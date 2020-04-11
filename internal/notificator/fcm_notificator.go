package notificator

import "fmt"

type FCMNotificator struct{}

func (fcmn FCMNotificator) send() error {
	fmt.Println("ewefwef")
	return nil
}
