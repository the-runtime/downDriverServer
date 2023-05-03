package fileController

import (
	"context"
	"fmt"
	"github.com/conduitio/bwlimit"
	"golang.org/x/oauth2"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	"io"
	"os"
)

//using writer interface to get progress bar and to implement

type UploadCounter struct {
	Total uint64
}

func (uc *UploadCounter) Write(p []byte) (int, error) {
	n := len(p)
	uc.Total += uint64(n)
	uc.UpdateProgress()
	return n, nil
}

func (uc *UploadCounter) UpdateProgress() {
	globalCurrentUser.ConsumedDataTransfer += uc.Total / uint64(2) // for not counting upload and download separately
}

func UploadFile(token *oauth2.Token, googleOauthConfig *oauth2.Config, filename string) {

	fmt.Printf("uploading file %s\n", filename)

	ctx := context.Background()
	driveService, err := drive.NewService(ctx, option.WithTokenSource(googleOauthConfig.TokenSource(ctx, token)))
	if err != nil {
		fmt.Println(err.Error())
	}
	f := &drive.File{Name: filename}
	driveFile := driveService.Files.Create(f)
	localFile, err := os.Open("workingDirectory/" + filename)
	if err != nil {
		println(err.Error())
	}

	//apply upload speed limit to google drive
	readLimit := bwlimit.Byte(globalCurrentUser.AllowedSpeed) * bwlimit.MiB
	fileLimited := bwlimit.NewReader(localFile, readLimit)

	defer println("Upload complete")
	defer localFile.Close()
	defer os.Remove("workingDirectory/" + filename)

	//For applying  transfer limit

	counter := &UploadCounter{}
	_, err = driveFile.Media(io.TeeReader(fileLimited, counter)).Do()
	if err != nil {
		fmt.Println(err.Error())
	}

}
