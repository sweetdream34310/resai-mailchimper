package libs

import (
	ctxSess "github.com/cloudsrc/api.awaymail.v1.go/src/shared/utils/context"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Error   bool        `json:"error"`
	Data    interface{} `json:"data"`
	Status  int         `json:"status"`
	Message string      `json:"message"`
}

type ResponseFile struct {
	Error   bool   `json:"error"`
	Data    []byte `json:"data"`
	Status  int    `json:"status"`
	Message string `json:"message"`
}

const salt = "dTj%/2}#+{5CC;zeK5-&"

// SendResponse : implements sending of response
func (res *Response) SendResponse(ctxSess *ctxSess.Context, ctx *gin.Context) {
	if res.Status != http.StatusOK {
		res.Error = true
	}
	if ctxSess != nil {
		//r, _ := json.Marshal(res)
		ctxSess.SetResponseCode(res.Status)
		//ctxSess.Lv4(utils.Encrypt(salt, string(r)))
		ctxSess.Lv4()
	}
	ctx.JSON(res.Status, &res)
}

func (res *ResponseFile) SendFleResponse(ctx *gin.Context, contentType string) {
	if res.Status != http.StatusOK {
		res.Error = true
	}
	ctx.Data(res.Status, contentType, res.Data)
}
