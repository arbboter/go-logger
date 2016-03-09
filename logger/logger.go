package logger

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path"
	"runtime"
	"strconv"
	"sync"
	"time"
)

const (
	_VER string = "1.0.3"
)

type LEVEL int32

var logLevel LEVEL = 1
var maxFileSize int64
var maxFileCount int32
var dailyRolling bool = true
var consoleAppender bool = true
var RollingFile bool = false
var logSet [OFF]*_FILE
var logFlag = 0

const DATEFORMAT = "2006-01-02"

type UNIT int64

const (
	_       = iota
	KB UNIT = 1 << (iota * 10)
	MB
	GB
	TB
)

const (
	ALL LEVEL = iota
	DEBUG
	ERROR
	KEY
	OFF
)

func (lv LEVEL) Tag() string {
	switch lv {
	case DEBUG:
		return "DEG"
	case ERROR:
		return "ERR"
	case KEY:
		return "KEY"
	default:
		return "LOG"
	}
}

func (lv LEVEL) String() string {
	switch lv {
	case DEBUG:
		return "debug"

	case ERROR:
		return "error"
	case KEY:
		return "key"
	default:
		return "log"
	}
}

func Init(appName string, logLevel LEVEL) {
	runtime.GOMAXPROCS(runtime.NumCPU())

	//指定是否控制台打印，默认为true
	SetConsole(false)

	//指定日志文件备份方式为日期的方式
	//第一个参数为日志文件存放目录
	//第二个参数为日志文件命名
	//logger.SetRollingDaily("F:/data/applog/main_test", "main_test.log")

	//指定日志级别  ALL，DEBUG，INFO，WARN，ERROR，FATAL，OFF 级别由低到高
	//一般习惯是测试阶段为debug，生成环境为info以上
	SetLevel(logLevel)

	//指定日志文件备份方式为文件大小的方式
	//第一个参数为日志文件存放目录
	//第二个参数为日志文件命名
	//第三个参数为备份文件最大数量
	//第四个参数为备份文件大小
	//第五个参数为文件大小的单位
	DefaultDir := "/data/applog/"
	if runtime.GOOS == "windows" {
		DefaultDir = "." + DefaultDir
	}
	SetRollingFile(DefaultDir+appName, appName, 50, 20, MB)

}

type _FILE struct {
	dir             string
	filename        string
	_suffix         int
	isCover         bool
	_date           *time.Time
	mu              *sync.RWMutex
	logfile         *os.File
	lg              *log.Logger
	logLevel        LEVEL
	maxFileSize     int64
	maxFileCount    int32
	dailyRolling    bool
	consoleAppender bool
	RollingFile     bool
}

func (f *_FILE) Print(lv LEVEL, v ...interface{}) {
	if lv < f.logLevel {
		return
	}

	if f.lg == nil {
		f.openWriteLogFile()
	}
	//fmt.Println(f.logLevel.Tag(), v)
	f.lg.Print(v...)

	lowerLevel := f.logLevel - 1
	if lowerLevel >= DEBUG {
		logObj := logSet[lowerLevel]
		logObj.Print(lv, v...)
	}
}

func (f *_FILE) isMustRename() bool {
	if f.dailyRolling {
		t, _ := time.Parse(DATEFORMAT, time.Now().Format(DATEFORMAT))
		if t.After(*f._date) {
			return true
		}
	} else {
		if f.maxFileCount > 1 {
			if f.lg != nil {
				//fmt.Println("isMustRename :", f.logFileName())
				if fileSize(f.logFileName()) >= f.maxFileSize {
					return true
				}
			}
		}
	}
	return false
}

func (f *_FILE) logFileName() string {
	return f.dir + "/" + f.filename + "_" + f.logLevel.String() + ".log"
}

