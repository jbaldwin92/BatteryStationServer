//The main webserver
package main

import (
  "fmt"
  "log"
  "net/http"
  )
  
func main() {

// Some Examples
//http.Handle("/foo", fooHandler)
//http.HandleFunc("/bar", func(w http.ResponseWriter, r *http.Request) {
//	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
//})

http.Handle("/", mainpage)


log.Fatal(http.ListenAndServe(":8080", nil))
}

func mainpage(w http.ResponseWriter, r *http.Request) {
 w.Write([]byte("Main Page"))
}
