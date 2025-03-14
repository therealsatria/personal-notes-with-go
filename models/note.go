package models

type Note struct {
	ID         string `json:"id"`
	Subject    string `json:"subject"`
	Content    string `json:"content"`
	Priority   string `json:"priority"`
	Tags       string `json:"tags"`
	CategoryID string `json:"category_id"`
}
