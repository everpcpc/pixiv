package pixiv

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAppPixivAPIUserDetail(t *testing.T) {
	r := require.New(t)
	_, err := Login("x", "x")
	r.Nil(err)
}
