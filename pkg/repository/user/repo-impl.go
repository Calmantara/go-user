package user

import (
	"context"

	"github.com/Calmantara/go-user/common/logger"
	creditcard "github.com/Calmantara/go-user/pkg/domain/credit-card"
	"github.com/Calmantara/go-user/pkg/domain/photo"
	"github.com/Calmantara/go-user/pkg/domain/user"

	confgorm "github.com/Calmantara/go-user/common/infra/gorm"
)

type UserRepoImpl struct {
	sugar logger.CustomLogger
	// conf config
	readCln  confgorm.PostgresConfig
	writeCln confgorm.PostgresConfig
}

func NewUserRepo(sugar logger.CustomLogger, readCln confgorm.PostgresConfig, writeCln confgorm.PostgresConfig) user.UserRepo {
	sugar.Logger().Info("init user repo. . .")
	return &UserRepoImpl{sugar: sugar, readCln: readCln, writeCln: writeCln}
}

func (u *UserRepoImpl) GetUsersRepo(ctx context.Context, userQuery user.UserQuery, users []*user.User) (err error) {
	u.sugar.WithContext(ctx).Infof("execute %T GetUsersRepo", u)
	defer u.sugar.WithContext(ctx).Infof("%T GetUsersRepo executed", u)

	// generate transaction
	txn := u.readCln.GenerateTransaction(ctx)
	txn.Model(user.User{}).
		Select("*").
		Find(users)
	if err = txn.Error; err != nil {
		u.sugar.WithContext(ctx).Errorf("error execute ReadWallet:%v", err.Error())
	}
	return err
}
func (u *UserRepoImpl) GetUserByIdRepo(ctx context.Context, userDet *user.User) (err error) {
	u.sugar.WithContext(ctx).Infof("execute %T GetUserByIdRepo", u)
	defer u.sugar.WithContext(ctx).Infof("%T GetUserByIdRepo executed", u)

	// generate transaction
	txn := u.readCln.GenerateTransaction(ctx)
	txn.Model(user.User{}).
		Select("*").
		Where("id = ?", userDet.Id).
		Find(userDet)
	if err = txn.Error; err != nil {
		u.sugar.WithContext(ctx).Errorf("error execute ReadWallet:%v", err.Error())
	}
	return err
}
func (u *UserRepoImpl) UpdateUserRepo(ctx context.Context, userDet *user.User) (err error) {
	u.sugar.WithContext(ctx).Infof("execute %T UpdateUserRepo", u)
	defer u.sugar.WithContext(ctx).Infof("%T UpdateUserRepo executed", u)
	// generate transaction
	txn := u.readCln.GenerateTransaction(ctx)
	txn.Model(user.User{}).
		Create(userDet)

	// dirty code
	// 1. input photos
	txn.Model(photo.Photo{}).
		Create(userDet.Photos)

	// 2. input credit card
	txn.Model(creditcard.CreditCard{}).
		Create(userDet.CreditCard)

	if err = txn.Error; err != nil {
		u.sugar.WithContext(ctx).Errorf("error execute ReadWallet:%v", err.Error())
	}
	return err
}
func (u *UserRepoImpl) InsertUserRepo(ctx context.Context, userDet *user.User) (err error) {
	u.sugar.WithContext(ctx).Infof("execute %T InsertUserRepo", u)
	defer u.sugar.WithContext(ctx).Infof("%T InsertUserRepo executed", u)
	// generate transaction
	txn := u.readCln.GenerateTransaction(ctx)
	txn.Model(user.User{}).
		Create(userDet).Debug()

	// dirty code
	// 1. input photos
	txn.Model(photo.Photo{}).
		Create(userDet.Photos)

	// 2. input credit card
	txn.Model(creditcard.CreditCard{}).
		Create(userDet.CreditCard)

	if err = txn.Error; err != nil {
		u.sugar.WithContext(ctx).Errorf("error execute ReadWallet:%v", err.Error())
	}
	return err
}
