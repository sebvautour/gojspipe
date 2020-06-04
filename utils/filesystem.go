package utils

import "io/ioutil"

type FileContent struct {
	Content interface{}
	Error   string
}

// ReadFile opens and reads a file
// Returns FileContent with Content in string, or Content nil with Error string
func (u *Utils) ReadFile(filename string) FileContent {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return FileContent{
			Content: nil,
			Error:   err.Error(),
		}
	}
	return FileContent{
		Content: string(b),
	}
}
