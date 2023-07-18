package pixiv

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseNextPageURL(t *testing.T) {
	r := require.New(t)
	next, err := parseNextPageOffset("https://app-api.pixiv.net/v1/user/bookmarks/illust?filter=for_ios&restrict=private&user_id=60984430&max_bookmark_id=21354656694", OffsetFieldMaxBookmarkID)
	r.Nil(err)
	r.Equal(21354656694, next)
}
