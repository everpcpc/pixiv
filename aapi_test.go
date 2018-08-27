package pixiv

import (
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var (
	tokenSet bool
	testUID  uint64
)

func initTest() (err error) {
	if tokenSet {
		return
	}
	LoadAuth(os.Getenv("TOKEN"), os.Getenv("REFRESH_TOKEN"), time.Time{})
	testUID, err = strconv.ParseUint(os.Getenv("TEST_UID"), 10, 0)
	tokenSet = true
	return err
}

func TestUserDetail(t *testing.T) {
	r := require.New(t)
	r.Nil(initTest())
	app := NewApp()
	detail, err := app.UserDetail(testUID)
	r.Nil(err)
	r.Equal(testUID, detail.User.ID)
}

func TestUserIllusts(t *testing.T) {
	r := require.New(t)
	r.Nil(initTest())
	app := NewApp()
	illusts, next, err := app.UserIllusts(490219, "illust", 0)
	r.Nil(err)
	r.Len(illusts, 30)
	r.Equal(30, next)
}

func TestUserBookmarksIllust(t *testing.T) {
	r := require.New(t)
	r.Nil(initTest())
	app := NewApp()
	illusts, _, err := app.UserBookmarksIllust(testUID, "public", 0, "")
	r.Nil(err)
	r.Equal(uint64(70095856), illusts[0].ID)

}

func TestIllustFollow(t *testing.T) {
	r := require.New(t)
	r.Nil(initTest())
	app := NewApp()
	illusts, next, err := app.IllustFollow("public", 0)
	r.Nil(err)
	r.Len(illusts, 30)
	r.Equal(30, next)
}

func TestIllustDetail(t *testing.T) {
	r := require.New(t)
	r.Nil(initTest())
	app := NewApp()
	illust, err := app.IllustDetail(70095856)
	r.Nil(err)
	r.Equal(uint64(70095856), illust.ID)
}

func TestDownload(t *testing.T) {
	r := require.New(t)
	r.Nil(initTest())
	app := NewApp()
	size, err := app.Download(70095856, ".")
	r.Nil(err)
	r.Equal(int64(993598), size)
}
