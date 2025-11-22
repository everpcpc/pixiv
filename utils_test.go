package pixiv

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseBookmarkNextPageURL(t *testing.T) {
	r := require.New(t)
	next, err := parseNextPageOffset("https://app-api.pixiv.net/v1/user/bookmarks/illust?filter=for_ios&restrict=private&user_id=60984430&max_bookmark_id=21354656694", OffsetFieldMaxBookmarkID)
	r.Nil(err)
	r.Equal(21354656694, next)
}

func TestParseOffsetNextPageURL(t *testing.T) {
	r := require.New(t)
	next, err := parseNextPageOffset("https://app-api.pixiv.net/v1/user/illusts?filter=for_ios&type=illust&user_id=490219&offset=30", OffsetFieldOffset)
	r.Nil(err)
	r.Equal(30, next)
}

func TestParseEmptyNextPageURL(t *testing.T) {
	r := require.New(t)
	next, err := parseNextPageOffset("", OffsetFieldOffset)
	r.Nil(err)
	r.Equal(0, next)
}

func TestParseInvalidNextPageURL(t *testing.T) {
	r := require.New(t)
	next, err := parseNextPageOffset("https://app-api.pixiv.net/v1/user/illusts?filter=for_ios&type=illust&user_id=490219&offset=30", "invalid")
	r.EqualError(err, "offset param omitted: invalid")
	r.Equal(0, next)
}

func TestUrlToPageNum(t *testing.T) {
	r := require.New(t)
	urls := [...]string{
		"https://i.pximg.net/c/360x360_10_webp/img-master/img/2023/01/22/10/02/20/104714396_p0_square1200.jpg",
		"https://i.pximg.net/c/600x1200_90_webp/img-master/img/2023/01/22/10/02/20/104714396_p1_master1200.jpg",
		"https://i.pximg.net/img-original/img/2023/01/22/10/02/20/104714396_p2.jpg",
	}
	for i, url := range urls {
		r.Equal(i, urlToPageNum(url))
	}
}

func TestIllustGetUrlsMulti(t *testing.T) {
	r := require.New(t)
	illust := Illust{}
	err := json.Unmarshal([]byte(TestIllustGetUrlsMultiJsonSample), &illust)
	r.Nil(err)

	expectedOriginal := []string{
		"https://i.pximg.net/img-original/img/2023/01/22/10/02/20/104714396_p0.jpg",
		"https://i.pximg.net/img-original/img/2023/01/22/10/02/20/104714396_p1.jpg",
		"https://i.pximg.net/img-original/img/2023/01/22/10/02/20/104714396_p2.jpg",
	}
	urlsOriginal := illustGetUrls(&illust, SIZE_ORIGINAL)
	r.Equal(expectedOriginal, urlsOriginal)

	expectedLarge := []string{
		"https://i.pximg.net/c/600x1200_90_webp/img-master/img/2023/01/22/10/02/20/104714396_p0_master1200.jpg",
		"https://i.pximg.net/c/600x1200_90_webp/img-master/img/2023/01/22/10/02/20/104714396_p1_master1200.jpg",
		"https://i.pximg.net/c/600x1200_90_webp/img-master/img/2023/01/22/10/02/20/104714396_p2_master1200.jpg",
	}
	urlsLarge := illustGetUrls(&illust, SIZE_LARGE)
	r.Equal(expectedLarge, urlsLarge)

	expectedMedium := []string{
		"https://i.pximg.net/c/540x540_70/img-master/img/2023/01/22/10/02/20/104714396_p0_master1200.jpg",
		"https://i.pximg.net/c/540x540_70/img-master/img/2023/01/22/10/02/20/104714396_p1_master1200.jpg",
		"https://i.pximg.net/c/540x540_70/img-master/img/2023/01/22/10/02/20/104714396_p2_master1200.jpg",
	}
	urlsMedium := illustGetUrls(&illust, SIZE_MEDIUM)
	r.Equal(expectedMedium, urlsMedium)

	expectedSquareMedium := []string{
		"https://i.pximg.net/c/360x360_10_webp/img-master/img/2023/01/22/10/02/20/104714396_p0_square1200.jpg",
		"https://i.pximg.net/c/360x360_10_webp/img-master/img/2023/01/22/10/02/20/104714396_p1_square1200.jpg",
		"https://i.pximg.net/c/360x360_10_webp/img-master/img/2023/01/22/10/02/20/104714396_p2_square1200.jpg",
	}
	urlsSquareMedium := illustGetUrls(&illust, SIZE_SQUARE_MEDIUM)
	r.Equal(expectedSquareMedium, urlsSquareMedium)
}

