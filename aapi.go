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
