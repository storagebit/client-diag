package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

func checkSystemTuning(){
	fmt.Println(formatBoldWhite("\nServer/Client System Tuning Settings:"))
	checkTuneD()
	checkIRQBalance()
	checkCstate()
}

func checkCstate(){
	if checkIfFileExists("/sys/module/intel_idle/parameters/max_cstate"){
		c, err := ioutil.ReadFile("/sys/module/intel_idle/parameters/max_cstate")
		if err != nil {
			log.Fatal(err)
		}
		sCstateSetting := strings.TrimRight(string(c), "\n\r")
		iCstateSetting, err := strconv.Atoi(sCstateSetting)

		if iCstateSetting > 0 {
			sWarning := "The CPU maximum cstate setting is at " + sCstateSetting + "! For best performance it is recommended to configure it to 0."
			fmt.Println("\tCPU maximum cstate setting: ", formatYellow(sWarning))
			mReport["CPU cstate"] = sWarning
			return
		} else {
			fmt.Println("\tCPU maximum cstate setting: " + sCstateSetting + " - " + formatGreen("OK"))
			return
		}
	}
}

func checkTuneD() {
	tunedReturnCode := getCommandReturnCode(strings.Fields("systemctl status tuned"))
	if tunedReturnCode == 4 {
		sWarning := "Not installed! To achieve best performance you should run tuned in the 'latency-performance' profile."
		fmt.Println("\tTuneD service: ", formatYellow(sWarning))
		mReport["TuneD"] = sWarning
		return
	}
	if tunedReturnCode != 0 {
		sWarning := "Service failed! Please check the tuned service."
		fmt.Println("\tTuneD service: ", formatYellow(sWarning))
		mReport["TuneD"] = sWarning
		return
	}
	if tunedReturnCode == 0 {
		tunedProfile, _ := runCommand(strings.Fields("tuned-adm active"))
		{
			if strings.Contains(tunedProfile, "latency-performance") {
				fmt.Println("\tTuneD service: ", formatGreen("OK"))
				return
			} else {
				sWarning := "TuneD is not running. To achieve best performance you should run tuned with the 'latency-performance' profile."
				fmt.Println("\tTuneD service: ", formatYellow(sWarning))
				mReport["TuneD"] = sWarning
				return
			}
		}
	}
}

func checkIRQBalance(){
	irqBalance := getCommandReturnCode(strings.Fields("systemctl status irqbalance"))
	if irqBalance == 4 {
		sWarning := "Not installed! To achieve best performance and user experience you should run IRQ Balance."
		fmt.Println("\tIRQ Balance: ", formatYellow(sWarning))
		mReport["IRQ Balance"] = sWarning
		return
	}
	if irqBalance != 0 {
		sWarning := "IRQBalance service is failed! Please check the IRQ Balance service."
		fmt.Println("\tIRQ Balance: ", formatYellow(sWarning))
		mReport["IRQ Balance"] = sWarning
		return
	}
	if irqBalance == 0{
		fmt.Println("\tIRQBalance: ", formatGreen("OK"))
		return
	}
}