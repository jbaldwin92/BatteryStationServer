//The main webserver
package main

import (
//  "fmt"
  "log"
  "github.com/jbaldwin92/bbb_io"
  "net/http"
  "strings"
  "strconv"
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

bbb_io.Analog_init()
v := bbb_io.AnalogReadN("P9_39",100) * 71.13
vs := strconv.FormatFloat(v,'f',2,64)
str1 := `<h1>Batt Server</h1>
<br>
Eventually, you can choose which power source to use.
<br>
Eventually, you can see the battery voltages and state of charge.
<br>
Eventually, you can see some power use plots, and see how much peak power has been saved.
<br>
Eventually, you can set the time when batteries are used, and the time when batteries are charged.
<br>
Voltage = [[VOLTAGE]]`

str1 = strings.Replace(str1,"[[VOLTAGE]]",vs,-1)



w.Write([]byte(str1))
}
