# pixiv

![test](https://github.com/everpcpc/pixiv/workflows/test/badge.svg)
[![codecov](https://codecov.io/gh/everpcpc/pixiv/branch/master/graph/badge.svg)](https://codecov.io/gh/everpcpc/pixiv)
[![Go Report Card](https://goreportcard.com/badge/github.com/everpcpc/pixiv)](https://goreportcard.com/report/github.com/everpcpc/pixiv)
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
