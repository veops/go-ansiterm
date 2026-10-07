package ansiterm

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"unicode"

	. "github.com/veops/go-ansiterm/const"
	. "github.com/veops/go-ansiterm/pkg"
)

const (
	Bell             = "bell"
	Backspace        = "backspace"
	Tab              = "tab"
	Linefeed         = "linefeed"
	CarriageReturn   = "carriage_return"
	ShiftOut         = "shift_out"
	ShiftIn          = "shift_in"
	Reset            = "reset"
	Index            = "index"
	ReverseIndex     = "reverse_index"
	SetTabStop       = "set_tab_stop"
	SaveCursor       = "save_cursor"
	RestoreCursor    = "restore_cursor"
	AlignmentDisplay = "alignment_display"
)

type Stream struct {
	Listener        *Screen
	Strict          bool
	UseUTF8         bool
	TakingPlainText bool
	Basic           map[string]struct{}
	Escape          map[string]struct{}
	Sharp           map[string]struct{}
	Csi             map[string]struct{}
	//Events          map[string]struct{} // or []string
	TextPattern *regexp.Regexp
	parser      *MyParser
	mu          sync.Mutex
	finished    <-chan struct{}
	closed      bool
	pending     []byte
}

func generateTextPattern() (*regexp.Regexp, error) {
	special := map[string]struct{}{
		RegBEL:   {},
		RegBS:    {},
		RegHT:    {},
		RegLF:    {},
		RegVT:    {},
		RegFF:    {},
		RegCR:    {},
		RegSO:    {},
		RegSI:    {},
		RegESC:   {},
		RegCSIC1: {},
		RegNUL:   {},
		RegDEL:   {},
		RegOSCC1: {},
	}
	//var specialChars []string
	//for k := range special {
	//	specialChars = append(specialChars, k)
	//	//specialChars = append(specialChars, regexp.QuoteMeta(k))
	//}

	var escaped string
	for s := range special {
		//escaped += regexp.QuoteMeta(s)
		escaped += s
	}

	// 创建正则表达式
	pattern := "[^" + escaped + "]+"

	//pattern := "[^" + strings.Join(specialChars, "") + "]+"
	return regexp.Compile(pattern)
}

func initializeStream(screen *Screen, strict bool) *Stream {
	textPattern, err := generateTextPattern()
	if err != nil {
		fmt.Println(err)
	}

	s := &Stream{
		Listener:        nil,
		Strict:          strict,
		UseUTF8:         true,
		TakingPlainText: false,
		TextPattern:     textPattern,
		Basic:           Basic,
		Escape:          Escape,
		Sharp:           Sharp,
		Csi:             Csi,
	}

	//if screen != nil {
	//	s.Attach(screen)
	//}

	return s
}

func (s *Stream) Attach(screen *Screen) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.stopParser()
	s.Listener = screen
	s.closed = false
	s.pending = nil
	s.initializeParser()
}

func (s *Stream) InitializeParser() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return
	}
	s.stopParser()
	s.initializeParser()
}

func (s *Stream) initializeParser() {
	parser := &MyParser{CharChan: make(chan string, 1), IsPlain: make(chan bool, 1)}
	finished := make(chan struct{})
	s.parser, s.finished = parser, finished
	s.TakingPlainText = true
	go func() {
		defer close(finished)
		s.parseFsm(parser)
	}()
}

func (s *Stream) stopParser() {
	if s.parser != nil {
		s.parser.Close()
		<-s.finished
		s.parser, s.finished = nil, nil
	}
}

// Close stops parsing and waits for the worker to exit.
func (s *Stream) Close() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.closed = true
	s.pending = nil
	s.stopParser()
}

func (s *Stream) Feed(data string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.feed(data)
}

