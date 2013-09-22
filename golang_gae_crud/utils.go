package golang_gae_crud

import "html/template"

var templateFuncMap = template.FuncMap{"eq": func(a, b string) bool { return a == b }}
func newTemplateFromFile(filename string) *template.Template{
    var temp = template.New(filename) 
    return template.Must(temp.Funcs(templateFuncMap).ParseFiles(filename))
}
