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
	"fmt"
	"github.com/spf13/pflag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	troubleReport                []string
	filesToArchive               []string
	bPlainOutput                 = false
	bAnswerYes                   = false
	bLustreInstalled             = false
	iInstalledLustrePackages     = 0
	iLoadedLustreModules         = 0
	bLnetLoaded                  = false
	bParseMellanoxLspciOutput    = false
	bLustreLoaded                = false
	cpuList                      string
	bSOSreportInstalled          = false
	bSOSreportRequested          = false
	bHelpRequested               = false
	bCreateClientDiagBundle      = false
	sWorkingDir                  = "/tmp"
	sTempDir                     = ""
	sClientDiagOutputFile        = ""
	bKeepDiagBundle              = false
	sSupportReference            = ""
	bQuietMode                   = false
	sClientDiagBundleOutputPath  = "/tmp"
	sClientDiagBundleArchiveName = ""
)

func init() { log.SetFlags(log.Lshortfile | log.LstdFlags) }

func main() {
	pflag.BoolVarP(&bPlainOutput, "plain-output", "p", false, "Plain output without colors or other formatters")
	pflag.BoolVarP(&bAnswerYes, "yes", "y", false, "Answer yes to all questions.")
	pflag.BoolVarP(&bSOSreportRequested, "sosreport", "s", false, "Create a sosreport and save it to /tmp/sosreport-<hostname>-<date>.tar.xz")
	pflag.BoolVarP(&bHelpRequested, "help", "h", false, "Show this help message")
	pflag.BoolVarP(&bCreateClientDiagBundle, "client-diag-bundle", "c", false, "Create a client diagnostic bundle and save it to /tmp/client-diag-bundle-<hostname>-<date>.tar.xz.\nThis bundle will also include the sosreport if the -s option is used.")
	pflag.StringVarP(&sWorkingDir, "working-dir", "w", "/tmp", "Working directory for sosreport and client diagnostic bundle creation.")
	//pflag.BoolVarP(&bKeepDiagBundle, "keep-diag-bundle", "k", false, "Keep the client-diag diagnostic bundle after the script has finished.\nDefault is to delete it. Only works if the -c option is used.")
	pflag.StringVarP(&sSupportReference, "support-reference", "r", "", "Support reference number or case for the client diagnostic and sosreport bundle.\nThis will be added to the filename of the bundle if provided.")
	pflag.BoolVarP(&bQuietMode, "quiet", "q", false, "Quiet mode. Only print errors and warnings. Only works if the -c option is used.")
	pflag.StringVarP(&sClientDiagBundleOutputPath, "client-diag-bundle-output-path", "o", "/tmp", "Output path for the client diagnostic bundle.\nOnly works if the -c option is used.")
	pflag.Parse()

	if bHelpRequested {
		pflag.Usage()
		os.Exit(0)
	}

	checkUser()

	validateCommandLineArgs()

	sHostname, _ := runCommand(strings.Fields("hostname"))
	sHostname = strings.TrimRight(sHostname, "\n\r")

	if bCreateClientDiagBundle {
		TempDir, err := os.MkdirTemp(sWorkingDir, "client-diag-bundle-*")
		if err != nil {
			log.Println("Error creating temporary directory for client diagnostic bundle creation: ", err.Error())
			os.Exit(1)
		} else {
			sTempDir = TempDir
			fmt.Println("Temporary directory for client diagnostic bundle creation: ", sTempDir)
		}

		defer func(path string) {
			err := os.RemoveAll(path)
			if err != nil {
				log.Println("Error deleting temporary directory for client diagnostic bundle creation:", err.Error())
			}
		}(sTempDir)
		sClientDiagOutputFile = filepath.Join(sTempDir, "/client-diag-output-"+sHostname+".txt")
	}

	writeOutLn(formatBoldWhite("Starting System Scan. Please wait..."))

	if checkExecutableExists("lshw") {
		strLshwOutput, _ := runCommand(strings.Fields("lshw -c system -quiet"))
		parseLshw(strLshwOutput)
	} else {
		writeOutLn(formatBoldWhite("Server/Hostname: " + sHostname))
	}

	if rootUser() {
		strDmiDecodeOutput, _ := runCommand(strings.Fields("dmidecode --no-sysfs -t baseboard -q"))
		parseDmiDecode(strDmiDecodeOutput)
	}
	if checkExecutableExists("free") {
		strFree, _ := runCommand(strings.Fields("free --si -h"))
		slcInstMem := strings.Split(strFree, "\n")
		writeOutLn(formatBoldWhite("Installed Memory: " + strings.Fields(slcInstMem[1])[1]))
	}
	writeOutLn(formatBoldWhite("\nOperating System:"))
	var strOSinfo string
	if checkExecutableExists("lsb_release") {
		strOSinfo, _ = runCommand(strings.Fields("lsb_release -a"))
	} else if checkExecutableExists("hostnamectl") {
		strOSinfo, _ = runCommand(strings.Fields("hostnamectl"))
	} else {
		strOSinfo, _ = runCommand(strings.Fields("cat /etc/os-release"))
	}
	strKernel, _ := runCommand(strings.Fields("uname -r"))
	parseOSinfo(strOSinfo, strKernel)

	writeOutLn(formatBoldWhite("\nCPU Information:"))
	strLsCpuOutput, _ := runCommand(strings.Fields("lscpu"))
	parseLsCpu(strLsCpuOutput)

	checkSystemTuning()

	if !checkFirewall() == false {
		sWarning := "Firewall is running! Make sure that all necessary TCP/IP ports for your environment are open."
		writeOutLn(formatBoldWhite("\nFirewall: " + formatYellow(sWarning)))
		troubleReport = append(troubleReport, "Firewall: "+sWarning)

	} else {
		writeOutLn(formatBoldWhite("\nFirewall:"), "No firewall found.")
	}

	writeOutLn(formatBoldWhite("\nInstalled Lustre packages:"))
	parseLustrePackages()

	if bLustreInstalled {
		writeOutLn(formatBoldWhite("\nLoaded Lustre Kernel modules:"))
		parseLoadedLustreKernelModules()

		writeOutLn(formatBoldWhite("\nLustre Kernel module configuration (\"/etc/modprobe.d/lustre.conf\"):"))
		parseLustreKernelModuleConfig()

		writeOutLn(formatBoldWhite("\nLustre Filesystem Targets:"))
		for _, line := range strings.Split(runCommand(strings.Fields("lfs check servers"))) {
			writeOutLn("\t", line)
		}

		writeOutLn(formatBoldWhite("\nLustre Filesystem OST and MST Information:"))
		if bLustreLoaded {
			parseLfsDf()
		} else {
			sWarning := "The Lustre kernel module is not loaded, cannot get OST and MST usage details."
			writeOutLn(formatYellow("\tWarning: " + sWarning))
			troubleReport = append(troubleReport, "Lustre OST and MST info: "+sWarning)
		}
		if rootUser() {
			writeOutLn(formatBoldWhite("\nLustre LNET Information:"))
			if bLnetLoaded {
				parseLnetInfo()
			} else {
				sWarning := "The lustre LNET kernel module is not loaded, cannot get LNET details.."
				writeOutLn(formatYellow("\tWarning: " + sWarning))
				troubleReport = append(troubleReport, "LNET info: "+sWarning)
			}
		}
		writeOutLn(formatBoldWhite("\nLustre Filesystem Client Mount Information:"))

		strMountOutput, _ := runCommand(strings.Fields("mount -t lustre -l"))
		slcMountOutput := strings.Split(strMountOutput, "\n")
		if len(slcMountOutput) > 1 {
			for _, line := range slcMountOutput {
				writeOutLn("\t", line)
			}
			writeOutLn(formatBoldWhite("\nClient lustre filesystem capacity information:"))

			strDfOutput, err := runCommand(strings.Fields("df -t lustre -H"))

			if len(err) > 0 {
				sWarning := "Cannot find any active lustre filesystem. Are all lustre resources mounted?"
				writeOutLn(formatYellow("\tWarning: " + sWarning))
				troubleReport = append(troubleReport, "Lustre filesystem: "+sWarning)
			} else {
				var slcDfOutput = strings.Split(strDfOutput, "\n")
				if len(slcDfOutput) > 0 {
					for _, line := range slcDfOutput {
						writeOutLn("\t", line)
					}
				} else {
					sWarning := "Cannot find any active lustre filesystem. Are all lustre resources mounted?"
					writeOutLn(formatYellow("\tWarning: " + sWarning))
					troubleReport = append(troubleReport, "Lustre filesystem: "+sWarning)
				}
			}

			writeOutLn(formatBoldWhite("\nLustre Filesystem Client Tuning Information:"))
			parseLustreFilesystemTuning()
		} else {
			sWarning := "Cannot find any active lustre device mounts. Are all lustre resources mounted?"
			writeOutLn(formatYellow("\tWarning" + sWarning))
			troubleReport = append(troubleReport, "Lustre devices: "+sWarning)
		}

	}
	writeOutLn(formatBoldWhite("\nInfiniband/Mellanox Device Information:"))

	if checkExecutableExists("ofed_info") {
		ofedVersion, _ := runCommand(strings.Fields("ofed_info -n"))
		writeOutLn(formatBoldWhite("\nMellanox OFED:"), strings.TrimRight(ofedVersion, "\n\r"))
	} else {
		writeOutLn(formatBoldWhite("\nMellanox OFED:"), "No OFED found.")
	}

	if checkExecutableExists("ibdev2netdev") {
		writeOutLn(formatBoldWhite("\nInfiniband device information (\"ibdev2netdev\" output):"))
		sIBDev, _ := runCommand(strings.Fields("ibdev2netdev -v"))
		slcIBNetdev := strings.Split(sIBDev, "\n")
		for _, line := range slcIBNetdev {
			writeOutLn("\t", line)
		}
	}

	if checkExecutableExists("ibv_devinfo") {
		parseIBDEVInfo()
	} else {
		bParseMellanoxLspciOutput = true
		writeOutLn("\nCannot find \"ibv_devinfo\" in $PATH. Will parse \"lspci -vvv\" output for Mellanox HCA information instead.")
	}

	if checkExecutableExists("ibnetdiscover") {
		writeOutLn(formatBoldWhite("\nInfiniband fabric peers information (\"ibnetdiscover\" output):"))
		sIBNetDiscover, _ := runCommand(strings.Fields("ibnetdiscover -H"))
		slcIBNetDiscover := strings.Split(sIBNetDiscover, "\n")
		for _, line := range slcIBNetDiscover {
			writeOutLn("\t", line)
		}
	} else {
		writeOutLn("\nCannot find \"ibnetdiscover\" in $PATH and therefore no IB fabric peer information will be available.")
	}

	strLspciOutput, _ := runCommand(strings.Fields("lspci -vvv"))
	parseLSPCI(strLspciOutput)

	writeOutLn(formatBoldWhite("\nIP Network Interface Information:"))

	for _, ipCommand := range []string{"ip -4 a s", "ip -s link"} {

		strIpOutput, _ := runCommand(strings.Fields(ipCommand))
		slcIpOutput := strings.Split(strIpOutput, "\n")
		for _, line := range slcIpOutput {
			writeOutLn("\t", line)
		}
	}

	if checkExecutableExists("netstat") {
		writeOutLn(formatBoldWhite("\nNetwork statistics (\"netstat -t\" output):"))
		strNetstatOutput, _ := runCommand(strings.Fields("netstat -t"))
		slcNetstatOutput := strings.Split(strNetstatOutput, "\n")
		for _, line := range slcNetstatOutput {
			writeOutLn("\t", line)
		}
	}

	writeOutLn(formatBoldWhite("\nTrouble Summary:"))
	if len(troubleReport) < 1 {
		writeOutLn(formatGreen("\tNo troubles found."))
	} else {
		for _, line := range troubleReport {
			writeOutLn(formatYellow("\t" + line))
		}
	}

	if bCreateClientDiagBundle {

		if bSOSreportRequested {
			writeOutLn(formatBoldWhite("\nSOS Report Information:"))
			strSosreportOutput, _ := runSosreport()
			slcSosreportOutput := strings.Split(strSosreportOutput, "\n")
			for _, line := range slcSosreportOutput {
				writeOutLn("\t", line)
			}
		}

		if len(sSupportReference) > 0 {
			sClientDiagBundleArchiveName = sSupportReference + "_client_diag_bundel_" + time.Now().Format("2006-01-02_15-04-05") + ".tar.gz"
		} else {
			sClientDiagBundleArchiveName = "client_diag_bundel_" + time.Now().Format("2006-01-02_15-04-05") + ".tar.gz"
		}
		writeOutLn(formatBoldWhite("\nCreating client diagnostic bundle archive: " + sClientDiagBundleOutputPath + "/" + sClientDiagBundleArchiveName))

		fileList, err := ioutil.ReadDir(sTempDir)
		if err != nil {
			log.Println(formatRed("Error reading temporary directory: " + sTempDir + " - " + err.Error()))
		}

		for _, file := range fileList {
			if !file.IsDir() {
				filesToArchive = append(filesToArchive, sTempDir+"/"+file.Name())
			}

			err = CreateTarball(sClientDiagBundleOutputPath+"/"+sClientDiagBundleArchiveName, filesToArchive)
			if err != nil {
				log.Println(formatRed("Error creating client diagnostic bundle archive: " + sClientDiagBundleArchiveName + " - " + err.Error()))
			}
		}
	}
}
