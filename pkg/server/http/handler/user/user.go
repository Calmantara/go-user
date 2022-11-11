package user

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Calmantara/go-user/common/infra/gorm/transaction"
	"github.com/Calmantara/go-user/common/logger"
	serviceutil "github.com/Calmantara/go-user/common/service/util"
	"github.com/Calmantara/go-user/pkg/domain/response"
	"github.com/Calmantara/go-user/pkg/domain/user"
	"github.com/gin-gonic/gin"
)

type UserHdlImpl struct {
	sugar       logger.CustomLogger
	readTrx     transaction.Transaction
	util        serviceutil.UtilService
	userUsecase user.UserUsecase
}

func NewUserHdl(sugar logger.CustomLogger, readTrx transaction.Transaction, userUsecase user.UserUsecase, util serviceutil.UtilService) user.UserHdl {
	return &UserHdlImpl{sugar: sugar, readTrx: readTrx, userUsecase: userUsecase, util: util}
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

func (u *UserHdlImpl) GetUserByIdHdl(ctx *gin.Context) {
	ctx.Set(transaction.TRANSACTION_KEY.String(), u.readTrx.GormBeginTransaction(ctx))
	u.sugar.WithContext(ctx).Info("%T-InsertUserHdl is invoked", u)
	defer func() {
		// close transaction
		if errTx := u.readTrx.GormEndTransaction(ctx); errTx != nil {
			u.sugar.WithContext(ctx).Errorf("error when process transaction:%v", errTx)
		}
		u.sugar.WithContext(ctx).Info("%T-InsertUserHdl executed", u)
	}()
	u.util.SetCorrelationIdFromHeader(ctx)

	// get user id
	userId := ctx.Param("user_id")
	if userId == "" {
		err := fmt.Sprintf(string(response.MISSING_FIELD_MSG), "user_id")
		// bad request
		u.sugar.WithContext(ctx).Errorf("error whengetting user_id %+v", err)
		ctx.AbortWithStatusJSON(int(response.MISSING_FIELD_CODE), map[string]any{
			"error": err,
		})
		ctx.Set(transaction.TRANSACTION_ERROR_KEY.String(), err)
		return
	}
	// calling usecase
	uid, _ := strconv.Atoi(userId)
	userDet := user.User{
		Id: uint64(uid),
	}
	if errResp := u.userUsecase.GetUserByIdSvc(ctx, &userDet); errResp.Error != "" {
		u.sugar.WithContext(ctx).Errorf("error when process service %+v", errResp)
		ctx.AbortWithStatusJSON(int(errResp.Code), errResp)
		ctx.Set(transaction.TRANSACTION_ERROR_KEY.String(), errResp)
		return
	}
	// success
	ctx.JSON(http.StatusOK, userDet)
}

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
	u.util.SetCorrelationIdFromHeader(ctx)

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
