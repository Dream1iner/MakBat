package main

import (
	"fmt"
	"image/color"
	"math"
	"os/exec"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"howett.net/plist"
)

var (
	myApp               fyne.App
	myWindow            fyne.Window
	batHealthV          string
	TemperatureV        string
	DesignCapacityV     string
	MaximumPackVoltageV string
	VoltageV            string
	MinimumPackVoltageV string
	DesignCapacitymAhV  string
	NominalCapaityV     string
	MaxCapacityV        string
	MaxCapacityV2       int
	CurrentCapacityV    string
	CurrentCapacityV2   int
	CurrentCapacityPV2  int
	WeightedRaV         string
	CycleCountV         string
	DesignCycleCountV   string
	QmaxV0              string
	QmaxV1              string
	QmaxV2              string
	CellVoltageV0       string
	CellVoltageV1       string
	CellVoltageV2       string
)

type BatteryData struct {
	LifetimeData LifetimeData `plist:"LifetimeData"`
	WeightedRa   int          `plist:"WeightedRa"` // Ohm
	Qmax         []int        `plist:"Qmax"`
	CellVoltage  []int        `plist:"CellVoltage"`
}

type LifetimeData struct {
	MaximumPackVoltage int `plist:"MaximumPackVoltage"` // V
	MinimumPackVoltage int `plist:"MinimumPackVoltage"` // V

}

type battery struct {
	BatteryData       BatteryData `plist:"BatteryData"`
	Voltage           int         `plist:"AppleRawBatteryVoltage"`  //  V
	CurrentCapacity   int         `plist:"AppleRawCurrentCapacity"` // mAh
	MaxCapacity       int         `plist:"AppleRawMaxCapacity"`     // mAh
	DesignCapacity    int         `plist:"MaxCapacity"`             // Wh
	NominalCapaity    int         `plist:"NominalChargeCapacity"`   // mAh
	DesignCapacitymAh int         `plist:"DesignCapacity"`          // Mah
	CurrentCapacityP  int         `plist:"CurrentCapacity"`         // % CellVoltage
	DeltaLimit        float32
	CellVoltageDealta int // calculate in %
	Temperature       int `plist:"Temperature"` // Celcius
	CycleCount        int `plist:"CycleCount"`
	DesignCycleCount  int `plist:"DesignCycleCount9C"`
}

func getData() ([]*battery, error) {
	out, err := exec.Command("ioreg", "-n", "AppleSmartBattery", "-r", "-a").Output()
	if err != nil {
		return nil, err
	}

	// fmt.Println(string(out)) // Print raw output for debug

	if len(out) == 0 {
		return nil, nil
	}

	var data []*battery
	if _, err := plist.Unmarshal(out, &data); err != nil {
		return nil, err
	}

	return data, nil
}

