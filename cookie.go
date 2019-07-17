package chromedp

import (
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
)

func CookieParamsFromCookies(cookies []*network.Cookie) []*network.CookieParam {
	ret := []*network.CookieParam{}
	for _, c := range cookies {
		expr := cdp.TimeSinceEpoch(time.Unix(int64(c.Expires), 0))

		ret = append(ret, &network.CookieParam{
			Name:     c.Name,
			Value:    c.Value,
			URL:      "",
			Domain:   c.Domain,
			Path:     c.Path,
			Secure:   c.Secure,
			HTTPOnly: c.HTTPOnly,
			SameSite: c.SameSite,
			Expires:  &expr,
		})
	}
	return ret
}
