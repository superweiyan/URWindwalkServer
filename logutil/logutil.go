package logutil

import (
	"log"
	"os"
)

var filePath string = "./log/test"
var file *os.File
var logger *log.Logger

func init() {

	var err error
	if isFileExist(filePath) == false {

		file, err = os.Create(filePath)
		if err != nil {
			log.Fatalln("fail to create test file!")
			panic("open log file fail")
		}
	}

	file, err = os.OpenFile("./log/test", os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		panic("open log file fail")
	}

	// log.SetOutput(file)
	// log.SetFlags(log.LstdFlags | log.Lshortfile)

	logger = log.New(file, "", log.LstdFlags|log.Lshortfile)
}

func isFileExist(filePath string) bool {
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false
	}
	return err != nil
}

func Writelog(info string) {
	if file == nil {
		return
	}

	logger.Printf(info)
}

func WriteLogInt(info int) {
	if file == nil {
		return
	}

	logger.Println(info)
}
