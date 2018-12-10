package console

import (
	"bufio"
	"fmt"
	"os"
	"sync"
	"time"

	//"github.com/raoulh/go-progress"
	"github.com/icoco/ixkit-cli/vendors/raoulh/go-progress"
)

type PromptTyping struct {
	tip      string
	onTyping OnTyping
	values   map[string]string

	next *PromptTyping
}

type OnTyping func(p *PromptTyping, value string) *PromptTyping

func NewPromptTyping(tip string, call OnTyping) *PromptTyping {
	p := &PromptTyping{tip: tip, onTyping: call}
	return p
}

func (p *PromptTyping) push(k string, v string) {
	if p.values == nil {
		p.values = make(map[string]string)
	}
	p.values[k] = v
}

func (p *PromptTyping) getValues(in map[string]string) map[string]string {
	dict := p.values
	if nil != dict {
		for k, v := range dict {
			in[k] = v
		}
	}
	if nil != p.next {
		p.next.getValues(in)
	}
	return in
}

//////////////////////////////////
func readTyping() string {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		value := scanner.Text()
		return value
	}

	if scanner.Err() != nil {
		// handle error.
	}

	return ""
}
func readTypingP(p *PromptTyping) string {

	fmt.Println(p.tip)
	v := readTyping()
	pv := p.onTyping(p, v)
	if nil == pv {
		return v
	}
	if pv != p {
		p.next = pv
	}

	return readTypingP(pv)
}

/////////////////
var (
	_quitChan = make(chan bool, 10) //chan bool
)

type Indictiator struct {
}

var instance *Indictiator
var once sync.Once

func TTYIndictiator() *Indictiator {
	once.Do(func() {
		//_quitChan =
		instance = &Indictiator{}
	})
	return instance
}

func (i *Indictiator) start() {
	go showProgress(_quitChan)
}

func showProgress(quitChan chan bool) {
	bar := progress.New(1)
	bar.Format = progress.ProgressFormats[1]
	for {

		select {
		case <-quitChan:
			bar.Finish()

			quitChan <- false
			bar.Cls()
			bar.Stop = true
			return
		default:
			if bar.Stop {
				return
			}
			bar.Busy(1)
			time.Sleep(time.Millisecond * 10)
		}
	}
}

func (i *Indictiator) stop() {
	_quitChan <- true
	time.Sleep(time.Millisecond * 10)
}
