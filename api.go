package pixiv

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/dghubble/sling"
)

const (
	clientID     = "MOBrBDS8blbauoSck0ZfDbtuzpyT"
	clientSecret = "lsACyCD94FhDUtGTXi3QzcFE2uU1hqtDaKeqrdwj"
)

var (
	_token, _refreshToken string
	_tokenDeadline        time.Time
)

type AccountProfileImages struct {
	Px16  string `json:"px_16x16"`
	Px50  string `json:"px_50x50"`
	Px170 string `json:"px_170x170"`
}

type Account struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	Account          string `json:"account"`
	MailAddress      string `json:"mail_address"`
	IsPremium        bool   `json:"is_premium"`
	XRestrict        int    `json:"x_restrict"`
	IsMailAuthorized bool   `json:"is_mail_authorized"`

	ProfileImage AccountProfileImages `json:"profile_image_urls"`
}

type authInfo struct {
	AccessToken  string   `json:"access_token"`
	ExpiresIn    int      `json:"expires_in"`
	TokenType    string   `json:"token_type"`
	Scope        string   `json:"scope"`
	RefreshToken string   `json:"refresh_token"`
	User         *Account `json:"user"`
	DeviceToken  string   `json:"device_token"`
}

type loginParams struct {
	GetSecureURL int    `url:"get_secure_url,omitempty"`
	ClientID     string `url:"client_id,omitempty"`
	ClientSecret string `url:"client_secret,omitempty"`
	GrantType    string `url:"grant_type,omitempty"`
	Username     string `url:"username,omitempty"`
	Password     string `url:"password,omitempty"`
	RefreshToken string `url:"refresh_token,omitempty"`
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

func auth(username, password string) (*authInfo, error) {
	s := sling.New().Base("https://oauth.secure.pixiv.net/").Set("User-Agent", "PixivAndroidApp/5.0.64 (Android 6.0)")
	params := &loginParams{
		GetSecureURL: 1,
		ClientID:     clientID,
		ClientSecret: clientSecret,
	}

	if (username != "") && (password != "") {
		params.GrantType = "password"
		params.Username = username
		params.Password = password
	} else {
		params.GrantType = "refresh_token"
		params.RefreshToken = _refreshToken
	}

	res := &loginResponse{
		Response: &authInfo{
			User: &Account{},
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
	_token = res.Response.AccessToken
	_tokenDeadline = time.Now().Add(time.Duration(res.Response.ExpiresIn) * time.Second)
	_refreshToken = res.Response.RefreshToken

	return res.Response, nil
}

func Login(username, password string) (*Account, error) {
	a, err := auth(username, password)
	if err != nil {
		return nil, err
	}
	return a.User, nil
}

func refreshToken() error {
	if _refreshToken == "" {
		return fmt.Errorf("missing refresh token")
	}
	_, err := auth("", "")
	if err != nil {
		return err
	}
	return nil
}

func CheckRefreshToken() error {
	if time.Now().Before(_tokenDeadline) {
		return nil
	}
	return refreshToken()
}

func setToken(t, rt string) {
	_token = t
	_refreshToken = rt
	_tokenDeadline = time.Time{}
}

func download(url, path, name string, replace bool) (int64, error) {
	if path == "" {
		return 0, fmt.Errorf("downloadpath needed")
	}
	if name == "" {
		name = filepath.Base(url)
	}
	fullPath := filepath.Join(path, name)

	if _, err := os.Stat(fullPath); err == nil {
		return 0, nil
	}

	output, err := os.Create(fullPath)
	if err != nil {
		return 0, err
	}
	defer output.Close()

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, err
	}
	req.Header.Add("Referer", apiBase)
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	n, err := io.Copy(output, resp.Body)
	if err != nil {
		return 0, err
	}
	return n, nil
}
