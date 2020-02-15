package scanner

import (
	"testing"
	"bitbucket.org/lewington/erosai-server/assist"
	"bitbucket.org/lewington/erosai-server/globals"
)
func TestImageSaver(t *testing.T) {
	s := imageSaver{}

	name, err := s.save("https://cache5.pbwstatic.com/thumbnail/quuI5XljFTI/180x135/1.jpg")
	assist.Check(err)

	if name != globals.ImageSaveDirectory + "temp.jpg" {
		t.Fatalf("expected file to be named %v, got %v",globals.ImageSaveDirectory + "temp.jpg", name)
	}
}