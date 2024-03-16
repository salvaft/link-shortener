package persistance

type Link struct {
	Href     string `json:"href"`
	B64_code string `json:"b64_code"`
	Url      string `json:"url"`
	Id       int    `json:"id"`
}
