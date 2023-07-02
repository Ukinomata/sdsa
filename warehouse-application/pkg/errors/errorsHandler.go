package errors

import "log"

func CheckError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func CheckWarning(err error) {
	if err != nil {
		log.Println(err)
	}
}
