package views
import (
"net/http"
"fmt"
"io/ioutil"
"database/sql"
 _ "github.com/go-sql-driver/mysql"
 "strconv"
"github.com/offical/system"
"github.com/offical/conf"
proto "github.com/golang/protobuf/proto"
)

func Home(w http.ResponseWriter,r *http.Request) {
    fmt.Fprintln(w,"")
}

func SignUp(w http.ResponseWriter,r *http.Request) {

}

func SignIn(w http.ResponseWriter,r *http.Request) {
     res,_:= ioutil.ReadAll(r.Body)
     user := &User{}
     err := proto.Unmarshal(res,user)

     if  err != nil {
     }else {
        config ,err :=conf.GetConfig()
        if (err !=nil) {
            panic("读取配置错误！")
        }
        sqlconf :=config.DataBase.User+":"+config.DataBase.Password+"@tcp("+config.DataBase.Host+":"+strconv.Itoa(config.DataBase.Port)+")/"+config.DataBase.Name+"?parseTime=true"
        db,err := sql.Open("mysql",sqlconf)
        if err!=nil {
                   system.MyLogger.LogToError()(err)
        }
        defer db.Close()
        sql := "select account,password from user where `account` ='"+ user.GetAccount()+ "'"
        rows,err :=db.Query(sql)
        if err !=nil {
             system.MyLogger.LogToError()(err)
        }
        var account string
        var password string
        for rows.Next() {
            err=rows.Scan(&account,&password)
            if (err!=nil) {
                 system.MyLogger.LogToError()(err)
            }
        }
        if password!=user.GetPassword()  {
             fmt.Fprintln(w,"密码错误")
        }else {
             sessionKey := r.Header.Get("session")
             session:=system.GlobalManage.SessionStore.SessionMap[sessionKey]
             session.SetSessionData("is_login",true)
             fmt.Fprintln(w,"登录成功")
        }
     }

}

func SignOut (w http.ResponseWriter,r *http.Request) {
             sessionKey := r.Header.Get("session")
             session:=system.GlobalManage.SessionStore.SessionMap[sessionKey]
             session.EmptySessionData()
             fmt.Fprintln(w,"退出成功")
}