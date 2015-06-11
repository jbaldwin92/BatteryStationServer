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
        var K []float64 = []float64{71.85,105.82,1,1,1,1,1}
        //number of battery cells, for help with the SOC
        var Cells []int64 = []int64{12,1,1,1,1,1,1}
        //GPIO pins for switching chargers and dischargers
        var SW []string = []string{"P9_11"}
        //Values
        var old_values []string = []string{}

//------Functions
func main() {
        //Initialize the pins
	bbb_io.Analog_init()
        bbb_io.PinMode(SW[0],"OUTPUT")
        //startup the datalogger (runs in parallel)
        go v_logger()
        //startup the on & off timer
        go charging_timer()

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
        volts_cell := make([]string,7)
        //read the voltages
        for i:=0; i<7; i++ { 
          volt[i] = bbb_io.AnalogReadN(AIN[i], 100)
          volts[i] = strconv.FormatFloat(volt[i] * K[i], 'f', 2, 64)  //K[i] is the conversion factor
          volts_cell[i] = strconv.FormatFloat(volt[i] * K[i]/float64(Cells[i]), 'f', 2, 64)  //number of volts per cell
        }
        //check to see if any actions need to be taken
//TODO: turn off discharge before turning on charger
        if r.FormValue("SW0")=="on" { bbb_io.DigitalWrite(SW[0],"HIGH") }
        if r.FormValue("SW0")=="off" { bbb_io.DigitalWrite(SW[0],"LOW") }
//TODO: turn off charger before turning on discharger
//TODO: discharge switch
	//Now write out the page
	str1 := `<html>
<head>
<script src="http://d3js.org/d3.v3.min.js"></script>
<script src="http://dimplejs.org/dist/dimple.v2.1.2.min.js"></script>
</head>
<body>
<h1>Batt Server</h1>
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
  <td>[[VOLTAGE1]]v<br>[[VOLTAGE1_CELL]]v/cell</td>
  <td>calibrate</td>
</tr>
<tr>
  <td>State of Charge</td>
  <td>[[SOC1]]%</td>
  <td>42.35V is about 100%</td>
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
<br>
<br>
[[OLD_VALUES]]
<div id="chartContainer">
<script type="text/javascript">
var svg = dimple.newSvg("#chartContainer",590,400);
var data = [
 {"X":"1","Y":"2.5" },
 {"X":"3.2","Y":"53" }
];
var chart = new dimple.chart(svg,data);
chart.setBounds(60,30,510,305)
chart.addCategoryAxis("x","X");
chart.addMeasureAxis("y","Y");
chart.addSeries(null,dimple.plot.line);
chart.draw();
</script>
</div>

</body>
</html>
`
        //calculate the SOC
        SOC1 = SOC("LiFePO4",Cells[0],volt[0]*K[0])
        SOC1s = strconv.FormatFloat(SOC1,'f',2,64)
        chargerSwitch := bbb_io.DigitalRead("P9_11") 
        //put the old values into a long string
        var old_values_list string
        for _,v := range(old_values) {
          old_values_list = old_values_list + v + "<br>\n"
        }	
	//do the string substitutions
	str1 = strings.Replace(str1, "[[VOLTAGE1]]", volts[0], -1)
        str1 = strings.Replace(str1, "[[VOLTAGE1_CELL]]", volts_cell[0], -1)
        str1 = strings.Replace(str1, "[[SOC1]]", SOC1s, -1)
   	if chargerSwitch == "LOW" {
          str1 = strings.Replace(str1, "[[CHARGE1]]", "Charger is Off", -1)
        } else {
          str1 = strings.Replace(str1, "[[CHARGE1]]", "Charger is On", -1)
        }
        str1 = strings.Replace(str1, "[[DISCHARGE1]]", "Supplying Power", -1)
        str1 = strings.Replace(str1, "[[TIME]]", t.Format(layout), -1)
        str1 = strings.Replace(str1, "[[OLD_VALUES]]", old_values_list, -1)
	w.Write([]byte(str1))
}

//Calculate the State of Charge
//Inputs: battery type, number of cells, and current voltage
//Returns a number between 0 and 100. Will return -1 on error.
func SOC(batType string, n int64, v float64) float64 {
  SOC := -1.0             //unless proven otherwise
  v_cell := v/float64(n)  //voltage per cell
  if batType=="LiFePO4" {
    if v_cell>3.8 {
      SOC=100.01
    } else if v_cell>3.53 {  //this is unique to my charger, which charges to 3.5v/cell
      SOC=99.99
    } else if v_cell > 3.3 {
      //SOC = 100 - (3.8 - v_cell) / (3.8-3.3) * (100.0-90.0)  //this is the general value
      SOC = 100 - (3.53 - v_cell) / (3.53-3.3) * (100 - 90.0)  //unique to my charger
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
  var voltage float64
  var voltages string
  for {
    voltage = bbb_io.AnalogReadN(AIN[0],100)*K[0]
    voltages = strconv.FormatFloat(voltage,'f',4,64)
    fmt.Println(voltage)
    old_values = append(old_values,voltages)
    time.Sleep(time.Second*15)
  }
}  


//Turn on and off the switch at a certain time
func charging_timer() {
  phase := "off"
  h_on := 1
  m_on := 0
  s_on := 0
  h_off := 9
  m_off := 0
  s_off :=0
  t := time.Now()
  h,m,s := t.Clock()
  for {
    t = time.Now()
    h,m,s = t.Clock()
    //is it time to turn on the charger?
    if phase == "off" {
      if h >= h_on {
      if m >= m_on {
      if s >= s_on {
        phase = "on"
        bbb_io.DigitalWrite(SW[0],"HIGH") //turn it on 
      }}}
    }
    //is it time to turn off the charger?      
    if phase == "on" {
      if h >= h_off {
      if m >= m_off {
      if s >= s_off {
        phase = "off"
        bbb_io.DigitalWrite(SW[0],"LOW") 
      }}}
    }
    time.Sleep(time.Second*1)
  }
} 
    
