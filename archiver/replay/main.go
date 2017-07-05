package main

import (
	"github.com/yuin/gopher-lua"
	"github.com/yuin/gopher-lua/parse"
	"gopkg.in/readline.v1"
	"log"
)

const (
	PS1 = "\033[1;31m> \033[0m"
	PS2 = "\033[1;31m>> \033[0m"
)

type REPL struct {
	L       *lua.LState // the lua virtual machine
	toolbox *ToolBox
	reader  *readline.Instance
}

func NewREPL() *REPL {
	r := new(REPL)
	r.L = lua.NewState()
	r.toolbox = NewToolBox("/data")
	if reader, err := readline.New(PS1); err == nil {
		r.reader = reader
	} else {
		log.Println(err)
		return nil
	}
	return r
}

func (r *REPL) Close() {
	r.toolbox.Close()
	r.reader.Close()
	r.L.Close()
}

// read/eval/print/loop
func (r *REPL) Start() {
	for {
		if str, err := r.loadline(); err == nil {
			r.toolbox.exec(str)
		} else {
			log.Println(err)
			return
		}
	}
}

func incomplete(err error) bool {
	if lerr, ok := err.(*lua.ApiError); ok {
		if perr, ok := lerr.Cause.(*parse.Error); ok {
			return perr.Pos.Line == parse.EOF
		}
	}
	return false
}

func (r *REPL) loadline() (string, error) {
	r.reader.SetPrompt(PS1)
	if line, err := r.reader.Readline(); err == nil {
		if _, err := r.L.LoadString("return " + line); err == nil { // try add return <...> then compile
			return line, nil
		} else {
			return r.multiline(line)
		}
	} else {
		return "", err
	}
}

func (r *REPL) multiline(ml string) (string, error) {
	for {
		if _, err := r.L.LoadString(ml); err == nil { // try compile
			return ml, nil
		} else if !incomplete(err) { // syntax error, but not EOF
			return ml, nil
		} else { // read next line
			r.reader.SetPrompt(PS2)
			if line, err := r.reader.Readline(); err == nil {
				ml = ml + "\n" + line
			} else {
				return "", err
			}
		}
	}
}

func main() {
	r := NewREPL()
	r.Start()
	r.Close()
}
