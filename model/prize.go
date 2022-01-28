package model

import (
	"github.com/jameycribbs/hare"
)

type PrizeList []Prize
func (self PrizeList)Len() int {return len(self)}
func (self PrizeList)Less(i, j int) bool {if self[i].ID < self[j].ID {return true}
	return false
}
func (self PrizeList)Swap(i, j int) {
	self[i], self[j] = self[j], self[i]
}

type Prize struct {
	ID int `json:"id"`
	Title string `json:"title"`
	Cover string `json:"cover"`
}

func (self *Prize)GetID()int{
	return self.ID
}

func (self *Prize)SetID(id int) {
	self.ID = id
}

func (self *Prize) AfterFind(db *hare.Database) error {
	*self = Prize(*self)
	return nil
}
