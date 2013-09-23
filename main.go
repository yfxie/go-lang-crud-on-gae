package golang_gae_crud

import (
	"appengine"
	"appengine/datastore"
	"appengine/user"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"
)

// 存入 Datastore 資料形態
type Greeting struct {
	Id      int64
	Author  string
	Content string
	Date    time.Time
}

// 用於 template 的 data
type Context struct {
	Greetings   []Greeting
	CurrentUser string
}

// template 無內建判斷等值的func, 須要手動實作
var templateFuncMap = template.FuncMap{"eq": func(a, b string) bool { return a == b }}

func render(w http.ResponseWriter, filename string, context interface{}) error {
	var temp = template.New(filename)
	return template.Must(temp.Funcs(templateFuncMap).ParseFiles("templates/"+filename)).Execute(w, context)
}

func guestbookKey(c appengine.Context) *datastore.Key {
	return datastore.NewKey(c, "Guestbook", "default_guestbook", 0, nil)
}

type HandleFuncType func(http.ResponseWriter, *http.Request)
type HandleFuncTemplateType func(http.ResponseWriter, *http.Request, appengine.Context, *user.User)

func handleFuncUserCheck(handlefunc HandleFuncTemplateType) HandleFuncType {
	outfunc := func(w http.ResponseWriter, r *http.Request) {
		c := appengine.NewContext(r)
		u := user.Current(c)
		if u == nil {
			fmt.Fprint(w, "permission denied")
			return
		}
		handlefunc(w, r, c, u)
	}

	return outfunc
}

func init() {
	router()
}

func router() {
	// hadndleFuncUserCheck(func) 會先檢查使用者才做func的事情
	http.HandleFunc("/", index)
	http.HandleFunc("/create", handleFuncUserCheck(create))
	http.HandleFunc("/destroy", handleFuncUserCheck(destroy))
	http.HandleFunc("/edit", handleFuncUserCheck(edit))
	http.HandleFunc("/update", handleFuncUserCheck(update))
}

func index(w http.ResponseWriter, r *http.Request) {
	/*
	   c有Debugf, Infof, Warningf, Errorf, and Criticalf..等method可用來產生log message
	   Example:
	       c.Debugf("%v",r)
	*/
	c := appengine.NewContext(r)

	// fore to login
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

	if err := render(w, "index.html", context); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func create(w http.ResponseWriter, r *http.Request, c appengine.Context, current_user *user.User) {
	g := Greeting{
		Author:  current_user.String(),
		Content: r.FormValue("content"),
		Date:    time.Now(),
	}
	key := datastore.NewIncompleteKey(c, "Greeting", guestbookKey(c))

	if _, err := datastore.Put(c, key, &g); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

func edit(w http.ResponseWriter, r *http.Request, c appengine.Context, current_user *user.User) {
	int_id, _ := strconv.Atoi(r.FormValue("id"))
	key := datastore.NewKey(c, "Greeting", "", int64(int_id), guestbookKey(c))
	greeting := new(Greeting)
	err := datastore.Get(c, key, greeting)
	greeting.Id = int64(int_id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else if current_user.String() != greeting.Author {
		fmt.Fprint(w, "permission denied")
		return
	}

	if err := render(w, "edit.html", greeting); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func update(w http.ResponseWriter, r *http.Request, c appengine.Context, current_user *user.User) {
	int_id, _ := strconv.Atoi(r.FormValue("id"))
	key := datastore.NewKey(c, "Greeting", "", int64(int_id), guestbookKey(c))
	greeting := new(Greeting)
	err := datastore.Get(c, key, greeting)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else if current_user.String() != greeting.Author {
		fmt.Fprint(w, "permission denied")
		return
	}

	greeting.Content = r.FormValue("content")

	if _, err := datastore.Put(c, key, greeting); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Location", "/")
	w.WriteHeader(http.StatusFound)
}

func destroy(w http.ResponseWriter, r *http.Request, c appengine.Context, current_user *user.User) {
	int_id, _ := strconv.Atoi(r.FormValue("id"))
	key := datastore.NewKey(c, "Greeting", "", int64(int_id), guestbookKey(c))
	greeting := new(Greeting)
	err := datastore.Get(c, key, greeting)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else if current_user.String() != greeting.Author {
		fmt.Fprint(w, "permission denied")
		return
	}

	if err := datastore.Delete(c, key); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Location", "/")
	w.WriteHeader(http.StatusFound)
}