func (f *_FILE) openWriteLogFile() {
	if f.lg != nil {
		return
	}

	logName := f.logFileName()
	lf, err := os.OpenFile(logName, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	fmt.Println("open file:", logName)

	if err != nil {
		fmt.Println("open file failed:", logName)
	}
	f.logfile = lf
	f.lg = log.New(f.logfile, "", logFlag)
}

func (f *_FILE) rename() {
	if f.dailyRolling {
		fn := f.dir + "/" + f.filename + "." + f._date.Format(DATEFORMAT)
		if !isExist(fn) && f.isMustRename() {
			if f.logfile != nil {
				f.logfile.Close()
			}
			err := os.Rename(f.dir+"/"+f.filename, fn)
			if err != nil {
				f.lg.Println("rename err", err.Error())
			}
			t, _ := time.Parse(DATEFORMAT, time.Now().Format(DATEFORMAT))
			f._date = &t
			f.logfile, _ = os.Create(f.dir + "/" + f.filename)
			f.lg = log.New(f.logfile, "", logFlag)
		}
	} else {
		f.coverNextOne()
	}
}

func (f *_FILE) nextSuffix() int {
	return int(f._suffix%int(maxFileCount) + 1)
}

// 滚动日志文件，返回受影响的文件数目
func (f *_FILE) rollLogFile() int {
	var moveNum int = 0

	// Del End File
	endFile := f.indexFile(int(f.maxFileCount))
	if isExist(endFile) {
		os.Remove(endFile)
		moveNum += 1
		fmt.Println("Del File:", endFile)
	}
	for i := int(maxFileCount) - 1; i > 0; i-- {
		// Move File i->i+1
		srcName := f.indexFile(i)
		if isExist(srcName) {
			dstName := f.indexFile(i + 1)
			os.Rename(srcName, dstName)
			moveNum += 1
			fmt.Println("Move File:", srcName, dstName)
		}
	}
	return moveNum
}

func (f *_FILE) indexFile(idx int) string {
	return f.logFileName() + "." + strconv.Itoa(int(idx))
}

func (f *_FILE) coverNextOne() {
	f._suffix = 1
	if f.logfile != nil {
		f.logfile.Close()
		f.logfile = nil
		f.lg = nil
	}

	f.rollLogFile()
	os.Rename(f.logFileName(), f.indexFile(f._suffix))

	f.openWriteLogFile()
}

func fileCheck(lv LEVEL) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()
	logObj := logSet[lv]
	//fmt.Println("fileCheck:", logObj.logFileName())
	if logObj != nil && logObj.isMustRename() {
		logObj.mu.Lock()
		defer logObj.mu.Unlock()
		logObj.rename()
	}
}

func SetConsole(isConsole bool) {
	consoleAppender = isConsole
}

func SetLevel(_level LEVEL) {
	logLevel = _level
}

func SetRollingFile(fileDir, fileName string, maxNumber int32, maxSize int64, _unit UNIT) {
	// Make Sure Dir Exist
	err := os.MkdirAll(fileDir, 0777)
	if err != nil {
		log.Printf("Create Dir Failed. Dir:%v\n", fileDir)
	}

	maxFileCount = maxNumber
	maxFileSize = maxSize * int64(_unit)
	RollingFile = true
	dailyRolling = false

	for lv := logLevel; lv < OFF; lv++ {
		logSet[lv] = &_FILE{dir: fileDir, filename: fileName, isCover: false, mu: new(sync.RWMutex), lg: nil}
		logSet[lv].maxFileCount = maxFileCount
		logSet[lv].maxFileSize = maxFileSize
		logSet[lv].RollingFile = RollingFile
		logSet[lv].dailyRolling = dailyRolling
		logSet[lv].logLevel = lv
		logSet[lv].consoleAppender = false
	}
	go fileMonitor()
}

func SetRollingDaily(fileDir, fileName string) {
	//	RollingFile = false
	//	dailyRolling = true
	//	t, _ := time.Parse(DATEFORMAT, time.Now().Format(DATEFORMAT))
	//	logObj = &_FILE{dir: fileDir, filename: fileName, _date: &t, isCover: false, mu: new(sync.RWMutex)}
	//	logObj.mu.Lock()
	//	defer logObj.mu.Unlock()

	//	if !logObj.isMustRename() {
	//		logObj.logfile, _ = os.OpenFile(fileDir+"/"+fileName, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	//		logObj.newLogFile()
	//	} else {
	//		logObj.rename()
	//	}
}

