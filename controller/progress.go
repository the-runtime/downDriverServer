package controller

var DataProgresses = make(map[string]*Progress)

type Progress struct {
	Filename    string `json:"filename"`
	UserId      string `json: "user_id"`
	Total       uint64 `json:"filesize"` //file size in bytes
	Transferred uint64 `json:"done"`
}

func NewProgress(filename, userid string, filesize uint64) *Progress {
	tempProgress := Progress{filename,
		userid,
		filesize,
		0,
	}

	DataProgresses[userid] = &tempProgress
	return &tempProgress
}

func GetDataProgress() *map[string]*Progress {
	return &DataProgresses
}
