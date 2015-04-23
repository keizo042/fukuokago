package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"unicode/utf8"
)

type stateFn func(*lexer) stateFn
type Pos int

type itemTypes int

const (
	itemEnd itemTypes = iota + 1
	itemIdentify
)

type item struct {
	typ itemTypes
	pos Pos
	val string
}

const eof = -1

type lexer struct {
	input string
	state stateFn
	start Pos
	pos   Pos
	width Pos
	items chan item
}

func lex(input string) *lexer {

	l := &lexer{
		input: input,
		start: 0,
		pos:   0,
		width: 0,
		items: make(chan item, 2),
	}

	go l.run()
	return l
}

func (l *lexer) run() {

	for l.state = lexText; l.state != nil; {

		l.state = l.state(l)

	}
	close(l.items)
}

func (l *lexer) peek() rune {
	r := l.next()
	l.backup()
	return r

}

func (l *lexer) next() rune {
	if int(l.pos) >= len(l.input) {
		l.width = 0
		return eof
	}
	r, w := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = Pos(w)
	l.pos += l.width
	return r
}

func (l *lexer) backup() {
	l.pos -= l.width
}

func (l *lexer) ignore() {
	l.start = l.pos
}

func (l *lexer) emit(i itemTypes) {
	l.items <- item{i, l.start, l.input[l.start:l.pos]}
	l.start = l.pos
}

func lexText(l *lexer) stateFn {

	if l.peek() == eof {
		l.emit(itemEnd)
		return nil
	}
	return lexIdentify
}

func lexIdentify(l *lexer) stateFn {
	if l.peek() == eof {
		l.emit(itemEnd)
		return nil
	}
	if l.next() != ' ' {
		return lexIdentify
	}
	l.emit(itemIdentify)
	return lexText
}

func (l *lexer) nextItem() item {
	i := <-l.items
	return i
}

func main() {
	f, _ := os.Open("input.txt")
	b, _ := ioutil.ReadAll(f)
	input := string(b)
	l := lex(input)

	for {
		switch item := l.nextItem(); item.typ {
		case itemEnd:
			return
		default:
			fmt.Println(item.val)
		}
	}

}
