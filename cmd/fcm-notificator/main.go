package main

import (
	"github.com/ddddddO/tag-mng/internal/notificator"
)

func main() {
	fcmNotificator := notificator.NewNotificator("fcm")
	notificator.Run(fcmNotificator)
}
