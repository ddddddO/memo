package notificator

import (
	"fmt"
)

type Notificator interface {
	send() error
}

type DefaultNotificator struct{}

func NewNotificator(to string) Notificator {
	switch to {
	case "fcm":
		return FCMNotificator{}
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
