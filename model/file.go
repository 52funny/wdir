package model

// FileInfo is description filename and ...
type FileInfo struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Date string `json:"date"`
	Size string `json:"size"`
}

// File is description filepath and fileinfo
type File struct {
	Path     string   `json:"path"`
	Fileinfo FileInfo `json:"fileinfo"`
}
