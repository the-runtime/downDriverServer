package fileController

import (
	"fmt"
	"github.com/dustin/go-humanize"
	"io"
	"net/http"
	"os"
	"serverFordownDrive/controller"
	"serverFordownDrive/model"
	"strings"
)

var GlobalCurrentUser *model.User
var globalProgresscounter *controller.Progress

//for implementing transfer limit on user

type WriteCounter struct {
	Total uint64
}

func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Total += uint64(n)
	wc.PrintProgress()
	//wc.UpdateProgress()
	return n, nil
}
func (wc *WriteCounter) PrintProgress() {
	fmt.Printf("\r%s", strings.Repeat(" ", 35))
	fmt.Printf("\rDownloading... %s complete", humanize.Bytes(wc.Total))
	globalProgresscounter.Done = wc.Total / uint64(2)
	//println("globalProgresscounter.Done   %d", globalProgresscounter.Done)
}

//func (wc *WriteCounter) UpdateProgress() {
//	GlobalCurrentUser.ConsumedDataTransfer += wc.Total / uint64(2) // for not counting upload and download separately
//	globalProgresscounter.Done += wc.Total / uint64(2)
//}

func StartDown(url string, CurrenUser *model.User) (string, int) {
	//client := http.Client{}

	GlobalCurrentUser = CurrenUser
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

	resp, err := http.Get(url)
	if err != nil {
		println(err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		println("file not found or file is inaccessible to this server")
		return "", 0
	}

	// to handle progress info

	globalProgresscounter = controller.NewProgress(filename, GlobalCurrentUser.UserId, uint64(resp.ContentLength))

	counter := &WriteCounter{}
	_, err = io.Copy(f, io.TeeReader(resp.Body, counter))

	//req, err := http.NewRequest(http.MethodGet, url, nil)
	//if err != nil {
	//	fmt.Println(err.Error())
	//	return "", 0
	//}
	//
	//res, err := client.Do(req)
	//if err != nil {
	//	fmt.Println(err.Error())
	//	return "", 0
	//}
	//
	//defer func(Body io.ReadCloser) {
	//	err := Body.Close()
	//	if err != nil {
	//		fmt.Println(err.Error())
	//		return
	//	}
	//}(res.Body)
	//
	//io.Copy(f, res.Body)
	//
	//body, err := io.ReadAll(res.Body)
	//if err != nil {
	//	fmt.Println(err.Error())
	//	return "", 0
	//}
	//
	//_, err = f.Write(body)
	//if err != nil {
	//	fmt.Println(err.Error())
	//	return "", 0
	//}

	return filename, 1

}
