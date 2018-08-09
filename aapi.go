package pixiv

import (
	"fmt"

	"github.com/dghubble/sling"
)

// AppPixivAPI -- App-API (6.x - app-api.pixiv.net)
type AppPixivAPI struct {
	sling *sling.Sling
}

func NewApp() *AppPixivAPI {
	s := sling.New().Base("https://app-api.pixiv.net/").Set("User-Agent", "PixivIOSApp/6.7.1 (iOS 10.3.1; iPhone8,1)").Set("App-Version", "6.7.1").Set("App-OS-VERSION", "10.3.1").Set("App-OS", "ios")
	return &AppPixivAPI{
		sling: s,
	}
}

type authParams struct {
}

func (a *AppPixivAPI) Auth(username, password string) {

	// url := "https://oauth.secure.pixiv.net/auth/token"
	// 'User-Agent': 'PixivAndroidApp/5.0.64 (Android 6.0)',
	// data = {
	// 	'get_secure_url': 1,
	// 	'client_id': self.client_id,
	// 	'client_secret': self.client_secret,
	// }
}

type userParams struct {
	userID string `url:"user_id,omitempty"`
	filter string `url:"filter,omitempty"`
}

func (a *AppPixivAPI) UserDetail(userID string) error {
	path := "v1/user/detail"
	params := &userParams{
		userID: userID,
		filter: "for_ios",
	}
	var user interface{}
	resp, err := a.sling.New().Get(path).QueryStruct(params).ReceiveSuccess(user)
	if err != nil {
		return err
	}
	fmt.Println("====", resp, user)
	return nil
}
