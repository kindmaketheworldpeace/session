package system
import (
"log"
"fmt"
"os"
)
type MyLog struct {

    logger *log.Logger
    File   *os.File
}
var MyLogger *MyLog
func init() {
     logger := log.New(nil, "", log.Lshortfile|log.Ldate|log.Ltime)
     logFile, _ := os.OpenFile("logs/info.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
     MyLogger = &MyLog{logger:logger,File:logFile}
}
func (this *MyLog) LogToError(v ...interface{}) func( ...interface{}) {
    this.File.Close()
    logFile, err := os.OpenFile("logs/error.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
    if err !=nil {
        fmt.Println(err)
    }
    this.logger.SetOutput(logFile)
     this.File = logFile
     this.logger.SetPrefix("[ERROR]")

     return   this.logger.Println


}
func (this *MyLog) LogToInfo(v ...interface{}) func( ...interface{}) {
    this.File.Close()
    logFile, err := os.OpenFile("logs/info.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
    if err !=nil {
        fmt.Println(err)
    }
    this.logger.SetOutput(logFile)
    this.File = logFile
    this.logger.SetPrefix("[INFO]")
    return this.logger.Println

}
