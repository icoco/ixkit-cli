package core

import (
	"fmt"
	wlog "github.com/dixonwille/wlog"
	"os"
	"strings"
)

type Tracer struct {
	ui    wlog.UI
	level int // 0 : info | 1 : debug
}

var TraceMe = &Tracer{ui: createUI(), level: 0}

func createUI() wlog.UI {
	var ui wlog.UI
	reader := strings.NewReader("\r\n") //Simulate user typing "User Input" then pressing [enter] when reading from os.Stdin
	ui = wlog.New(reader, os.Stdout, os.Stdout)
	ui = wlog.AddPrefix("?", wlog.Cross, " ", "", "", "~", wlog.Check, "!", ui)
	ui = wlog.AddConcurrent(ui)
	return ui
}

func (t *Tracer) SetDebugMode(d bool) *Tracer {
	if d {
		t.level = 1
	} else {
		t.level = 0
	}
	return t
}

func (t *Tracer) Info(format string, args ...interface{}) {
	s := fmt.Sprintf("[ixkit] "+format, args...)
	t.ui.Info(s)
}

func (t *Tracer) Debug(format string, args ...interface{}) {
	if t.level < 1 {
		return
	}
	s := fmt.Sprintf("[ixkit] "+format, args...)
	t.ui.Info(s)
}

func (t *Tracer) Error(format string, args ...interface{}) {
	s := fmt.Sprintf("[ixkit] "+format, args...)
	t.ui.Error(s)
}
