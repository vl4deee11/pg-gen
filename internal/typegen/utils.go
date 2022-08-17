package typegen

import "math/rand"

func randomSelectFromSlice(s []interface{}) interface{} {
	return s[rand.Intn(len(s))]
}
