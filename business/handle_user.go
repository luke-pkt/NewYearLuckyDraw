package business

import (
	"NewYearLuckyDraw/libray/storage"
	"NewYearLuckyDraw/model"
	"errors"
	"fmt"
)

func (self *LuckyDraw) ExistUser(uid int) (bool, error) {
	ids, err := storage.Database.IDs(self.Key(model.UserTable))
	if err != nil {
		return false, err
	}
	for _, id := range ids { if id == uid { return true, nil}}
	return false, nil
}

func (self *LuckyDraw) ExistUserName(name string) (bool, error) {
	users, err := self.QueryUsers(true)
	if err != nil {
		return false, err
	}
	for _, user := range users {
		if user.Name == name {
			return true, nil
		}
	}
	return false, nil
}

func (self *LuckyDraw)QueryByUserName(name string) (model.User, error) {
	users, err := self.QueryUsers(true)
	if err != nil {
		return model.User{}, err
	}
	for _, user := range users {
		if user.Name == name {
			return user, nil
		}
	}

	return model.User{}, errors.New("用户不存在")
}

func (self *LuckyDraw) CountUser() (int, error){
	ids, err := storage.Database.IDs(self.Key(model.UserTable))
	if err != nil {
		return 0, err
	}
	return len(ids), nil
}

func (self *LuckyDraw) QueryUsers(isAll bool) ([]model.User, error) {
	var Users []model.User
	ids, err := storage.Database.IDs(self.Key(model.UserTable))
	if err != nil {
		return Users, err
	}
	for _, id := range ids {

		if isAll {
			var user model.User
			if err := storage.Database.Find(self.Key(model.UserTable), id, &user); err != nil {
				return Users, err
			}
			Users = append(Users, user)

		} else {
			if exist, err := self.ExistUserRecord(id); err != nil {
				return Users, err

			} else if err == nil && !exist {
				var user model.User
				if err := storage.Database.Find(self.Key(model.UserTable), id, &user); err != nil {
					return Users, err
				}
				Users = append(Users, user)
			}
		}
	}

	return Users, err
}

func (self *LuckyDraw) AppendUser(name string) error {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	exist, err := self.ExistUserName(name)
	if err != nil {
		return err
	}
	if exist {
		return errors.New(fmt.Sprintf("用户已存在"))
	}
	if _, err := storage.Database.Insert(self.Key(model.UserTable), &model.User{Name: name}); err != nil {
		return err
	}
	return nil
}

func (self *LuckyDraw) DeleteUser(id int) error {
	if err := storage.Database.Delete(self.Key(model.UserTable), id); err != nil {
		return err
	}
	return nil
}

