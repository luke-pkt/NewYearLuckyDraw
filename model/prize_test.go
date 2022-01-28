package model

import (
	"fmt"
	"sort"
	"testing"
)

func TestPrizeList(t *testing.T) {
	var prizes PrizeList

	prizes = append(prizes, Prize{ID: 2, Title: "2"}, Prize{ID: 1, Title: "1"}, Prize{ID: 3, Title: "3"})
	sort.Sort(prizes)

	if prizes[0].ID != 1 {t.Fail()}
	fmt.Printf("%+v", prizes)
}