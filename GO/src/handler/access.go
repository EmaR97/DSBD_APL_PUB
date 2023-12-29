package handler

import (
	"CamMonitoring/src/entity"
	"CamMonitoring/src/repository"
	"CamMonitoring/src/utility"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"time"
)

const (
	tokenLength     = 32 // Adjust the length of the token as needed
	cookieTokenName = "token"
	cookieUserName  = "user"
)

// AccessHandler handles user authentication, registration, and token management
type AccessHandler struct {
	userRepository  *repository.MongoDBRepository[entity.User]
	tokenRepository *repository.TokenRepository
}

// NewAccessHandler creates a new AccessHandler with the provided MqttClient
func NewAccessHandler(
	userRepository *repository.MongoDBRepository[entity.User],
	tokenRepository *repository.MongoDBRepository[entity.Token],
) *AccessHandler {
	//userRepository.Add(entity.User{Id: "username", Password: "password", Email: "email"})

	return &AccessHandler{
		userRepository:  userRepository,
		tokenRepository: &repository.TokenRepository{MongoDBRepository: *tokenRepository},
	}
}

// LoginPost handles user authentication
func (h *AccessHandler) LoginPost(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	camId := c.PostForm("cam_id")

	id, err := h.authenticateUser(username, password, camId)
	if err != nil {
		utility.WarningLog().Printf("Login failed for user %s: %s", username, err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token, err := h.generateAndSetToken(c, username)

	if err != nil {
		utility.ErrorLog().Printf("Error generating and setting token for user %s: %s", username, err)
		c.Status(http.StatusInternalServerError)
		return
	}
	utility.InfoLog().Printf("Login successful for user %s, token: %s", id, token.Id)
	c.JSON(http.StatusOK, gin.H{"token": token.Id})
}

// authenticateUser performs user authentication
func (h *AccessHandler) authenticateUser(username, password, camID string) (string, error) {
	user, err := h.userRepository.GetByID(username)
	if err != nil {
		return "", fmt.Errorf("user not found")
	}
	toCheck := user.Password
	idToReturn := username

	if camID != "" {
		toCheck = user.CamPassword
		idToReturn = camID
	}

	if password != toCheck {
		return "", fmt.Errorf("invalid password")
	}

	return idToReturn, nil
}

// SignupPost handles user registration
func (h *AccessHandler) SignupPost(c *gin.Context) {
	username := c.PostForm("username")
	email := c.PostForm("email")
	password := c.PostForm("password")

	// Validate input (replace this with your actual validation logic)
	if username == "" || email == "" || password == "" {
		utility.WarningLog().Printf("Invalid input for user registration: %s", username)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	user := entity.User{
		Id:          username,
		Email:       email,
		Password:    password,
		CamPassword: uuid.New().String(),
	}
	if err := h.userRepository.Add(user); err != nil {
		utility.ErrorLog().Printf("Error adding user %s: %s", username, err)
		c.Status(http.StatusInternalServerError)
		return
	}
	token, err := h.generateAndSetToken(c, username)

	if err != nil {
		utility.ErrorLog().Printf("Error generating and setting token for user %s: %s", username, err)
		c.Status(http.StatusInternalServerError)
		return
	}
	utility.InfoLog().Printf("User registered successfully: %s, token: %s", username, token.Id)
	c.JSON(http.StatusOK, gin.H{"Cam_pass": user.CamPassword})
}

// generateAndSetToken generates a token and stores it in the database
func (h *AccessHandler) generateAndSetToken(c *gin.Context, username string) (*entity.Token, error) {
	// Generate a token and store it in the database
	token, err := entity.NewToken(username, tokenLength)
	if err != nil {
		return nil, fmt.Errorf("error generating token: %s", err)
	}
	if err := h.tokenRepository.StoreToken(*token); err != nil {
		return nil, fmt.Errorf("error storing token: %s", err)
	}

	h.setCookies(c, token.Id, username)

	return token, nil
}

// setCookies sets the necessary cookies for the user
func (h *AccessHandler) setCookies(c *gin.Context, token, username string) {
	// Set cookies with one year expiration
	cookieOptions := http.Cookie{
		Path:     "/",
		Domain:   "localhost",
		Secure:   false,
		HttpOnly: true,
		SameSite: http.SameSiteDefaultMode,
	}

	// Set token cookie
	cookieOptions.Name = cookieTokenName
	cookieOptions.Value = token
	cookieOptions.Expires = time.Now().AddDate(1, 0, 0)
	c.SetCookie(
		cookieOptions.Name, cookieOptions.Value, int(cookieOptions.Expires.Unix()), cookieOptions.Path,
		cookieOptions.Domain, cookieOptions.Secure, cookieOptions.HttpOnly,
	)

	// Set username cookie
	cookieOptions.Name = cookieUserName
	cookieOptions.Value = username
	c.SetCookie(
		cookieOptions.Name, cookieOptions.Value, int(cookieOptions.Expires.Unix()), cookieOptions.Path,
		cookieOptions.Domain, cookieOptions.Secure, cookieOptions.HttpOnly,
	)
}

// Logout handles user logout
func (h *AccessHandler) Logout(c *gin.Context) {
	token, err := c.Cookie(cookieTokenName)
	if err != nil {
		utility.WarningLog().Printf("Logout failed: %s", err)
		c.String(http.StatusNotFound, "Cookie not found")
		return
	}
	// Remove the token from the database
	if err := h.tokenRepository.RemoveToken(token); err != nil {
		utility.ErrorLog().Printf("Error removing token for user: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error removing token"})
		return
	}

	h.clearCookies(c)
	utility.InfoLog().Println("User logged out successfully")
	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}

// clearCookies clears the cookies for the user
func (h *AccessHandler) clearCookies(c *gin.Context) {
	// Clear token cookie
	c.SetCookie(cookieTokenName, "", -1, "/", "localhost", false, true)

	// Clear username cookie
	c.SetCookie(cookieUserName, "", -1, "/", "localhost", false, true)
}

// VerifyToken checks the validity of the token
func (h *AccessHandler) VerifyToken(c *gin.Context) {
	// Retrieve the token from the request body
	var requestBody struct {
		Token string `form:"token" binding:"required"`
	}

	if err := c.ShouldBind(&requestBody); err != nil {
		utility.WarningLog().Printf("Token not provided: %s", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token not provided"})
		return
	}

	// Validate the token against the database
	if isValid, err := h.tokenRepository.IsValidToken(requestBody.Token); err != nil || !isValid {
		utility.WarningLog().Printf("Invalid or expired token: %s", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		return
	}

	utility.InfoLog().Println("Token is valid")
	c.JSON(http.StatusOK, gin.H{"message": "Token is valid"})
}
