package pixiv

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dghubble/sling"
	"github.com/pkg/errors"
)

const (
	clientID         = "MOBrBDS8blbauoSck0ZfDbtuzpyT"
	clientSecret     = "lsACyCD94FhDUtGTXi3QzcFE2uU1hqtDaKeqrdwj"
	clientHashSecret = "28c1fdd170a5204386cb1313c7077b34f83e4aaf4aa829ce78c231e05b0bae2c"
)

var (
	_token, _refreshToken string
	_tokenDeadline        time.Time
	authHook              func(string, string, time.Time) error
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

type authParams struct {
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

func genClientHash(clientTime string) string {
	h := md5.New()
	io.WriteString(h, clientTime)
	io.WriteString(h, clientHashSecret)
	return hex.EncodeToString(h.Sum(nil))
}

func auth(params *authParams) (*authInfo, error) {
	clientTime := time.Now().Format(time.RFC3339)
	s := sling.New().Base("https://oauth.secure.pixiv.net/").Set("User-Agent", "PixivAndroidApp/5.0.115 (Android 6.0)").Set("X-Client-Time", clientTime).Set("X-Client-Hash", genClientHash(clientTime))

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
		return nil, errors.Wrap(err, "auth request failed")
	}
	if loginErr.HasError {
		for k, v := range loginErr.Errors {
			return nil, fmt.Errorf("Login %s error: %s", k, v.Message)
		}
	}
	_token = res.Response.AccessToken
	_refreshToken = res.Response.RefreshToken
	_tokenDeadline = time.Now().Add(time.Duration(res.Response.ExpiresIn) * time.Second)

	if authHook != nil {
		err = authHook(_token, _refreshToken, _tokenDeadline)
	}

	return res.Response, err
}

// HookAuth add a hook with (token, refreshToken, tokenDeadline) after a successful auth.
// Prividing a way to store the latest token.
func HookAuth(f func(string, string, time.Time) error) {
	authHook = f
}

func Login(username, password string) (*Account, error) {
	params := &authParams{
		GetSecureURL: 1,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		GrantType:    "password",
		Username:     username,
		Password:     password,
	}
	a, err := auth(params)
	if err != nil {
		return nil, errors.Wrap(err, "login failed")
	}
	return a.User, nil
}

func LoadAuth(token, refreshToken string, tokenDeadline time.Time) (*Account, error) {
	_token = token
	_refreshToken = refreshToken
	_tokenDeadline = tokenDeadline
	return refreshAuth(true)
}

func refreshAuth(force bool) (*Account, error) {
	if !force && time.Now().Before(_tokenDeadline) {
		return nil, nil
	}
	if _refreshToken == "" {
		return nil, fmt.Errorf("missing refresh token")
	}
	params := &authParams{
		GetSecureURL: 1,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		GrantType:    "refresh_token",
		RefreshToken: _refreshToken,
	}
	a, err := auth(params)
	if err != nil {
		return nil, errors.Wrap(err, "refresh token")
	}
	return a.User, nil
}

// download image to file (use 6.0 app-api)
func download(client *http.Client, url, path, name, tmpdir string, replace bool) (int64, error) {
	if path == "" {
		return 0, fmt.Errorf("download path needed")
	}
	if name == "" {
		name = filepath.Base(url)
	}
	fullPath := filepath.Join(path, name)

	if !replace {
		if _, err := os.Stat(fullPath); err == nil {
			return 0, nil
		}
	}

	tmpFile, err := os.CreateTemp(tmpdir, "pixiv-*")
	if err != nil {
		return 0, errors.Wrap(err, "create temp file failed")
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, errors.Wrap(err, "create request failed")
	}
	req.Header.Add("Referer", apiBase)
	resp, err := client.Do(req)
	if err != nil {
		return 0, errors.Wrap(err, "request failed")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("download failed: %s", resp.Status)
	}

	n, err := io.Copy(tmpFile, resp.Body)
	if err != nil {
		return 0, errors.Wrap(err, "download failed")
	}

	err = tmpFile.Close()
	if err != nil {
		return 0, errors.Wrap(err, "close temp file failed")
	}

	err = os.Rename(tmpFile.Name(), fullPath)
	if err != nil {
		return 0, errors.Wrap(err, "rename temp file failed")
	}

	return n, nil
}

func downloadImage(ctx context.Context, client *http.Client, url string) (Image, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return Image{}, errors.Wrap(err, "create request failed")
	}
	req.Header.Add("Referer", apiBase)
	resp, err := client.Do(req)
	if err != nil {
		return Image{}, errors.Wrap(err, "request failed")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Image{}, fmt.Errorf("download failed: %s", resp.Status)
	}

	// image/xxx
	contentType := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image") {
		return Image{}, fmt.Errorf("invalid content type: %s, %s", contentType, url)
	}
	contentTypeSplits := strings.Split(contentType, "/")
	if len(contentTypeSplits) != 2 || len(contentTypeSplits[1]) == 0 {
		return Image{}, fmt.Errorf("invalid content type: %s, %s", contentType, url)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return Image{}, errors.Wrap(err, "download failed")
	}
	return Image{
		Data:    b,
		P:       urlToPageNum(url),
		Type:    parseImageType(contentTypeSplits[1]),
		TypeRaw: contentTypeSplits[1],
	}, nil
}
