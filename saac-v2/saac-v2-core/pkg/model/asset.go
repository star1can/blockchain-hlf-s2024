package model

type Asset struct {
	ID             int    `json:"id"`
	AppraisedValue int    `json:"appraisedValue"`
	Color          string `json:"color"`
	Size           int    `json:"size"`
	Type           string `json:"type"`
	Owner          string `json:"owner"`
}
