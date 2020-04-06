package pixiv

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAuth(t *testing.T) {
	r := require.New(t)
	_, err := Login("username", "password")
	if nil != err{
		r.EqualError(err, "Login system error: 103:pixiv ID、またはメールアドレス、パスワードが正しいかチェックしてください。")
	}
}
