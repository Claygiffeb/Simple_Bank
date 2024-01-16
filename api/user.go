package api

import (
	"net/http"
	"time"

	db "github.com/Clayagiffeb/Simple_Bank/db/sqlc"
	"github.com/Clayagiffeb/Simple_Bank/util"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type createUser struct {
	Username string `json:"username" binding:"required,alphanum"` // username is required not to compose special characters
	Password string `json:"password" binding:"required,min=6"`    // minimum password length
	Fullname string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"` // check email to make sure that it's a valid email
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
	userResponse := createUserResponse{
		Username:        user.Username,
		FullName:        user.FullName,
		Email:           user.Email,
		PasswordChanged: user.PasswordChanged,
		CreatedAt:       user.CreatedAt,
	}
	ctx.JSON(http.StatusOK, userResponse)
}

type createUserResponse struct { // since we will not response the hashedpassword for the users, we will create a response struct (it kinda like dto for MVC Spring)
	Username        string    `json:"username"`
	FullName        string    `json:"full_name"`
	Email           string    `json:"email"`
	PasswordChanged time.Time `json:"password_changed"`
	CreatedAt       time.Time `json:"created_at"`
}
