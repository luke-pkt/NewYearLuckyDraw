package business

import (
	"NewYearLuckyDraw/libray/logger"
	"NewYearLuckyDraw/libray/storage"
	"NewYearLuckyDraw/model"
	"errors"
	"fmt"
	"math/rand"
	"sort"
	"sync"
	"time"
)

type LuckyDraw struct {
	mutex *sync.Mutex

	start bool
	round int // 当前抽奖者

	settings model.Settings
}


func NewLuckyDraw(settings model.Settings) LuckyDraw {
	lucky := LuckyDraw{mutex: new(sync.Mutex), settings: settings}

	if !storage.Database.TableExists(lucky.Key(model.UserTable)) {
		storage.Database.CreateTable(lucky.Key(model.UserTable))
	}

	if !storage.Database.TableExists(lucky.Key(model.PrizeTable)) {
		storage.Database.CreateTable(lucky.Key(model.PrizeTable))
	}

	if !storage.Database.TableExists(lucky.Key(model.RecordTable)) {
		storage.Database.CreateTable(lucky.Key(model.RecordTable))
	}

	return lucky
}

func (self *LuckyDraw) Settings() model.Settings {
	return self.settings
}

func (self *LuckyDraw) Key(tb model.TableName) string {
	return fmt.Sprintf("%v_%v", tb, self.settings.ID)
}

func (self *LuckyDraw) Start() {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	self.start = true
}

func (self *LuckyDraw) Started() bool {
	return self.start
}

func (self *LuckyDraw)GetRound() (model.User, error) {
	if self.round != 0 {
		users, err := self.QueryUsers(true)
		if err != nil {
			return model.User{}, err
		}
		for _, user := range users {
			if user.ID == self.round {
				return user, nil
			}
		}
	}
	return model.User{}, nil
}

func (self *LuckyDraw)SetLotteryUserNumber(num int) {
	self.settings.LotteryUserNumber = num
}

func (self *LuckyDraw)SetRound(uid int) error {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	if !self.start {
		return errors.New("抽奖未开始")
	}

	if self.round != 0 {
		return errors.New("已有抽奖者")
	}
	if exist, err := self.ExistUser(uid); err != nil {
		return err

	} else if err == nil && !exist {
		return errors.New("用户不存在")
	}

	exist, err := self.ExistUserRecord(uid)
	if err != nil {
		return err
	}
	if exist {
		return errors.New("用户已抽奖")
	}
	self.round = uid
	return nil
}

func (self *LuckyDraw)Lottery(uid int) (model.Record, error) {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	var record model.Record

	if !self.start {
		return record, errors.New("抽奖未开始")
	}

	if self.round != uid {
		return record, errors.New("未轮到用户")
	}

	if exist, err := self.ExistUser(uid); err != nil {
		return record, err

	} else if err == nil && !exist {
		return record, errors.New("用户不存在")
	}

	if exist, err := self.ExistUserRecord(uid); err != nil {
		return record, err

	} else if err == nil && exist {
		return record, errors.New("用户已抽奖")
	}

	prizes, err := self.QueryPrizes(false)
	if err != nil {
		return record, err
	}

	i := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(len(prizes))
	record.UID = uid
	record.PID = prizes[i].ID

	rid, err := self.AppendRecord(uid, prizes[i].ID)
	if err != nil {
		return record, err
	}

	if err := storage.Database.Find(self.Key(model.RecordTable), rid, &record); err != nil {
		return record, err
	}

	self.round = 0
	
	return record, nil
}

func (self *LuckyDraw)Report() (string, error) {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	var report string
	records, err := self.QueryRecords()
	if err != nil {
		return report, err
	}

	prizes, err := self.QueryPrizes(true)
	if err != nil {
		return "", err
	}
	sort.Sort(prizes)

	for _, prize := range prizes {

		var userName string
		for _, record := range records {
			if record.PID == prize.ID {
				var user model.User
				if err := storage.Database.Find(self.Key(model.UserTable), record.UID, &user); err != nil {
					logger.ErrorF("[Report]Index(%v) find uid(%v) error: %v", self.settings.ID, record.UID, err)
				}
				userName = user.Name
			}
		}
		report += fmt.Sprintf("%v\t->\t%v\n", prize.Title, userName)
	}
	return report, nil
}