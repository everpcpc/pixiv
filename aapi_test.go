package pixiv

import (
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

// appMockTest if one of the env not defined
func appMockTest() bool {
	if os.Getenv("TOKEN") == "" {
		return true
	}
	if os.Getenv("REFRESH_TOKEN") == "" {
		return true
	}
	if os.Getenv("TEST_UID") == "" {
		return true
	}
	return false
}

func setupAPPMockTest() uint64 {
	httpmock.Activate()
	resp, _ := getMockedResponse("auth.json")
	httpmock.RegisterResponder("POST", "https://oauth.secure.pixiv.net/auth/token",
		httpmock.NewStringResponder(200, resp))
	LoadAuth("fake_token", "fake_refresh_token", time.Time{})
	return 12345678
}

func setupAPPRealTest() uint64 {
	LoadAuth(os.Getenv("TOKEN"), os.Getenv("REFRESH_TOKEN"), time.Time{})
	testUID, _ := strconv.ParseUint(os.Getenv("TEST_UID"), 10, 0)
	return testUID
}

func setupAPPTest() (uint64, bool, *AppPixivAPI) {
	var testUID uint64
	var mock bool
	if appMockTest() {
		testUID = setupAPPMockTest()
		mock = true
	} else {
		testUID = setupAPPRealTest()
	}
	os.MkdirAll("tmp", 0755)
	return testUID, mock, NewApp().WithTmpdir("tmp")
}

func getMockedResponse(file string) (string, error) {
	f, err := os.ReadFile("fixtures/response_" + file)
	return string(f), err
}

func TestUserDetail(t *testing.T) {
	testUID, mock, app := setupAPPTest()
	if mock {
		resp, _ := getMockedResponse("user_detail.json")
		httpmock.RegisterResponder("GET", "https://app-api.pixiv.net/v1/user/detail?filter=for_ios&user_id=12345678",
			httpmock.NewStringResponder(200, resp))
	}

	r := require.New(t)
	detail, err := app.UserDetail(testUID)
	r.Nil(err)
	r.Equal(testUID, detail.User.ID)

	if mock {
		httpmock.DeactivateAndReset()
	}
}

func TestUserIllusts(t *testing.T) {
	_, mock, app := setupAPPTest()
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

	if mock {
		httpmock.DeactivateAndReset()
	}
}

func TestUserBookmarksIllust(t *testing.T) {
	testUID, mock, app := setupAPPTest()
	if mock {
		resp, _ := getMockedResponse("user_bookmarks_illust.json")
		httpmock.RegisterResponder("GET", "https://app-api.pixiv.net/v1/user/bookmarks/illust?filter=for_ios&restrict=public&user_id=12345678",
			httpmock.NewStringResponder(200, resp))
	}

	r := require.New(t)
	illusts, _, err := app.UserBookmarksIllust(testUID, "public", 0, "")
	r.Nil(err)
	r.Equal(uint64(70095856), illusts[0].ID)

	if mock {
		httpmock.DeactivateAndReset()
	}
}

func TestIllustFollow(t *testing.T) {
	_, mock, app := setupAPPTest()
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

	if mock {
		httpmock.DeactivateAndReset()
	}
}

func TestIllustDetail(t *testing.T) {
	_, mock, app := setupAPPTest()
	if mock {
		resp, _ := getMockedResponse("illust_detail.json")
		httpmock.RegisterResponder("GET", "https://app-api.pixiv.net/v1/illust/detail?illust_id=68943534",
			httpmock.NewStringResponder(200, resp))
	}

	r := require.New(t)
	illust, err := app.IllustDetail(68943534)
	r.Nil(err)
	r.Equal(uint64(68943534), illust.ID)

	tt, err := time.Parse(time.RFC3339, "2018-05-27T12:14:11+09:00")
	r.Nil(err)
	r.Equal(tt, illust.CreateDate)

	if mock {
		httpmock.DeactivateAndReset()
	}
}

func TestDownload(t *testing.T) {
	_, mock, app := setupAPPTest()
	if mock {
		resp, _ := getMockedResponse("illust_detail.json")
		httpmock.RegisterResponder("GET", "https://app-api.pixiv.net/v1/illust/detail?illust_id=68943534",
			httpmock.NewStringResponder(200, resp))
		p0, _ := os.ReadFile("fixtures/68943534_p0.jpg")
		httpmock.RegisterResponder("GET", "https://i.pximg.net/img-original/img/2018/05/27/12/14/11/68943534_p0.jpg",
			httpmock.NewBytesResponder(200, p0))
		p1, _ := os.ReadFile("fixtures/68943534_p1.jpg")
		httpmock.RegisterResponder("GET", "https://i.pximg.net/img-original/img/2018/05/27/12/14/11/68943534_p1.jpg",
			httpmock.NewBytesResponder(200, p1))
		p2, _ := os.ReadFile("fixtures/68943534_p2.jpg")
		httpmock.RegisterResponder("GET", "https://i.pximg.net/img-original/img/2018/05/27/12/14/11/68943534_p2.jpg",
			httpmock.NewBytesResponder(200, p2))
	}

	r := require.New(t)
	sizes, err := app.Download(68943534, ".")
	r.Nil(err)
	r.Len(sizes, 3)

	r.Equal(int64(2742014), sizes[0])
	r.Equal(int64(2008306), sizes[1])
	r.Equal(int64(560294), sizes[2])

	if mock {
		httpmock.DeactivateAndReset()
	}
}
