package controllers

import (
	"context"
	"time"

	"idealist/app/crud"
	"idealist/app/models"
	"idealist/pkg/repository"
	"idealist/pkg/utils"
	"idealist/platform/cache"
	"idealist/platform/database"

	"github.com/gofiber/fiber/v2"
)

// UserSignUp method to create a new user.
// @Description Create a new user given username, email and password
// @Summary Create a new user
// @Tags User
// @Accept json
// @Produce json
// @Param username body string true "User Username"
// @Param email body string true "User Email"
// @Param password body string true "User Password"
// @Success 201 {string} status "ok"
// @Router /v1/user/sign/up [post]
func UserSignUp(c *fiber.Ctx) error {
	signUp := &crud.SignUp{}

	// Checking received data from JSON body.
	if err := c.BodyParser(signUp); err != nil {
		// Return status 400 and error message.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"code":  repository.INVALID_DATA,
			"msg":   err.Error(),
		})
	}

	// Create database connection.
	db := database.DB

	// Generate password hash.
	hashedPassword := utils.GeneratePassword(signUp.Password)

	// Create a new user.
	user := &models.User{
		Username: signUp.Username,
		Password: hashedPassword,
		Email:    signUp.Email,
		Verified: false,
		Member:   models.Member{},
	}

	// Create a new user in database.
	err := db.Create(&user).Error
	if err != nil {
		// Return status 500 and database connection error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"code":  repository.DATABASE_ERROR,
			"msg":   err.Error(),
		})
	}

	// Return status 201 and created user.
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"error": false,
		"code":  repository.REGISTERED,
		"msg":   nil,
		"user":  user,
	})
}

// UserSignIn method to auth user and return access and refresh tokens.
// @Description Log In user and return access and refresh token.
// @Summary User Sign In
// @Tags User
// @Accept json
// @Produce json
// @Param username body string true "User's username"
// @Param password body string true "User's password'"
// @Success 200 {string} status "ok"
// @Router /v1/user/sign/in [post]
func UserSignIn(c *fiber.Ctx) error {
	// Create a new user auth struct.
	signIn := &crud.SignIn{}

	// Checking received data from JSON body.
	if err := c.BodyParser(signIn); err != nil {
		// Return status 400 and error message.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"code":  repository.INVALID_DATA,
			"msg":   err.Error(),
		})
	}

	// Create database connection.
	db := database.DB

	// Get user by username.
	var foundedUser models.User
	err := db.Model(models.User{Username: signIn.Username}).First(&foundedUser).Error
	if err != nil {
		// Return, if user not found.
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": true,
			"code":  repository.USER_NOT_FOUND,
			"msg":   "User with the given name is not found",
		})
	}

	// Compare given user password with stored in found user.
	compareUserPassword := utils.ComparePasswords(foundedUser.Password, signIn.Password)
	if !compareUserPassword {
		// Return, if password is not compare to stored in database.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"code":  repository.USER_NOT_FOUND,
			"msg":   "Incorrect credentials",
		})
	}

	// Generate a new pair of access and refresh tokens.
	tokens, err := utils.GenerateNewTokens(foundedUser.ID)
	if err != nil {
		// Return status 500 and token generation error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"code":  repository.ERROR_RETREIVING_TOKEN,
			"msg":   err.Error(),
		})
	}

	// Define user ID.
	userUsername := foundedUser.Username

	// Create a new Redis connection.
	connRedis, err := cache.RedisConnection()
	if err != nil {
		// Return status 500 and Redis connection error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"code":  repository.CACHE_ERROR,
			"msg":   err.Error(),
		})
	}

	// Save refresh token to Redis.
	errSaveToRedis := connRedis.Set(context.Background(), userUsername, tokens.Refresh, 0).Err()
	if errSaveToRedis != nil {
		// Return status 500 and Redis connection error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"code":  repository.CACHE_ERROR,
			"msg":   errSaveToRedis.Error(),
		})
	}

	// Return status 200 OK.
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error": false,
		"msg":   "Successfully signed in",
		"tokens": fiber.Map{
			"access":  tokens.Access,
			"refresh": tokens.Refresh,
		},
	})
}

