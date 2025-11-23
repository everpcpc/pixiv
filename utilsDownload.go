package pixiv

import (
	"context"
	"net/http"
	"path"
	"strconv"
	"strings"
	"sync"

	"github.com/pkg/errors"
)

type Size int

const (
	SIZE_ORIGINAL Size = iota
	SIZE_LARGE
	SIZE_MEDIUM
	SIZE_SQUARE_MEDIUM
)

// func (i *Illust) urls() []string
func illustGetUrls(illust *Illust, size Size) []string {
	isSingle := illust.MetaSinglePage != nil && illust.MetaSinglePage.OriginalImageURL != ""
	var urls []string

	if isSingle {
		switch size {
		case SIZE_ORIGINAL:
			if illust.MetaSinglePage != nil {
				urls = []string{illust.MetaSinglePage.OriginalImageURL}
			}
		case SIZE_LARGE:
			urls = []string{illust.Images.Large}
		case SIZE_MEDIUM:
			urls = []string{illust.Images.Medium}
		case SIZE_SQUARE_MEDIUM:
			urls = []string{illust.Images.SquareMedium}
		}
	} else {
		urls = make([]string, 0, len(illust.MetaPages))
		for _, img := range illust.MetaPages {
			switch size {
			case SIZE_ORIGINAL:
				urls = append(urls, img.Images.Original)
			case SIZE_LARGE:
				urls = append(urls, img.Images.Large)
			case SIZE_MEDIUM:
				urls = append(urls, img.Images.Medium)
			case SIZE_SQUARE_MEDIUM:
				urls = append(urls, img.Images.SquareMedium)
			}
		}
	}
	return urls
}

var numberTestMap = [256]bool{
	'0': true, '1': true, '2': true, '3': true, '4': true,
	'5': true, '6': true, '7': true, '8': true, '9': true,
}

// urlToPageNum returns -1 when failed
func urlToPageNum(url string) int {
	name := path.Base(url) // "104714396_p0.jpg"
	pi := strings.Index(name, "_p")
	// extract numbers after "_p"
	start := pi + 2
	if start >= len(name) {
		return -1
	}

	end := start
	for numberTestMap[name[end]] {
		end++
	}
	if end == start { // at least one digit
		return -1
	}

	p, err := strconv.Atoi(name[start:end])
	if err != nil {
		return -1
	}

	return p
}

type imageType int

const (
	IMAGE_TYPE_UNKNOWN imageType = iota

	IMAGE_TYPE_WEBP
	IMAGE_TYPE_JPEG
	IMAGE_TYPE_PNG
	// [TODO] ugoria?
)

func (it imageType) String() string {
	switch it {
	case IMAGE_TYPE_UNKNOWN:
		return "unknown"
	case IMAGE_TYPE_WEBP:
		return "webp"
	case IMAGE_TYPE_JPEG:
		return "jpeg"
	case IMAGE_TYPE_PNG:
		return "png"
	default:
		return ""
	}
}

func parseImageType(s string) imageType {
	switch strings.ToLower(s) {
	case "webp":
		return IMAGE_TYPE_WEBP
	case "jpeg", "jpg":
		return IMAGE_TYPE_JPEG
	case "png":
		return IMAGE_TYPE_PNG
	default:
		return IMAGE_TYPE_UNKNOWN
	}
}

type Image struct {
	Data    []byte
	P       int
	Type    imageType
	TypeRaw string
}

func (i *Image) String() string {
	return "image/" + i.Type.String() + " (" + strconv.Itoa(len(i.Data)) + " bytes)"
}

type downloadItem struct {
	url string
	img *Image
	err chan error
}

func newDownloadItems(urls []string) (dis []*downloadItem) {
	dis = make([]*downloadItem, len(urls))
	for i, url := range urls {
		dis[i] = &downloadItem{
			url: url,
			img: &Image{},
			err: make(chan error, 1),
		}
	}
	return dis
}

type downloader struct {
	ctx     context.Context
	cancel  context.CancelFunc
	client  *http.Client
	threads int
	items   []*downloadItem
}

func newDownloader(ctx context.Context, client *http.Client, threads int, dis []*downloadItem) *downloader {
	ctx, cancel := context.WithCancel(ctx)
	return &downloader{
		ctx:     ctx,
		cancel:  cancel,
		client:  client,
		threads: threads,
		items:   dis,
	}
}

func (dl *downloader) startBackground() {
	go func() {
		limiter := newLimiter(dl.threads)
		defer limiter.close()

		for i, item := range dl.items {
			i, item := i, item // [TODO] <=go1.21
			select {
			case <-dl.ctx.Done():
				return
			case limiter.acquire() <- struct{}{}:
			}

			go func() {
				defer limiter.release()

				img, err := downloadImage(dl.ctx, dl.client, item.url)
				if err == nil {
					if img.P == -1 {
						img.P = i // [urlToPageNum] failed, fallback
					}
					*item.img = img
				}
				item.err <- err // mark an item as done
			}()
		}
	}()
}

func (dl *downloader) downloadIter() <-chan DownloadResult {
	dl.startBackground()
	ch := make(chan DownloadResult)
	go func() {
		defer close(ch)
		for _, item := range dl.items {
			var err error
			var ok bool
			select {
			case <-dl.ctx.Done():
				return
			case err, ok = <-item.err:
				if !ok {
					continue
				}
			}

			if err != nil {
				ch <- DownloadResult{Err: errors.Wrapf(err, "download url %s failed", item.url)}
			} else {
				ch <- DownloadResult{Img: *item.img}
			}
		}
	}()
	return ch
}

func (dl *downloader) downloadTo(f func(int, DownloadResult)) {
	dl.startBackground()
	for i, item := range dl.items {
		i, item := i, item // [TODO] <=go1.21
		go func() {
			// there may be over a hundred goroutines(depend on illust)
			// waiting for item.err at the same time,
			// but should not be a big deal
			var err error
			var ok bool
			select {
			case <-dl.ctx.Done():
				return
			case err, ok = <-item.err:
				if !ok {
					return
				}
			}

			if err != nil {
				f(i, DownloadResult{Err: errors.Wrapf(err, "download url %s failed", item.url)})
			} else {
				f(i, DownloadResult{Img: *item.img})
			}
		}()
	}
}

type limiter struct {
	sem  chan struct{}
	once sync.Once
}

func newLimiter(threads int) *limiter {
	return &limiter{
		sem: make(chan struct{}, threads),
	}
}

func (l *limiter) acquire() chan<- struct{} {
	return l.sem
}

func (l *limiter) release() {
	<-l.sem
}

func (l *limiter) close() {
	l.once.Do(func() {
		close(l.sem)
	})
}
