/*The MIT License (MIT)
Copyright © 2020 StorageBIT.ch
Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated
documentation files (the “Software”), to deal in the Software without restriction, including without limitation the
rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit
persons to whom the Software is furnished to do so, subject to the following conditions:
The above copyright notice and this permission notice shall be included in all copies or substantial portions of the
Software.
THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO
THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.IN NO EVENT SHALL THE AUTHORS
OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR
OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package main

import (
	"strings"
)

func checkFirewallType(d string) int {
	returnCode := getCommandReturnCode(strings.Fields("systemctl status " + d))
	return returnCode
}

func checkFirewall() bool {

	var firewallRunning bool

	//Checking if Firewalld is running on the host/server. A return code of zero tells us that the firewalld is running
	if checkFirewallType("firewalld") == 0 {
		firewallRunning = true
	}
	//Checking if Ubuntu Firewalld is running on the host/server
	if checkFirewallType("ufw") == 0 {
		firewallRunning = true
	}
	//checking if SLES Susefirewall2 is running
	if checkFirewallType("SuSEfirewall2") == 0 {
		firewallRunning = true
	}

	return firewallRunning
}
