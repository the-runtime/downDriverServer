package fileController

import (
	"fmt"
	"github.com/conduitio/bwlimit"
	"io"
	"net"
	"net/http"
	"os"
	"serverFordownDrive/controller"
	"serverFordownDrive/model"
	"strings"
	"time"
)

//var GlobalCurrentUser *model.User
//var globalProgresscounter *controller.Progress

//for implementing transfer limit on user

type WriteCounter struct {
	Total    uint64
	progress *controller.Progress
}

func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Total += uint64(n) / 2
	//wc.PrintProgress()
	wc.UpdateProgress()
	return n, nil
}
func (wc *WriteCounter) PrintProgress() {
	//fmt.Printf("\r%s", strings.Repeat(" ", 35))
	//fmt.Printf("\rDownloading... %s complete", humanize.Bytes(wc.Total))
	wc.progress.Transferred = wc.Total
	//fmt.Printf("id is %d transferd data is %d ", wc.progress.ProcessId, wc.progress.Transferred)
	//println("filename : %s\n", globalProgresscounter.Filename, "Downloaded: %d\n", globalProgresscounter.Transferred)
	//println("globalProgresscounter.Transferred   %d", globalProgresscounter.Transferred)
}

func (wc *WriteCounter) UpdateProgress() {
	wc.progress.Transferred = wc.Total
}

func StartDown(url string, CurrenUser *model.User, progressId int) (string, int) {
	//client := http.Client{}

	//GlobalCurrentUser = CurrenUser
	urlSplit := strings.Split(url, "/")
	filename := urlSplit[len(urlSplit)-1]
	//res, err := http.Head(url)
	//if err != nil {
	//	fmt.Println(err.Error())
	//	return "", 0
	//}

	f, err := os.Create("workingDirectory/" + filename)
	if err != nil {
		fmt.Println(err.Error())
		return "", 0
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			fmt.Println(err.Error())
			return
		}

	}(f)

	//writeLimit := bwlimit.Byte(CurrenUser.AllowedSpeed) * bwlimit.MiB
	//readLimit := bwlimit.Byte(CurrenUser.AllowedSpeed) * bwlimit.MiB
	writeLimit := bwlimit.Byte(CurrenUser.AllowedSpeed) * bwlimit.MiB
	readLimit := bwlimit.Byte(CurrenUser.AllowedSpeed) * bwlimit.MiB
	//fileLimited := bwlimit.NewWriter(f, writeLimit)
	dialer := bwlimit.NewDialer(&net.Dialer{
		Timeout:   30 * time.Minute,
		KeepAlive: 30 * time.Minute,
	}, writeLimit, readLimit)

	http.DefaultTransport.(*http.Transport).DialContext = dialer.DialContext

	resp, err := http.DefaultClient.Get(url)
	if err != nil {
		println(err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		println("file not found or file is inaccessible to this server")
		return "", 0
	}

	// to handle progress info

	//globalProgresscounter = controller.NewProgress(filename, GlobalCurrentUser.UserId, uint64(resp.ContentLength))
	//dataprogress := *controller.GetDataProgress()
	//globalProgresscounter = controller.GetProgressById(CurrenUser.UserId, progressId)
	//globalProgresscounter.Filename = filename
	//globalProgresscounter.Total = uint64(resp.ContentLength)
	tempProgress := controller.GetProgressById(CurrenUser.UserId, progressId)
	tempProgress.Filename = filename
	tempProgress.Total = uint64(resp.ContentLength)
	counter := &WriteCounter{0, tempProgress}
	_, err = io.Copy(f, io.TeeReader(resp.Body, counter))
	println("IsOn is :", tempProgress.IsOn, "\nProgressId is", tempProgress.ProcessId)
	//println("progress id  is ", progressId)

	return filename, 1

}
