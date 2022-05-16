//go:build !containers && !disable_notification

package app

import "github.com/adamwoolhether/cliApps/distributing/notify"

func send_notification(msg string) {
	n := notify.New("Pomodoro", msg, notify.SeverityNormal)

	n.Send()
}
