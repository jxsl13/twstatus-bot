package model

import "github.com/diamondburned/arikawa/v3/discord"

type UserTarget struct {
	Target
	UserID discord.UserID
}

type PlayerCountNotification struct {
	UserTarget
	Threshold int
}

func (p PlayerCountNotification) Notify(change *ChangedServerStatus) bool {
	return len(change.Curr.Clients) >= p.Threshold
}

type PlayerCountNotifications []PlayerCountNotification

type ByPlayerCountNotificationIDs []PlayerCountNotification

func (a ByPlayerCountNotificationIDs) Len() int      { return len(a) }
func (a ByPlayerCountNotificationIDs) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByPlayerCountNotificationIDs) Less(i, j int) bool {
	return a[i].Target.Less(a[j].Target)
}
