package models

type Blog struct {
	Title   string   `json:"title"`
	Tags    []string `json:"tags"`
	Content string   `json:"content"`
	Likes   int      `json:"likes"`
}
