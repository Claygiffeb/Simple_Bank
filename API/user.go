package api

import (
	"database/sql"
	"net/http"
	"time"

	db "github.com/Clayagiffeb/Simple_Bank/db/sqlc"
	"github.com/Clayagiffeb/Simple_Bank/util"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type createUser struct {
	Username string `json:"username" binding:"required,alphanum"` // username is required not to compose special characters
	Password string `json:"password" binding:"required,min=6"`    // minimum password length
	Fullname string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"` // check email to make sure that it's a valid email
}

type userResponse struct {
	Username        string    `json:"username"`
	FullName        string    `json:"full_name"`
	Email           string    `json:"email"`
	PasswordChanged time.Time `json:"password_changed_at"`
	CreatedAt       time.Time `json:"created_at"`
}

func newUserResponse(user db.User) userResponse {
	return userResponse{Username: user.Username,
		FullName:        user.FullName,
		Email:           user.Email,
		PasswordChanged: user.PasswordChanged,
		CreatedAt:       user.CreatedAt}
}
func (server *Server) CreateUser(ctx *gin.Context) {
	var req createUser
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	HashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.CreateUserParams{
		Username:       req.Username,
		HashedPassword: HashedPassword,
		FullName:       req.Fullname,
		Email:          req.Email,
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() { // handle the http status code in case of confliting by creating many currencies account for the same owner
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	userResponse := newUserResponse(user)
	ctx.JSON(http.StatusOK, userResponse)
}

type createUserResponse struct { // since we will not response the hashedpassword for the users, we will create a response struct (it kinda like dto for MVC Spring)
	Username        string    `json:"username"`
	FullName        string    `json:"full_name"`
	Email           string    `json:"email"`
	PasswordChanged time.Time `json:"password_changed"`
	CreatedAt       time.Time `json:"created_at"`
}
type loginUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"` // username is required not to compose special characters
	Password string `json:"password" binding:"required,min=6"`    // minimum password length
}

type loginUserResponse struct {
	SessionID             uuid.UUID    `json:"session_id"`
	AcessToken            string       `json:"acces_token"`
	AccessTokenExpiresAt  time.Time    `json:"access_token_expires_at"`
	RefreshToken          string       `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time    `json:"refresh_token_expires_at"`
	User                  userResponse `json:"user"`
}

func (server *Server) loginUser(ctx *gin.Context) {
	var req loginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUser(ctx, req.Username)

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	// check for the password
	err = util.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(
		user.Username,
		server.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(user.Username, server.config.RefreshTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	session, err := server.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshPayload.ID,
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    ctx.Request.UserAgent(),
		ClientIp:     ctx.ClientIP(),
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpiresAt.Time,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	respone := loginUserResponse{
		SessionID:             session.ID,
		AcessToken:            accessToken,
		AccessTokenExpiresAt:  accessPayload.ExpiresAt.Time,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshPayload.ExpiresAt.Time,
		User:                  newUserResponse(user),
	}

	ctx.JSON(http.StatusOK, respone)
}
