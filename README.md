# client-diag
A tool collecting hardware and software information allowing you to take a quick inventory or
to start troubleshooting a lustre filesystem client.

The tool collects and analyses information of the following components:
* Platform hardware
* Operating system
* CPU config and settings
* Memory
* Operating system and hardware tuning
* Network interfaces and configuration
* Mellanox specific information and configuration, OFED and hardware 
* Installed lustre packages
* Lustre kernel modules and module configuration
* Lustre LNET information
* Lustre filesystem information
* Lustre device mount information
* Lustre filesystem tuning information
* Collects sos report and creates a client diagnostic bundle

## Installation
Quite simple actually. 
Just download the binary from here and run it on your Lustre clients.

Or `git clone https://github.com/storagebit/client-diag/` and `cd` into the `bin` directory where you find the binary or build and compile it from the source in the `src` directory.
The choice is yours.

## How to use it
Also, quite simple.
```
Usage of ./client_diag:
  -c, --client-diag-bundle                      Create a client diagnostic bundle and save it to /tmp/client-diag-bundle-<hostname>-<date>.tar.xz.
                                                This bundle will also include the sosreport if the -s option is used.
  -o, --client-diag-bundle-output-path string   Output path for the client diagnostic bundle.
                                                Only works if the -c option is used. (default "/tmp")
  -h, --help                                    Show this help message
  -p, --plain-output                            Plain output without colors or other formatters
  -q, --quiet                                   Quiet mode. Only print errors and warnings. Only works if the -c option is used.
  -s, --sosreport                               Create a sosreport and save it to /tmp/sosreport-<hostname>-<date>.tar.xz
  -r, --support-reference string                Support reference number or case for the client diagnostic and sosreport bundle.
                                                This will be added to the filename of the bundle if provided.
  -w, --working-dir string                      Working directory for sosreport and client diagnostic bundle creation. (default "/tmp")
  -y, --yes                                     Answer yes to all questions.

```
## Note
I started this project as I was in need of a very simple to use tool which doesn't have 3rd party package or other software dependencies and can be easily distributed(just copy the binary).
The code might look a bit clunky here and there in its first iteration, but it does the job for me, and I'll see where I can improve it, if required.
I'll also add more code documentation as I work on it and time allows.
Please feel free to contribute if you feel the need.