func throwData() {
	batteries, err := getData()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	for _, batt := range batteries {

		DesignCapacityV = ("Designed energy capacity: " + strconv.Itoa(batt.DesignCapacity) + "Wh")
		MaximumPackVoltageV = ("Maximum pack voltage : " + strconv.FormatFloat(float64(batt.BatteryData.LifetimeData.MaximumPackVoltage)/1000, 'f', -1, 64) + "V")
		VoltageV = ("Current voltage:" + strconv.FormatFloat(float64(batt.Voltage)/1000.0, 'f', -1, 64) + "V")
		MinimumPackVoltageV = ("Minimum pack voltage: " + strconv.FormatFloat(float64(batt.BatteryData.LifetimeData.MinimumPackVoltage)/1000, 'f', -1, 64) + "V")

		DesignCapacitymAhV = ("Design capacity: " + strconv.Itoa(batt.DesignCapacitymAh) + "mAh")
		NominalCapaityV = ("Nominal capacity: " + strconv.Itoa(batt.NominalCapaity) + "mAh")
		MaxCapacityV = ("Max capacity: " + strconv.Itoa(batt.MaxCapacity) + "mAh")
		MaxCapacityV2 = batt.MaxCapacity
		CurrentCapacityV = ("Current capacity: " + strconv.Itoa(batt.CurrentCapacity) + "mAh")
		CurrentCapacityV2 = batt.CurrentCapacity
		CurrentCapacityPV2 = batt.CurrentCapacityP

		WeightedRaV = ("Measured resistance: " + strconv.Itoa(batt.BatteryData.WeightedRa) + " Ohm")
		CycleCountV = ("Cycles: " + strconv.Itoa(batt.CycleCount))
		DesignCycleCountV = ("Designed cycles: " + strconv.Itoa(batt.DesignCycleCount))
		QmaxV0 = ("Cell 1 Charge: " + strconv.Itoa(batt.BatteryData.Qmax[0]) + "mAh")
		QmaxV1 = ("Cell 2 Charge: " + strconv.Itoa(batt.BatteryData.Qmax[1]) + "mAh")
		QmaxV2 = ("Cell 3 Charge: " + strconv.Itoa(batt.BatteryData.Qmax[2]) + "mAh")

		CellVoltageV0 = "Cell 1 voltage: " + strconv.FormatFloat(float64(batt.BatteryData.CellVoltage[0])/1000, 'f', -1, 64) + "V"
		CellVoltageV1 = "Cell 2 voltage: " + strconv.FormatFloat(float64(batt.BatteryData.CellVoltage[1])/1000, 'f', -1, 64) + "V"
		CellVoltageV2 = "Cell 3 voltage: " + strconv.FormatFloat(float64(batt.BatteryData.CellVoltage[2])/1000, 'f', -1, 64) + "V"

		TemperatureV = ("Temperature: " + strconv.FormatFloat(float64(batt.Temperature)/100, 'f', -1, 64) + " °C")
		batHealthV = "Battery health: " + strconv.Itoa(int(math.Round(float64(batt.MaxCapacity)/float64(batt.DesignCapacitymAh)*100))) + "%"

		// In case you have some unfulfilled desire to print values to output

		//fmt.Println("Maximum pack voltage:", float64(batt.BatteryData.LifetimeData.MaximumPackVoltage)/1000, "V")
		//fmt.Println("Design energy сapacity:", batt.DesignCapacity, "Wh")
		//fmt.Println("Current voltage:", float64(batt.Voltage)/1000.0, "V")
		//fmt.Println("Minimum pack voltage:", float64(batt.BatteryData.LifetimeData.MinimumPackVoltage)/1000, "V")
		//fmt.Println("Cell 1 voltage:", float64(batt.BatteryData.CellVoltage[0])/1000, "V")
		//fmt.Println("Cell 2 voltage:", float64(batt.BatteryData.CellVoltage[1])/1000, "V")
		//fmt.Println("Cell 3 voltage:", float64(batt.BatteryData.CellVoltage[2])/1000, "V")
		//fmt.Println("Design capacity:", batt.DesignCapacitymAh, "mAh")
		//fmt.Println("Nominal capacity:", batt.NominalCapaity, "mAh")
		//fmt.Println("Max capacity:", batt.MaxCapacity, "mAh")
		//fmt.Println("Current capacity:", batt.CurrentCapacity, "mAh")
		//fmt.Println("Current charge level:", batt.CurrentCapacityP, "%")
		//fmt.Println("Temperature:", float64(batt.Temperature)/100, "°C")
		//fmt.Println("Measured resistance:", batt.BatteryData.WeightedRa)
		//fmt.Println("Cycles counter:", batt.CycleCount)
		//fmt.Println("Design cycle count:", batt.DesignCycleCount)
		//fmt.Println("Cell 1 Charge:", batt.BatteryData.Qmax[0], "mAh")
		//fmt.Println("Cell 2 Charge:", batt.BatteryData.Qmax[1], "mAh")
		//fmt.Println("Cell 3 Charge:", batt.BatteryData.Qmax[2], "mAh")

		// calculating delta V per 1->2, 1->3, 2->3

		/*
			if float64(batt.BatteryData.CellVoltage[0]-batt.BatteryData.CellVoltage[1])/100.0 <= 0.1 {
				fmt.Println("Delta voltage between Cell 1 and Cell 2 is OK. ")
			} else {
				fmt.Println("Potential issue with voltage delta between Cell 2 and Cell 3")
			}

			if float64(batt.BatteryData.CellVoltage[0]-batt.BatteryData.CellVoltage[2])/100.0 <= 0.1 {
				fmt.Println("Delta voltage between Cell 1 and Cell 3 is OK. ")
			} else {
				fmt.Println("Potential issue with voltage delta between Cell 2 and Cell 3")
			}

			if float64(batt.BatteryData.CellVoltage[1]-batt.BatteryData.Qmax[2])/100.0 <= 0.1 {
				fmt.Println("Delta voltage between Cell 2 and Cell 3 is OK. ")
			} else {
				fmt.Println("Potential issue with voltage delta between Cell 2 and Cell 3")
			}

			// calculating delta mAh per 1->2, 1->3, 2->3
			if float64(batt.BatteryData.Qmax[0]-batt.BatteryData.Qmax[1])/100.0 <= 0.5 {
				fmt.Println("Charge delta between Cell 1 and Cell 2 is OK.")
			} else {
				fmt.Println("Potential issue with charge delta between Cell 1 and Cell 2")
			}

			if float64(batt.BatteryData.Qmax[0]-batt.BatteryData.Qmax[2])/100.0 <= 0.5 {
				fmt.Println("Charge delta between Cell 1 and Cell 3 is OK.")
			} else {
				fmt.Println("Potential issue with charge delta between Cell 1 and Cell 3")
			}

			if float64(batt.BatteryData.Qmax[1]-batt.BatteryData.Qmax[2])/100.0 <= 0.5 {
				fmt.Println("Charge delta between Cell 2 and Cell 3 is OK.")
			} else {
				fmt.Println("Potential issue with charge delta between Cell 2 and Cell 3")
			}
		*/

	}

}

