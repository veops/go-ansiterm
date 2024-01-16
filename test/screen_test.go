package test

import (
	"fmt"
	. "github.com/veops/go-ansiterm"
	"strings"
	"testing"
)

func TestScreen_Display(t *testing.T) {
	t.Run("test", func(t *testing.T) {
		input := []byte{0xd, 0xa, 0x1b, 0x5b, 0x30, 0x3b, 0x33, 0x31, 0x6d, 0x20, 0x79, 0x6f, 0x75, 0x20, 0x68, 0x61, 0x76, 0x65, 0x20, 0x6e, 0x6f, 0x20, 0x70, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x20, 0x74, 0x6f, 0x20, 0x65, 0x78, 0x65, 0x63, 0x75, 0x74, 0x65, 0x20, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x3a, 0x20, 0x1b, 0x5b, 0x30, 0x6d, 0x20, 0x20, 0x1b, 0x5b, 0x30, 0x3b, 0x33, 0x31, 0x6d, 0x20, 0x1b, 0x5b, 0x30, 0x6d, 0xa, 0x72, 0x6f, 0x6f, 0x74, 0x40, 0x32, 0x36, 0x66, 0x36, 0x34, 0x63, 0x39, 0x63, 0x38, 0x65, 0x33, 0x36, 0x3a, 0x7e, 0x23, 0x20, 0x70, 0x77, 0x64, 0x64, 0x8, 0x1b, 0x5b, 0x4b}
		screen := NewScreen(80, 24)
		stream := InitByteStream(screen, false)

		stream.Attach(screen)
		for i := 0; i < 1; i++ {
			fmt.Println("loop:", i)
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

//     def parse_output(self, data, sep='\n'):
//        output = []
//        if not isinstance(data, bytes):
//            data = data.encode('utf-8', 'ignore')
//
//        self.stream.feed(data)
//
//        for line in self.screen.display:
//            print("line:", line)
//            if line.strip():
//                output.append(line)
//        self.screen.reset()
//        return sep.join(output[0:-1])
//
//    def parse_input(self, data):
//        command = []
//        if not isinstance(data, bytes):
//            data = data.encode('utf-8', 'ignore')
//
//        self.stream.feed(data)
//        for line in self.screen.display:
//            line = line.strip()
//            if line:
//                command.append(line)
//        if command:
//            command = command[-1]
//        else:
//            command = ''
//        self.screen.reset()
//        command = self.clean_ps1_etc(command)
//        return command
