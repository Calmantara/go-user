package main

import (
	"encoding/json"

	"github.com/Calmantara/go-user/common/infra/gorm/transaction"
	"github.com/Calmantara/go-user/common/logger"
	"github.com/Calmantara/go-user/common/setup/config"
	"github.com/Calmantara/go-user/pkg/domain/user"
	"go.uber.org/dig"

	ginrouter "github.com/Calmantara/go-user/common/infra/gin/router"
	confgorm "github.com/Calmantara/go-user/common/infra/gorm"
	"github.com/Calmantara/go-user/common/middleware/cors"
	serviceutil "github.com/Calmantara/go-user/common/service/util"
	userrepo "github.com/Calmantara/go-user/pkg/repository/user"
	"github.com/Calmantara/go-user/pkg/server/http/handler/auth"
	userhdl "github.com/Calmantara/go-user/pkg/server/http/handler/user"
	userrouter "github.com/Calmantara/go-user/pkg/server/http/router/v1/user"
	userusecase "github.com/Calmantara/go-user/pkg/usecase/user"
)

// initiate all grouped DI
func commonDependencies() []any {
	return []any{logger.NewCustomLogger, config.NewConfigSetup,
		serviceutil.NewUtilService, ginrouter.NewGinRouter}
}

func svcDependencies() []any {
	return []any{userusecase.NewUserUsecase}
}

func handlerDependencies() []any {
	return []any{userhdl.NewUserHdl}

}

func routerDependencies() []any {
	return []any{userrouter.NewUserRouter}
}

func BuildRepoDependencies(sugar logger.CustomLogger, conf config.ConfigSetup) (transaction.Transaction, user.UserRepo) {
	readCln := confgorm.NewPostgresConfig(sugar, conf, confgorm.WithPostgresMode("read"))
	writeCln := confgorm.NewPostgresConfig(sugar, conf, confgorm.WithPostgresMode("write"))
	return transaction.NewTransaction(sugar, readCln), userrepo.NewUserRepo(sugar, readCln, writeCln)
}

func BuildInRuntime() (serviceConf map[string]any, ginRouter ginrouter.GinRouter, err error) {
	c := dig.New()
	// define all generic
	var constructor []any
	constructor = append(constructor, BuildRepoDependencies)
	constructor = append(constructor, commonDependencies()...)
	constructor = append(constructor, svcDependencies()...)
	constructor = append(constructor, handlerDependencies()...)
	constructor = append(constructor, routerDependencies()...)

	// provide all generic
	for _, service := range constructor {
		if err := c.Provide(service); err != nil {
			return nil, nil, err
		}
	}
	if err = c.Invoke(func(config config.ConfigSetup, gn ginrouter.GinRouter, userRouter userrouter.UserRouter) {
		// service information
		app, _ := json.Marshal(config.GetRawConfig()["service"])
		// unmarshal config
		json.Unmarshal(app, &serviceConf)
		// init server
		ginRouter = gn
		// init middleware
		gn.USE(cors.NewCorsMiddleware().Cors, auth.AuthStatic)
		// init routers
		userRouter.Routers()
	}); err != nil {
		panic(err)
	}
	return serviceConf, ginRouter, err
}