func main() {
	myApp = app.New()
	myWindow = myApp.NewWindow("MakBat")
	myWindow.Resize(fyne.NewSize(640, 480))
	myWindow.SetFixedSize(true)

	throwData()

	// labels for battery data
	DesignCapacityVLabel := widget.NewLabel(DesignCapacityV)
	MaximumPackVoltageVLabel := widget.NewLabel(MaximumPackVoltageV)
	VoltageVlabel := widget.NewLabel(VoltageV)
	MinimumPackVoltageVLabel := widget.NewLabel(MinimumPackVoltageV)
	DesignCapacitymAhVLabel := widget.NewLabel(DesignCapacitymAhV)
	NominalCapaityVlabel := widget.NewLabel(NominalCapaityV)
	MaxCapacityVLabel := widget.NewLabel(MaxCapacityV)
	CurrentCapacityVLabel := widget.NewLabel(CurrentCapacityV)
	CycleCountVLabel := widget.NewLabel(CycleCountV)
	DesignCycleCountVLabel := widget.NewLabel(DesignCycleCountV)
	healthLabel := widget.NewLabel(batHealthV)
	temperatureLabel := widget.NewLabel(TemperatureV)
	// detailed data
	WeightedRaVLabel := widget.NewLabel(WeightedRaV)
	QmaxV0Label := widget.NewLabel(QmaxV0)
	QmaxV1Label := widget.NewLabel(QmaxV1)
	QmaxV2Label := widget.NewLabel(QmaxV2)
	CellVoltageV0Label := widget.NewLabel(CellVoltageV0)
	CellVoltageV1Label := widget.NewLabel(CellVoltageV1)
	CellVoltageV2Label := widget.NewLabel(CellVoltageV2)

	// in fyne you can only use once each instance of separator
	separator := widget.NewSeparator()
	separator0 := widget.NewSeparator()
	separator1 := widget.NewSeparator()
	separator2 := widget.NewSeparator()
	separator3 := widget.NewSeparator()

	// yep, that sounds wierd
	vSeparator := canvas.NewRectangle(color.Black)
	vSeparator.Resize(fyne.NewSize(20, 50))

	vSeparator1 := canvas.NewRectangle(color.Black)
	vSeparator1.Resize(fyne.NewSize(20, 50))

	vSeparator2 := canvas.NewRectangle(color.Black)
	vSeparator2.Resize(fyne.NewSize(20, 50))

	vSeparator3 := canvas.NewRectangle(color.Black)
	vSeparator3.Resize(fyne.NewSize(20, 50))

	vSeparator4 := canvas.NewRectangle(color.Black)
	vSeparator4.Resize(fyne.NewSize(20, 50))

	vSeparator5 := canvas.NewRectangle(color.Black)
	vSeparator5.Resize(fyne.NewSize(20, 50))

	vSeparator6 := canvas.NewRectangle(color.Black)
	vSeparator6.Resize(fyne.NewSize(20, 50))

	vSeparator7 := canvas.NewRectangle(color.Black)
	vSeparator7.Resize(fyne.NewSize(20, 50))

	// Create a spacer with a custom size
	customSpacer := canvas.NewRectangle(color.Transparent)
	customSpacer.Resize(fyne.NewSize(350, 480))

	// Progress bar for charge level in %
	// Value reported by BMS
	pbExpalinerlabel1 := widget.NewLabel("BMS reported value")
	currentCapacityProgressBar1 := widget.NewProgressBar()
	currentCapacityProgressBar1.Max = 100
	currentCapacityProgressBar1.Min = 0
	currentCapacityProgressBar1.Value = float64(CurrentCapacityPV2)

	// Progress bar for charge level in %
	// Calulated value
	pbExpalinerlabel2 := widget.NewLabel("Calculated value")
	currentCapacityProgressBar2 := widget.NewProgressBar()
	currentCapacityProgressBar2.Max = 100
	currentCapacityProgressBar2.Min = 0
	currentCapacityProgressBar2.Value = 0
	currentCapacityProgressBar2.Value = float64((CurrentCapacityV2 * 100) / MaxCapacityV2)

	batteryDataContainerH1 := container.NewHBox(DesignCapacityVLabel)

	batteryDataContainer1 := container.NewHBox(
		MaximumPackVoltageVLabel,
		VoltageVlabel,
		MinimumPackVoltageVLabel)

	batteryDataContainerH2 := container.NewHBox(
		DesignCapacitymAhVLabel,
		vSeparator7,
		NominalCapaityVlabel)

	batteryDataContainerH3 := container.NewHBox(
		vSeparator5,
		MaxCapacityVLabel,
		vSeparator6,
		CurrentCapacityVLabel,
	)

	batteryDataContainerH4 := container.NewHBox(
		healthLabel,
		vSeparator2,
		CycleCountVLabel,
		vSeparator3,
		DesignCycleCountVLabel,
		vSeparator4,
		temperatureLabel,
	)

	batteryDataContainer := container.NewVBox(
		customSpacer,
		pbExpalinerlabel1,
		currentCapacityProgressBar1,
		pbExpalinerlabel2,
		currentCapacityProgressBar2,
		batteryDataContainer1,
		separator0,
		batteryDataContainerH1,
		separator,
		batteryDataContainerH2,
		separator2,
		batteryDataContainerH3,
		separator3,
		batteryDataContainerH4,
	)

	detailedDataContainer0 := container.NewVBox(QmaxV0Label, CellVoltageV0Label)
	detailedDataContainer1 := container.NewVBox(QmaxV1Label, CellVoltageV1Label)
	detailedDataContainer2 := container.NewVBox(QmaxV2Label, CellVoltageV2Label)

	detailedDataContainerH := container.NewHBox(
		detailedDataContainer0,
		vSeparator,
		detailedDataContainer1,
		vSeparator1,
		detailedDataContainer2)

	detailedDataContainer := container.NewVBox(
		WeightedRaVLabel,
		separator1,
		detailedDataContainerH,
		separator,
	)

	// Content for tabs
	tabs := container.NewAppTabs(
		container.NewTabItem("Battery data", batteryDataContainer),
		container.NewTabItem("Detailed data", detailedDataContainer),
		container.NewTabItem("About dev", widget.NewLabel("Murat")),
		// Add more widgets as needed to display additional information
	)

	// Exit button
	buttonExit := widget.NewButton("Exit", func() {
		myApp.Quit()
	})

	// Create a container for buttons with custom spacers
	buttonsContainer := container.New(layout.NewVBoxLayout(),
		layout.NewSpacer(),
		customSpacer,
		customSpacer,
		customSpacer,
		customSpacer,
		customSpacer,
		customSpacer,
		customSpacer,
		customSpacer,
		buttonExit,
	)

	// Combine topBar and appsGrid into a single container
	content := container.New(layout.NewVBoxLayout(), tabs, buttonsContainer)

	// Creating go routine to handle UI updates
	go func() {
		for {
			time.Sleep(time.Second * 5)
			throwData()
			currentCapacityProgressBar1.SetValue((float64(CurrentCapacityPV2)))
			//fmt.Println("new % level: " + strconv.Itoa(CurrentCapacityPV2)) //debug message
			currentCapacityProgressBar2.Value = float64((CurrentCapacityV2 * 100) / MaxCapacityV2)
			myWindow.Canvas().Refresh(content)
		}
	}()

	// Set the content container as the content of the window
	myWindow.SetContent(content)

	myWindow.ShowAndRun()
}
