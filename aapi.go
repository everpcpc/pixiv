package pixiv

import (
	"fmt"
	"path/filepath"

	"github.com/dghubble/sling"
)

const (
	apiBase = "https://app-api.pixiv.net/"
)

// AppPixivAPI -- App-API (6.x - app-api.pixiv.net)
type AppPixivAPI struct {
	sling *sling.Sling
}

func NewApp() *AppPixivAPI {
	s := sling.New().Base(apiBase).Set("User-Agent", "PixivIOSApp/7.6.2 (iOS 12.2; iPhone9,1)").Set("App-Version", "7.6.2").Set("App-OS-VERSION", "12.2").Set("App-OS", "ios")
	return &AppPixivAPI{sling: s}
}

type UserImages struct {
	Medium string `json:"medium"`
}
type User struct {
	ID         uint64 `json:"id"`
	Name       string `json:"name"`
	Account    string `json:"account"`
	Comment    string `json:"comment"`
	IsFollowed bool   `json:"is_followed"`

	ProfileImages UserImages `json:"profile_image_urls"`
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
type illustResponse struct {
	Illust Illust `json:"illust"`
}

func (a *AppPixivAPI) request(path string, params, data interface{}, auth bool) (err error) {
	if auth {
		if _, err := refreshAuth(); err != nil {
			return fmt.Errorf("refresh token failed: %v", err)
		}
		_, err = a.sling.New().Get(path).Set("Authorization", "Bearer "+_token).QueryStruct(params).ReceiveSuccess(data)
	} else {
		_, err = a.sling.New().Get(path).QueryStruct(params).ReceiveSuccess(data)
	}
	return err
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
func (a *AppPixivAPI) UserIllusts(uid uint64, _type string, offset int) ([]Illust, int, error) {
	path := "v1/user/illusts"
	params := &userIllustsParams{
		UserID: uid,
		Filter: "for_ios",
		Type:   _type,
		Offset: offset,
	}
	data := &illustsResponse{}
	if err := a.request(path, params, data, true); err != nil {
		return nil, 0, err
	}
	next, err := parseNextPageOffset(data.NextURL)
	return data.Illusts, next, err
}

type userBookmarkIllustsParams struct {
	UserID        uint64 `url:"user_id,omitempty"`
	Restrict      string `url:"restrict,omitempty"`
	Filter        string `url:"filter,omitempty"`
	MaxBookmarkID int    `url:"max_bookmark_id,omitempty"`
	Tag           string `url:"tag,omitempty"`
}

// UserBookmarksIllust restrict: [public, private]
func (a *AppPixivAPI) UserBookmarksIllust(uid uint64, restrict string, maxBookmarkID int, tag string) ([]Illust, int, error) {
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
		return nil, 0, err
	}
	next, err := parseNextPageOffset(data.NextURL)
	return data.Illusts, next, err
}

type illustFollowParams struct {
	Restrict string `url:"restrict,omitempty"`
	Offset   int    `url:"offset,omitempty"`
}

// IllustFollow restrict: [public, private]
func (a *AppPixivAPI) IllustFollow(restrict string, offset int) ([]Illust, int, error) {
	path := "v2/illust/follow"
	params := &illustFollowParams{
		Restrict: restrict,
		Offset:   offset,
	}
	data := &illustsResponse{}
	if err := a.request(path, params, data, true); err != nil {
		return nil, 0, err
	}
	next, err := parseNextPageOffset(data.NextURL)
	return data.Illusts, next, err
}

type illustDetailParams struct {
	IllustID uint64 `url:"illust_id,omitemtpy"`
}

// IllustDetail get a detailed illust with id
func (a *AppPixivAPI) IllustDetail(id uint64) (*Illust, error) {
	path := "v1/illust/detail"
	data := &illustResponse{}
	params := &illustDetailParams{
		IllustID: id,
	}
	if err := a.request(path, params, data, true); err != nil {
		return nil, err
	}
	return &data.Illust, nil
}

// Download a specific picture from pixiv id
func (a *AppPixivAPI) Download(id uint64, path string) ([]int64, []error) {
	illust, err := a.IllustDetail(id)
	if err != nil {
		return []int64{0}, []error{err}
	}
	var urls []string
	if illust.MetaSinglePage.OriginalImageURL == "" {
		for _, img := range illust.MetaPages {
			urls = append(urls, img.Images.Original)
		}
	} else {
		urls = append(urls, illust.MetaSinglePage.OriginalImageURL)
	}
	var sizes []int64
	var errs []error
	for _, u := range urls {
		size, err := download(u, path, filepath.Base(u), false)
		sizes = append(sizes, size)
		errs = append(errs, err)
	}

	return sizes, errs
}
