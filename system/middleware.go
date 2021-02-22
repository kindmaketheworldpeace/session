package system
import (
"net/http"
"net/url"
"io"
"bytes"
)
type Middleware func(http.Handler) http.Handler
func Chain(f http.Handler,mmap ...Middleware) http.Handler {
    for _,m := range mmap {
        f=m(f)
    }
    return f
}
type responsewriter struct {
    w http.ResponseWriter
    buf bytes.Buffer
    code int
}
func (rw *responsewriter) Header() http.Header {
    return rw.w.Header()
}
func (rw *responsewriter) WriteHeader (statusCode int) {
    rw.code =statusCode
}
func (rw *responsewriter) Write(data []byte ) (int,error) {

    return rw.buf.Write(data)
}
func (rw *responsewriter) Done() (int64,error) {
    if rw.code >0 {
        rw.w.WriteHeader(rw.code)
    }
    return io.Copy(rw.w,&rw.buf)

}


func SessionMiddleware() Middleware {
    return func(f http.Handler) http.Handler {
        return http.HandlerFunc(func (w http.ResponseWriter,r *http.Request){
                rw := &responsewriter{w:w}
                cookie,err:= r.Cookie(GlobalManage.cookieName)
                var session *Session
                if (err!=nil) {
                    MyLogger.LogToError()(err)
                    sessionKey:=GlobalManage.createCach()
                    session = NewSession(sessionKey,true)
                }else {
                    session = NewSession(cookie.Value,false)
                }
                GlobalManage.SessionStore.Push(session.sessionKey,session)
                r.Header.Set("session",session.sessionKey)
                f.ServeHTTP(rw,r)
                if isEmpty := session.isEmpty();isEmpty {
                           //若会话没有操作session,则清空cookies
                           cookieP:=http.Cookie{Name:GlobalManage.cookieName,Path:"/",HttpOnly:true, MaxAge: -1}
                           http.SetCookie(rw,&cookieP)
                }else {

                    if (session.modify) {
                           session.Save()
                           cookieP :=http.Cookie{Name:GlobalManage.cookieName,Value:url.QueryEscape(session.sessionKey),Path:"/",HttpOnly:true,MaxAge:int(GlobalManage.maxlifetime)}
                           http.SetCookie(rw,&cookieP)
                    }
                }
                if _, err := rw.Done(); err != nil {
                           GlobalManage.SessionStore.Delete(session.sessionKey)
                }
        })
    }
}


func LogMiddleware() Middleware {
        return func (f http.Handler) http.Handler {
            return  http.HandlerFunc(func (w http.ResponseWriter,r *http.Request) {

                  MyLogger.LogToInfo()(r.URL.Path)
            })
        }
}