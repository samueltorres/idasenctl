package notification

import (
	"github.com/gen2brain/beeep"
)

type Notifier struct{}

func NewNotifier() *Notifier {
	return &Notifier{}
}

func (n *Notifier) SendNotification(title, message string) error {
	return beeep.Notify(title, message, "")
}
