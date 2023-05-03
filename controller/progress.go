package controller

var DataProgresses = make(map[string]*Progress)

type Progress struct {
	Filename string
	UserId   string
	Total    uint64 //file size in bytes
	Done     uint64
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
