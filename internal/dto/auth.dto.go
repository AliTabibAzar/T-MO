package dto

type (
	SignUp struct {
		Email    string `bson:"email" json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
		FullName string `bson:"fullname" json:"fullname" validate:"required"`
		Username string `bson:"username" json:"username" validate:"required"`
	}
	SignIn struct {
		UsernameOrEmail string `json:"usernameOrEmail" validate:"required"`
		Password        string `json:"password" validate:"required"`
	}

	ForgotPassword struct {
		Email string `bson:"email" json:"email" validate:"required,email"`
	}

	ChangePassword struct {
		Email    string `bson:"email" json:"email" validate:"required,email"`
		Password string `bson:"password" json:"password" validate:"required,min=8"`
	}
)
