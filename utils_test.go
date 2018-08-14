package pixiv

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseNextPageURL(t *testing.T) {
	r := require.New(t)
	next, err := parseNextPageOffset("https://app-api.pixiv.net/v2/illust/follow?restrict=public&offset=30")
	r.Nil(err)
	r.Equal(30, next)
}
