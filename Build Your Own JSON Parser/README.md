# JSON Parser Implementation

A custom JSON parser written in Go that validates and parses JSON documents according to the [JSON specification](https://www.json.org/). This implementation was created as part of the [Coding Challenges](https://codingchallenges.fyi/challenges/challenge-json-parser) series.

## Features

- Tokenizes and parses JSON input
- Supports all JSON data types:
  - Objects
  - Arrays
  - Strings (with escape sequences)
  - Numbers (including scientific notation)
  - Booleans
  - Null
- Validates JSON syntax and structure
- Handles nested structures
- Detailed error reporting
- Comprehensive test suite

## Implementation Details

The parser is implemented in two main stages:

### 1. Lexical Analysis (Tokenization)

The `ConvertToToken` function breaks down the input JSON string into tokens:

- Handles structural characters ({, }, [, ], :, ,)
- Processes string literals with escape sequences
- Parses numbers (including scientific notation)
- Recognizes boolean values and null
- Skips whitespace

### 2. Parsing

The parser converts the token stream into an Abstract Syntax Tree (AST) using these node types:

- `ObjectNode`: For JSON objects
- `ArrayNode`: For JSON arrays
- `StringNode`: For string values
- `NumberNode`: For numeric values
- `BooleanNode`: For true/false values
- `NullNode`: For null values

## Error Handling

The parser includes comprehensive error checking for:

- Unterminated strings
- Invalid escape sequences
- Unexpected characters
- Trailing commas
- Invalid numbers
- Malformed objects/arrays
- Unexpected end of input

## Usage

To use the parser to validate JSON:

```go
valid, err := checkValidJson(jsonString)
if err != nil {
    // Handle error
}
if valid {
    // JSON is valid
}
```

## Testing

The implementation includes a test runner (`RunAllTests`) that processes test files organized in step directories. Test files should be placed in the following structure:

```
tests/
    step1/
        valid1.json
        invalid1.json
    step2/
        valid2.json
        invalid2.json
    ...
```

## Credits

This implementation was inspired by:

- [Coding Challenges JSON Parser Challenge](https://codingchallenges.fyi/challenges/challenge-json-parser)
- [Write Your Own JSON Parser](https://ogzhanolguncu.com/blog/write-your-own-json-parser/)
