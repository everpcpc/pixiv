package pixiv

import (
	"fmt"
	"net/url"
	"strconv"
)

func parseNextPageOffset(s string) (int, error) {
	if s == "" {
		return 0, nil
	}
	u, err := url.Parse(s)
	if err != nil {
		return 0, fmt.Errorf("parse next_url error: %s {%s}", s, err)
	}

	m, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		return 0, fmt.Errorf("parse next_url raw query error: %s {%s}", s, err)
	}

	offsetParam := m.Get("max_bookmark_id")
	if offsetParam == "" {
		return 0, nil
	}

	offset, err := strconv.Atoi(offsetParam)
	if err != nil {
		return 0, fmt.Errorf("getting offset from url: %s {%s}", s, err)
	}
	return offset, nil
}
