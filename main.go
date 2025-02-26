package main

import (
	"fmt"
	"os"
	"time"

	"github.com/ankon07/gui-cli-reminder/utils" // Corrected import path
)

const (
	markName    = "GOLANG_CLI_REMINDER"
	markValue   = "1"
	storageFile = "storage/reminders.json"
)

type Reminder struct {
	Time    time.Time `json:"time"`
	Message string    `json:"message"`
}

func main() {
	if os.Getenv(markName) == markValue {
		// Run reminder notification
		utils.RunReminder()
		return
	}

	// Load existing reminders
	reminders := utils.LoadReminders(storageFile)

	fmt.Println("===== CLI Reminder System =====")
	fmt.Println("1. Add Reminder")
	fmt.Println("2. Show Reminders")
	fmt.Println("3. Exit")

	var choice int
	fmt.Print("Select an option: ")
	fmt.Scan(&choice)

	switch choice {
	case 1:
		if len(reminders) >= 5 {
			fmt.Println("❌ You can only set up to 5 reminders at a time!")
			return
		}
		utils.AddReminder(storageFile)
	case 2:
		utils.ShowReminders(storageFile)
	case 3:
		fmt.Println("Exiting...")
		return
	default:
		fmt.Println("Invalid option!")
	}
}
