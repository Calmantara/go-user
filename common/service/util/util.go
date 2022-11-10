//go:generate mockgen -source util.go -destination mock/util_mock.go -package mock

package serviceutil

import (
	"context"

	"github.com/gin-gonic/gin"
)

type UtilService interface {
	// correlation
	SetCorrelationIdFromHeader(ctx *gin.Context)
	UpsertCorrelationId(ctx context.Context, corrUid ...string) (ctxResult context.Context)
	// Mapping object can to struct
	ObjectMapper(source, destination any) (err error)
	// Generate new context background with correlation id
	ContextBackground(ctx context.Context) (result context.Context)
}
