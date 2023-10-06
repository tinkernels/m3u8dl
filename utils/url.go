package utils

import "net/url"

func IsUrlSchemeHttp(urlP *url.URL) bool {
	return urlP.Scheme == "http" || urlP.Scheme == "https"
}
