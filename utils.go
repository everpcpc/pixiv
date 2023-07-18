package pixiv

import (
	"fmt"
	"net/url"
	"strconv"
)

// parseNextPageOffset parses next_url and returns the offset
//
// field is either "max_bookmark_id" or "offset"
func parseNextPageOffset(s, field string) (int, error) {
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

	offsetParam := m.Get(field)
	if offsetParam == "" {
		return 0, fmt.Errorf("offset param omitted: %s", field)
	}

	offset, err := strconv.Atoi(offsetParam)
	if err != nil {
		return 0, fmt.Errorf("getting offset from url: %s {%s}", s, err)
	}
	return offset, nil
}
