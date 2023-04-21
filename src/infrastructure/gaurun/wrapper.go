package gaurun

import ctxSess "github.com/cloudsrc/api.awaymail.v1.go/src/shared/utils/context"

type Wrapper interface {
	Send(ctxSess *ctxSess.Context, payload Notification) (err error)
}
