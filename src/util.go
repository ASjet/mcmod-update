package src

import (
	"io"
	"net/http"
	"os"
)

func Get(url, name string) error {
	if len(url) == 0 {
		return nil
	}

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	f, err := os.Create(name)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, resp.Body)
	return err
}

// Check if a file exist
func FileExist(path string) bool {
	// Skip empty path
	if len(path) == 0 {
		return true
	}
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}
