package services

import (
	"os"

	"github.com/AliTr404/T-MO/internal/dto"
	"github.com/AliTr404/T-MO/internal/http/exception"
	model "github.com/AliTr404/T-MO/internal/models"
	repository "github.com/AliTr404/T-MO/internal/repositories"
	"github.com/AliTr404/T-MO/pkg/derrors"
	"github.com/AliTr404/T-MO/pkg/dtime"
	"github.com/AliTr404/T-MO/pkg/jwt"
	"github.com/AliTr404/T-MO/pkg/mail"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	AuthService interface {
		SignUp(user *model.User) (*dto.User, error)
		SignIn(signInDTO *dto.SignIn) (*dto.User, error)
		ForgotPassword(email string) error
		ChangePassword(changePasswordDTO *dto.ChangePassword) error
	}
	authService struct {
		userRepository repository.UserRepository
	}
)

func NewAuthService(userRepository repository.UserRepository) AuthService {
	return &authService{
		userRepository: userRepository,
	}
}

func (s *authService) SignUp(user *model.User) (*dto.User, error) {
	if _, err := s.userRepository.FindByUsername(user.Username); err == nil {
		return nil, derrors.New(derrors.KindConflict, "نام کاربری وارد شده از قبل در دیتابیس موجود است !")
	}

	if _, err := s.userRepository.FindByEmail(user.Email); err == nil {
		return nil, derrors.New(derrors.KindConflict, "ایمیل وارد شده از قبل در دیتابیس موجود است !")
	}

	password, err := user.HashPassword()
	if err != nil {
		return nil, derrors.New(derrors.KindInvalid, "مشکلی رخ داده !")
	}
	user.Password = string(password)

	id, err := s.userRepository.Insert(user)

	if err != nil {
		return nil, derrors.New(derrors.KindInvalid, "مشکلی رخ داده !")
	}

	return &dto.User{
		ID:       id,
		SID:      user.SID,
		FullName: user.FullName,
		Username: user.Username,
		Email:    user.Email,
	}, nil
}

func (s *authService) SignIn(signInDTO *dto.SignIn) (*dto.User, error) {
	user, err := s.userRepository.FindByEmailOrUsername(signInDTO.UsernameOrEmail)
	if err != nil {
		return nil, derrors.New(derrors.KindUnauthorized, "مشخصات ورود غیر معتبر است !")
	}
	if err := user.ValidatePassword(signInDTO.Password); err != nil {
		return nil, derrors.New(derrors.KindUnauthorized, "مشخصات ورود غیر معتبر است !")
	}

	return &dto.User{
		ID:       user.ID,
		SID:      user.SID,
		FullName: user.FullName,
		Username: user.Username,
		Email:    user.Email,
	}, nil
}

func (s *authService) ForgotPassword(email string) error {
	user, err := s.userRepository.FindByEmail(email)
	if err != nil {
		return exception.NotFoundException("یمیل وارد شده در دیتابیس موجود نمی باشد.")
	}

	token := jwt.GenerateForgotPasswordToken(user.ID.Hex())

	mail.NewMail([]string{email}, "فراموشی رمز عبور | Rythmik").Send("templates/forgot-password.html", map[string]string{"Url": os.Getenv("FRONTEND_HOST") + "/auth/change-password/?email=" + user.Email + "&token=" + token})

	return nil
}

func (s *authService) ChangePassword(changePasswordDTO *dto.ChangePassword) error {
	res, err := s.userRepository.FindByEmail(changePasswordDTO.Email)
	if err != nil {
		return exception.NotFoundException([]dto.ApiError{{Param: "Email", Message: "ایمیل وارد شده در دیتابیس موجود نمی باشد."}})
	}
	res.Password = changePasswordDTO.Password
	pass, err := res.HashPassword()

	if err != nil {
		return exception.InternalServerException()
	}
	err = s.userRepository.UpdateOne(bson.M{"email": changePasswordDTO.Email}, bson.D{
		primitive.E{Key: "$set", Value: bson.D{primitive.E{Key: "password", Value: string(pass)}, primitive.E{Key: "updated_at", Value: dtime.Now()}}},
	})
	if err != nil {
		return exception.InternalServerException()
	}

	return nil
}
