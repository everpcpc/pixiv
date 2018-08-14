# pixiv

[![build](https://travis-ci.org/everpcpc/pixiv.svg)](https://travis-ci.org/everpcpc/pixiv)
[![godoc](https://img.shields.io/badge/godoc-reference-5272B4.svg)](https://godoc.org/github.com/everpcpc/pixiv)

Pixiv API for Golang (with Auth supported)

Inspired by [pixivpy](https://github.com/upbit/pixivpy)

## example

```golang
account, err := pixiv.Login("username", "password")
app := pixiv.NewApp()
user, err := app.UserDetail(uid)
illusts, next, err := app.UserIllusts(uid, "illust", 0)
illusts, next, err := app.UserBookmarksIllust(uid, "public", 0, "")
illusts, next, err := app.IllustFollow("public", 0)
```
