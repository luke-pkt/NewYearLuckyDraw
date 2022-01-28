package model

import "github.com/jameycribbs/hare"

type User struct {
	ID int `json:"id"`
	Name string `json:"name"`
}

func (self *User)GetID()int{
	return self.ID
}

func (self *User)SetID(id int) {
	self.ID = id
}

func (self *User) AfterFind(db *hare.Database) error {
	*self = User(*self)
	return nil
}