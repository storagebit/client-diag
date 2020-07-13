package main

import (
	"fmt"
	"log"
	"os/user"
)

func rootUser() bool {
	currentUser, userError := user.Current()

	if userError != nil {
		log.Fatalf("Error while retrieving current user details: '%s'", userError)
	}

	if currentUser.Uid == "0" {
		return true
	} else {
		return false
	}
}

func checkUser() {

	currentUser, err := user.Current()

	if err != nil {
		log.Fatalf("Error while retrieving current user details: '%s'", err)
	}

	fmt.Println("Running client-diag as user: " + currentUser.Username)

	if !rootUser() {
		fmt.Println("Executing client-diag without root privileges or sudo will limit the dianostic/reporting capabilities.\n" +
			"Run as root or sudo if you want to see more.")
	} else {

		fmt.Println("client-diag is being executed with elevated/root privileges.")
	}

}
