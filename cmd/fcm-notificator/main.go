package main

import (
	"github.com/ddddddO/tag-mng/notificator"
)

func main() {
	fcmNotificator := notificator.NewNotificator("fcm")
	notificator.Run(fcmNotificator)
}
