package notify

import (
	"time"

	"github.com/getlantern/golog"
)

var (
	log = golog.LoggerFor("flashlight/notify")
)

// notifications is an internal representation of the Notifier interface.
type notifications struct {
	ui UISender
}

// Notifier is an interface for sending notifications to the user.
type Notifier interface {
	// Notify sends a notification to the user.
	Notify(*Notification)
}

// Notification contains data for notifying the user about something. This
// is directly modeled after Chrome notifications, as detailed at:
// https://developer.chrome.com/apps/notifications
type Notification struct {
	ID                 string
	Type               string
	Title              string
	Message            string
	IconURL            string
	RequireInteraction bool
	IsClickable        bool
	Buttons            []*Button
}

// Button is a button for a notification.
type Button struct {
	Title    string
	IconURL  string
	ClickURL string
}

// UISender is an interface for allowing this class to send thing to the UI.
type UISender interface {
	Send(interface{})
}

// NewNotifications creates a new Notifier that can notify the user about stuff.
func NewNotifications(register func(string) (UISender, error)) (Notifier, error) {
	uiSender, err := register("notification")
	if err != nil {
		log.Errorf("Could not register with UI? %v", err)
		return nil, err
	}
	n := &notifications{ui: uiSender}

	go func() {
		for {
			buttons := []*Button{&Button{Title: "OK", ClickURL: "https://www.getlantern.org"}}
			not := &Notification{
				ID:                 "backendtesting",
				Type:               "basic",
				Title:              "Data Cap",
				Message:            "Please fix this",
				IconURL:            "lantern_logo.png",
				RequireInteraction: true,
				IsClickable:        true,
				Buttons:            buttons,
			}
			n.Notify(not)
			time.Sleep(30 * time.Second)
		}
	}()

	return n, nil
}

// Notify sends a notification to the user.
func (n *notifications) Notify(msg *Notification) {
	go func() {
		n.ui.Send(msg)
	}()
}
