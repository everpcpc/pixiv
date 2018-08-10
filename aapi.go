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

type UserProfileImages struct {
	Medium string `json:"medium"`
}
type User struct {
	ID         uint64 `json:"id"`
	Name       string `json:"name"`
	Account    string `json:"account"`
	Comment    string `json:"comment"`
	IsFollowed bool   `json:"is_followed"`

	ProfileImages UserProfileImages `json:"profile_image_urls"`
}
type UserDetail struct {
	User *User `json:"user"`
	// TODO:
	// Profile
	// ProfilePublicity
	// Workspace
}

type Tag struct {
	Name string `json:"name"`
}
type Images struct {
	SquareMedium string `json:"square_medium"`
	Medium       string `json:"medium"`
	Large        string `json:"large"`
	Original     string `json:"original"`
}
type MetaSinglePage struct {
	OriginalImageURL string `json:"original_image_url"`
}
type MetaPage struct {
	Images Images `json:"image_urls"`
}
type Illust struct {
	ID          uint64   `json:"id"`
	Title       string   `json:"title"`
	Type        string   `json:"type"`
	Images      Images   `json:"image_urls"`
	Caption     string   `json:"caption"`
	Restrict    int      `json:"restrict"`
	User        User     `json:"user"`
	Tags        []Tag    `json:"tags"`
	Tools       []string `json:"tools"`
	CreateData  string   `json:"create_data"`
	PageCount   int      `json:"page_count"`
	Width       int      `json:"width"`
	Height      int      `json:"height"`
	SanityLevel int      `json:"sanity_level"`
	// TODO:
	// Series `json:"series"`
	MetaSinglePage MetaSinglePage `json:"meta_single_page"`
	MetaPages      []MetaPage     `json:"meta_pages"`
	TotalView      int            `json:"total_view"`
	TotalBookmarks int            `json:"total_bookmarks"`
	IsBookmarked   bool           `json:"is_bookmarked"`
	Visible        bool           `json:"visible"`
	IsMuted        bool           `json:"is_muted"`
	TotalComments  int            `json:"total_comments"`
}

type illustsResponse struct {
	Illusts []Illust `json:"illusts"`
	NextURL string   `json:"next_url"`
}

func (a *AppPixivAPI) request(path string, params, data interface{}, auth bool) (err error) {
	if auth {
		if err := CheckRefreshToken(); err != nil {
			return fmt.Errorf("refresh token failed: %v", err)
		}
		_, err = a.sling.New().Get(path).Set("Authorization", "Bearer "+_token).QueryStruct(params).ReceiveSuccess(data)
	} else {
		_, err = a.sling.New().Get(path).QueryStruct(params).ReceiveSuccess(data)
	}

	if err != nil {
		return err
	}
	return nil
}

type userDetailParams struct {
	UserID uint64 `url:"user_id,omitempty"`
	Filter string `url:"filter,omitempty"`
}

func (a *AppPixivAPI) UserDetail(uid uint64) (*UserDetail, error) {
	path := "v1/user/detail"
	params := &userDetailParams{
		UserID: uid,
		Filter: "for_ios",
	}
	detail := &UserDetail{
		User: &User{},
	}
	if err := a.request(path, params, detail, true); err != nil {
		return nil, err
	}
	return detail, nil
}

type userIllustsParams struct {
	UserID uint64 `url:"user_id,omitempty"`
	Filter string `url:"filter,omitempty"`
	Type   string `url:"type,omitempty"`
	Offset int    `url:"offset,omitempty"`
}

// UserIllusts type: [illust, manga]
func (a *AppPixivAPI) UserIllusts(uid uint64, _type string, offset int) ([]Illust, error) {
	path := "v1/user/illusts"
	params := &userIllustsParams{
		UserID: uid,
		Filter: "for_ios",
		Type:   _type,
		Offset: offset,
	}
	data := &illustsResponse{}
	if err := a.request(path, params, data, true); err != nil {
		return nil, err
	}
	return data.Illusts, nil
}

type userBookmarkIllustsParams struct {
	UserID        uint64 `url:"user_id,omitempty"`
	Restrict      string `url:"restrict,omitempty"`
	Filter        string `url:"filter,omitempty"`
	MaxBookmarkID int    `url:"max_bookmark_id,omitempty"`
	Tag           string `url:"tag,omitempty"`
}

// UserBookmarksIllust restrict: [public, private]
func (a *AppPixivAPI) UserBookmarksIllust(uid uint64, restrict string, maxBookmarkID int, tag string) ([]Illust, error) {
	path := "v1/user/bookmarks/illust"
	params := &userBookmarkIllustsParams{
		UserID:        uid,
		Restrict:      "public",
		Filter:        "for_ios",
		MaxBookmarkID: maxBookmarkID,
		Tag:           tag,
	}
	data := &illustsResponse{}
	if err := a.request(path, params, data, true); err != nil {
		return nil, err
	}
	return data.Illusts, nil
}

type illustFollowParams struct {
	Restrict string `url:"restrict,omitempty"`
	Offset   int    `url:"offset,omitempty"`
}

// IllustFollow restrict: [public, private]
func (a *AppPixivAPI) IllustFollow(restrict string, offset int) ([]Illust, error) {
	path := "v2/illust/follow"
	params := &illustFollowParams{
		Restrict: restrict,
		Offset:   offset,
	}
	data := &illustsResponse{}
	if err := a.request(path, params, data, true); err != nil {
		return nil, err
	}
	return data.Illusts, nil
}

// func (a *AppPixivAPI) testResponse() error {
// 	if err := CheckRefreshToken(); err != nil {
// 		return fmt.Errorf("refresh token failed: %v", err)
// 	}
// 	path := "v2/illust/follow"
// 	params := &illustFollowParams{
// 		Restrict: "public",
// 	}
// 	req, err := a.sling.New().Get(path).Set("Authorization", "Bearer "+_token).QueryStruct(params).Request()
// 	client := &http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return err
// 	}
// 	body, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		return err
// 	}
// 	fmt.Println("======", string(body))
// 	return fmt.Errorf("OK")
// }
