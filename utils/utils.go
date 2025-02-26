package utils

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/gen2brain/beeep"
	"github.com/olebedev/when"
	"github.com/olebedev/when/rules/common"
	"github.com/olebedev/when/rules/en"
)

type Reminder struct {
	Time    time.Time `json:"time"`
	Message string    `json:"message"`
}

const (
	markName    = "GOLANG_CLI_REMINDER"
	markValue   = "1"
	storageFile = "storage/reminders.json"
)

// Load existing reminders from JSON file
func LoadReminders(filename string) []Reminder {
	var reminders []Reminder
	file, err := os.ReadFile(filename)
	if err == nil {
		json.Unmarshal(file, &reminders)
	}
	return reminders
}

// Save reminders to JSON file
func SaveReminders(filename string, reminders []Reminder) {
	data, _ := json.Marshal(reminders)
	os.WriteFile(filename, data, 0644)
}

// Add a new reminder
func AddReminder(filename string) {
	var timeInput, messageInput string

	fmt.Print("Enter time (e.g., '5 minutes', '3pm', '15:30'): ")
	fmt.Scanln(&timeInput)
	fmt.Print("Enter reminder message: ")
	reader := bufio.NewReader(os.Stdin)
	messageInput, _ = reader.ReadString('\n')
	messageInput = strings.TrimSpace(messageInput)

	now := time.Now()
	w := when.New(nil)
	w.Add(en.All...)
	w.Add(common.All...)

	t, err := w.Parse(timeInput, now)
	if err != nil || t == nil {
		fmt.Println("Invalid time format. Try '5 minutes', 'tomorrow 9am', etc.")
		return
	}

	if now.After(t.Time) {
		fmt.Println("❌ Please enter a future time!")
		return
	}

	diff := t.Time.Sub(now)

	reminders := LoadReminders(filename)
	newReminder := Reminder{Time: t.Time, Message: messageInput}
	reminders = append(reminders, newReminder)

	SaveReminders(filename, reminders)

	cmd := exec.Command(os.Args[0])
	cmd.Env = append(os.Environ(), fmt.Sprintf("%s=%s", markName, markValue))
	cmd.Start()

	fmt.Printf("✅ Reminder set for %s (%s from now) - \"%s\"\n",
		t.Time.Format("Mon Jan 2 15:04:05"), diff.Round(time.Second), messageInput)
}

// Show reminders
func ShowReminders(filename string) {
	reminders := LoadReminders(filename)

	if len(reminders) == 0 {
		fmt.Println("📭 No reminders set.")
		return
	}

	fmt.Println("📅 Upcoming Reminders:")
	for i, r := range reminders {
		fmt.Printf("[%d] %s - \"%s\"\n", i+1, r.Time.Format("Mon Jan 2 15:04:05"), r.Message)
	}
}

// Run the reminder notification
func RunReminder() {
	reminders := LoadReminders(storageFile)
	if len(reminders) == 0 {
		return
	}

	// Sort reminders by time
	now := time.Now()
	for _, r := range reminders {
		if now.After(r.Time) {
			continue
		}
		diff := r.Time.Sub(now)
		time.Sleep(diff)

		err := beeep.Alert("🔔 Reminder", r.Message, "assets/information.png")
		if err != nil {
			fmt.Println("Error displaying notification:", err)
		}
	}

	// Remove past reminders
	newReminders := []Reminder{}
	for _, r := range reminders {
		if now.Before(r.Time) {
			newReminders = append(newReminders, r)
		}
	}

	SaveReminders(storageFile, newReminders)
}
