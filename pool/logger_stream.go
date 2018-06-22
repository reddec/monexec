package pool

import (
	"io"
	"bufio"
)

type LogInterface interface {
	Println(v ...interface{})
}

func NewLoggerStream(logger LogInterface, prefix string) (io.WriteCloser) {
	reader, writer := io.Pipe()
	go func() {
		scanner := bufio.NewReader(reader)
		for {
			line, _, err := scanner.ReadLine()
			if err != nil {
				break
			}
			logger.Println(prefix, string(line))
		}
	}()
	return writer
}
