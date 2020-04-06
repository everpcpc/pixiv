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
		return 0, fmt.Errorf("parse next_url error: %s {%s}", s, err)
	}

	offset, err := strconv.Atoi(m.Get("offset"))

	if err != nil {
		fmt.Errorf("parse next_url error: %s {%s}\n", s, err)
		offset, err = strconv.Atoi(m.Get("max_bookmark_id"))
	}//add max_bookmark_id check

	if nil != err{
		return 0, fmt.Errorf("parse next_url error: %s {%s}", s, err)
	}
	return offset, nil
}
