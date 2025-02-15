````markdown
# Word Count Tool (ccwc)

This is an implementation of the `wc` command line tool as a solution to the [Coding Challenges project](https://codingchallenges.fyi/challenges/challenge-wc/). It builds a clone of the Unix command line tool `wc`.

## Overview

The `ccwc` tool counts bytes, lines, words, and characters in text files, replicating the core functionality of the Unix `wc` command.

## Installation

1. Clone the repository:

```bash
git clone <your-repository-url>
cd ccwc
```
````

2. Build and install:

```bash
go build -o ccwc main.go
sudo mv ccwc /usr/local/bin/
```

## Features

- `-c` : Count bytes in file
- `-l` : Count lines in file
- `-w` : Count words in file
- `-m` : Count characters with locale awareness
- Default behavior (no flags): Displays line, word, and byte counts
- Supports reading from files or standard input using `-`

## Usage Examples

Count lines, words, and bytes (default behavior):

```bash
ccwc file.txt
```

Count only bytes:

```bash
ccwc -c file.txt
```

Count only lines:

```bash
ccwc -l file.txt
```

Count only words:

```bash
ccwc -w file.txt
```

Count characters (with locale awareness):

```bash
ccwc -m file.txt
```

Read from standard input:

```bash
cat file.txt | ccwc
# or
cat file.txt | ccwc -
```

## Implementation Details

The tool is implemented in Go and features:

- Efficient buffered I/O using 32KB buffer size
- Single-pass counting for better performance
- Proper handling of multibyte characters
- Error handling for file operations
- Support for standard input and file input

## Testing

Compare the output with the original `wc` command:

```bash
# Compare default output
wc test.txt
ccwc test.txt

# Compare line counts
wc -l test.txt
ccwc -l test.txt

# Compare word counts
wc -w test.txt
ccwc -w test.txt

# Compare byte counts
wc -c test.txt
ccwc -c test.txt
```

## License

[MIT License](LICENSE)

## Contributing

Feel free to submit issues and enhancement requests.

## References

- [Coding Challenges - WC Tool](https://codingchallenges.fyi/challenges/challenge-wc/)
- [Digital Ocean - How to Build and Install Go Programs](https://www.digitalocean.com/community/tutorials/how-to-build-and-install-go-programs)

```

```
