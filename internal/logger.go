package internal

import (
	"fmt"
	"github.com/sqweek/dialog"
	"log"
	"os"
	"time"
)

func Set_log_file() *os.File {

	log_file_location, log_file_selection_err :=
		dialog.Directory().Title("Select log file storage location").Browse()

	fmt.Println(log_file_selection_err)

	log_file_name := log_file_location + "\\" + time.Now().Format("2006-01-02_15_04_05") + "_info.log"

	file, err := os.OpenFile(log_file_name, os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	//defer file.Close()
	log.SetOutput(file)

	return file

}
