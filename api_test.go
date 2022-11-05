package pixiv

import (
	"os"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

func setupAPIMockTest(code int, responseFile string) {
	httpmock.Activate()
	resp, _ := getMockedResponse(responseFile)
	httpmock.RegisterResponder("POST", "https://oauth.secure.pixiv.net/auth/token",
		httpmock.NewStringResponder(code, resp))
}

func TestAuth(t *testing.T) {
	username := os.Getenv("USERNAME")
	password := os.Getenv("PASSWORD")
	testUID := os.Getenv("TEST_UID")
	if username == "" || password == "" || testUID == "" {
		t.Log("No username or password found, mock TestAuth")
		testUID = "12345678"
		setupAPIMockTest(200, "auth.json")
	}

	r := require.New(t)
	account, err := Login(username, username)
	r.Nil(err)
	r.Equal(testUID, account.ID)

	httpmock.DeactivateAndReset()
}

func TestLoadAuth(t *testing.T) {
	token := os.Getenv("TOKEN")
	refreshToken := os.Getenv("REFRESH_TOKEN")
	testUID := os.Getenv("TEST_UID")
	if token == "" || refreshToken == "" {
		t.Log("No token or refresh token found, mock TestLoadAuth")
		setupAPIMockTest(200, "auth_refresh_token.json")
		token = "xxxxxxxx"
		refreshToken = "xxxxxxxxxxxxxx"
		testUID = "12345678"
	}

	r := require.New(t)
	account, err := LoadAuth(token, refreshToken, time.Time{})
	r.Nil(err)
	r.Equal(testUID, account.ID)

	httpmock.DeactivateAndReset()
}

func TestLoginFail(t *testing.T) {
	username := os.Getenv("USERNAME")
	password := os.Getenv("PASSWORD")
	if username == "" || password == "" {
		t.Log("No username or password found, mock TestLoginFail")
		setupAPIMockTest(400, "auth_invalid_password.json")
		username = "fake_username"
		password = "fake_password"
	}

	r := require.New(t)
	account, err := Login(username[:5], password[:5])
	r.Nil(account)
	r.EqualError(err, "Login system error: 103:pixiv ID、またはメールアドレス、パスワードが正しいかチェックしてください。")

	httpmock.DeactivateAndReset()
}

func TestRefreshTokenFail(t *testing.T) {
	token := os.Getenv("TOKEN")
	refreshToken := os.Getenv("REFRESH_TOKEN")
	if token == "" || refreshToken == "" {
		t.Log("No token or refresh token found, mock TestRefreshTokenFail")
		setupAPIMockTest(400, "auth_invalid_token.json")
		token = "xxxxxxxx"
		refreshToken = "xxxxxxxxxxxxxx"
	}

	r := require.New(t)
	account, err := LoadAuth(token, refreshToken[:10], time.Time{})
	r.Nil(account)
	r.EqualError(err, "refresh token: Login system error: Invalid refresh token")

	httpmock.DeactivateAndReset()
}
