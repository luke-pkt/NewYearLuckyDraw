package business

import (
	"NewYearLuckyDraw/libray/storage"
	"NewYearLuckyDraw/model"
	"time"
)

func (self *LuckyDraw)ExistUserRecord(uid int) (bool, error) {
	ids, err := storage.Database.IDs(self.Key(model.RecordTable))
	if err != nil {
		return false, err
	}

	for _, id := range ids {
		var record model.Record
		if err := storage.Database.Find(self.Key(model.RecordTable), id, &record); err != nil {
			return false, err
		}
		if record.UID == uid {
			return true, nil
		}
	}
	return false, nil
}

func (self *LuckyDraw)ExistPrizeRecord(pid int) (bool, error) {
	ids, err := storage.Database.IDs(self.Key(model.RecordTable))
	if err != nil {
		return false, err
	}

	for _, id := range ids {
		var record model.Record
		if err := storage.Database.Find(self.Key(model.RecordTable), id, &record); err != nil {
			return false, err
		}
		if record.PID == pid {
			return true, nil
		}
	}
	return false, nil
}

func (self *LuckyDraw)AppendRecord(uid, pid int) (int, error) {
	id, err := storage.Database.Insert(self.Key(model.RecordTable), &model.Record{UID: uid, PID: pid, Timestamp: time.Now().Unix()})
	if err != nil {
		return 0, err
	}
	return id, nil
}


func (self *LuckyDraw) QueryRecords() ([]model.Record, error) {
	var records []model.Record
	ids, err := storage.Database.IDs(self.Key(model.RecordTable))
	if err != nil {
		return records, err
	}
	for _, id := range ids {
		var record model.Record
		if err := storage.Database.Find(self.Key(model.RecordTable), id, &record); err != nil {
			return records, err
		}
		records = append(records, record)
	}

	return records, err
}