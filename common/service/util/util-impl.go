package serviceutil

import (
	"context"
	"encoding/json"

	"github.com/Calmantara/go-user/common/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UtilServiceImpl struct{}

func NewUtilService() UtilService {
	return &UtilServiceImpl{}
}

func (c *UtilServiceImpl) ContextBackground(ctx context.Context) (result context.Context) {
	result = context.Background()
	result = context.WithValue(result, logger.CorrelationKey.String(), ctx.Value(logger.CorrelationKey.String()))
	return result
}

func (c *UtilServiceImpl) SetCorrelationIdFromHeader(ctx *gin.Context) {
	corr := ctx.GetHeader(logger.CorrelationKey.String())
	if corr == "" {
		corr = uuid.New().String()
		val, exist := ctx.Get(logger.CorrelationKey.String())
		if exist {
			corr = val.(string)
		}
	}
	ctx.Set(logger.CorrelationKey.String(), corr)
}

func (c *UtilServiceImpl) UpsertCorrelationId(ctx context.Context, corrUid ...string) (ctxResult context.Context) {
	// check correlation id
	corr := (ctx).Value(logger.CorrelationKey.String())
	if corr == nil && len(corrUid) <= 0 {
		corr = uuid.New().String()
	} else if len(corrUid) > 0 {
		corr = corrUid[0]
	}
	(ctxResult) = context.WithValue(ctx, logger.CorrelationKey.String(), corr)
	return ctxResult
}

func (c *UtilServiceImpl) ObjectMapper(source interface{}, destination interface{}) (err error) {
	// encode
	byteObject, err := json.Marshal(&source)
	if err != nil {
		return err
	}
	// decode to object
	err = json.Unmarshal(byteObject, &destination)
	return err
}