func TestIllustGetUrlsSingle(t *testing.T) {
	r := require.New(t)
	illust := Illust{}
	err := json.Unmarshal([]byte(TestIllustGetUrlsSingleJsonSample), &illust)
	r.Nil(err)

	expectedOriginal := []string{
		"https://i.pximg.net/img-original/img/2023/01/31/20/06/17/104968552_p0.jpg",
	}
	urlsOriginal := illustGetUrls(&illust, SIZE_ORIGINAL)
	r.Equal(expectedOriginal, urlsOriginal)

	expectedLarge := []string{
		"https://i.pximg.net/c/600x1200_90_webp/img-master/img/2023/01/31/20/06/17/104968552_p0_master1200.jpg",
	}
	urlsLarge := illustGetUrls(&illust, SIZE_LARGE)
	r.Equal(expectedLarge, urlsLarge)

	expectedMedium := []string{
		"https://i.pximg.net/c/540x540_70/img-master/img/2023/01/31/20/06/17/104968552_p0_master1200.jpg",
	}
	urlsMedium := illustGetUrls(&illust, SIZE_MEDIUM)
	r.Equal(expectedMedium, urlsMedium)

	expectedSquareMedium := []string{
		"https://i.pximg.net/c/540x540_10_webp/img-master/img/2023/01/31/20/06/17/104968552_p0_square1200.jpg",
	}
	urlsSquareMedium := illustGetUrls(&illust, SIZE_SQUARE_MEDIUM)
	r.Equal(expectedSquareMedium, urlsSquareMedium)
}

const TestIllustGetUrlsMultiJsonSample = `{
    "image_urls": {
        "square_medium": "https://i.pximg.net/c/540x540_10_webp/img-master/img/2023/01/22/10/02/20/104714396_p0_square1200.jpg",
        "medium": "https://i.pximg.net/c/540x540_70/img-master/img/2023/01/22/10/02/20/104714396_p0_master1200.jpg",
        "large": "https://i.pximg.net/c/600x1200_90_webp/img-master/img/2023/01/22/10/02/20/104714396_p0_master1200.jpg",
        "original": ""
    },
    "meta_single_page": {
        "original_image_url": ""
    },
    "meta_pages": [
        {
            "image_urls": {
                "square_medium": "https://i.pximg.net/c/360x360_10_webp/img-master/img/2023/01/22/10/02/20/104714396_p0_square1200.jpg",
                "medium": "https://i.pximg.net/c/540x540_70/img-master/img/2023/01/22/10/02/20/104714396_p0_master1200.jpg",
                "large": "https://i.pximg.net/c/600x1200_90_webp/img-master/img/2023/01/22/10/02/20/104714396_p0_master1200.jpg",
                "original": "https://i.pximg.net/img-original/img/2023/01/22/10/02/20/104714396_p0.jpg"
            }
        },
        {
            "image_urls": {
                "square_medium": "https://i.pximg.net/c/360x360_10_webp/img-master/img/2023/01/22/10/02/20/104714396_p1_square1200.jpg",
                "medium": "https://i.pximg.net/c/540x540_70/img-master/img/2023/01/22/10/02/20/104714396_p1_master1200.jpg",
                "large": "https://i.pximg.net/c/600x1200_90_webp/img-master/img/2023/01/22/10/02/20/104714396_p1_master1200.jpg",
                "original": "https://i.pximg.net/img-original/img/2023/01/22/10/02/20/104714396_p1.jpg"
            }
        },
        {
            "image_urls": {
                "square_medium": "https://i.pximg.net/c/360x360_10_webp/img-master/img/2023/01/22/10/02/20/104714396_p2_square1200.jpg",
                "medium": "https://i.pximg.net/c/540x540_70/img-master/img/2023/01/22/10/02/20/104714396_p2_master1200.jpg",
                "large": "https://i.pximg.net/c/600x1200_90_webp/img-master/img/2023/01/22/10/02/20/104714396_p2_master1200.jpg",
                "original": "https://i.pximg.net/img-original/img/2023/01/22/10/02/20/104714396_p2.jpg"
            }
        }
    ]
}`

const TestIllustGetUrlsSingleJsonSample = `{
    "image_urls": {
        "square_medium": "https://i.pximg.net/c/540x540_10_webp/img-master/img/2023/01/31/20/06/17/104968552_p0_square1200.jpg",
        "medium": "https://i.pximg.net/c/540x540_70/img-master/img/2023/01/31/20/06/17/104968552_p0_master1200.jpg",
        "large": "https://i.pximg.net/c/600x1200_90_webp/img-master/img/2023/01/31/20/06/17/104968552_p0_master1200.jpg",
        "original": ""
    },
    "meta_single_page": {
        "original_image_url": "https://i.pximg.net/img-original/img/2023/01/31/20/06/17/104968552_p0.jpg"
    },
    "meta_pages": []
}`