func console(lv LEVEL, s ...interface{}) {
	needWriter := logSet[lv].consoleAppender
	if needWriter {
		_, file, line, _ := runtime.Caller(2)
		short := file
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				short = file[i+1:]
				break
			}
		}
		file = short
		log.Println(file, strconv.Itoa(line), s)
	}
}

func catchError() {
	if err := recover(); err != nil {
		log.Println("err", err)
	}
}

// 日志信息
func getLineInfo(levelName string, calldepth int) string {
	var buffer bytes.Buffer

	tNow := time.Now()
	buffer.WriteString(tNow.Format("[2006-01-02 15:04:05.999]"))
	buffer.WriteString(fmt.Sprintf(" [PID:%v] #%v# ", os.Getpid(), levelName))
	funcName, file, line, ok := runtime.Caller(calldepth)
	if ok {
		buffer.WriteString(fmt.Sprintf("FILE:%v LN:%v FUNC:%v", path.Base(file), line, runtime.FuncForPC(funcName).Name()))
	} else {
		buffer.WriteString("FILE:nil LN:nil FUNC:nil")
	}
	return buffer.String()
}

func logByLevelln(lv LEVEL, v ...interface{}) {
	defer catchError()
	logObj := logSet[lv]
	if logObj != nil {
		logObj.mu.RLock()
		defer logObj.mu.RUnlock()
	}

	if logLevel <= lv {
		var buffer bytes.Buffer
		buffer.WriteString(fmt.Sprintf("%v EM:", getLineInfo(lv.Tag(), 3)))
		buffer.WriteString(fmt.Sprint(v...))
		if logObj != nil {
			logObj.Print(lv, buffer.String())
		}
		console(lv, buffer.String())
	}
}

func logByLevelf(lv LEVEL, arg string, v ...interface{}) {
	defer catchError()
	logObj := logSet[lv]
	if logObj != nil {
		logObj.mu.RLock()
		defer logObj.mu.RUnlock()
	}

	if logLevel <= lv {
		var buffer bytes.Buffer
		buffer.WriteString(fmt.Sprintf("%v EM:", getLineInfo(lv.Tag(), 3)))
		buffer.WriteString(fmt.Sprintf(arg, v...))
		if logObj != nil {
			logObj.Print(lv, buffer.String())
		}
		console(lv, buffer.String())
	}
}

// 行输出日志
func Debug(v ...interface{}) {
	logByLevelln(DEBUG, v...)
}

func Error(v ...interface{}) {
	logByLevelln(ERROR, v...)
}
func Key(v ...interface{}) {
	logByLevelln(KEY, v...)
}
func Header() {
	var buffer bytes.Buffer

	funcName, _, _, ok := runtime.Caller(1)
	buffer.WriteString("\n======== ")
	if ok {
		buffer.WriteString(fmt.Sprintf("%v ", runtime.FuncForPC(funcName).Name()))
	} else {
		buffer.WriteString("nil ")
	}

	tNow := time.Now()
	buffer.WriteString(tNow.Format("REQUEST START cptime:[2006-01-02 15:04:05.999] ========"))

	logObj := logSet[DEBUG]
	if logObj.lg == nil {
		logObj.openWriteLogFile()
	}
	logObj.lg.Print(buffer.String())

}

// 格式化日志
func Debugf(arg string, v ...interface{}) {
	logByLevelf(DEBUG, arg, v...)
}
func Errorf(arg string, v ...interface{}) {
	logByLevelf(ERROR, arg, v...)
}
func Keyf(arg string, v ...interface{}) {
	logByLevelf(KEY, arg, v...)
}

func fileSize(file string) int64 {
	f, e := os.Stat(file)
	if e != nil {
		fmt.Println(e.Error())
		return 0
	}
	return f.Size()
}

func isExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

func fileMonitor() {
	timer := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-timer.C:
			for lv := logLevel; lv < OFF; lv++ {
				fileCheck(lv)
			}
		}
	}
}
