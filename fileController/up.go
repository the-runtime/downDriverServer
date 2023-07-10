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
	"serverFordownDrive/controller"
	"serverFordownDrive/database"
	"serverFordownDrive/model"
	"time"
)

// using writer interface to get progress bar and to implement
//var downloaded_value uint64

type UploadCounter struct {
	Total    uint64
	Progress *controller.Progress
}

func (uc *UploadCounter) Write(p []byte) (int, error) {
	n := len(p)
	uc.Total += uint64(n) / 2
	uc.PrintProgress()
	//uc.UpdateProgress()
	return n, nil
}
func (uc *UploadCounter) PrintProgress() {
	//GlobalCurrentUser.ConsumedDataTransfer = uc.Total // for not counting upload and download separately
	uc.Progress.Transferred = uc.Total + uc.Progress.Total/2
	//println("flobalProgress %d", globalProgresscounter.Transferred)
	//fmt.Printf("\r%s", strings.Repeat(" ", 35))
	//fmt.Printf("\rUploading... %s complete", humanize.Bytes(uc.Total))
}

//func (uc *UploadCounter) UpdateProgress() {
//	GlobalCurrentUser.ConsumedDataTransfer += uc.Total / uint64(2) // for not counting upload and download separately
//	globalProgresscounter.Transferred += uc.Total / uint64(2)
//}

func UploadFile(token *oauth2.Token, googleOauthConfig *oauth2.Config, filename string, tempUser *model.User, progressId int) {

	//GlobalCurrentUser = tempUser

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
	readLimit := bwlimit.Byte(tempUser.AllowedSpeed) * bwlimit.MiB
	fileLimited := bwlimit.NewReader(localFile, readLimit)

	defer localFile.Close()
	defer os.Remove("workingDirectory/" + filename)

	//For applying  transfer limit

	tempProgress := controller.GetProgressById(tempUser.UserId, progressId)

	counter := &UploadCounter{0, tempProgress}
	_, err = driveFile.Media(io.TeeReader(fileLimited, counter)).Do()
	if err != nil {
		fmt.Println(err.Error())
	}
	println("Upload complete")
	tempProgress.IsOn = false
	println("IsOn is :", tempProgress.IsOn, "\nProgressId is", tempProgress.ProcessId)
	println("progress id  is ", progressId)

	userdb, err := database.NewUserDb()
	if err != nil {
		println(err.Error())
	}
	historyDb, err := database.NewHistoryDb()
	if err != nil {
		println(err.Error())
		return
	}

	///consumedDataUser := model.User{}

	history := model.SingleHistory{
		UserId:     tempUser.UserId,
		Filename:   filename,
		Filesize:   tempProgress.Total,
		Finishedat: time.Now(),
	}

	historyDb.Model(&model.SingleHistory{}).Create(&history)
	userdb.Model(&model.User{}).Where("user_id=?", tempUser.UserId).Update("consumed_data_transfer", tempProgress.Total+tempUser.ConsumedDataTransfer)

	//remove the progress for progressTable
	//controller.EndProgress(tempProgress.UserId, tempProgress.ProcessId)
	//	will implement this latter

}
