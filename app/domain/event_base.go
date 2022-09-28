package domain

import "time"

type EventBase struct {
	Id       string    `gorm:"column:id;type:char(24);not null;primarykey" json:"id"`
	SourceId string    `gorm:"column:source_id;type:char(30);not null" json:"source_id"`
	Action   string    `gorm:"column:action;type:text;not null" json:"action"`
	Time     time.Time `gorm:"column:time;type:datetime;not null" json:"time"`
}

func (this *EventBase) TableName() string {
	return "event"
}
