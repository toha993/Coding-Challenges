package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
)

func countCharacters(reader io.Reader) (int, error) {
	r := bufio.NewReader(reader)
	count := 0

	for {
		_, size, err := r.ReadRune()

		if err == io.EOF {
			break
		}
		if err != nil {
			return count, err
		}
		if size > 0 {
			count++
		}
	}
	return count, nil
}

type counts struct {
	lines int64
	words int64
	bytes int64
}

func countAll(reader io.Reader) (counts, error) {
	var c counts
	var inWord bool
	buf := make([]byte, 32*1024)

	for {
		n, err := reader.Read(buf)
		if err != nil && err != io.EOF {
			return c, err
		}
		c.bytes += int64(n)

		for i := 0; i < n; i++ {
			if buf[i] == '\n' {
				c.lines++
			}

			if buf[i] == ' ' || buf[i] == '\n' || buf[i] == '\t' || buf[i] == '\r' {
				inWord = false
			} else {
				if !inWord {
					c.words++
					inWord = true
				}
			}
		}

		if err == io.EOF {
			break
		}
	}

	return c, nil
}

func main() {
	countBytesFlag := flag.Bool("c", false, "count bytes in file")
	countLinesFlag := flag.Bool("l", false, "count lines in file")
	countWordsFlag := flag.Bool("w", false, "count words in file")
	countCharactersFlag := flag.Bool("m", false, "count characters in file")

	flag.Parse()

	args := flag.Args()

	var input io.Reader
	fileName := ""

	if len(args) > 0 {
		fileName = args[0]
		file, err := os.Open(fileName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Cannot open file %s\n", fileName)
			os.Exit(1)
		}
		defer file.Close()
		input = file
	} else {
		input = os.Stdin
	}

	if !*countBytesFlag && !*countLinesFlag && !*countWordsFlag && !*countCharactersFlag {
		counts, err := countAll(input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to count bytes\n")
			os.Exit(1)
		}

		if fileName != "" {
			fmt.Printf("%d %d %d %s\n", counts.lines, counts.words, counts.bytes, fileName)
		} else {
			fmt.Printf("%d %d %d\n", counts.lines, counts.words, counts.bytes)
		}

	}

	if *countBytesFlag {
		counts, err := countAll(input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to count bytes\n")
			os.Exit(1)
		}

		if fileName != "" {
			fmt.Printf("%d %s\n", counts.bytes, fileName)
		} else {
			fmt.Printf("%d\n", counts.bytes)
		}
	}

	if *countLinesFlag {
		counts, err := countAll(input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to count lines\n")
			os.Exit(1)
		}

		if fileName != "" {
			fmt.Printf("%d %s\n", counts.lines, fileName)
		} else {
			fmt.Printf("%d\n", counts.lines)
		}
	}

	if *countWordsFlag {
		counts, err := countAll(input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to count words\n")
			os.Exit(1)
		}

		if fileName != "" {
			fmt.Printf("%d %s\n", counts.words, fileName)
		} else {
			fmt.Printf("%d\n", counts.words)
		}
	}

	if *countCharactersFlag {
		characters, err := countCharacters(input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to count words\n")
			os.Exit(1)
		}

		if fileName != "" {
			fmt.Printf("%d %s\n", characters, fileName)
		} else {
			fmt.Printf("%d\n", characters)
		}
	}
}
