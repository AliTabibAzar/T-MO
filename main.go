package main

import (
	"log"
	"os"

	"github.com/AliTr404/T-MO/internal/config"
	"github.com/AliTr404/T-MO/internal/db"
	"github.com/AliTr404/T-MO/internal/http/controllers"
	"github.com/AliTr404/T-MO/internal/http/exception"
	mw "github.com/AliTr404/T-MO/internal/http/middleware"
	"github.com/AliTr404/T-MO/internal/repositories"
	"github.com/AliTr404/T-MO/internal/routes"
	"github.com/AliTr404/T-MO/internal/services"
	"github.com/AliTr404/T-MO/pkg/jwt"
	"github.com/AliTr404/T-MO/pkg/validation"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	configuration config.Config = config.New(".env")
	mongoClient   *mongo.Client = db.SetupMongo()

	// setup repositories
	userRepository    repositories.UserRepository    = repositories.NewUserRepository(mongoClient)
	videoRepository   repositories.VideoRepository   = repositories.NewVideoRepository(mongoClient)
	commentRepository repositories.CommentRepository = repositories.NewCommentRepository(mongoClient)
	// setup services
	authService   services.AuthService   = services.NewAuthService(userRepository)
	videoService  services.VideoService  = services.NewVideoService(videoRepository, commentRepository)
	uploadService services.UploadService = services.NewUploadService(userRepository, videoRepository)
	homeService   services.HomeService   = services.NewHomeService(videoRepository)

	// setup controllers
	authController   controllers.AuthController   = controllers.NewAuthController(authService)
	videoController  controllers.VideoController  = controllers.NewVideoController(videoService)
	uploadController controllers.UploadController = controllers.NewUploadController(uploadService)
	homeController   controllers.HomeController   = controllers.NewHomeController(homeService)
)

func main() {
	logFile, err := os.Create("./logs/.log")
	if err != nil {
		log.Fatal("cannot create the log file")
	}
	defer func() {
		err := logFile.Close()
		if err != nil {
			log.Printf("failed to load log file: %s ", err)
		}
	}()
	server := echo.New()
	server.Static("/files", "upload/public")

	// api middlewares
	server.Use(middleware.CORS())
	server.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Output: logFile,
	}))

	server.Validator = validation.NewValidation()
	server.HTTPErrorHandler = exception.CustomHTTPErrorHandler

	// jwt managers
	accessJwt := jwt.NewJwtManager(os.Getenv("JWT_ACCESS_SECRET_KEY"))
	// api routes
	routes.AuthRoutes(server, authController)
	routes.VideoRoutes(server, videoController, mw.Authorization(accessJwt))
	routes.UploadRoutes(server, uploadController, mw.Authorization(accessJwt))
	routes.HomeRoutes(server, homeController)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	server.Logger.Fatal(server.Start(":" + port))
}