// UserSignOut De-authorize user and delete refresh token from cache.
// @Description De-authorize user and delete refresh token from cache.
// @Summary De-authorize user
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {string} status "ok"
// @Security ApiKeyAuth
// @Router /v1/user/sign/out [post]
func UserSignOut(c *fiber.Ctx) error {
	// Get claims from JWT.
	claims, err := utils.ExtractTokenMetadata(c)
	if err != nil {
		// Return status 500 and JWT parse error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"code":  repository.INVALID_DATA,
			"msg":   err.Error(),
		})
	}

	// Define user ID.
	userID := claims.UserID

	// Create a new Redis connection.
	connRedis, err := cache.RedisConnection()
	if err != nil {
		// Return status 500 and Redis connection error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"code":  repository.CACHE_ERROR,
			"msg":   err.Error(),
		})
	}

	// Save refresh token to Redis.
	errDelFromRedis := connRedis.Del(context.Background(), userID).Err()
	if errDelFromRedis != nil {
		// Return status 500 and Redis deletion error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"code":  repository.CACHE_ERROR,
			"msg":   errDelFromRedis.Error(),
		})
	}

	// Return status 200, sign out successfull.
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error": false,
		"code":  repository.LOGGED_OUT,
		"msg":   "User logged out",
	})
}

// RenewTokens method for renew access and refresh tokens.
// @Description Renew access token
// @Summary Renew access and refresh tokens
// @Tags Token
// @Accept json
// @Produce json
// @Param refresh_token body string true "Refresh token"
// @Success 200 {string} status "ok"
// @Security ApiKeyAuth
// @Router /v1/token/renew [post]
func RenewTokens(c *fiber.Ctx) error {
	// Get now time.
	now := time.Now().Unix()

	// Get claims from JWT.
	claims, err := utils.ExtractTokenMetadata(c)
	if err != nil {
		// Return status 500 and JWT parse error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"code":  repository.ERROR_RETREIVING_TOKEN,
			"msg":   err.Error(),
		})
	}

	// Set expiration time from JWT data of current user.
	expiresAccessToken := claims.Expires

	// Checking, if now time greather than Access token expiration time.
	if now > expiresAccessToken {
		// Return status 401 and unauthorized error message.
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"code":  repository.ERROR_RETREIVING_TOKEN,
			"msg":   "Unauthorized, check expiration time of your token",
		})
	}

	// Create a new renew refresh token struct.
	renew := &crud.Renew{}

	// Checking received data from JSON body.
	if err := c.BodyParser(renew); err != nil {
		// Return, if JSON data is not correct.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"code":  repository.INVALID_DATA,
			"msg":   err.Error(),
		})
	}

	// Set expiration time from Refresh token of current user.
	expiresRefreshToken, err := utils.ParseRefreshToken(renew.RefreshToken)
	if err != nil {
		// Return status 400 and error message.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"code":  repository.ERROR_RETREIVING_TOKEN,
			"msg":   err.Error(),
		})
	}
	// Checking, if now time greather than Refresh token expiration time.
	if now >= expiresRefreshToken {

		// Return status 401 and unauthorized error message.
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"code":  repository.ERROR_RETREIVING_TOKEN,
			"msg":   "unauthorized, your session was ended earlier",
		})
	}

	// Define user ID.
	userID := claims.UserID

	// Create database connection.
	db := database.DB

	// Get user by ID.
	var foundedUser models.User
	err = db.Model(models.User{ID: userID}).First(&foundedUser).Error
	if err != nil {
		// Return, if user not found.
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": true,
			"code":  repository.USER_NOT_FOUND,
			"msg":   "User with the given ID is not found",
		})
	}

	// Generate JWT Access & Refresh tokens.
	tokens, err := utils.GenerateNewTokens(userID)
	if err != nil {
		// Return status 500 and token generation error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"code":  repository.ERROR_RETREIVING_TOKEN,
			"msg":   err.Error(),
		})
	}

	// Create a new Redis connection.
	connRedis, err := cache.RedisConnection()
	if err != nil {
		// Return status 500 and Redis connection error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"code":  repository.CACHE_ERROR,
			"msg":   err.Error(),
		})
	}

	// Save refresh token to Redis.
	errRedis := connRedis.Set(context.Background(), userID, tokens.Refresh, 0).Err()
	if errRedis != nil {
		// Return status 500 and Redis connection error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"code":  repository.CACHE_ERROR,
			"msg":   errRedis.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"error": false,
		"code":  repository.TOKEN_RENEWED,
		"msg":   "Success regenerating token data",
		"tokens": fiber.Map{
			"access":  tokens.Access,
			"refresh": tokens.Refresh,
		},
	})
}
