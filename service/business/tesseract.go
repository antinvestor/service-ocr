package business

import (
	"context"
	gosseract2 "github.com/otiai10/gosseract/v2"
	"os"
)

type tesseract struct {
}

func (ts *tesseract) Recognize(ctx context.Context, image *os.File) (string, error) {

	localClient := gosseract2.NewClient()
	defer localClient.Close()

	err := localClient.SetImage(image.Name())
	if err != nil {
		return "", err
	}

	result, err := localClient.Text()
	if err != nil {
		return "", err
	}

	return result, err

}
