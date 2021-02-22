package system
import (
"sync"
"crypto/md5"
"encoding/hex"
"crypto/rand"
"database/sql"
 _ "github.com/go-sql-driver/mysql"
   "github.com/json-iterator/go"
 "encoding/base64"
 "time"
 "strconv"
"github.com/offical/conf"
"io"
"strings"
)
type Session struct {
  sessionKey string
  sessionData map[interface{}]interface{}
  expireDate time.Time
  modify   bool
  isCach   bool
}
type SessionStore struct {
        lock sync.Mutex
        SessionMap map[string]*Session
}
type Manager struct {
    storeType string
    cookieName string
    lock sync.Mutex
    SessionStore *SessionStore
    maxlifetime int64

}
var GlobalManage *Manager

func init() {
    sessionMap := make(map[string]*Session)
    var GlobalSessionStore *SessionStore
    GlobalSessionStore =&SessionStore{SessionMap:sessionMap}
    GlobalManage,_=NewManager("database","sessionId",60*60*24,GlobalSessionStore)
}

func (this *SessionStore) Push(sessionKey string ,session *Session) {

        this.lock.Lock()

        defer this.lock.Unlock()

        this.SessionMap[sessionKey] = session
}
func (this *SessionStore) Delete(sessionKey string ) {

        this.lock.Lock()
        defer this.lock.Unlock()
        delete(this.SessionMap,sessionKey)
}

func NewManager(sessionStoreName,cookieName string,maxlifetime int64,sessionStore *SessionStore) (*Manager,error) {


    return &Manager{storeType:sessionStoreName,SessionStore:sessionStore,cookieName:cookieName,maxlifetime:maxlifetime},nil
}
func NewSession (sid string,isCach bool ) (*Session) {
    sessionData :=make(map[interface{}]interface{})
    strInt64 := strconv.FormatInt(GlobalManage.maxlifetime, 10)
    id16 ,_ := strconv.Atoi(strInt64)
    expireDate :=time.Now().Add(time.Second*time.Duration(id16))
    session := &Session{sessionKey:sid,sessionData:sessionData,expireDate:expireDate,modify:false,isCach:isCach}
    return session

}
func (session *Session) isEmpty ()  bool {
   return (session.sessionKey=="")||(session.sessionData==nil)
}
func (session *Session) EmptySessionData() {
   session.sessionData=nil
}
func (session *Session) SetSessionData(key interface{},value interface{}) {
         session.modify = true
         session.sessionData[key] =value

}
func (session *Session) Save() error {
        if (session.sessionKey ==""||session.isCach) {
                return session.Create()
        }
        sessionDataS:=session.Encode(session.GetSessionData())


        sqlconf :=conf.MyConfig.DataBase.User+":"+conf.MyConfig.DataBase.Password+"@tcp("+conf.MyConfig.DataBase.Host+":"+strconv.Itoa(conf.MyConfig.DataBase.Port)+")/"+conf.MyConfig.DataBase.Name+"?parseTime=true"
        db,err := sql.Open("mysql",sqlconf)
        if err!=nil {

               MyLogger.LogToError()(err)
            return err
        }
        defer db.Close()
        sql:=`replace INTO session (session_key, session_data, expire_date) VALUES ("` + session.sessionKey+`","`+sessionDataS+`","` +session.expireDate.Format("2006-01-02 15:04:05")+`")`
         rows, err := db.Query(sql)
        defer rows.Close()
        if err !=nil {
                      MyLogger.LogToError()(err)
                    return err
        }
        //使用数据库存储

        return nil
}

func (session *Session) GetSessionData()  map[interface{}]interface{} {

    if (session.sessionData==nil) {
        if (session.sessionKey=="") {
            sessionData := make(map[interface{}]interface{})
            session.sessionData = sessionData

        }else {

            //从数据库获取


            sqlconf :=conf.MyConfig.DataBase.User+":"+conf.MyConfig.DataBase.Password+"@tcp("+conf.MyConfig.DataBase.Host+":"+strconv.Itoa(conf.MyConfig.DataBase.Port)+")/"+conf.MyConfig.DataBase.Name+"?parseTime=true"
            db,err := sql.Open("mysql",sqlconf)
            if err!=nil {
                      MyLogger.LogToError()(err)
                    sessionData := make(map[interface{}]interface{})
                    session.sessionData = sessionData
                    return session.sessionData
             }
             defer db.Close()
                sql:= "select * from session where  `expire_date`>CURRENT_TIME()  AND `session_key`=" +session.sessionKey
                rows,err :=db.Query(sql)
                if err !=nil {
                      MyLogger.LogToError()(err)
                }
                var currSessionData map[interface{}]interface{}
                for rows.Next() {
                    var sessionKey string
                    var sessionDataS string
                    var expireDate string
                    err=rows.Scan(&sessionKey,&sessionDataS,&expireDate)
                    if (err!=nil) {
                            MyLogger.LogToError()(err)
                    }
                    currSessionData=session.Decode(sessionDataS)
                }
                session.sessionData =currSessionData
        }
    }
    return session.sessionData



}
func (session *Session) Encode(sessionData map[interface{}]interface{}) string {
    data,err :=  jsoniter.Marshal(sessionData)
    h :=md5.New()
    h.Write([]byte(data))
    cipherStr:=h.Sum(nil)
    s:=hex.EncodeToString(cipherStr)+":" + string(data)
    if (err!=nil ) {
                       MyLogger.LogToError()(err)
                    return ""
    } else {
                 sessionDataS:= base64.StdEncoding.EncodeToString([]byte(s))
                 return sessionDataS
    }
}

func (session *Session) Decode(sessionDataS string) map[interface{}]interface{} {
    data,_:=base64.StdEncoding.DecodeString(sessionDataS)
    arr:=strings.Split(string(data),":")
    var sessionData map[interface{}]interface{}
    err :=  jsoniter.Unmarshal([]byte(arr[1]),&sessionData)

    if err!=nil {
             MyLogger.LogToError()(err)
          sessionData = make(map[interface{}]interface{})

    }
    return sessionData


}

func (session *Session) Create() error {
    for {
        session.sessionKey= GlobalManage.createCach()
        err:=session.Save()
        if (err!=nil ){
           continue
        }else {
            break
        }
    }
    return nil

}
func(manager *Manager) createCach() string {
     //该函数仅在cach中检查是否重复,即使不在数据库中查重也没关系,因为存入数据库后会创建一个真正的session
     var sessionKey string
     for {
         sessionKey= GlobalManage.sessionId()
         _,ok :=GlobalManage.SessionStore.SessionMap[sessionKey]
         if ok {
             continue
         }else {
            break
         }
     }
     return sessionKey
}
func (manager *Manager) sessionId() string {
    b:=make([]byte,32)
    if _,err := io.ReadFull(rand.Reader,b);err !=nil {
           MyLogger.LogToError()(err)
        return ""
    }else {
        h :=md5.New()
        h.Write(b)
        cipherStr:=h.Sum(nil)
        return hex.EncodeToString(cipherStr)
    }
}
