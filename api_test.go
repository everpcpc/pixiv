package pixiv

import (
	"fmt"
	"os"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

func TestAuth(t *testing.T) {
	username := os.Getenv("USERNAME")
	password := os.Getenv("PASSWORD")
	testUID := os.Getenv("TEST_UID")
	if username == "" || password == "" || testUID == "" {
		testUID = "12345678"
		fmt.Println("=== RUNNING mock tests for api")
		httpmock.Activate()
		resp, _ := getMockedResponse("auth.json")
		httpmock.RegisterResponder("POST", "https://oauth.secure.pixiv.net/auth/token",
			httpmock.NewStringResponder(200, resp))
	}

	r := require.New(t)
	account, err := Login(username, username)
	r.Nil(err)
	r.Equal(testUID, account.ID)
}
