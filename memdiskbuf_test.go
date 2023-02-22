package memdiskbuf_test

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/pschou/go-memdiskbuf"
)

func ExampleNewBuffer() {
	buf := memdiskbuf.NewBuffer("/nowhere/unused_file", 10<<10, 8<<10)
	{
		str := "This is a test\n"
		n, err := buf.Write([]byte(str))
		fmt.Printf("Wrote %q (%d) err: %s\n", str, n, err)
	}

	{
		n, err := io.Copy(os.Stdout, buf)
		fmt.Println("Copied", n, "err:", err)
	}
	// Output:
	// Wrote "This is a test\n" (15) err: %!s(<nil>)
	// This is a test
	// Copied 15 err: <nil>
}

func ExampleBuffer_Rewind() {
	buf := memdiskbuf.NewBuffer("/nowhere/unused_file", 10<<10, 8<<10)
	{
		str := "This is a test\n"
		n, err := buf.Write([]byte(str))
		fmt.Printf("Wrote %q (%d) err: %s\n", str, n, err)
	}

	{
		n, err := io.Copy(os.Stdout, buf)
		fmt.Println("Copied", n, "err:", err)
	}

	buf.Rewind()

	{
		n, err := io.Copy(os.Stdout, buf)
		fmt.Println("Copied", n, "err:", err)
	}
	// Output:
	// Wrote "This is a test\n" (15) err: %!s(<nil>)
	// This is a test
	// Copied 15 err: <nil>
	// This is a test
	// Copied 15 err: <nil>
}

func ExampleBuffer_Write_twiceError() {
	buf := memdiskbuf.NewBuffer("/nowhere/unused_file", 10<<10, 8<<10)

	str := "This is a test\n"
	n, err := buf.Write([]byte(str))
	fmt.Printf("Wrote %q (%d) err: %s\n", str, n, err)

	io.Copy(io.Discard, buf)

	n, err = buf.Write([]byte(str))
	fmt.Printf("Second write %q (%d) err: %s\n", str, n, err)
	// Output:
	// Wrote "This is a test\n" (15) err: %!s(<nil>)
	// Second write "This is a test\n" (0) err: Already in read mode
}

func ExampleBuffer_Write_toDisk() {
	buf := memdiskbuf.NewBuffer("./data.tmp", 5, 10)
	str := "This is a test of a long test data string as a dump to file.\n"
	{
		n, err := buf.Write([]byte(str))
		fmt.Printf("Wrote %q (%d) err: %s\n", str, n, err)
	}

	var out strings.Builder
	{
		n, err := io.Copy(&out, buf)
		fmt.Println("Copied", n, "err:", err)
	}

	if out.String() == str {
		fmt.Println("Strings match!")
	}
	buf.Reset()
	// Output:
	// Wrote "This is a test of a long test data string as a dump to file.\n" (61) err: %!s(<nil>)
	// Copied 61 err: <nil>
	// Strings match!
}

func ExampleBuffer_Write_withReset() {
	buf := memdiskbuf.NewBuffer("./abe.tmp", 5, 10)
	str := `Four score and seven years ago our fathers brought forth on this continent, a new nation, conceived in Liberty, and dedicated to the proposition that all men are created equal.

Now we are engaged in a great civil war, testing whether that nation, or any nation so conceived and so dedicated, can long endure. We are met on a great battle-field of that war. We have come to dedicate a portion of that field, as a final resting place for those who here gave their lives that that nation might live. It is altogether fitting and proper that we should do this.

But, in a larger sense, we can not dedicate—we can not consecrate—we can not hallow—this ground. The brave men, living and dead, who struggled here, have consecrated it, far above our poor power to add or detract. The world will little note, nor long remember what we say here, but it can never forget what they did here. It is for us the living, rather, to be dedicated here to the unfinished work which they who fought here have thus far so nobly advanced. It is rather for us to be here dedicated to the great task remaining before us—that from these honored dead we take increased devotion to that cause for which they gave the last full measure of devotion—that we here highly resolve that these dead shall not have died in vain—that this nation, under God, shall have a new birth of freedom—and that government of the people, by the people, for the people, shall not perish from the earth.

—Abraham Lincoln`
	{
		n, err := buf.Write([]byte(str))
		fmt.Println("Wrote", n, "err:", err)
	}

	var out strings.Builder
	{
		n, err := io.Copy(&out, buf)
		fmt.Println("Copied", n, "err:", err)
	}
	if out.String() == str {
		fmt.Println("Strings match!")
	}

	out.Reset()
	buf.Reset()

	{
		n, err := buf.Write([]byte("Hello"))
		fmt.Println("Wrote", n, "err:", err)
	}

	{
		n, err := io.Copy(&out, buf)
		fmt.Println("Copied", n, "err:", err)
	}
	if out.String() == "Hello" {
		fmt.Println("Strings match!")
	}

	// Output:
	// Wrote 1488 err: <nil>
	// Copied 1488 err: <nil>
	// Strings match!
	// Wrote 5 err: <nil>
	// Copied 5 err: <nil>
	// Strings match!
}

func ExampleBuffer_ReadFrom() {
	str := `Four score and seven years ago our fathers brought forth on this continent, a new nation, conceived in Liberty, and dedicated to the proposition that all men are created equal.

Now we are engaged in a great civil war, testing whether that nation, or any nation so conceived and so dedicated, can long endure. We are met on a great battle-field of that war. We have come to dedicate a portion of that field, as a final resting place for those who here gave their lives that that nation might live. It is altogether fitting and proper that we should do this.

But, in a larger sense, we can not dedicate—we can not consecrate—we can not hallow—this ground. The brave men, living and dead, who struggled here, have consecrated it, far above our poor power to add or detract. The world will little note, nor long remember what we say here, but it can never forget what they did here. It is for us the living, rather, to be dedicated here to the unfinished work which they who fought here have thus far so nobly advanced. It is rather for us to be here dedicated to the great task remaining before us—that from these honored dead we take increased devotion to that cause for which they gave the last full measure of devotion—that we here highly resolve that these dead shall not have died in vain—that this nation, under God, shall have a new birth of freedom—and that government of the people, by the people, for the people, shall not perish from the earth.

—Abraham Lincoln`

	buf := memdiskbuf.NewBuffer("./abe2.tmp", 5, 10)
	{
		n, err := buf.ReadFrom(strings.NewReader(str))
		fmt.Println("Wrote", n, "err:", err)
	}

	var out strings.Builder
	{
		n, err := io.Copy(&out, buf)
		fmt.Println("Copied", n, "err:", err)
	}
	if out.String() == str {
		fmt.Println("Strings match!")
	}
	out.Reset()
	buf.Reset()

	// Output:
	// Wrote 1488 err: <nil>
	// Copied 1488 err: <nil>
	// Strings match!
}
