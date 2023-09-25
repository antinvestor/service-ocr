package business

import (
	"context"
	"encoding/base64"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/vision/v1"
	"io/ioutil"
	"os"
)

type googleCloudVision struct {
}

func (gcv *googleCloudVision) Recognize(ctx context.Context, image *os.File) (string, error) {

	// Authenticate to generate a vision service
	client, err := google.DefaultClient(ctx, vision.CloudPlatformScope)
	if err != nil {
		return "", err
	}

	service, err := vision.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return "", err
	}

	imageContent, err := ioutil.ReadAll(image)
	if err != nil {
		return "", err
	}

	// We now have a Vision API service with which we can make API calls.
	imageRequestList := make([]*vision.AnnotateImageRequest, 0)

	// Construct a text request, encoding the image in base64.
	req := &vision.AnnotateImageRequest{
		// Apply image which is encoded by base64
		Image: &vision.Image{
			Content: base64.StdEncoding.EncodeToString(imageContent),
		},
		// Apply features to indicate what type of image detection
		Features: []*vision.Feature{
			{
				Type: "TEXT_DETECTION",
			},
		},
	}

	imageRequestList = append(imageRequestList, req)

	batch := &vision.BatchAnnotateImagesRequest{
		Requests: imageRequestList,
	}

	res, err := service.Images.Annotate(batch).Do()
	if err != nil {
		return "", err
	}

	result := ""
	for _, resp := range res.Responses {
		result = result + resp.FullTextAnnotation.Text + " \n"
	}
	return result, nil
}
