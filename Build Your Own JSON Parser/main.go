package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unicode"
)

type TokenType int
type ASTNode interface{}

const (
	BraceOpen TokenType = iota
	BraceClose
	BracketOpen
	BracketClose
	String
	Number
	Comma
	Colon
	True
	False
	Null
)

type ObjectNode struct {
	Type  string             `json:"type"`
	Value map[string]ASTNode `json:"value"`
}

type ArrayNode struct {
	Type  string    `json:"type"`
	Value []ASTNode `json:"value"`
}

type StringNode struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type NumberNode struct {
	Type  string  `json:"type"`
	Value float64 `json:"value"`
}

type BooleanNode struct {
	Type  string `json:"type"`
	Value bool   `json:"value"`
}

type NullNode struct {
	Type string `json:"type"`
}

type Token struct {
	TokenType TokenType
	Value     string
}

func isNumber(value string) bool {

	if len(value) > 1 && value[0] == '0' && !strings.Contains(value[1:], ".") && !strings.ContainsAny(value[1:], "eE") {
		return false
	}

	_, err := strconv.ParseFloat(value, 64)
	return err == nil
}

func isBooleanTrue(s string) bool {
	return s == "true"
}

func isBooleanFalse(s string) bool {
	return s == "false"
}

func isNull(s string) bool {
	return s == "null"
}

func NewObjectNode(value map[string]ASTNode) ObjectNode {
	return ObjectNode{Type: "Object", Value: value}
}

func NewArrayNode(value []ASTNode) ArrayNode {
	return ArrayNode{Type: "Array", Value: value}
}

func NewStringNode(value string) StringNode {
	return StringNode{Type: "String", Value: value}
}

func NewNumberNode(value float64) NumberNode {
	return NumberNode{Type: "Number", Value: value}
}

func NewBooleanNode(value bool) BooleanNode {
	return BooleanNode{Type: "Boolean", Value: value}
}

func NewNullNode() NullNode {
	return NullNode{Type: "Null"}
}

func ConvertToToken(content string) ([]Token, error) {
	var tokens []Token

	for i := 0; i < len(content); i++ {
		char := content[i]

		switch char {
		case '{':
			tokens = append(tokens, Token{TokenType: BraceOpen, Value: "char"})
			break
		case '}':
			tokens = append(tokens, Token{TokenType: BraceClose, Value: "char"})
			break
		case '[':
			tokens = append(tokens, Token{TokenType: BracketOpen, Value: "char"})
			break
		case ']':
			tokens = append(tokens, Token{TokenType: BracketClose, Value: "char"})
			break
		case ':':
			tokens = append(tokens, Token{TokenType: Colon, Value: "char"})
			break
		case ',':
			tokens = append(tokens, Token{TokenType: Comma, Value: "char"})
			break
		case '"':
			var builder strings.Builder
			builder.WriteString("")
			i += 1
			for i < len(content) && content[i] != '"' {
				if content[i] == '\\' {
					i++
					if i >= len(content) {
						return nil, fmt.Errorf("unterminated escape sequence")
					}

					switch content[i] {
					case '"', '\\', '/', 'b', 'f', 'n', 'r', 't', 'u':
						builder.WriteByte('\\')
						builder.WriteByte(content[i])
					default:
						return nil, fmt.Errorf("illegal backslash escape: \\%c", content[i])
					}

				} else if content[i] < 0x20 {
					// Control characters must be escaped
					return nil, fmt.Errorf("unescaped control character in string")
				} else {
					builder.WriteByte(content[i])
				}
				i++
			}

			if i >= len(content) {
				return nil, fmt.Errorf("unterminated string")
			}
			tokens = append(tokens, Token{TokenType: String, Value: builder.String()})
			break
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 't', 'f', 'n', '-', '.', 'e', 'E':
			var builder strings.Builder
			for i < len(content) && (unicode.IsDigit(rune(content[i])) || unicode.IsLetter(rune(content[i]))) ||
				content[i] == '.' ||
				content[i] == '-' ||
				content[i] == '+' ||
				content[i] == 'e' ||
				content[i] == 'E' {
				builder.WriteByte(content[i])
				i++
			}
			i--

			value := builder.String()
			if isNumber(value) {
				tokens = append(tokens, Token{TokenType: Number, Value: value})
			} else if isBooleanTrue(value) {
				tokens = append(tokens, Token{TokenType: True, Value: value})
			} else if isBooleanFalse(value) {
				tokens = append(tokens, Token{TokenType: False, Value: value})
			} else if isNull(value) {
				tokens = append(tokens, Token{TokenType: Null, Value: value})
			} else {
				return nil, fmt.Errorf("unexpected value: %s", value)
			}
		case ' ', '\t', '\n', '\r':
			continue
		default:
			return nil, fmt.Errorf("unexpected char: %c", char)
		}
	}

	return tokens, nil
}

