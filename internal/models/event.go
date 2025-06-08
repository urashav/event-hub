package models

import "time"

type Event struct {
	BaseModel
	ID       int
	Name     string
	Date     time.Time
	Address  string
	Status   string
	Capacity int
}

type UserEvent struct {
	BaseModel
	User_id  int
	Event_id int
}
