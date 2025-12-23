package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
)

// limiting the function's access by mengerucutkan type ke io.ReadCloser
func getLinesChannel(f io.ReadCloser) <-chan string {
	out := make(chan string, 1)

	go func() {
		defer close(out)
		defer f.Close()

		str := ""
		for {
			data := make([]byte, 8)
			n, err := f.Read(data)
			if err == io.EOF {
				break
			}

			data = data[:n]
			if i := bytes.IndexByte(data, '\n'); i != -1 {
				str += string(data[:i])
				data = data[1+i:]
				out <- str
				str = ""
			}

			str += string(data)
		}

		if len(str) != 0 {
			out <- str
		}

	}()

	return out
}

func main() {
	f, err := os.Open("./message.txt")
	if err != nil {
		log.Fatal("error", "error", err)
	}

	lines := getLinesChannel(f)
	for line := range lines {
		fmt.Printf("read: %s\n", line)
	}

}
