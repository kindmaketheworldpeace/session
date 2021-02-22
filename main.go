package main
import (
"net/http"
"strconv"
"github.com/offical/conf"
"github.com/offical/system"
_ "github.com/offical/router"
)

func main() {

    defer system.MyLogger.File.Close()


    http.ListenAndServe(":"+strconv.Itoa(conf.MyConfig.System.Port),nil)
}