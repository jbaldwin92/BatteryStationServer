package main

import (
  "fmt"
  "io/ioutil"
  "strconv"
  "strings"
  "time"
)
//Global Variables
//***************************
//(Thanks to gobot.io for these!)
var pins = map[string]int{
	"P8_3":  38,
	"P8_4":  39,
	"P8_5":  34,
	"P8_6":  35,
	"P8_7":  66,
	"P8_8":  67,
	"P8_9":  69,
	"P8_10": 68,
	"P8_11": 45,
	"P8_12": 44,
	"P8_13": 23,
	"P8_14": 26,
	"P8_15": 47,
	"P8_16": 46,
	"P8_17": 27,
	"P8_18": 65,
	"P8_19": 22,
	"P8_20": 63,
	"P8_21": 62,
	"P8_22": 37,
	"P8_23": 36,
	"P8_24": 33,
	"P8_25": 32,
	"P8_26": 61,
	"P8_27": 86,
	"P8_28": 88,
	"P8_29": 87,
	"P8_30": 89,
	"P8_31": 10,
	"P8_32": 11,
	"P8_33": 9,
	"P8_34": 81,
	"P8_35": 8,
	"P8_36": 80,
	"P8_37": 78,
	"P8_38": 79,
	"P8_39": 76,
	"P8_40": 77,
	"P8_41": 74,
	"P8_42": 75,
	"P8_43": 72,
	"P8_44": 73,
	"P8_45": 70,
	"P8_46": 71,
	"P9_11": 30,
	"P9_12": 60,
	"P9_13": 31,
	"P9_14": 50,
	"P9_15": 48,
	"P9_16": 51,
	"P9_17": 5,
	"P9_18": 4,
	"P9_19": 13,
	"P9_20": 12,
	"P9_21": 3,
	"P9_22": 2,
	"P9_23": 49,
	"P9_24": 15,
	"P9_25": 117,
	"P9_26": 14,
	"P9_27": 115,
	"P9_28": 113,
	"P9_29": 111,
	"P9_30": 112,
	"P9_31": 110,
}

var pwmPins = map[string]string{
	"P9_14": "P9_14",
	"P9_21": "P9_21",
	"P9_22": "P9_22",
	"P9_29": "P9_29",
	"P9_42": "P9_42",
	"P8_13": "P8_13",
	"P8_34": "P8_34",
	"P8_45": "P8_45",
	"P8_46": "P8_46",
}

var analogPins = map[string]string{
	"P9_39": "AIN0",
	"P9_40": "AIN1",
	"P9_37": "AIN2",
	"P9_38": "AIN3",
	"P9_33": "AIN4",
	"P8_36": "AIN5",
	"P8_35": "AIN6",
}

//Functions
//****************************************

func main() {
  led_init()
  analog_init()
  pinMode("P9_31","INPUT")
  pwm_init("P8_13")
  var value float64
  var value2 string
  for i:=0; i<2; i++ {
    led_off("1")
    led_on("0")
    time.Sleep(time.Second*1)
    led_on("1")
    led_off("0")
    time.Sleep(time.Second*1)
    value = analogRead("P9_39")
    value2 = digitalRead("P9_31")
    fmt.Println(value)
    fmt.Println(value2)
  }
  led_off("1")
  analogWrite("P8_13",0,5000000)
}
//Initialize pulse width modulation
func pwm_init(pinName string) {
 if pwmPins[pinName] == pinName { //this checks to see if the pinName is a special pwm pin
   ioutil.WriteFile("/sys/devices/bone_capemgr.*/slots",[]byte("am33xx_pwm"),0444)
   ioutil.WriteFile("/sys/devices/bone_capemgr.*/slots",[]byte("bone_pwm_"+pinName),0444)
 } else {
   //TODO: error
   fmt.Println("pin name was not a pwm pin")
 }
}


//Initialize a gpio pin (e.g., P9_20)
//direction can be "OUTPUT" or "INPUT" 
func pinMode(pinName string, direction string) {
 pin := pins[pinName]  //gets the integer from the pin list
 pin_str := strconv.Itoa(pin)
 //export the pin
 ioutil.WriteFile("/sys/class/gpio/export",[]byte(pin_str),0444)
 //set the direction
 var dir string
 if direction=="OUTPUT" {
   dir = "low"  //this is the same as "out" but starts low
 } else if direction=="INPUT" {
   dir = "in"  //TODO: what about INPUT_PULLUP??
 } else {
   //TODO: make this an error
 }
 ioutil.WriteFile("/sys/class/gpio/gpio"+pin_str+"/direction",[]byte(dir),0444)
}

