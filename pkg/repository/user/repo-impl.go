package user

import (
	"context"
	"fmt"

	"github.com/Calmantara/go-user/lib/logger"
	creditcard "github.com/Calmantara/go-user/pkg/domain/credit-card"
	"github.com/Calmantara/go-user/pkg/domain/photo"
	"github.com/Calmantara/go-user/pkg/domain/user"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	confgorm "github.com/Calmantara/go-user/lib/infra/gorm"
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

func (u *UserRepoImpl) GetUsersRepo(ctx context.Context, userQuery user.UserQuery, users *[]*user.User) (err error) {
	u.sugar.WithContext(ctx).Infof("execute %T GetUsersRepo", u)
	defer u.sugar.WithContext(ctx).Infof("%T GetUsersRepo executed", u)

	// generate transaction
	txn := u.readCln.GenerateTransaction(ctx)
	if userQuery.Q != "" {
		txn = txn.Where("name like ? or email like ?",
			userQuery.Q, userQuery.Q)
	}

	txn.Model(user.User{}).
		Select("id, name, email, address").
		Preload("Photos",
			func(db *gorm.DB) *gorm.DB {
				return db.Select("name, user_id")
			}).
		Preload("CreditCardToken",
			func(db *gorm.DB) *gorm.DB {
				return db.Select("token, user_id")
			}).Limit(userQuery.Lt).
		Offset(userQuery.Of).
		// TODO: fmt.sprintf can be replaced with string builder
		Order(fmt.Sprintf("%v %v", userQuery.Ob, userQuery.Sb)).
		Find(&users)
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
		Select("id, name, email, password, address").
		Preload("Photos", func(db *gorm.DB) *gorm.DB {
			return db.Select("name, user_id")
		}).
		Preload("CreditCardToken", func(db *gorm.DB) *gorm.DB {
			return db.Select("token, user_id")
		}).
		Where("id = ?", userDet.Id).
		Find(userDet)
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
		Clauses(
			clause.OnConflict{
				Columns:   []clause.Column{{Name: "id"}},
				UpdateAll: true,
			},
		).
		Create(userDet)

	// slightly bad approach
	txn.Delete(&creditcard.CreditCardToken{}, "user_id = ? and id <> ?",
		userDet.Id, userDet.CreditCardToken.Id).Debug()

	// delete photo
	if userDet.Photos != nil {
		var photoId []uint64
		for _, val := range userDet.Photos {
			photoId = append(photoId, val.ID)
		}
		txn.Delete(&photo.Photo{}, "user_id = ? and id not in ?",
			userDet.Id, photoId).Debug()
	}

	if err = txn.Error; err != nil {
		u.sugar.WithContext(ctx).Errorf("error execute upsert user:%v", err.Error())
	}
	return err
}
