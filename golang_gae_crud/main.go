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

func init() {
	router()
}

func router() {
	http.HandleFunc("/", index)
	http.HandleFunc("/create", create)
	http.HandleFunc("/destroy", destroy)
	http.HandleFunc("/edit", edit)
	http.HandleFunc("/update", update)
}

func guestbookKey(c appengine.Context) *datastore.Key {
	return datastore.NewKey(c, "Guestbook", "default_guestbook", 0, nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	/*
		c有Debugf, Infof, Warningf, Errorf, and Criticalf..等method可用來產生log message
		Example:
			c.Debugf("%v",r)
	*/
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

var guestbookTemplate = template.Must(template.New("index").Funcs(template.FuncMap{"eq": func(a, b string) bool { return a == b }}).Parse(guestbookTemplateHTML))

const guestbookTemplateHTML = `
<html>
  <body>
      <form action="/create" method="post">
      <div><textarea name="content" rows="3" cols="60"></textarea></div>
      <div><input type="submit" value="create"></div>
    </form>
  	{{ $CurrentUser := .CurrentUser}}
  	<ul>
    {{range .Greetings}}
    	<li>
      	{{with .Author}}
        	<b>{{.}}</b> said:
      	{{else}}
        	Guest said:
      	{{end}}
      	<pre>{{.Content}}</pre>
      	<p>{{ if eq $CurrentUser .Author}} <a href="/edit?id={{.Id}}">Edit</a> / <a href="/destroy?id={{.Id}}">Delete</a> {{end}}</p>
    	</li>
    {{end}}
    </ul>

  </body>
</html>
`

func create(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	u := user.Current(c)
	if u == nil {
		fmt.Fprint(w, "permission denied")
		return
	}
	g := Greeting{
		Author:  u.String(),
		Content: r.FormValue("content"),
		Date:    time.Now(),
	}

	key := datastore.NewIncompleteKey(c, "Greeting", guestbookKey(c))
	_, err := datastore.Put(c, key, &g)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func edit(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	u := user.Current(c)
	if u == nil {
		fmt.Fprint(w, "permission denied")
		return
	}
	id, err := strconv.Atoi(r.FormValue("id"))
	key := datastore.NewKey(c, "Greeting", "", int64(id), guestbookKey(c))
	greeting := new(Greeting)
	err = datastore.Get(c, key, greeting)
	greeting.Id = key.IntID()
	if u.String() != greeting.Author {
		fmt.Fprint(w, "permission denied")
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := editTemplate.Execute(w, greeting); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

var editTemplate = template.Must(template.New("edit").Funcs(template.FuncMap{"eq": func(a, b string) bool { return a == b }}).Parse(editTemplateHTML))

const editTemplateHTML = `
<html>
	<body>
		<form action="/update" method="post">
			<div><textarea name="content" rows="3" cols="60">{{.Content}}</textarea></div>
      		<div><input type="submit" value="update"><input type="hidden" name="id" value={{.Id}} /></div>
		</form>
	</body>
</html>
`

func update(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	u := user.Current(c)
	if u == nil {
		fmt.Fprint(w, "permission denied")
		return
	}

	id, err := strconv.Atoi(r.FormValue("id"))
	key := datastore.NewKey(c, "Greeting", "", int64(id), guestbookKey(c))
	greeting := new(Greeting)
	greeting.Id = key.IntID()
	err = datastore.Get(c, key, greeting)

	if u.String() != greeting.Author {
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

func destroy(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	u := user.Current(c)

	if u == nil {
		fmt.Fprint(w, "permission denied")
	}

	id, err := strconv.Atoi(r.FormValue("id"))
	key := datastore.NewKey(c, "Greeting", "", int64(id), guestbookKey(c))
	err = datastore.Delete(c, key)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Location", "/")
	w.WriteHeader(http.StatusFound)
}
