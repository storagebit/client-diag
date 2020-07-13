package main

import (
	"strings"
)

func checkFirewallType(d string) int{
	returnCode := getCommandReturnCode(strings.Fields("systemctl status " + d))
	return returnCode
}

func checkFirewall() bool {

	var firewallRunning bool

	//Checking if Firewalld is running on the host/server. A return code of zero tells us that the firewalld is running
	if checkFirewallType("firewalld") == 0{
		firewallRunning = true
	}
	//Checking if Ubuntu Firewalld is running on the host/server
	if checkFirewallType("ufw") == 0{
		firewallRunning = true
		}
	//checking if SLES Susefirewall2 is running
	if checkFirewallType("SuSEfirewall2") == 0{
		firewallRunning = true
	}

	return firewallRunning
}
