package user

import (
	gingroup "github.com/Calmantara/go-user/lib/infra/gin/group"
	ginrouter "github.com/Calmantara/go-user/lib/infra/gin/router"
	"github.com/Calmantara/go-user/pkg/domain/user"
)

type UserRouter interface {
	Routers()
}

type UserRouterImpl struct {
	group   gingroup.GinGroup
	userhdl user.UserHdl
}

func NewUserRouter(ginRouter ginrouter.GinRouter, userhdl user.UserHdl) UserRouter {
	group := ginRouter.GROUP("api/v1/user")
	return &UserRouterImpl{group: group, userhdl: userhdl}
}

func (w *UserRouterImpl) get() {
	w.group.GET("/:user_id", w.userhdl.GetUserByIdHdl)
	w.group.GET("/list", w.userhdl.GetUsersHdl)
}

func (w *UserRouterImpl) post() {
	w.group.POST("/register", w.userhdl.InsertUserHdl)
	w.group.POST("", w.userhdl.UpdateUserHdl)
}

func (w *UserRouterImpl) Routers() {
	w.get()
	w.post()
}
