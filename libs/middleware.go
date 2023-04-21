package libs

import (
	"github.com/gin-gonic/gin"

	Logger "github.com/cloudsrc/api.awaymail.v1.go/src/shared/logger"
	"github.com/cloudsrc/api.awaymail.v1.go/src/shared/utils"
	"github.com/cloudsrc/api.awaymail.v1.go/src/shared/utils/context"
)

func RequestIdMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		reqId := c.Request.Header.Get("X-Request-Id")
		if len(reqId) == 0 {
			reqId = utils.GenerateThreadId()
		}
		ctxSess := context.New(Logger.GetLogger()).
			SetXRequestID(reqId).
			SetAppName("api.awaymail.v1.go").
			SetAppVersion("0.0").
			SetPort(8080).
			SetSrcIP(c.Request.Host).
			SetURL(c.Request.URL.Path).
			SetMethod(c.Request.Method).
			SetHeader(c.Request.Header)

		c.Set(context.AppSession, ctxSess)

		c.Next()
	}
}
