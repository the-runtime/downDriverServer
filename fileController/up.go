package fileController

import (
	"context"
	"fmt"
	"golang.org/x/oauth2"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	"os"
)

func UploadFile(token *oauth2.Token, googleOauthConfig *oauth2.Config, filename string) {

	fmt.Printf("uploading file %s\n", filename)

	ctx := context.Background()
	driveService, err := drive.NewService(ctx, option.WithTokenSource(googleOauthConfig.TokenSource(ctx, token)))
	if err != nil {
		fmt.Println(err.Error())
	}
	f := &drive.File{Name: filename}
	driveFile := driveService.Files.Create(f)
	localFile, err := os.Open(filename)
	if err != nil {
		println(err.Error())
	}
	defer localFile.Close()

	_, err = driveFile.Media(localFile).Do()
	if err != nil {
		fmt.Println(err.Error())
	}
}
