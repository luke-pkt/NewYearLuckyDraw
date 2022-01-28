package model

type TableName string

func (self TableName)String() string {
	return string(self)
}

const (
	UserTable TableName = "UserTable"
	PrizeTable TableName = "PrizeTable"
	RecordTable TableName = "RecordTable"
)

type Settings struct {
	ID string `json:"id"` //
	LotteryUserNumber int `json:"lottery_user_number"` // 本次抽奖人数
}

type Prizes []struct{
	Title string `json:"title"`
	CoverUrl string `json:"cover_url"`
}