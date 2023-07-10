package controller

var dataProgresses = make(map[string][]*Progress)

//var dataProgresses = make(map[string]map[int]Progress)

type Progress struct {
	Filename    string `json:"filename"`
	UserId      string `json:"user_id"`
	ProcessId   int    `json:"process_id"`
	Total       uint64 `json:"filesize"` //file size in bytes
	Transferred uint64 `json:"done"`
	IsOn        bool   `json:"state"`
}

func NewProgress(filename, userid string, filesize uint64) int {

	userProcesses := dataProgresses[userid]
	next := len(userProcesses)

	tempProgress := Progress{filename,
		userid,
		next,
		filesize,
		0,
		true,
	}

	if next == 0 {
		dataProgresses[userid] = []*Progress{&tempProgress}
		println("0 is the progress id")
		return 0

	} else {
		dataProgresses[userid] = append(dataProgresses[userid], &tempProgress)
		println(next, " is the progress id")
		return next
	}

}

//func GetDataProgress() *map[string][]Progress {
//	return &dataProgresses
//}

func GetProgressById(userId string, progressId int) *Progress {
	return dataProgresses[userId][progressId]
}

func GetProgressList(userId string) []Progress {
	var listProgress []Progress
	for _, r := range dataProgresses[userId] {
		listProgress = append(listProgress, *r)
	}
	return listProgress
}

//need some fundamental changes make it work
//func EndProgress(userId string, progressId int) {
//	place := progressId
//	dataProgresses[userId] = append(dataProgresses[userId][:place-1], dataProgresses[userId][place:]...)
//
//}
