package log

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"sync"
	"time"
)

const (
	//PATH  LOGPATH/time.Now().Format(FORMAT)/*.log
	PATH = "log/"
	//FORMAT .
	FORMAT = "20060102"
	//LineFeed 换行
	LineFeed = "\r\n"
)

var (
	Trace   *log.Logger // 记录所有日志
	Info    *log.Logger // 重要的信息
	Warning *log.Logger // 需要注意的信息
	Error   *log.Logger // 非常严重的问题
	mutex   sync.Mutex  // 写入锁
)

func Init() {

	Trace = log.New(ioutil.Discard,
		"TRACE: ",
		log.Ldate|log.Ltime)

	Info = log.New(os.Stdout,
		"INFO: ",
		log.Ldate|log.Ltime)

	Warning = log.New(os.Stdout,
		"WARNING: ",
		log.Ldate|log.Ltime)

	Error = log.New(os.Stdout,
		"ERROR: ",
		log.Ldate|log.Ltime)
}

func Writer(l *log.Logger, mess interface{}) error {
	mutex.Lock()
	var path = PATH + time.Now().Format(FORMAT)
	if !isExist(PATH) {
		return createDir(PATH)
	}

	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println(err)
	}

	defer file.Close()

	_, fileName, line, _ := runtime.Caller(1)

	l.SetOutput(file)
	l.Println(fileName, line, mess)
	fmt.Println(fileName, line, mess)

	InsertLog(fileName, line, mess)

	mutex.Unlock()
	return nil
}

func logError(l *log.Logger, mess interface{}) error {
	mutex.Lock()
	var path = PATH + time.Now().Format(FORMAT)
	if !isExist(PATH) {
		return createDir(PATH)
	}

	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println(err)
	}

	defer file.Close()

	_, fileName, line, _ := runtime.Caller(1)

	l.SetOutput(file)
	l.Println(fileName, line, mess)
	fmt.Println(fileName, line, mess)

	mutex.Unlock()
	return nil
}

//createDir  文件夹创建
func createDir(path string) error {
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}
	os.Chmod(path, os.ModePerm)
	return nil
}

//isExist  判断文件夹/文件是否存在  存在返回 true
func isExist(f string) bool {
	_, err := os.Stat(f)
	return err == nil || os.IsExist(err)
}
