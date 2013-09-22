package golang_gae_crud

import (
    "net/http"
    "appengine"
    "appengine/user"
    "fmt"
)

type HandleFuncType func (http.ResponseWriter, *http.Request)
type HandleFuncTemplateType func (http.ResponseWriter, *http.Request, appengine.Context, *user.User)

func handleFuncUserCheckTemplate(handlefunc HandleFuncTemplateType) HandleFuncType {
    outfunc := func (w http.ResponseWriter, r *http.Request) {
        c := appengine.NewContext(r)
        u := user.Current(c)
        if u == nil {
            fmt.Fprint(w, "permission denied")
            return
        }
        
        handlefunc(w,r,c,u)
    }

    return outfunc
}
