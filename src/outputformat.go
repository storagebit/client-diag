package main

func formatBoldWhite(s string) string {
	if plainOutput == true {
		return s
	} else {
		out := "\033[1;37m" + s + "\033[0m"
		return out
	}
}

func formatYellow(s string) string {
	if plainOutput == true {
		return s
	} else {
		out := "\033[0;33m" + s + "\033[0m"
		return out
	}
}

func formatGreen(s string) string {
	if plainOutput == true {
		return s
	} else {
		out := "\033[0;32m" + s + "\033[0m"
		return out
	}
}

func statusFormat(returnCode int) string {
	var status string

	if returnCode == 0 {
		status = formatGreen("OK")
		return status
	} else {
		status = formatYellow("WARNING")
		return status
	}
}
