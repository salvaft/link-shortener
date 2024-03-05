package persistance

import "net/url"

type Link struct {
	Href     *url.URL
	B64_code string
	Url      string
	Id       int
}
