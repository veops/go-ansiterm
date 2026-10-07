package ansiterm

import (
	"fmt"
	"strings"
	"unicode/utf8"

	. "github.com/veops/go-ansiterm/pkg"
)

type ByteStream struct {
	*Stream
	utf8Decoder func(data []byte) (string, error)
}

func (b *ByteStream) Feed(data []byte) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.closed {
		return
	}
	var dataStr string
	if b.UseUTF8 {
		if len(b.pending) > 0 {
			data = append(b.pending, data...)
			b.pending = nil
		}
		end := len(data)
		for i := len(data) - 1; i >= max(0, len(data)-utf8.UTFMax); i-- {
			if utf8.RuneStart(data[i]) {
				if !utf8.FullRune(data[i:]) {
					end = i
					b.pending = append([]byte(nil), data[i:]...)
				}
				break
			}
		}
		if utf8.Valid(data[:end]) {
			dataStr = string(data[:end])
		} else {
			var err error
			dataStr, err = b.utf8Decoder(data[:end])
			if err != nil {
				fmt.Println(err)
			}
		}
	} else {
		dataStr = BytesToString(data)
	}
	b.feed(dataStr)
}

func (b *ByteStream) selectOtherCharset(code string) {
	if code == "@" {
		b.UseUTF8 = false
	} else if strings.Contains("G8", code) {
		b.UseUTF8 = true
	}
}

func InitByteStream(screen *Screen, strict bool) *ByteStream {
	stream := initializeStream(screen, strict)
	bs := &ByteStream{
		Stream:      stream,
		utf8Decoder: DecodeUTF8WithReplacement,
	}
	return bs
}
