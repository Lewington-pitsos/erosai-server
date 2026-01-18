package scanner

import (
	"io"
	"net/http"
	"os"

	"bitbucket.org/lewington/erosai-server/globals"
)

type imageSaver struct{}

func (s *imageSaver) save(URL string) (string, error) {
	// Ensure temp directory exists
	os.MkdirAll(globals.ImageSaveDirectory, 0755)

	response, err := http.Get(URL)

	if err != nil {
		return "", err
	}

	defer response.Body.Close()

	file, err := os.Create(globals.ImageSaveDirectory + "temp.jpg")

	if err != nil {
		return "", err
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		return "", err
	}

	return globals.ImageSaveDirectory + "temp.jpg", nil
}

func(s *imageSaver) delete() error {
	return os.Remove(globals.ImageSaveDirectory + "temp.jpg")
}