package typegen

import (
	"fmt"
	"math/rand"
)

func GenPK(nullable bool, ps ...interface{}) string {
	return nullWrap(nullable, fmt.Sprintf("(SELECT %v FROM %v WHERE random() > 0.01 LIMIT 1)", ps[1], ps[0]))
}

func GenPKInMemory(nullable bool, ps ...interface{}) string {
	rows := ps[0].([][]string)
	fIdx := ps[1].(int)
	val := rows[rand.Intn(len(rows))][fIdx]
	return nullWrap(nullable, val)
}

var tableIdx2Offset = map[int]map[int]int{}

func GenUniqPKInMemory(nullable bool, ps ...interface{}) string {
	rows := ps[0].([][]string)
	fIdx := ps[1].(int)
	srcTIdx := ps[2].(int)
	trgTIdx := ps[3].(int)

	if nullable && rand.Intn(2) < 1 {
		return "null"
	}

	offsetMap, ok := tableIdx2Offset[srcTIdx]
	if !ok {
		tableIdx2Offset[srcTIdx] = make(map[int]int)
		offsetMap = tableIdx2Offset[srcTIdx]
	}
	offset := offsetMap[trgTIdx]

	if offset >= len(rows) {
		return "-1"
	}

	val := rows[offset][fIdx]

	offsetMap[trgTIdx]++

	return val
}
