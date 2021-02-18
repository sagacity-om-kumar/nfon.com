package helper

import (
	"path"

	ghelper "nfon.com/helper"
)

const queryFileName string = "dbTemplateQueries.json"

func ReadDBQueryFile() (bool, []byte) {
	newPath := path.Join("queries", queryFileName)
	return ghelper.ReadFileContent(newPath)
}

func IsQueryFileMatched(fileName string) bool {

	if queryFileName == fileName {
		return true
	} else {
		return false

	}
}
