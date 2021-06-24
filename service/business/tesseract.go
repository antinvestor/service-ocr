package business

import (
	"context"
	"github.com/otiai10/gosseract"
	"os"
)

type tesseract struct {
}

func (ts *tesseract) Recognize(ctx context.Context, image *os.File) (string, error) {

	localClient := gosseract.NewClient()
	defer localClient.Close()

	err := localClient.SetImage(image.Name())
	if err != nil {
		return "", err
	}

	return localClient.Text()

}
