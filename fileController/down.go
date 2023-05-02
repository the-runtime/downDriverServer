package fileController

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func StartDown(url string) (string, int) {
	client := http.Client{}

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

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		fmt.Println(err.Error())
		return "", 0
	}

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		return "", 0
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}(res.Body)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err.Error())
		return "", 0
	}

	_, err = f.Write(body)
	if err != nil {
		fmt.Println(err.Error())
		return "", 0
	}

	return filename, 1

}
