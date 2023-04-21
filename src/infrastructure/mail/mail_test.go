package mail

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSendEmail(t *testing.T) {
	e := Email{}
	_, _, err := e.SendEmail("vicky.dwi.k@icloud.com", "2nd January 2023", "https://am.dev.cloudsrc.com/v2/breakthrough/amp/notification?email=vicky@cloudsource.io&id=01GP1FBD3VPSJH27SFKY7PSAPSt", "https://am.dev.cloudsrc.com/v2/breakthrough/notification?email=vicky@cloudsource.io&id=01GP1FBD3VPSJH27SFKY7PSAPSt", "vicky@cloudsource.io", "Vicky")
	assert.Empty(t, err)
}