func (s *Stream) feed(data string) {
	if s.closed {
		return
	}
	//matchText := s.TextPattern.MatchString
	matchText := s.TextPattern.FindStringSubmatchIndex
	takingPlainText := s.TakingPlainText
	if s.Listener == nil {
		panic("Listener is nil")
	}

	length := len(data)
	offset := 0
	if !s.parser.Running() {
		s.parser.Start()
		s.parser.GetPlain()
	}
	for offset < length {
		if takingPlainText {
			matches := matchText(data[offset:])

			if matches != nil && matches[0] == 0 {
				start, end := matches[0]+offset, matches[1]+offset
				s.Listener.Draw(data[start:end])
				offset = end
			} else {
				takingPlainText = false
			}
		} else {
			// Feed and Close share s.mu, so the input channel stays open here.
			s.parser.CharChan <- data[offset : offset+1]
			takingPlainText = s.parser.GetPlain()
			offset++
		}
	}
	s.TakingPlainText = takingPlainText
}

func (s *Stream) parseFsm(parser *MyParser) {
	if s.Listener == nil {
		panic("listener is nil")
	}

	SpOrGt := SP + ">"
	NulOrDel := NUL + DEL
	CanOrSub := CAN + SUB
	AllowedInCsi := BEL + BS + HT + LF + VT + FF + CR
	OscTermINATORS := map[string]struct{}{
		STC0: {},
		STC1: {},
		BEL:  {},
	}

	var char string
	defer func() {
		parser.Close()
	}()
	for {
		parser.SetPlain(true)
		char = parser.Next()
		if char == "" {
			return
		}
		if char == ESC {
			parser.SetPlain(false)

			char = parser.Next()
			if char == "" {
				return
			}
			if char == "[" {
				char = CSIC1
			} else if char == "]" {
				char = OSCC1
			} else {
				if char == "#" {
					parser.SetPlain(false)
					code := parser.Next()
					if code == "" {
						return
					}
					s.HandleSharp(code)
				} else if char == "%" {
					parser.SetPlain(false)
					code := parser.Next()
					if code == "" {
						return
					}
					s.selectOtherCharset(code)
				} else if char == "(" || char == ")" {
					parser.SetPlain(false)
					code := parser.Next()
					if code == "" {
						return
					}
					if s.UseUTF8 {
						continue
					}
					s.Listener.defineCharset(code, char)
				} else {
					s.HandleEscape(char)
					s.HandleEscape(char)
				}
				continue
			}
		}
		if _, ok := s.Basic[char]; ok {
			if (char == SI || char == SO) && s.UseUTF8 {
				continue
			}
			s.HandleBasic(char)
		} else if char == CSIC1 {
			var params []int
			current := ""
			private := false
			for {
				parser.SetPlain(false)
				char = parser.Next()
				if char == "" {
					return
				}
				if char == "?" {
					private = true
				} else if strings.Contains(AllowedInCsi, char) {
					s.HandleBasic(char)
				} else if strings.Contains(SpOrGt, char) {
				} else if strings.Contains(CanOrSub, char) {
					s.Listener.Draw(char)
					break
				} else if unicode.IsDigit(rune(char[0])) {
					current += char
				} else if char == "$" {
					parser.SetPlain(false)
					char = parser.Next()
					if char == "" {
						return
					}
					break
				} else {
					num, _ := strconv.Atoi(current)
					params = append(params, int(math.Min(float64(num), 9999)))
					if char == ";" {
						current = ""
					} else {
						if private {
							s.HandleCSI(char, params, map[string]any{"private": true})
						} else {
							s.HandleCSI(char, params, nil)
						}
						break
					}
				}
			}
		} else if char == OSCC1 {
			parser.SetPlain(false)
			code := parser.Next()
			if code == "" {
				return
			}
			switch code {
			case "R", "P":
				continue
			}
			param := ""
			for {
				parser.SetPlain(false)
				char = parser.Next()
				if char == "" {
					return
				}
				if char == ESC {
					parser.SetPlain(false)
					next := parser.Next()
					if next == "" {
						return
					}
					char += next
				}
				if _, ok := OscTermINATORS[char]; ok {
					break
				} else {
					param += char
				}
			}
			param = strings.TrimPrefix(param, ";")
			if strings.Contains("01", code) {
				s.Listener.setIconName(param)
			}
			if strings.Contains("02", code) {
				s.Listener.setTitle(param)
			}
		} else if strings.Contains(NulOrDel, char) {
			s.Listener.Draw(char)
		}
	}

}

func (s *Stream) selectOtherCharset(code string) {

}
