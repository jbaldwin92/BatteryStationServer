package main

import (
  "fmt"
  "time"
  "github.com/jbaldwin92/bbb_io"   //bbb io functions
)
func main() {
  bbb_io.LED_init()
  bbb_io.Analog_init()
  var value1, value2 float64
  for i:=0; i<5000; i++ {
    bbb_io.LED_off("1")
    bbb_io.LED_on("0")
    time.Sleep(time.Second*2)
    bbb_io.LED_on("1")
    bbb_io.LED_off("0")
    time.Sleep(time.Second*2)
    value1 = bbb_io.AnalogReadN("P9_39",200) * 70.82  //36v conversion factor to volts
    //TODO: Fix this conversion factor
    value2 = bbb_io.AnalogReadN("P9_37",200) * 105.82  //12v
    fmt.Printf("LiFePO4 Battery: %6.2fv %6.f %\n", value1, LiFePO4_SOC(value1/12) )
    //fmt.Println(value1)
    //fmt.Println(LiFePO4_SOC(value1/12) )
    fmt.Printf("\nPb Battery:      %4.2fv\n",value2)
    fmt.Println("------")
  }
  bbb_io.LED_off("1")
  bbb_io.LED_off("0")
}

//input the voltage (per cell), the output is 0-100% state of charge
func LiFePO4_SOC(v float64) float64 {
  var SOC float64
  if v>3.8 {
    SOC=100.0
  } else if v > 3.3 {
    SOC = 100.0 - (3.8 - v) / (3.8 - 3.3) * (100.0-90.0)
  } else if v > 3.2 {
    SOC = 90.0 - (3.3-v) / (3.3-3.2) * (90-20)
  } else if v > 2.0 {
    SOC = 20 - (3.2-v) / (3.2 - 2.0) * (20-0)
  } else {
    SOC = 0.00
  }
  return SOC
}

