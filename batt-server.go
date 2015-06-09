//The main webserver
package main

import (
	"fmt"
	"github.com/jbaldwin92/bbb_io"
	"log"
	"net/http"
	"strconv"
	"strings"
        "time"
)
//------Global Variables
        //analog pins
        var AIN []string = []string{"P9_39","P9_40","P9_37","P9_38","P9_33","P9_36","P9_35"}
        //calibration factors: AIN x K = Output
        var K []float64 = []float64{70.82,105.82,1,1,1,1,1}
        //GPIO pins for switching chargers and dischargers
        var SW []string = []string{"P9_11"}


//------Functions
func main() {
        //Initialize the pins
	bbb_io.Analog_init()
        bbb_io.PinMode(SW[0],"OUTPUT")
        //startup the datalogger (runs in parallel)
        go v_logger()

	// Some Examples
	//http.Handle("/foo", fooHandler)
	//http.HandleFunc("/bar", func(w http.ResponseWriter, r *http.Request) {
	//	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	//})

	http.HandleFunc("/", mainpage)

	log.Fatal(http.ListenAndServe(":1721", nil))
}

//-----------------------------------------------------
func mainpage(w http.ResponseWriter, r *http.Request) {
        //Set up the time
        const layout="Jan 2, 2006 at 3:04pm (MST)"
        t := time.Now()        


	SOC1 := 0.00
        SOC1s := ""
        volt := make([]float64, 7)
        volts:= make([]string, 7)
        //read the voltages
        for i:=0; i<7; i++ { 
          volt[i] = bbb_io.AnalogReadN(AIN[i], 100)
          volts[i] = strconv.FormatFloat(volt[i] * K[i], 'f', 2, 64)  //K[i] is the conversion factor
        }
        //check to see if any actions need to be taken
//TODO: turn off discharge before turning on charger
        if r.FormValue("SW0")=="on" { bbb_io.DigitalWrite(SW[0],"HIGH") }
        if r.FormValue("SW0")=="off" { bbb_io.DigitalWrite(SW[0],"LOW") }
//TODO: turn off charger before turning on discharger
//TODO: discharge switch
	//Now write out the page
	str1 := `<h1>Batt Server</h1>
<br>
Eventually, you can see some power use plots, and see how much peak power has been saved.
<br>
Eventually, you can set the time when batteries are used, and the time when batteries are charged.
<br>
<br>
<table>
<tr>
  <td>Time</td>
  <td>[[TIME]]</td>
  <td></td>
<tr>
<th COLSPAN=3>Battery 1</th></tr>
<tr>
  <td>LiFePO4</td>
  <td>12 cells</td>
  <td>edit</td>
</tr>
<tr>
  <td>Voltage</td>
  <td>[[VOLTAGE1]]</td>
  <td>calibrate</td>
</tr>
<tr>
  <td>State of Charge</td>
  <td>[[SOC1]]%</td>
  <td></td>
</tr>
<tr>
  <td>Charge</td>
  <td>[[CHARGE1]]</td>
  <td><a href="/?SW0=on">on</a> <a href="/?SW0=off">off</a></td>
</tr>
<tr>
  <td>Discharge</td>
  <td>[[DISCHARGE1]]</td>
  <td><a href="/?SW1=on">on</a> <a href="/?SW1=off">off</a></td>
</tr>
</table>
<a href="/">refresh</a>
`
        //calculate the SOC
        SOC1 = SOC("LiFePO4",12,volt[0]*K[0])
        SOC1s = strconv.FormatFloat(SOC1,'f',2,64)
	//do the string substitutions
	str1 = strings.Replace(str1, "[[VOLTAGE1]]", volts[0], -1)
        str1 = strings.Replace(str1, "[[SOC1]]", SOC1s, -1)
   	str1 = strings.Replace(str1, "[[CHARGE1]]", "Charger Off", -1)
        str1 = strings.Replace(str1, "[[DISCHARGE1]]", "Supplying Power", -1)
        str1 = strings.Replace(str1, "[[TIME]]", t.Format(layout), -1)

	w.Write([]byte(str1))
}

//Calculate the State of Charge
//Inputs: battery type, number of cells, and current voltage
//Returns a number between 0 and 100. Will return -1 on error.
func SOC(batType string, n int, v float64) float64 {
  SOC := -1.0             //unless proven otherwise
  v_cell := v/float64(n)  //voltage per cell
  if batType=="LiFePO4" {
    if v_cell>3.8 {
      SOC=100.0
    } else if v_cell > 3.3 {
      SOC = 100 - (3.8 - v_cell) / (3.8-3.3) * (100.0-90.0)
    } else if v_cell > 3.2 {
      SOC = 90 - (3.3 - v_cell) / (3.3-3.2) * (90.0 - 20.0)
    } else if v_cell > 2.0 {
      SOC = 20 - (3.2 - v_cell) / (3.2-2.0) * (20.0 - 0.00)
    } else {
      SOC = 0.00
    }
  } else if batType=="Pb" {
     SOC = 50.0
  }
  return SOC
}

//Data logger
func v_logger() {
  for {
    fmt.Println(bbb_io.AnalogReadN(AIN[0],100)*K[0])
    time.Sleep(time.Second*2)
  }
}  
