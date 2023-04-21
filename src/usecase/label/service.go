package label

import (
	"google.golang.org/api/gmail/v1"

	ctxSess "github.com/cloudsrc/api.awaymail.v1.go/src/shared/utils/context"
)

type Service interface {
	CreateLabel(ctxSess *ctxSess.Context, req *CreateLabelReq) (resp *gmail.Label, err error)
	FindAllLabels(ctxSess *ctxSess.Context) (resp []*gmail.Label, err error)
	GetLabel(ctxSess *ctxSess.Context, labelId string) (resp *gmail.Label, err error)
	DeleteLabel(ctxSess *ctxSess.Context, labelID string) (err error)
	PatchLabel(ctxSess *ctxSess.Context, labelID string, req *PatchLabelReq) (resp *gmail.Label, err error)
}
