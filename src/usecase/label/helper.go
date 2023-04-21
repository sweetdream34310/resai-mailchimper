package label

import (
	"google.golang.org/api/gmail/v1"

	"github.com/cloudsrc/api.awaymail.v1.go/src/shared/constants"
	ctxSess "github.com/cloudsrc/api.awaymail.v1.go/src/shared/utils/context"
)

func (s *service) checkLabelArchive(ctxSess *ctxSess.Context) (archive *gmail.Label, err error) {
	ls, err := s.googleWrapper.GetLabelList(ctxSess)
	if err != nil {
		err = constants.ErrorGeneral
		return
	}
	var flagArchive bool
	for _, eachLabel := range ls {
		if eachLabel.Name == constants.LABEL_AwayARCHIVE {
			flagArchive = true
			archive = eachLabel
			break
		}
	}

	if !flagArchive {
		archive, err = s.googleWrapper.CreateLabel(ctxSess, &gmail.Label{Name: constants.LABEL_AwayARCHIVE})
		if err != nil {
			ctxSess.ErrorMessage = err.Error()
			err = constants.ErrorGeneral
			return
		}
	}

	return
}
