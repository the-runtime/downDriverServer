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
	"serverFordownDrive/database"
	"serverFordownDrive/model"
)

// using writer interface to get progress bar and to implement
var downloaded_value uint64

type UploadCounter struct {
	Total uint64
}

func (uc *UploadCounter) Write(p []byte) (int, error) {
	n := len(p)
	uc.Total += uint64(n)
	uc.PrintProgress()
	//uc.UpdateProgress()
	return n, nil
}
func (uc *UploadCounter) PrintProgress() {
	GlobalCurrentUser.ConsumedDataTransfer = uc.Total // for not counting upload and download separately
	globalProgresscounter.Transferred = downloaded_value + uc.Total/uint64(2)
	//println("flobalProgress %d", globalProgresscounter.Transferred)
	//fmt.Printf("\r%s", strings.Repeat(" ", 35))
	//fmt.Printf("\rUploading... %s complete", humanize.Bytes(uc.Total))
}

//func (uc *UploadCounter) UpdateProgress() {
//	GlobalCurrentUser.ConsumedDataTransfer += uc.Total / uint64(2) // for not counting upload and download separately
//	globalProgresscounter.Transferred += uc.Total / uint64(2)
//}

func UploadFile(token *oauth2.Token, googleOauthConfig *oauth2.Config, filename string, tempUser *model.User, progressId int) {

	GlobalCurrentUser = tempUser

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
	readLimit := bwlimit.Byte(GlobalCurrentUser.AllowedSpeed) * bwlimit.MiB
	fileLimited := bwlimit.NewReader(localFile, readLimit)

	defer localFile.Close()
	defer os.Remove("workingDirectory/" + filename)

	//For applying  transfer limit

	downloaded_value = globalProgresscounter.Total
	counter := &UploadCounter{}
	_, err = driveFile.Media(io.TeeReader(fileLimited, counter)).Do()
	if err != nil {
		fmt.Println(err.Error())
	}
	println("Upload complete")
	globalProgresscounter.IsOn = false
	println("IsOn is :", globalProgresscounter.IsOn, "\nProgressId is", globalProgresscounter.ProcessId)
	userdb, err := database.NewUserDb()
	if err != nil {
		println(err.Error())
	}

	userdb.Model(&model.User{}).Where("user_id=?", GlobalCurrentUser.UserId).Update("consumed_data_transfer", GlobalCurrentUser.ConsumedDataTransfer)
	println("updated consumed data transfer in database %d", GlobalCurrentUser.ConsumedDataTransfer)

}
