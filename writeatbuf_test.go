package memdiskbuf_test

import (
	"fmt"
	"log"
	"os"

	"github.com/pschou/go-memdiskbuf"
)

func ExampleNewWriterAtBufWithBlockSize() {
	fh, err := os.OpenFile("example_file", os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		log.Fatal(err)
	}
	fh.Truncate(100)

	wab := memdiskbuf.NewWriterAtBufWithBlockSize(fh, 40, 10)
	wab.WriteAt([]byte("world,"), 6)
	wab.WriteAt([]byte("hello "), 0)
	wab.WriteAt([]byte(" will "), 12)
	wab.WriteAt([]byte("you be"), 18)
	wab.WriteAt([]byte("iend?\n"), 30)
	wab.WriteAt([]byte(" my fr"), 24)
	err = wab.Flush()
	if err != nil {
		fmt.Println("err:", err)
	}

	dat := make([]byte, 36)
	fh.ReadAt(dat, 0)
	fmt.Printf("%q\n", dat)
	fmt.Println("written:", wab.Written())

	fh.Close()
	os.Remove("example_file")
	// Output:
	// "hello world, will you be my friend?\n"
	// written: 36
}