func parser(tokens []Token) (ASTNode, error) {
	if len(tokens) == 0 {
		return nil, fmt.Errorf("Invalid empty Json")
	}

	current := 0

	var advance = func() Token {
		current++
		if current >= len(tokens) {
			return Token{}
		}
		return tokens[current]
	}
	var parseValue func() (ASTNode, error)
	var parseObject func() (ASTNode, error)
	var parseArray func() (ASTNode, error)

	parseValue = func() (ASTNode, error) {
		token := tokens[current]
		switch token.TokenType {
		case String:
			return NewStringNode(token.Value), nil
		case Number:
			num, err := strconv.ParseFloat(token.Value, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid number: %s", token.Value)
			}
			return NewNumberNode(num), nil
		case True:
			return NewBooleanNode(true), nil
		case False:
			return NewBooleanNode(false), nil
		case Null:
			return NewNullNode(), nil
		case BraceOpen:
			return parseObject()
		case BracketOpen:
			return parseArray()
		case Comma:
			return nil, fmt.Errorf("Value expected")
		default:
			return nil, fmt.Errorf("unexpected token type: %v", token.TokenType)
		}
	}

	parseObject = func() (ASTNode, error) {
		obj := make(map[string]ASTNode)
		token := advance()
		if token == (Token{}) {
			return nil, fmt.Errorf("Unexpected eof")
		}

		for token.TokenType != BraceClose {
			if token.TokenType != String {
				return nil, fmt.Errorf("expected String key in object. Token type: %v", token.TokenType)
			}

			key := token.Value
			token = advance()
			if token == (Token{}) {
				return nil, fmt.Errorf("Unexpected eof")
			}

			if token.TokenType != Colon {
				return nil, fmt.Errorf("expected : in key-value pair")
			}

			token = advance()
			if token == (Token{}) {
				return nil, fmt.Errorf("Unexpected eof")
			}

			value, err := parseValue()
			if err != nil {
				return nil, err
			}
			obj[key] = value

			token = advance()
			if token == (Token{}) {
				return nil, fmt.Errorf("Unexpected eof")
			}

			if token.TokenType == Comma {
				token = advance()
				if token == (Token{}) {
					return nil, fmt.Errorf("Unexpected eof")
				}

				if token.TokenType == BraceClose {
					return nil, fmt.Errorf("Trailing comma")

				}
			}
		}

		return NewObjectNode(obj), nil
	}

	parseArray = func() (ASTNode, error) {
		var arr []ASTNode
		token := advance()
		if token == (Token{}) {
			return nil, fmt.Errorf("Unexpected eof")
		}

		for token.TokenType != BracketClose {
			value, err := parseValue()
			if err != nil {
				return nil, err
			}
			arr = append(arr, value)

			token = advance()
			if token == (Token{}) {
				return nil, fmt.Errorf("Unexpected eof")
			}

			if token.TokenType == Comma {
				token = advance()
				if token == (Token{}) {
					return nil, fmt.Errorf("Unexpected eof")
				}

				if token.TokenType == BracketClose {
					return nil, fmt.Errorf("Trailing comma")

				}
			}
		}

		return NewArrayNode(arr), nil
	}

	AST, err := parseValue()
	if err != nil {
		return nil, err
	}

	if current < len(tokens)-1 {
		return nil, fmt.Errorf("End of file expected")
	}

	return AST, nil
}

func checkValidJson(content string) (bool, error) {
	tokens, err := ConvertToToken(content)

	if err != nil {
		return false, err
	}

	_, err = parser(tokens)
	if err != nil {
		return false, err
	}
	return true, nil
}

func RunAllTests() error {
	baseDir := "tests"

	stepDirs, err := os.ReadDir(baseDir)
	if err != nil {
		return fmt.Errorf("error reading base directory: %v", err)
	}

	for _, stepDir := range stepDirs {
		if stepDir.IsDir() && strings.HasPrefix(stepDir.Name(), "step") {
			fmt.Printf("\n=== Running tests for %s ===\n", stepDir.Name())

			stepPath := filepath.Join(baseDir, stepDir.Name())
			files, err := os.ReadDir(stepPath)
			if err != nil {
				return fmt.Errorf("error reading directory %s: %v", stepDir.Name(), err)
			}

			for _, file := range files {
				if !file.IsDir() && strings.HasSuffix(file.Name(), ".json") {
					filePath := filepath.Join(stepPath, file.Name())
					content, err := os.ReadFile(filePath)
					if err != nil {
						return fmt.Errorf("error reading file %s: %v", file.Name(), err)
					}

					fmt.Printf("\nTesting file: %s/%s\n", stepDir.Name(), file.Name())

					_, err = checkValidJson(string(content))
					if err != nil {
						if strings.HasPrefix(file.Name(), "invalid") {
							fmt.Printf("INVALID JSON - %v\n", err)
							continue
						}
						return fmt.Errorf("error processing %s/%s: %v", stepDir.Name(), file.Name(), err)
					}

					if strings.HasPrefix(file.Name(), "valid") {
						fmt.Printf("VALID JSON\n")
					}
				}
			}
			fmt.Printf("\n=== Completed tests for %s ===\n", stepDir.Name())
		}
	}

	return nil
}

func main() {

	err := RunAllTests()
	if err != nil {
		log.Fatalf("Tests failed: %v", err)
	}
	fmt.Println("\nAll tests completed successfully!")

}
