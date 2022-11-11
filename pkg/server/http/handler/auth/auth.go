package auth

import (
	"github.com/Calmantara/go-user/pkg/domain/response"
	"github.com/gin-gonic/gin"
)

const (
	staticKey string = "HiJhvL$T27@1u^%u86g"
)

func AuthStatic(c *gin.Context) {
	// check header
	key := c.GetHeader("key")
	// check if key exist or not
	if key == "" {
		c.AbortWithStatusJSON(int(response.MISSING_API_CODE), response.ErrorResponse{
			Error: response.INVALID_API_KEY_MSG,
		})
		return
	}
	// check if same or not
	if key != staticKey {
		c.AbortWithStatusJSON(int(response.INVALID_API_KEY_CODE), response.ErrorResponse{
			Error: response.INVALID_API_KEY_MSG,
		})
		return
	}
	c.Next()
}
