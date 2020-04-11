package notificator

import (
	"fmt"
	"os"
)

type Notificator interface {
	send() error
}

type DefaultNotificator struct{}

func NewNotificator(to string) Notificator {
	switch to {
	case "fcm":
		return FCMNotificator{
			endpoint: "https://fcm.googleapis.com/fcm/send",
			to:       os.Getenv("FCM_TO"),
			authKey:  os.Getenv("FCM_AUTH_KEY"),
		}
	default:
		return DefaultNotificator{}
	}
}

func Run(n Notificator) {
	if err := n.send(); err != nil {
		panic(err)
	}
}

func (dn DefaultNotificator) send() error {
	fmt.Println("not implemented")
	return nil
}
