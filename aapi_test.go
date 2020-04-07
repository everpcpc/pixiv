package pixiv

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

var (
	testUID uint64
	mock    bool
	app     *AppPixivAPI
)

func init() {
	if mockTest() {
		mock = true
		fmt.Println("=== RUNNING mock tests")
		LoadAuth("fake_token", "fake_refresh_token", time.Time{})
		testUID = 12345678
		httpmock.Activate()

		resp, _ := getMockedResponse("auth.json")
		httpmock.RegisterResponder("POST", "https://oauth.secure.pixiv.net/auth/token",
			httpmock.NewStringResponder(200, resp))
	} else {
		fmt.Println("=== RUNNING real tests")
		LoadAuth(os.Getenv("TOKEN"), os.Getenv("REFRESH_TOKEN"), time.Time{})
		testUID, _ = strconv.ParseUint(os.Getenv("TEST_UID"), 10, 0)
	}
	app = NewApp()
}

// mockTest if one of the env not defined
func mockTest() bool {
	if os.Getenv("TOKEN") == "" {
		fmt.Println("TOKEN not set")
		return true
	}
	if os.Getenv("REFRESH_TOKEN") == "" {
		fmt.Println("REFRESH_TOKEN not set")
		return true
	}
	if os.Getenv("TEST_UID") == "" {
		fmt.Println("TEST_UID not set")
		return true
	}
	return false
}

func getMockedResponse(file string) (string, error) {
	f, err := ioutil.ReadFile("fixtures/response_" + file)
	return string(f), err
}

func TestUserDetail(t *testing.T) {
	if mock {
		resp, _ := getMockedResponse("user_detail.json")
		httpmock.RegisterResponder("GET", "https://app-api.pixiv.net/v1/user/detail?filter=for_ios&user_id=12345678",
			httpmock.NewStringResponder(200, resp))
	}

	r := require.New(t)
	detail, err := app.UserDetail(testUID)
	r.Nil(err)
	r.Equal(testUID, detail.User.ID)
}

func TestUserIllusts(t *testing.T) {
	if mock {
		resp, _ := getMockedResponse("user_illusts.json")
		httpmock.RegisterResponder("GET", "https://app-api.pixiv.net/v1/user/illusts?filter=for_ios&type=illust&user_id=490219",
			httpmock.NewStringResponder(200, resp))
	}

	r := require.New(t)
	illusts, next, err := app.UserIllusts(490219, "illust", 0)
	r.Nil(err)
	r.Len(illusts, 30)
	r.Equal(30, next)
}

func TestUserBookmarksIllust(t *testing.T) {
	if mock {
		resp, _ := getMockedResponse("user_bookmarks_illust.json")
		httpmock.RegisterResponder("GET", "https://app-api.pixiv.net/v1/user/bookmarks/illust?filter=for_ios&restrict=public&user_id=12345678",
			httpmock.NewStringResponder(200, resp))
	}

	r := require.New(t)
	illusts, _, err := app.UserBookmarksIllust(testUID, "public", 0, "")
	r.Nil(err)
	r.Equal(uint64(70095856), illusts[0].ID)
}

func TestIllustFollow(t *testing.T) {
	if mock {
		resp, _ := getMockedResponse("illust_follow.json")
		httpmock.RegisterResponder("GET", "https://app-api.pixiv.net/v2/illust/follow?restrict=public",
			httpmock.NewStringResponder(200, resp))
	}

	r := require.New(t)
	illusts, next, err := app.IllustFollow("public", 0)
	r.Nil(err)
	r.Len(illusts, 30)
	r.Equal(30, next)
}

func TestIllustDetail(t *testing.T) {
	if mock {
		resp, _ := getMockedResponse("illust_detail.json")
		httpmock.RegisterResponder("GET", "https://app-api.pixiv.net/v1/illust/detail?illust_id=68943534",
			httpmock.NewStringResponder(200, resp))
	}

	r := require.New(t)
	illust, err := app.IllustDetail(68943534)
	r.Nil(err)
	r.Equal(uint64(68943534), illust.ID)
}

func TestDownload(t *testing.T) {
	if mock {
		resp, _ := getMockedResponse("illust_detail.json")
		httpmock.RegisterResponder("GET", "https://app-api.pixiv.net/v1/illust/detail?illust_id=68943534",
			httpmock.NewStringResponder(200, resp))
		p0, _ := ioutil.ReadFile("fixtures/68943534_p0.jpg")
		httpmock.RegisterResponder("GET", "https://i.pximg.net/img-original/img/2018/05/27/12/14/11/68943534_p0.jpg",
			httpmock.NewBytesResponder(200, p0))
		p1, _ := ioutil.ReadFile("fixtures/68943534_p1.jpg")
		httpmock.RegisterResponder("GET", "https://i.pximg.net/img-original/img/2018/05/27/12/14/11/68943534_p1.jpg",
			httpmock.NewBytesResponder(200, p1))
		p2, _ := ioutil.ReadFile("fixtures/68943534_p2.jpg")
		httpmock.RegisterResponder("GET", "https://i.pximg.net/img-original/img/2018/05/27/12/14/11/68943534_p2.jpg",
			httpmock.NewBytesResponder(200, p2))
	}

	r := require.New(t)
	sizes, errs := app.Download(68943534, ".")
	r.Len(sizes, 3)
	for i := range errs {
		r.Nil(errs[i])
	}
	r.Equal(int64(2748932), sizes[0])
	r.Equal(int64(2032716), sizes[1])
	r.Equal(int64(600670), sizes[2])
}
