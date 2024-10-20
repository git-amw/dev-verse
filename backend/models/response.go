package models

type TagSearchResponse struct {
	TagId        uint   `json:"id"`
	TagValue     string `json:"tag_value"`
	BlogsWithTag int    `json:"blogs_with_tag"`
}

type BlogSearchResponse struct {
	ID       uint               `json:"id"`
	Title    string             `json:"title"`
	Content  string             `json:"content"`
	Likes    int                `json:"likes"`
	BlogTags []BlogTagsResponse `json:"blogTags"`
}

type BlogTagsResponse struct {
	TagId int `json:"tag_id"`
}
