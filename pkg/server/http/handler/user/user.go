package user

import (
	"net/http"

	"github.com/Calmantara/go-user/common/infra/gorm/transaction"
	"github.com/Calmantara/go-user/common/logger"
	"github.com/Calmantara/go-user/pkg/domain/response"
	"github.com/Calmantara/go-user/pkg/domain/user"
	"github.com/gin-gonic/gin"
)

type UserHdlImpl struct {
	sugar       logger.CustomLogger
	readTrx     transaction.Transaction
	userUsecase user.UserUsecase
}

func NewUserHdl(sugar logger.CustomLogger, readTrx transaction.Transaction, userUsecase user.UserUsecase) user.UserHdl {
	return &UserHdlImpl{sugar: sugar, readTrx: readTrx, userUsecase: userUsecase}
}

func (u *UserHdlImpl) GetUsersHdl(ctx *gin.Context) {
	ctx.Set(transaction.TRANSACTION_KEY.String(), u.readTrx.GormBeginTransaction(ctx))
	u.sugar.WithContext(ctx).Info("%T-GetUsersHdl is invoked", u)
	defer func() {
		// close transaction
		if errTx := u.readTrx.GormEndTransaction(ctx); errTx != nil {
			u.sugar.WithContext(ctx).Errorf("error when process transaction:%v", errTx)
		}
		u.sugar.WithContext(ctx).Info("%T-GetUsersHdl executed", u)
	}()
	// calling service
	var users []*user.User
	if errResp := u.userUsecase.GetUsersSvc(ctx, users); errResp.Error != "" {
		u.sugar.WithContext(ctx).Errorf("error when process service %+v", errResp)
		ctx.AbortWithStatusJSON(int(errResp.Code), errResp)
		ctx.Set(transaction.TRANSACTION_ERROR_KEY.String(), errResp)
		return
	}

	// success
	ctx.JSON(http.StatusOK, map[string]any{
		"count": len(users),
		"rows":  users,
	})
}

func (u *UserHdlImpl) GetUserByIdHdl(ctx *gin.Context) {}

func (u *UserHdlImpl) UpdateUserHdl(ctx *gin.Context) {}

func (u *UserHdlImpl) InsertUserHdl(ctx *gin.Context) {
	ctx.Set(transaction.TRANSACTION_KEY.String(), u.readTrx.GormBeginTransaction(ctx))
	u.sugar.WithContext(ctx).Info("%T-InsertUserHdl is invoked", u)
	defer func() {
		// close transaction
		if errTx := u.readTrx.GormEndTransaction(ctx); errTx != nil {
			u.sugar.WithContext(ctx).Errorf("error when process transaction:%v", errTx)
		}
		u.sugar.WithContext(ctx).Info("%T-InsertUserHdl executed", u)
	}()
	// binding
	var userDet user.User
	if err := ctx.ShouldBind(&userDet); err != nil {
		u.sugar.WithContext(ctx).Errorf("error when biding %+v", err)
		ctx.AbortWithStatusJSON(int(response.MISSING_FIELD_CODE), map[string]any{
			"error": err.Error(),
		})
		ctx.Set(transaction.TRANSACTION_ERROR_KEY.String(), err)
		return
	}

	// calling service
	if errResp := u.userUsecase.InsertUserSvc(ctx, &userDet); errResp.Error != "" {
		u.sugar.WithContext(ctx).Errorf("error when process service %+v", errResp)
		ctx.AbortWithStatusJSON(int(errResp.Code), errResp)
		ctx.Set(transaction.TRANSACTION_ERROR_KEY.String(), errResp)
		return
	}

	// success
	ctx.JSON(http.StatusOK, map[string]any{
		"user_id": userDet.Id,
	})
}
