package model

import "github.com/jameycribbs/hare"

type Record struct {
	ID int `json:"id"`
	UID int `json:"uid"`
	PID int `json:"pid"`
	Timestamp int64 `json:"timestamp"`
}

func (self *Record)GetID()int{
	return self.ID
}

func (self *Record)SetID(id int) {
	self.ID = id
}

func (self *Record) AfterFind(db *hare.Database) error {
	*self = Record(*self)
	return nil
}