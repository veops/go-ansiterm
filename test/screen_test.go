package test

import (
	"fmt"
	. "github.com/veops/go-ansiterm"
	"strings"
	"testing"
)

func TestScreen_Display(t *testing.T) {
	t.Run("test", func(t *testing.T) {
		input := []byte{
			76, 105, 110, 117, 120, 32, 50, 54, 102, 54, 52, 99, 57, 99, 56, 101, 51, 54, 32, 53, 46, 52, 46, 49, 55, 45, 50, 49, 51, 54, 46, 51, 49, 56, 46, 55, 46, 49, 46, 101, 108, 55, 117, 101, 107, 46, 120, 56, 54, 95, 54, 52, 32, 35, 50, 32, 83, 77, 80, 32, 84, 104, 117, 32, 65, 112, 114, 32, 49, 51, 32, 49, 55, 58, 53, 49, 58, 49, 49, 32, 80, 68, 84, 32, 50, 48, 50, 51, 32, 120, 56, 54, 95, 54, 52, 13, 10, 13, 10, 84, 104, 101, 32, 112, 114, 111, 103, 114, 97, 109, 115, 32, 105, 110, 99, 108, 117, 100, 101, 100, 32, 119, 105, 116, 104, 32, 116, 104, 101, 32, 68, 101, 98, 105, 97, 110, 32, 71, 78, 85, 47, 76, 105, 110, 117, 120, 32, 115, 121, 115, 116, 101, 109, 32, 97, 114, 101, 32, 102, 114, 101, 101, 32, 115, 111, 102, 116, 119, 97, 114, 101, 59, 13, 10, 116, 104, 101, 32, 101, 120, 97, 99, 116, 32, 100, 105, 115, 116, 114, 105, 98, 117, 116, 105, 111, 110, 32, 116, 101, 114, 109, 115, 32, 102, 111, 114, 32, 101, 97, 99, 104, 32, 112, 114, 111, 103, 114, 97, 109, 32, 97, 114, 101, 32, 100, 101, 115, 99, 114, 105, 98, 101, 100, 32, 105, 110, 32, 116, 104, 101, 13, 10, 105, 110, 100, 105, 118, 105, 100, 117, 97, 108, 32, 102, 105, 108, 101, 115, 32, 105, 110, 32, 47, 117, 115, 114, 47, 115, 104, 97, 114, 101, 47, 100, 111, 99, 47, 42, 47, 99, 111, 112, 121, 114, 105, 103, 104, 116, 46, 13, 10, 13, 10, 68, 101, 98, 105, 97, 110, 32, 71, 78, 85, 47, 76, 105, 110, 117, 120, 32, 99, 111, 109, 101, 115, 32, 119, 105, 116, 104, 32, 65, 66, 83, 79, 76, 85, 84, 69, 76, 89, 32, 78, 79, 32, 87, 65, 82, 82, 65, 78, 84, 89, 44, 32, 116, 111, 32, 116, 104, 101, 32, 101, 120, 116, 101, 110, 116, 13, 10, 112, 101, 114, 109, 105, 116, 116, 101, 100, 32, 98, 121, 32, 97, 112, 112, 108, 105, 99, 97, 98, 108, 101, 32, 108, 97, 119, 46, 13, 10, 76, 97, 115, 116, 32, 108, 111, 103, 105, 110, 58, 32, 84, 104, 117, 32, 74, 97, 110, 32, 32, 52, 32, 49, 48, 58, 50, 57, 58, 51, 50, 32, 50, 48, 50, 52, 32, 102, 114, 111, 109, 32, 49, 57, 50, 46, 49, 54, 56, 46, 54, 53, 46, 55, 52, 13, 13, 10, 27, 91, 63, 50, 48, 48, 52, 104, 114, 111, 111, 116, 64, 50, 54, 102, 54, 52, 99, 57, 99, 56, 101, 51, 54, 58, 126, 35, 32,
		}
		screen := NewScreen(80, 24)
		stream := InitByteStream(screen, false)

		stream.Attach(screen)
		stream.Feed(input)
		got := screen.Display()
		s := ParseOutput(got, "\n")
		t.Log(s)
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
