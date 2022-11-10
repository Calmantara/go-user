package user

import "github.com/gin-gonic/gin"

type UserHdl interface {
	GetUsersHdl(ctx *gin.Context)
	GetUserByIdHdl(ctx *gin.Context)
	UpdateUserHdl(ctx *gin.Context)
	InsertUserHdl(ctx *gin.Context)
}
