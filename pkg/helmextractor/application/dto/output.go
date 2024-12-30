package dto

type ImageDetails struct {
	Name   string `json:"name"`
	Size   int64  `json:"size"`
	Layers int    `json:"layers"`
}

type Response struct {
	Images []ImageDetails `json:"images"`
}
