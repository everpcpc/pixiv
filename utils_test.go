package pixiv

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseBookmarkNextPageURL(t *testing.T) {
	r := require.New(t)
	next, err := parseNextPageOffset("https://app-api.pixiv.net/v1/user/bookmarks/illust?filter=for_ios&restrict=private&user_id=60984430&max_bookmark_id=21354656694", OffsetFieldMaxBookmarkID)
	r.Nil(err)
	r.Equal(21354656694, next)
}

func TestParseOffsetNextPageURL(t *testing.T) {
	r := require.New(t)
	next, err := parseNextPageOffset("https://app-api.pixiv.net/v1/user/illusts?filter=for_ios&type=illust&user_id=490219&offset=30", OffsetFieldOffset)
	r.Nil(err)
	r.Equal(30, next)
}

func TestParseEmptyNextPageURL(t *testing.T) {
	r := require.New(t)
	next, err := parseNextPageOffset("", OffsetFieldOffset)
	r.Nil(err)
	r.Equal(0, next)
}

func TestParseInvalidNextPageURL(t *testing.T) {
	r := require.New(t)
	next, err := parseNextPageOffset("https://app-api.pixiv.net/v1/user/illusts?filter=for_ios&type=illust&user_id=490219&offset=30", "invalid")
	r.EqualError(err, "offset param omitted: invalid")
	r.Equal(0, next)
}
