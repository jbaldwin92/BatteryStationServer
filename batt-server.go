//The main webserver
package main

import (
//  "fmt"
  "log"
  "net/http"
  )
  
func main() {

// Some Examples
//http.Handle("/foo", fooHandler)
//http.HandleFunc("/bar", func(w http.ResponseWriter, r *http.Request) {
//	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
//})

http.HandleFunc("/", mainpage)


log.Fatal(http.ListenAndServe(":1721", nil))
}

func mainpage(w http.ResponseWriter, r *http.Request) {
str1 := `<b>Batt Server</b>
<br>
Eventually, you can choose which power source to use.
<br>
Eventually, you can see the battery voltages and state of charge.
<br>
Eventually, you can see some power use plots, and see how much peak power has been saved.
<br>
Eventually, you can set the time when batteries are used, and the time when batteries are charged.`

w.Write([]byte(str1))
}
