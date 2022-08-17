package typegen

import (
	"github.com/vl4deee11/pg-gen/internal/util"
)

const bQuote byte = 39

func GenString(nullable bool, ps ...interface{}) string {
	b := make([]byte, 102)
	b[0] = bQuote
	for i := 1; i < 102; i++ {
		b[i] = byte(randIntInRange(65, 122))
	}
	b[len(b)-1] = bQuote
	return nullWrap(nullable, util.B2S(b))
}

func GenUniqString(nullable bool, ps ...interface{}) string {
	x := GenUniqInt(false)
	b := make([]byte, 102+len(x))
	b[0] = bQuote
	for i := 1; i < 102; i++ {
		b[i] = byte(randIntInRange(65, 122))
	}

	end := len(b) - 1
	for i := 101; i < end; i++ {
		b[i] = x[i-101]
	}
	b[len(b)-1] = bQuote
	return nullWrap(nullable, util.B2S(b))
}

func GenStringN(nullable bool, ps ...interface{}) string {
	b := make([]byte, ps[0].(int)+2)
	b[0] = bQuote
	for i := 1; i < ps[0].(int)+2; i++ {
		b[i] = byte(randIntInRange(65, 122))
	}
	b[len(b)-1] = bQuote
	return nullWrap(nullable, util.B2S(b))
}
