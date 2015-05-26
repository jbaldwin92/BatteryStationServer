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
  bbb_io.Led_init()
  bbb_io.Analog_init()
  var value1, value2 float64
  var values1, values2 string
  for i:=0; i<5000; i++ {
    bbb_io.Led_off("1")
    bbb_io.Led_on("0")
    time.Sleep(time.Second*2)
    bbb_io.Led_on("1")
    bbb_io.Led_off("0")
    time.Sleep(time.Second*2)
    value1 = bbb_io.AnalogReadN("P9_39",200) * 70.82  //36v conversion factor to volts
    //TODO: Fix this conversion factor
    value2 = bbb_io.AnalogReadN("P9_37",200) * 120  //12v
    values1 = strconv.FormatFloat(value1,'f',2,64)
    values2 = strconv.FormatFloat(value2,'f',2,64)
    fmt.Println(values1)
    fmt.Println(values2)
    fmt.Println("------")
  }
  bbb_io.Led_off("1")
  bbb_io.Led_off("0")
}
