package model

type NodeSelector struct {
	Selector string `json:"selector"`
	Parent   bool   `json:"parent"`
	Many     bool   `json:"many"`
}
