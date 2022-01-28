package business

import (
	"NewYearLuckyDraw/libray/storage"
	"NewYearLuckyDraw/model"
)

func (self *LuckyDraw) ExistPrize(pid int) (bool, error) {
	ids, err := storage.Database.IDs(self.Key(model.PrizeTable))
	if err != nil {
		return false, err
	}
	for _, id := range ids {
		if id == pid {
			return true, nil
		}
	}
	return false, nil
}

func (self *LuckyDraw) QueryPrizeByID(pid int) (model.Prize, error) {
	var prize model.Prize
	if err := storage.Database.Find(self.Key(model.PrizeTable), pid, &prize); err != nil {
		return prize, err
	}
	return prize, nil
}

func (self *LuckyDraw) QueryPrizes(isAll bool) (model.PrizeList, error) {
	var prizes model.PrizeList
	ids, err := storage.Database.IDs(self.Key(model.PrizeTable))
	if err != nil {
		return prizes, err
	}
	for _, id := range ids {
		if isAll {
			var prize model.Prize
			if err := storage.Database.Find(self.Key(model.PrizeTable), id, &prize); err != nil {
				return prizes, err
			}
			prizes = append(prizes, prize)

		} else {
			if exist, err := self.ExistPrizeRecord(id); err != nil {
				return prizes, err

			} else if err == nil && !exist {
				var prize model.Prize
				if err := storage.Database.Find(self.Key(model.PrizeTable), id, &prize); err != nil {
					return prizes, err
				}
				prizes = append(prizes, prize)
			}
		}
	}

	return prizes, err
}

func (self *LuckyDraw) AppendPrize(title, cover string, count int) error {
	if count < 1 {
		count = 1
	}

	for i := 0; i < count; i++ {
		if _, err := storage.Database.Insert(self.Key(model.PrizeTable), &model.Prize{Title: title, Cover: cover}); err != nil {
			return err
		}
	}
	return nil
}

func (self *LuckyDraw) DeletePrize(id int) error {
	if err := storage.Database.Delete(self.Key(model.PrizeTable), id); err != nil {
		return err
	}

	return nil
}
