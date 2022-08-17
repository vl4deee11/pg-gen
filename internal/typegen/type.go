package typegen

import (
	"math/rand"
)

type PGType string

const (
	IntT        PGType = "int"
	BigIntT     PGType = "bigint"
	UniqIntT    PGType = "uniq_int"
	UniqBigIntT PGType = "uniq_bigint"
	TextT       PGType = "text"
	UniqTextT   PGType = "uniq_text"
	PK          PGType = "pk"
	UniqPK      PGType = "uniq_pk"
	NTextT      PGType = "n_text"
	TimestampTZ PGType = "timestamptz"
	TimeTZ      PGType = "timetz"
)

var GenMapInsertT = map[PGType]func(nullable bool, ps ...interface{}) string{
	UniqIntT:    GenUniqInt,
	UniqBigIntT: GenUniqBigInt,
	IntT:        GenInt,
	BigIntT:     GenBigInt,
	TextT:       GenString,
	UniqTextT:   GenUniqString,
	PK:          GenPK,
	NTextT:      GenStringN,
	TimestampTZ: GenTimestampTZ,
	TimeTZ:      GenTimeTZ,
}

var GenMapMemoryT = map[PGType]func(nullable bool, ps ...interface{}) string{
	UniqIntT:    GenUniqInt,
	UniqBigIntT: GenUniqBigInt,
	IntT:        GenInt,
	BigIntT:     GenBigInt,
	TextT:       GenString,
	UniqTextT:   GenUniqString,
	PK:          GenPKInMemory,
	UniqPK:      GenUniqPKInMemory,
	NTextT:      GenStringN,
	TimestampTZ: GenTimestampTZ,
	TimeTZ:      GenTimeTZ,
}

func nullWrap(nullable bool, v string) string {
	if !nullable {
		return v
	}

	if rand.Intn(2) >= 1 {
		return v
	}
	return "null"
}
