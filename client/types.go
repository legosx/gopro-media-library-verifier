package client

type page struct {
	Pages    pages    `json:"_pages"`
	Embedded embedded `json:"_embedded"`
}

type pages struct {
	CurrentPage int `json:"current_page"`
	PerPage     int `json:"per_page"`
	TotalItems  int `json:"total_items"`
	TotalPages  int `json:"total_pages"`
}

type embedded struct {
	Errors []interface{} `json:"errors"`
	Media  []media       `json:"media"`
}

type media struct {
	FileName string `json:"filename"`
	FileSize int64  `json:"file_size"`
}
