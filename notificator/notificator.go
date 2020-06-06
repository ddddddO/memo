package notificator

import (
	"fmt"
	"os"

	"github.com/ddddddO/tag-mng/notificator/fcm"
)

type Notificator interface {
	detect() error
	notify() error
}

type DefaultNotificator struct{}

func NewNotificator(to string) Notificator {
	switch to {
	case "fcm":
		return fcm.FCMNotificator{
			endpoint: "https://fcm.googleapis.com/fcm/send",
			token:    os.Getenv("FCM_TOKEN"),
			authKey:  os.Getenv("FCM_AUTH_KEY"),
			dsn:      os.Getenv("DBDSN"),
		}
	default:
		return DefaultNotificator{}
	}
}

func Run(n Notificator) {
	if err := n.detect(); err != nil {
		panic(err)
	}

	if err := n.notify(); err != nil {
		panic(err)
	}
}

func (dn DefaultNotificator) detect() error {
	fmt.Println("not implemented")
	return nil
}

func (dn DefaultNotificator) notify() error {
	fmt.Println("not implemented")
	return nil
}
