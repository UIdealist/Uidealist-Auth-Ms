package crud

// SignIn struct to describe login user.
type SignInCredentials struct {
	Username string `json:"username" validate:"required,username,lte=255"`
	Password string `json:"password" validate:"required,lte=255"`
}

// SignUp struct to describe register user.
type SignUpCredentials struct {
	Username string `json:"username" validate:"required,username,lte=255"`
	Password string `json:"password" validate:"required,lte=255"`
	Email    string `json:"email" validate:"required,email,lte=255"`
}

// Renew struct to describe refresh token object.
type Renew struct {
	RefreshToken string `json:"refresh_token"`
}
