package pixiv

import (
	"fmt"

	"github.com/dghubble/sling"
)

const (
	clientID     = "MOBrBDS8blbauoSck0ZfDbtuzpyT"
	clientSecret = "lsACyCD94FhDUtGTXi3QzcFE2uU1hqtDaKeqrdwj"
)

var (
	token, refreshToken string
)

type UserProfileImages struct {
	Px16  string `json:"px_16x16"`
	Px50  string `json:"px_50x50"`
	Px170 string `json:"px_170x170"`
}

type User struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	Account          string `json:"account"`
	MailAddress      string `json:"mail_address"`
	IsPremium        bool   `json:"is_premium"`
	XRestrict        int    `json:"x_restrict"`
	IsMailAuthorized bool   `json:"is_mail_authorized"`

	ProfileImageURLs UserProfileImages `json:"profile_image_urls"`
}

type authInfo struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	RefreshToken string `json:"refresh_token"`
	User         *User  `json:"user"`
	DeviceToken  string `json:"device_token"`
}

type loginParams struct {
	GetSecureURL int    `url:"get_secure_url,omitempty"`
	ClientID     string `url:"client_id,omitempty"`
	ClientSecret string `url:"client_secret,omitempty"`
	GrantType    string `url:"grant_type,omitempty"`
	Username     string `url:"username,omitempty"`
	Password     string `url:"password,omitempty"`
}

type loginResponse struct {
	Response *authInfo `json:"response"`
}
type loginError struct {
	HasError bool              `json:"has_error"`
	Errors   map[string]Perror `json:"errors"`
}
type Perror struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func Login(username, password string) (*User, error) {
	s := sling.New().Base("https://oauth.secure.pixiv.net/").Set("User-Agent", "PixivAndroidApp/5.0.64 (Android 6.0)")
	params := &loginParams{
		GetSecureURL: 1,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		GrantType:    "password",
		Username:     username,
		Password:     password,
	}
	res := &loginResponse{
		Response: &authInfo{
			User: &User{},
		},
	}
	loginErr := &loginError{
		Errors: map[string]Perror{},
	}
	_, err := s.New().Post("auth/token").BodyForm(params).Receive(res, loginErr)
	if err != nil {
		return nil, err
	}
	if loginErr.HasError {
		for k, v := range loginErr.Errors {
			return nil, fmt.Errorf("Login %s error: %s", k, v.Message)
		}
	}
	token = res.Response.AccessToken
	refreshToken = res.Response.RefreshToken

	return res.Response.User, nil
}

func SetToken(t, rt string) {
	token = t
	refreshToken = rt
}
