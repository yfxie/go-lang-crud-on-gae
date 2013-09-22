package golang_gae_crud

import (
    "appengine"
    "appengine/datastore"
    "appengine/user"
    "fmt"
    "net/http"
    "strconv"
    "time"
)

type Greeting struct {
    Id      int64
    Author  string
    Content string
    Date    time.Time
}

type Context struct {
    Greetings   []Greeting
    CurrentUser string
}

func NewGreeting (u *user.User, content string, c appengine.Context) (*Greeting, *datastore.Key) {
    g := &Greeting{
        Author:  u.String(),
        Content: content,
        Date:    time.Now(),
    }

    key := datastore.NewIncompleteKey(c, "Greeting", guestbookKey(c))
    return g,key
}

func GetGreeting (id_str string, c appengine.Context) (*Greeting, *datastore.Key, error) {
    id, err := strconv.Atoi(id_str)
    key := datastore.NewKey(c, "Greeting", "", int64(id), guestbookKey(c))
    
    greeting := new(Greeting)
    err = datastore.Get(c, key, greeting)
    greeting.Id = key.IntID()

    return greeting, key, err
}

var guestbookTemplate = newTemplateFromFile("index.html")
var editTemplate = newTemplateFromFile("edit.html")

func init() {
    router()
}

func router() {
    http.HandleFunc("/", index)
    http.HandleFunc("/create", handleFuncUserCheckTemplate(create))
    http.HandleFunc("/destroy", handleFuncUserCheckTemplate(destroy))
    http.HandleFunc("/edit", handleFuncUserCheckTemplate(edit))
    http.HandleFunc("/update", handleFuncUserCheckTemplate(update))
}

func guestbookKey(c appengine.Context) *datastore.Key {
    return datastore.NewKey(c, "Guestbook", "default_guestbook", 0, nil)
}

func index(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
    if u := user.Current(c); u == nil {
        url, err := user.LoginURL(c, r.URL.String())
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        w.Header().Set("Location", url)
        w.WriteHeader(http.StatusFound)
        return
    }
    q := datastore.NewQuery("Greeting").Ancestor(guestbookKey(c)).Order("-Date").Limit(10)
    greetings := make([]Greeting, 0, 10)
    keys, err := q.GetAll(c, &greetings)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    for i := 0; i < len(greetings); i++ {
        greetings[i].Id = keys[i].IntID()
    }
    context := Context{Greetings: greetings, CurrentUser: user.Current(c).String()}
    if err := guestbookTemplate.Execute(w, context); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}


func create(w http.ResponseWriter, r *http.Request, c appengine.Context, u *user.User) {
    g,key := NewGreeting(u,r.FormValue("content"),c)
    _, err := datastore.Put(c, key, g)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    http.Redirect(w, r, "/", http.StatusFound)
}

func edit(w http.ResponseWriter, r *http.Request, c appengine.Context, u *user.User) {
    greeting, _, err := GetGreeting(r.FormValue("id"),c)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }else if u.String() != greeting.Author {
        fmt.Fprint(w, "permission denied")
        return 
    }

    if err := editTemplate.Execute(w, greeting); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func update(w http.ResponseWriter, r *http.Request, c appengine.Context, u *user.User) {
    greeting, key, err := GetGreeting(r.FormValue("id"),c)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }else if u.String() != greeting.Author {
        fmt.Fprint(w, "permission denied")
        return
    }

    greeting.Content = r.FormValue("content")
    key, err = datastore.Put(c, key, greeting)

    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Location", "/")
    w.WriteHeader(http.StatusFound)
}

func destroy(w http.ResponseWriter, r *http.Request, c appengine.Context, u *user.User) {
    greeting, key, err := GetGreeting(r.FormValue("id"),c)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return 
    }else if u.String() != greeting.Author {
        fmt.Fprint(w, "permission denied")
        return 
    }

    err = datastore.Delete(c, key)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Location", "/")
    w.WriteHeader(http.StatusFound)
}
