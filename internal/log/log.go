package log

var Logger *LVLWrap

func MakeLogger(_lvl string) {
	Logger = newLvlWarp(_lvl)
}
