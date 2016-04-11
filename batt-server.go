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
        var Cells []int64 = []int64{10,1,1,1,1,1,1}
        //Battery Type
        var batteryType []string = []string{"Li-Ion","Pb","","","","",""}
        //GPIO pins for switching chargers and dischargers
        var SW []string = []string{"P9_11","P9_13"}  //charge, discharge
        //Values
        var old_values []string = []string{}
        //Time
        var time_list []string = []string{}
        //Manual On or Off
        //only the hours and minutes are used (not the year or month or day)
        var ontime_hour int = 13
        var ontime_minute int = 00
        var offtime_hour int = 21
        var offtime_minute int =00
        var chargeon_hour int = 21
        var chargeon_minute int = 2
        var chargeoff_hour int = 6
        var chargeoff_minute int = 0 
        //Minimum State of Charge
        var StayAboveSOC float64 = 41  //this keeps it above about 50%

//------Functions
func main() {
        //Initialize the pins
	bbb_io.Analog_init()
        bbb_io.PinMode(SW[0], "OUTPUT" )  //charger
        bbb_io.PinMode(SW[1], "OUTPUT" )  //discharger
        //startup the datalogger (runs in parallel)
        go v_logger()
        //startup the on & off timer
        //go charging_timer()

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
        if r.FormValue("SW1")=="on" {bbb_io.DigitalWrite(SW[1], "HIGH") }
        if r.FormValue("SW1")=="off" {bbb_io.DigitalWrite(SW[1],"LOW") }
	//Now write out the page
	str1 := `<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<meta http-equiv="X-UA-Compatible" content="IE=edge">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>Home</title>


<script src="http://cdnjs.cloudflare.com/ajax/libs/dygraph/1.1.1/dygraph-combined.js"></script>
<script type="text/javascript">
function reloadfunction() {
  window.setTimeout(function() {location.assign("/");},60000);
}
</script> 
</head>
<body onload="reloadfunction()">
<h1>Batt Server</h1>
<br>
<div id="chartContainer2" style="width:600px; height:300px; border:1px;"></div>
<script type="text/javascript">
g = new Dygraph(
  document.getElementById("chartContainer2"),
  [
    [[PLOT_DATA]]
  ],
   { }
);
</script>


<br>
<table>
<tr>
  <td>Time</td>
  <td>[[TIME]]</td>
  <td></td>
<tr>
<th COLSPAN=3>Battery 1</th></tr>
<tr>
  <td>[[BAT_TYPE]]</td>
  <td>[[NCELLS]] cells</td>
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
  <td>42.35V is about 100%<br>When discharging, the real SOC is 7% higher than shown.</td>
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
Don't forget: use 'screen' to keep the web server running. Here's the format:
<br>
screen -r [PID]
just type 'screen -r' to see what the PID number is
<br>
Eventually, you can set the time when batteries are used, and the time when batteries are charged.
<br>


</body>
</html>
`
        //calculate the SOC
        SOC1 = SOC("Li-Ion",Cells[0],volt[0]*K[0])
        SOC1s = strconv.FormatFloat(SOC1,'f',2,64)
        chargerSwitch := bbb_io.DigitalRead( SW[0] ) 
        dischargerSwitch := bbb_io.DigitalRead( SW[1] )
        //put the old values into a long string
        var plot_data string
        for i,v := range(old_values) {
          plot_data = plot_data + "[ new Date(\"" + time_list[i]+ "\") ," +v+"]"
          if i!=len(old_values) -1 {
            plot_data = plot_data + ","
          } 
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
        if dischargerSwitch == "LOW" {
          str1 = strings.Replace(str1, "[[DISCHARGE1]]", "Not Supplying Power", -1)
        } else {
          str1 = strings.Replace(str1, "[[DISCHARGE1]]", "Supplying Power", -1)
        }
        str1 = strings.Replace(str1, "[[TIME]]", t.Format(layout), -1)
        str1 = strings.Replace(str1, "[[PLOT_DATA]]", plot_data, -1)
	str1 = strings.Replace(str1, "[[BAT_TYPE]]", batteryType[0], -1)
        str1 = strings.Replace(str1, "[[NCELLS]]", strconv.Itoa(int(Cells[0])), -1)
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
      SOC = 100 - (3.8 - v_cell) / (3.8-3.3) * (100.0-90.0)  //this is the general value
    } else if v_cell > 3.2 {
      SOC = 90 - (3.3 - v_cell) / (3.3-3.2) * (90.0 - 20.0)
    } else if v_cell > 2.0 {
      SOC = 20 - (3.2 - v_cell) / (3.2-2.0) * (20.0 - 0.00)
    } else {
      SOC = 0.00
    }
  } else if batType=="Li-Ion" {
    if v_cell>4.2 {
      SOC=100.01
    } else if v_cell>3.95  {  //this is unique to my charger, which charges to 3.5v/cell
      SOC = 100 - (4.2 - v_cell) / (4.2-3.95) * (100-70) 
    } else if v_cell > 3.6 {
      SOC = 70 - (3.95 - v_cell) / (3.95-3.6) * (70.0-10)  //this is the general value
    } else if v_cell > 2.9 {
      SOC = 10 - (3.6 - v_cell) / (3.6-2.9) * (10-0)
    }
  } else if batType=="Pb" {
     SOC = 50.0
  }
  return SOC
}

