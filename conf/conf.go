package conf
import (
    "gopkg.in/gcfg.v1"
    "fmt"
)
type SYSTEM struct {
    Port int
}
type DATABASE struct {
    Name string
    User string
    Password string
    Host string
    Port int
}
type Config struct  {
    System SYSTEM
    DataBase DATABASE
}
var MyConfig Config

func init() {
    var err error
    MyConfig,err = GetConfig()
    if err!=nil {

        fmt.Println("读取配置文件错误")
    }


}
func GetConfig() (Config,error ) {
        config :=Config{}
        err := gcfg.ReadFileInto(&config,"conf/conf.ini")
        if err!=nil {
            return  config,err
        }
        return config,nil
}