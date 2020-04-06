package pixiv

import (
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

func init() {
	resp, _ := getMockedResponse("auth.json")
	httpmock.RegisterResponder("POST", "https://oauth.secure.pixiv.net/auth/token",
		httpmock.NewStringResponder(200, resp))
}

func TestAuth(t *testing.T) {
	httpmock.Activate()

	r := require.New(t)
	account, err := Login("x", "x")
	r.Nil(err)
	r.Equal("12345678", account.ID)
}
