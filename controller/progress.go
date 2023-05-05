package controller

var DataProgresses = make(map[string][]Progress)

type Progress struct {
	Filename    string `json:"filename"`
	UserId      string `json:"user_id"`
	ProcessId   int    `json:"process_id"`
	Total       uint64 `json:"filesize"` //file size in bytes
	Transferred uint64 `json:"done"`
	IsOn        bool   `json:"state"`
}

func NewProgress(filename, userid string, filesize uint64) int {

	listProcess := DataProgresses[userid]
	next := len(listProcess)

	tempProgress := Progress{filename,
		userid,
		next,
		filesize,
		0,
		true,
	}
	if next == 0 {
		//tempFirstProcess := [&tempProcess]
		DataProgresses[userid] = []Progress{tempProgress}
		println("0 is the progress id")
		return 0
	} else {
		DataProgresses[userid] = append(DataProgresses[userid], tempProgress)
		println(next, " is the progress id")
		return next
	}

}

func GetDataProgress() *map[string][]Progress {
	return &DataProgresses
}
