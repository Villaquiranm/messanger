package stdin

import (
	"bufio"
	"fmt"
	"os"
)

// Reader reads messages from Stdin
type Reader struct {
	channel chan string
}

// NewReader returns and initializes a new Stdin reader
func NewReader() *Reader {
	r := &Reader{}
	r.channel = make(chan string, 5)
	go r.listenStandardInput()
	return r
}

func (r *Reader) listenStandardInput() {
	for {
		reader := bufio.NewReader(os.Stdin)
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("error: %s", err.Error())
			break
		}
		if line != "\n" {
			r.channel <- line
		}
	}
}

// Next Gets the next message from the channel
func (r *Reader) Read() string {
	return <-r.channel
}
