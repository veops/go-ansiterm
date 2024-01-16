package test

import (
	"fmt"
	. "github.com/veops/go-ansiterm"
	"strings"
	"testing"
)

func TestScreen_Display(t *testing.T) {
	t.Run("test", func(t *testing.T) {
		input := []byte{}
		screen := NewScreen(80, 24)
		stream := InitByteStream(screen, false)

		stream.Attach(screen)
		for i := 0; i < 2; i++ {

			stream.Feed(input)
			got := screen.Display()
			screen.Reset()
			s := ParseOutput(got, "\n")

			t.Log(s)
		}

	})
}

func ParseOutput(data []string, sep string) string {
	var output []string
	for _, line := range data {
		if strings.TrimSpace(line) != "" {
			fmt.Println("line:", strings.TrimSpace(line))
			output = append(output, strings.TrimSpace(line))
		}
	}
	return strings.Join(output, sep)
}