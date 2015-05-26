package main

import (
  "fmt"
  "io/ioutil"
  "strconv"
  "strings"
  "time"
  "github.com/jbaldwin92/bbb_io"   //bbb io functions
)
func main() {
  led_init()
  analog_init()
  var value1, value2 float64
  var values1, values2 string
  for i:=0; i<5000; i++ {
    led_off("1")
    led_on("0")
    time.Sleep(time.Second*2)
    led_on("1")
    led_off("0")
    time.Sleep(time.Second*2)
    value1 = analogReadN("P9_39",200) * 70.82  //36v conversion factor to volts
    value2 = analogReadN("P9_37",200) * 1  //12v
    values1 = strconv.FormatFloat(value1,'f',2,64)
    values2 = strconv.FormatFloat(value2,'f',2,64)
    fmt.Println(values1)
    fmt.Println(values2)
    fmt.Println("------")
  }
  led_off("1")
  led_off("0")
}
