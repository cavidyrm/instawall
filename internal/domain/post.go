package domain

type Post struct {
	ID    uint64 `json:"id"`
	Title string `json:"title"`
	Body  string `json:"body"`
	Meta  string `json:"meta"`
}
