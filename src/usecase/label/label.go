package label

import (
	"strings"

	"google.golang.org/api/gmail/v1"

	googleWrapper "github.com/cloudsrc/api.awaymail.v1.go/src/infrastructure/google"
	"github.com/cloudsrc/api.awaymail.v1.go/src/shared/constants"
	ctxSess "github.com/cloudsrc/api.awaymail.v1.go/src/shared/utils/context"
)

type service struct {
	googleWrapper googleWrapper.Wrapper
}

func New(googleWrapper googleWrapper.Wrapper) Service {
	return &service{
		googleWrapper: googleWrapper,
	}
}

func (s *service) CreateLabel(ctxSess *ctxSess.Context, req *CreateLabelReq) (resp *gmail.Label, err error) {
	resp, err = s.googleWrapper.CreateLabel(ctxSess, &gmail.Label{Name: req.Name})
	if err != nil {
		err = constants.ErrorGeneral
		return
	}
	return
}

func (s *service) DeleteLabel(ctxSess *ctxSess.Context, labelID string) (err error) {
	err = s.googleWrapper.DeleteLabel(ctxSess, labelID)
	if err != nil {
		err = constants.ErrorGeneral
		return
	}
	return
}

func (s *service) PatchLabel(ctxSess *ctxSess.Context, labelID string, req *PatchLabelReq) (resp *gmail.Label, err error) {
	resp, err = s.googleWrapper.PatchLabel(ctxSess, labelID, &gmail.Label{Name: req.Name})
	if err != nil {
		err = constants.ErrorGeneral
		return
	}
	return
}

func (s *service) GetLabel(ctxSess *ctxSess.Context, labelId string) (resp *gmail.Label, err error) {
	if strings.ToLower(labelId) == strings.ToLower(constants.LABEL_ARCHIVE) {
		archive, errs := s.checkLabelArchive(ctxSess)
		if errs != nil {
			return
		}

		labelId = archive.Id
	}

	resp, err = s.googleWrapper.GetUserLabel(ctxSess, labelId)
	if err != nil {
		err = constants.ErrorGeneral
		return
	}
	return
}

func (s *service) FindAllLabels(ctxSess *ctxSess.Context) (resp []*gmail.Label, err error) {
	ls, err := s.googleWrapper.GetLabelList(ctxSess)
	if err != nil {
		err = constants.ErrorGeneral
		return
	}
	var tempLabel []*gmail.Label
	for _, eachLabel := range ls {
		if eachLabel.Type == constants.LabelTypeUser {
			tempLabel = append(tempLabel, eachLabel)
		} else if eachLabel.MessageListVisibility == "" || eachLabel.LabelListVisibility == "labelShow" {
			tempLabel = append(tempLabel, eachLabel)
		}
	}
	for {
		for _, each := range tempLabel {
			flagBreak := false
			if len(resp) == 2 {
				resp = append(resp, &gmail.Label{
					Id:   constants.LABEL_SNOOZED,
					Name: constants.LABEL_SNOOZED,
					Type: "system",
				})
			} else if len(resp) == 5 {
				for _, v := range tempLabel {
					if v.Name == constants.LABEL_AwayARCHIVE {
						v.Name = constants.LABEL_ARCHIVE
						resp = append(resp, v)
						flagBreak = true
						break
					}
				}
				if flagBreak {
					break
				}
			}
			switch each.Name {
			case "INBOX":
				if len(resp) == 0 {
					resp = append(resp, each)
					flagBreak = true
				}
			case "STARRED":
				if len(resp) == 1 {
					resp = append(resp, each)
					flagBreak = true
				}
			case "SENT":
				if len(resp) == 3 {
					resp = append(resp, each)
					flagBreak = true
				}
			case "DRAFT":
				if len(resp) == 4 {
					resp = append(resp, each)
					flagBreak = true
				}
			case constants.LABEL_AwayARCHIVE, constants.LABEL_ARCHIVE:
			default:
				resp = append(resp, each)
			}
			if flagBreak {
				break
			}
		}
		if len(resp) > len(tempLabel) {
			break
		}
	}

	return
}