//Data logger and watchdog
func v_logger() {
  var voltage, percent float64
  var voltages, percents string
  ontime := ontime_hour * 60 + ontime_minute
  offtime:= offtime_hour*60 + offtime_minute
  chargeontime := chargeon_hour * 60 + chargeon_minute
  chargeofftime := chargeoff_hour * 60 + chargeoff_minute  
  StayAboveSOC1 := StayAboveSOC //this leaves the global variable untouched
  var nowtime int
  var t time.Time
  for {
    voltage = bbb_io.AnalogReadN(AIN[0],200)*K[0]
    voltages = strconv.FormatFloat(voltage,'f',4,64)
    percent = SOC("Li-Ion", Cells[0], voltage)
    percents = strconv.FormatFloat(percent,'f',2,64)
    fmt.Println(voltages+", "+percents+"%")
    old_values = append(old_values,percents)
    time_list = append(time_list,time.Now().Format("2006-01-02 15:04:05"))
    if len(old_values)>1500 {
      old_values = old_values[1:]  //this keeps the file from getting too big
      time_list = time_list[1:]
    }
    //check to see if it's time to turn on or off
    t = time.Now()
    nowtime = t.Hour() * 60 + t.Minute()  //minutes since midnight
    if percent > StayAboveSOC1 {  //then see if it should be turned on 
      if nowtime < ontime {  //too early, but check to see if it's time to turn off the charger
         StayAboveSOC1 = StayAboveSOC  //resets this at midnight, if needed
         bbb_io.DigitalWrite(SW[1],"LOW")
         if nowtime > chargeofftime {
           bbb_io.DigitalWrite(SW[0],"LOW")  //or else just leave it on
         }
      } else if nowtime > offtime {  //too late, but check to see if it's time to turn on the charger 
         bbb_io.DigitalWrite(SW[1],"LOW")
         if nowtime > chargeontime {  //turn on the charger
           bbb_io.DigitalWrite(SW[0],"HIGH")
         }
      } else {  //nowtime is between ontime and offtime, turn it on
         bbb_io.DigitalWrite(SW[0],"LOW")  //make sure charger is off
         time.Sleep(1*time.Second)
         bbb_io.DigitalWrite(SW[1],"HIGH") //turn on discharger
      }
    } else {  //the SOC is too low, just shut it off
      StayAboveSOC1 = 200  //this keeps the system from bouncing on and off at the lower limit
      if nowtime > ontime {  
        StayAboveSOC1 = 200 //make sure it doesn't turn back on
      } else {  //then it's after midnight but before the ontime
        StayAboveSOC1 = StayAboveSOC  //reset the variable to the global variable value
      }
      bbb_io.DigitalWrite(SW[1], "LOW")  //turn off the discharger for the rest of the day
      if nowtime > chargeontime {  //but turn on the charger if it's time
           bbb_io.DigitalWrite(SW[0],"HIGH")  //turn on the charger
      }        
    }
    fmt.Println(nowtime)
    fmt.Println(ontime)
    fmt.Println(offtime)
    fmt.Println(StayAboveSOC1)
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
    

