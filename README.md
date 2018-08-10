# pixiv [![build](https://travis-ci.org/everpcpc/pixiv.svg)](https://travis-ci.org/everpcpc/pixiv)

Pixiv API for Golang (with Auth supported)

Inspired by [pixivpy](https://github.com/upbit/pixivpy)

## example

```golang
account, err := pixiv.Login("username", "password")
app := pixiv.NewApp()
user, err := app.UserDetail(uid)
illusts, err := app.UserIllusts(uid, "illust", 0)
illusts, err := app.UserBookmarksIllust(uid, "public", 0, "")
illusts, err := app.IllustFollow("public", 0)
```
