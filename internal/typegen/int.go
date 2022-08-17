package typegen

import (
	"math/rand"
	"strconv"
	"time"
)

const (
	PGINTMAX = 2147483647
	PGINTMIN = -2147483648
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func GenInt(nullable bool, ps ...interface{}) string {
	var (
		min = PGINTMIN
		max = PGINTMAX
	)
	if len(ps) == 2 {
		min = ps[0].(int)
		max = ps[1].(int)
	}
	return nullWrap(nullable, strconv.Itoa(randIntInRange(min, max)))
}

func GenBigInt(nullable bool, ps ...interface{}) string {
	return nullWrap(nullable, strconv.FormatInt(rand.Int63(), 10))
}

func GenUniqInt(nullable bool, ps ...interface{}) string {
	return nullWrap(nullable, strconv.Itoa(int(time.Now().UTC().UnixNano()%PGINTMAX)))
}

func GenUniqBigInt(nullable bool, ps ...interface{}) string {
	return nullWrap(nullable, strconv.FormatInt(time.Now().UTC().UnixNano(), 10))
}

func randIntInRange(min, max int) int {
	return rand.Intn(max-min) + min
}