//Make a pin go HIGH or LOW. The pinName should be (for example) "P9_20".
//TODO: check to make sure the pin was exported first
func digitalWrite(pinName string, value string) {
  pin := pins[pinName] 
  pin_str := strconv.Itoa(pin)
  var val string
  if value=="HIGH" {
    val = "1"
  } else if value=="LOW" {
    val = "0"
  } else {
    //TODO Error
  }
  ioutil.WriteFile("/sys/class/gpio/gpio"+pin_str+"/value",[]byte(val),0444)
}

//Read a HIGH or LOW result, for a pin that's an input. pinName should be (for example) "P9_20".
func digitalRead(pinName string) string {
  pin := pins[pinName]
  pin_str := strconv.Itoa(pin)
  val,_ := ioutil.ReadFile("/sys/class/gpio/gpio"+pin_str+"/value")
  vals := byteArrayToString(val)
  var val_str string
  if vals=="1" {
    val_str = "HIGH"
  } else if vals=="0" {
    val_str = "LOW"
  } else {
    val_str = "ERROR"
    //TODO err
  }
  return val_str  //returns "HIGH" or "LOW"
}

//PWM. Inputs are the pinName, duty cyle % (some number beween 0 and 1), and freq (Hz)
//These get translated before being written to the bbb
//note: max pwm freq is around 9 MHz
func analogWrite(pinName string, dc float64, freq float64) {
  //what's the period? In nanoseconds
  period := 1/freq*1000000000.0  //float64
  duty := dc * period
fmt.Println(int(period))
fmt.Println(int(duty))
  //This function assumes a positive polarity. Voltage starts at 3.3v and stays until the end of the duty time
  //duty and period are in nanoseconds. duty must be less than period
  ioutil.WriteFile("/sys/devices/ocp.3/pwm_test_"+pinName+".16/duty",[]byte("0"),0444)
  ioutil.WriteFile("/sys/devices/ocp.3/pwm_test_"+pinName+".16/period",[]byte("0"),0444)
  ioutil.WriteFile("/sys/devices/ocp.3/pwm_test_"+pinName+".16/run",[]byte("1"),0444)
  ioutil.WriteFile("/sys/devices/ocp.3/pwm_test_"+pinName+".16/polarity",[]byte("1"),0444)
  ioutil.WriteFile("/sys/devices/ocp.3/pwm_test_"+pinName+".16/period",[]byte(strconv.Itoa(int(period))),0444)
  ioutil.WriteFile("/sys/devices/ocp.3/pwm_test_"+pinName+".16/duty",[]byte(strconv.Itoa(int(duty))),0444)
}

//configure for analog out
func analog_init() {
  ioutil.WriteFile("/sys/devices/bone_capemgr.*/slots",[]byte("cape-bone-iio"),0444)
}
//returns a value that is 0-1.8v
func analogRead(pinName string) float64 {
  pin_str := analogPins[pinName] 
  val,_ := ioutil.ReadFile("/sys/devices/ocp.3/helper.15/"+pin_str)
  //transform val into other types
  vals := byteArrayToString(val)
  valint,err := strconv.Atoi(vals)
  if err != nil { 
    fmt.Println(err)
  }
  val64 := float64(valint)/1000.0
  return val64
}

func byteArrayToString(input []byte) string {
  n := -1
  for i,b := range input {
    if b == 0 {
      break
    }
    n=i
 }
  return strings.TrimSpace(string(input[:n+1]))
}


//initialize the leds; turn them all off and set the trigger to "none"
func led_init() {
  var s string
  for i:=0; i<4; i++ {
    s = strconv.Itoa(i)
    trigger(s,"none")
    led_off(s)
  }
}

//set the trigger of the led. "id" is a number "0" thru "3" Acceptable Triggers: 
//none nand-disk mmc0 mmc1 timer oneshot heartbeat backlight gpio cpu0 default-on transient
func trigger(id string, trig string) {
  ioutil.WriteFile("/sys/class/leds/beaglebone:green:usr"+id+"/trigger",[]byte(trig),0444)
}

//turn the led on
func led_on(id string) { 
//turn led on
  ioutil.WriteFile("/sys/class/leds/beaglebone:green:usr"+id+"/brightness",[]byte("255"),0444)
}

//Turn the led off. "id" is a number "0" thru "3"
func led_off(id string) {
  //turn led off
  ioutil.WriteFile("/sys/class/leds/beaglebone:green:usr"+id+"/brightness",[]byte("0"),0444)
}

