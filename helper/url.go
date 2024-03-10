package helper

import "net/url"

func IsValidURL(urlString string) bool {
	url, err := url.ParseRequestURI(urlString)
	if err != nil {
		return false
	}

	switch url.Scheme {
	case "http", "https":
	default:
		return false
	}

	return true
}
