package router
import (
"net/http"
"github.com/offical/views"
"github.com/offical/system"
)
func init() {
    http.Handle("/",system.Chain(http.HandlerFunc(views.Home),system.LogMiddleware(),system.SessionMiddleware()))

}