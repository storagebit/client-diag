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

func formatBoldWhite(s string) string {
	if bPlainOutput {
		return s
	} else {
		out := "\033[1;37m" + s + "\033[0m"
		return out
	}
}

func formatYellow(s string) string {
	if bPlainOutput {
		return s
	} else {
		out := "\033[0;33m" + s + "\033[0m"
		return out
	}
}

func formatGreen(s string) string {
	if bPlainOutput {
		return s
	} else {
		out := "\033[0;32m" + s + "\033[0m"
		return out
	}
}

func formatRed(s string) string {
	if bPlainOutput {
		return s
	} else {
		out := "\033[0;31m" + s + "\033[0m"
		return out
	}
}

func statusFormat(returnCode int) string {
	var status string

	if returnCode == 0 {
		if bPlainOutput {
			status = "OK"
		} else {
			status = formatGreen("OK")
		}
		return status
	} else {
		if bPlainOutput {
			status = "WARNING"
		} else {
			status = formatYellow("WARNING")
		}
		return status
	}
}
