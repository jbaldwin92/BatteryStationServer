package main

import (
  "io/ioutil"
  "time"
)

func main() {
  for i:=0; i<10; i++ {
    led_off("1")
    led_on("0")
    time.Sleep(time.Second*1)
    led_on("1")
    led_off("0")
    time.Sleep(time.Second*1)
  }
  led_off("1")
}

func led_on(id string) { 
//turn led on
  ioutil.WriteFile("/sys/class/leds/beaglebone:green:usr"+id+"/brightness",[]byte("255"),044)
}

func led_off(id string) {
  //turn led off
  ioutil.WriteFile("/sys/class/leds/beaglebone:green:usr"+id+"/brightness",[]byte("0"),0444)
}
