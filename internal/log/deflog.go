package log

import (
	"fmt"
	"log"
	"os"
	"strings"
	"unsafe"
)

const (
	trace = iota
	debug
	info
	warn
	err
)

var lvl2str = [5]string{"TRACE", "DEBUG", "INFO", "WARN", "ERROR"}
var str2lvl = map[string]int{
	"debug": debug,
	"info":  info,
	"warn":  warn,
	"error": err,
	"trace": trace,
}

type LVLWrap struct {
	*log.Logger
	l int
}

func (l *LVLWrap) canPrint(lvl int) bool {
	return lvl >= l.l
}

func (l *LVLWrap) lvlFormat(lvl int, format string, v ...interface{}) string {
	buf := make([]byte, 0, len(lvl2str[lvl])+len(format)+3)
	buf = append(buf, 91)
	buf = append(buf, lvl2str[lvl]...)
	buf = append(buf, 93, 32)
	buf = append(buf, format...)
	msg := fmt.Sprintf(*(*string)(unsafe.Pointer(&buf)), v...)
	return msg
}

func (l *LVLWrap) Tracef(format string, v ...interface{}) {
	l.printf(trace, format, v...)
}

func (l *LVLWrap) Trace(v ...interface{}) {
	l.println(trace, v...)
}

func (l *LVLWrap) Debugf(format string, v ...interface{}) {
	l.printf(debug, format, v...)
}

func (l *LVLWrap) Debug(v ...interface{}) {
	l.println(debug, v...)
}

func (l *LVLWrap) Infof(format string, v ...interface{}) {
	l.printf(info, format, v...)
}

func (l *LVLWrap) Info(v ...interface{}) {
	l.println(info, v...)
}

func (l *LVLWrap) Warnf(format string, v ...interface{}) {
	l.printf(warn, format, v...)
}

func (l *LVLWrap) Warn(v ...interface{}) {
	l.println(info, v...)
}

func (l *LVLWrap) Errorf(format string, v ...interface{}) {
	l.printf(err, format, v...)
}

func (l *LVLWrap) Error(v ...interface{}) {
	l.println(err, v...)
}

func (l *LVLWrap) printf(lvl int, format string, v ...interface{}) {
	if !l.canPrint(lvl) {
		return
	}
	l.Logger.Println(l.lvlFormat(lvl, format, v...))
}

func (l *LVLWrap) println(lvl int, v ...interface{}) {
	if !l.canPrint(lvl) {
		return
	}
	l.Logger.Println(l.lvlFormat(lvl, "%s", v...))
}

func newLvlWarp(lvl string) *LVLWrap {
	lvl = strings.ToLower(lvl)

	lo := new(LVLWrap)
	lo.Logger = log.New(os.Stdout, "", log.LstdFlags)
	if l, ok := str2lvl[lvl]; ok {
		lo.l = l
	} else {
		lo.l = 2
	}

	return lo
}
