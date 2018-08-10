package pixiv

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAuth(t *testing.T) {
	r := require.New(t)
	_, err := Login("x", "x")
	r.EqualError(err, "Login system error: 103:pixiv ID、またはメールアドレス、パスワードが正しいかチェックしてください。")

	// NOTE: ignore login success test
	// user, err := Login("", "")
	// r.Nil(err)
	// r.Equal("", user.ID)

	// r.Nil(refreshToken())
}
