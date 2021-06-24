package business

import (
	"context"
	"os"
	"strings"
	"testing"
)

func SkpTest_tesseract_Recognize(t *testing.T) {

	type args struct {
		ctx  context.Context
		file *os.File
		lang string
	}

	//Found on google : https://www.google.com/url?sa=i&url=https%3A%2F%2Ftwitter.com%2Fmirimuneuer16%2Fstatus%2F913049416935845893&psig=AOvVaw3Hvy7yp1c-zh-_xCgdwNyb&ust=1624522512079000&source=images&cd=vfe&ved=0CAoQjRxqFwoTCJD8xIyorfECFQAAAAAdAAAAABAI
	testFile1, err := os.Open("../../tests_runner/test_id_card.jpg")
	if err != nil {
		getwd, _ := os.Getwd()

		t.Errorf("useGosseractForOCR() error = %v, could not get file in %v", err, getwd)
		return
	}

	// Found on : https://nairobinews.nation.co.ke/editors-picks/man-beheaded-in-broad-daylight-in-kampala
	testFile2, err := os.Open("../../tests_runner/test2_id_card.jpg")
	if err != nil {
		t.Errorf("useGosseractForOCR() error = %v, could not get file", err)
		return
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{name: "Normal image",
			args: args{
				ctx:  context.Background(),
				file: testFile2,
				lang: "eng",
			},
			want: []string{"SIGNATURE", "ACHAN", "BEKI", "NIN", "CF95071101J9JH", "CARD NO", "001985238"},
		},

		{name: "Another Normal image",
			args: args{
				ctx:  context.Background(),
				file: testFile1,
				lang: "eng",
			},
			want: []string{"NATIONAL", "FRED", "MUGISHA", "NIN", "CM65106105MFLD", "CARD NO", "018565889"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := &tesseract{}
			got, err := ts.Recognize(tt.args.ctx, tt.args.file)
			if (err != nil) != tt.wantErr {
				t.Errorf("useGosseractForOCR() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			for _, w := range tt.want {
				if !strings.Contains(got, w) {
					t.Errorf("useGosseractForOCR() got = %v, does not contain : %v", got, w)
				}
			}

		})
	}
}
