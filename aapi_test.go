package pixiv

import (
	"fmt"
	"testing"
)

func TestAppPixivAPIUserDetail(t *testing.T) {
	aapi := NewApp()
	fmt.Println(aapi.UserDetail("6101418"))
	panic("fdsa")
}
