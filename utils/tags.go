package utils

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
)

const (
	EMPTYSTR = " \t\n\r"
)

func skipEmpty(str []byte, offset int) int {
	i := offset
	for i < len(str) && strings.IndexByte(EMPTYSTR, str[i]) >= 0 {
		i++
	}
	return i
}

func Unquote(str string) string {
	ret, _ := findString([]byte(str), 0)
	return ret
}

func findString(str []byte, offset int) (string, int) {
	return _findWord(str, offset, "\n\r")
}

func findWord(str []byte, offset int) (string, int) {
	return _findWord(str, offset, " :,\t\n}]")
}

func _findWord(str []byte, offset int, sepChars string) (string, int) {
	var buffer bytes.Buffer
	i := skipEmpty(str, offset)
	if i >= len(str) {
		return "", i
	}
	var endstr string
	quote := false
	if str[i] == '"' {
		quote = true
		endstr = "\""
		i++
	} else if str[i] == '\'' {
		quote = true
		endstr = "'"
		i++
	} else {
		// endstr = " :,\t\n\r}]"
		endstr = sepChars
	}
	for i < len(str) {
		if quote && str[i] == '\\' {
			if i+1 < len(str) {
				i++
				switch str[i] {
				case 'n':
					buffer.WriteByte('\n')
				case 'r':
					buffer.WriteByte('\r')
				case 't':
					buffer.WriteByte('\t')
				default:
					buffer.WriteByte(str[i])
				}
				i++
			} else {
				break
			}
		} else if strings.IndexByte(endstr, str[i]) >= 0 { // end
			if quote {
				i++
			}
			break
		} else {
			buffer.WriteByte(str[i])
			i++
		}
	}
	return buffer.String(), i
}

func FindWords(str []byte, offset int) []string {
	words := make([]string, 0)
	for offset < len(str) {
		word, i := findWord(str, offset)
		words = append(words, word)
		i = skipEmpty(str, i)
		if i < len(str) {
			if str[i] == ',' {
				offset = i + 1
			} else {
				panic(fmt.Sprintf("Malformed multi value string: %s", string(str[offset:])))
			}
		} else {
			offset = i
		}
	}
	return words
}

func TagMap(tag reflect.StructTag) map[string]string {
	ret := make(map[string]string)
	str := []byte(tag)
	i := 0
	for i < len(str) {
		var k, val string
		k, i = findWord(str, i)
		if len(k) == 0 {
			break
		}
		i = skipEmpty(str, i)
		if i >= len(str) || strings.IndexByte(EMPTYSTR, str[i]) >= 0 {
			val = ""
		} else if str[i] != ':' {
			panic(fmt.Sprintf("Invalid structTag: %s", tag))
		} else {
			i++
			val, i = findWord(str, i)
		}
		ret[k] = val
		i = skipEmpty(str, i)
	}
	return ret
}

func TagPop(m map[string]string, key string) (map[string]string, string, bool) {
	val, ok := m[key]
	if ok {
		delete(m, key)
	}
	return m, val, ok
}
