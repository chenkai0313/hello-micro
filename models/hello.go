package models

import "time"

type Hello struct {
	Id        int
	Name      string
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (Hello) TableName() string {
	return "hello"
}
