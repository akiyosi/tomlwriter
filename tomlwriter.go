package tomlwriter

import (
	"fmt"
	"io/ioutil"
	"strings"
	"unicode"
	"unsafe"
)

func replaceLinebreak(str, code string) string {
	return strings.NewReplacer("\r\n", code, "\r", code, "\n", code).Replace(str)
}

// WriteValue takes new value, file path, table name, key name, old value,
// replace old value with new value.
func WriteValue(newvalue interface{}, b []byte, table interface{}, keyname interface{}, oldvalue interface{}) []byte {
	v := fmt.Sprintf("%v", newvalue)
	t := fmt.Sprintf("%v", table)
	k := fmt.Sprintf("%v", keyname)
	o := fmt.Sprintf("%v", oldvalue)

	var matchTable bool
	var inMultiline bool
	var inMultilineString bool
	var inMultilineLiteral bool
	var inMultilineArray bool
	var isLineEndingBackSlash bool
	var isMultilineEnd bool
	var isglobalkey bool
	isglobalkey = true
	var key, value, multilinevaluebuffer, multilinevalue, parsedvalue string
	var writestring string
	var writebytes []byte

	// convert line break to "\n" and split with "\n"
	lines := strings.Split(replaceLinebreak(string(b), "\n"), "\n")

	for i := 0; i < len(lines); i++ {
		line := lines[i]

		// Toml Comment
		// A hash symbol marks the rest of the line as a comment.
		var vline, cline string
		if strings.Contains(line, "#") {
			s := strings.Index(line, "#")
			vline = line[:s]
			cline = line[s:]
		} else {
			vline = line
			cline = ""
		}

		// if Toml Table t is nil
		if t == "" {
		} else {

			// Toml Inline Table
			// Inline tables are enclosed in curly braces { and }
			if strings.Contains(vline, "{") && strings.Contains(vline, "}") && !strings.Contains(vline, "=") {
				table := strings.Replace(strings.Split(string(vline), "=")[0], "\x20", "", -1)
				if t != table {
					matchTable = false
				} else {
					matchTable = true
				}
				isglobalkey = false
			}

			// Toml Array of Tables
			if strings.Contains(vline, "[[") && strings.Contains(vline, "]]") && !strings.Contains(vline, "=") {
				table := strings.Replace(strings.Split(string(strings.Split(string(vline), "]]")[0]), "[[")[1], "\x20", "", -1)
				if t != table {
					matchTable = false
				} else {
					matchTable = true
				}
				isglobalkey = false
			}

			// Toml Table
			// They appear in square brackets on a line by themselves.
			if strings.Contains(vline, "[") && strings.Contains(vline, "]") && !strings.Contains(vline, "=") {
				table := strings.Replace(strings.Split(string(strings.Split(string(vline), "]")[0]), "[")[1], "\x20", "", -1)
				if t != table {
					matchTable = false
				} else {
					matchTable = true
				}
				isglobalkey = false
			}
		}

		// if !strings.Contains(line, "=")
		if strings.Contains(line, "=") {
			// key = strings.Replace(strings.Split(string(vline), "=")[0], "\x20", "", -1)
			key = strings.TrimRight(strings.Split(string(vline), "=")[0], "\x20")
			value = strings.Split(string(vline), "=")[1]
			if unicode.IsSpace([]rune(value)[0]) {
				value = value[1:]
			}
		} else {
			value = vline
		}

		// Toml multi line string
		if len(value) > 2 && strings.Contains(value, `"""`) {
			switch inMultilineString {
			case true:
				if value == `"""` {
					value = ""
				} else {
					value = value[:len(value)-3]
				}
				inMultilineString = false
				isMultilineEnd = true
			case false:
				if value == `"""` {
					value = ""
				} else {
					value = value[3:]
				}
				inMultiline = true
				inMultilineString = true
			}
		}

		// Toml multi line literal
		if len(value) > 2 && strings.Contains(value, `'''`) {
			// Toml multi line string
			switch inMultilineLiteral {
			case true:
				if value == `'''` {
					value = ""
				} else {
					value = value[:len(value)-3]
				}
				inMultilineLiteral = false
				isMultilineEnd = true
			case false:
				if value == `'''` {
					value = ""
				} else {
					value = value[3:]
				}
				inMultiline = true
				inMultilineLiteral = true
			}
		}

		// Toml multi line array
		if strings.Contains(value, `[`) && !strings.Contains(value, `]`) && !inMultilineArray {
			inMultiline = true
			inMultilineArray = true
		}
		if !strings.Contains(value, `[`) && strings.Contains(value, `]`) && inMultilineArray {
			inMultilineArray = false
			isMultilineEnd = true
		}

		if !inMultiline {
			parsedvalue = value
		} else {

			if strings.Contains(string(vline), "=") && strings.Contains(string(vline), `"""`) {
				multilinevaluebuffer += `"""` + value + "\n"
			} else if !strings.Contains(string(vline), "=") && strings.Contains(string(vline), `"""`) && !inMultilineString {
				multilinevaluebuffer += value + `"""`
			} else if strings.Contains(string(vline), "=") && strings.Contains(string(vline), `'''`) {
				multilinevaluebuffer += `'''` + value + "\n"
			} else if !strings.Contains(string(vline), "=") && strings.Contains(string(vline), `'''`) && !inMultilineLiteral {
				multilinevaluebuffer += value + `'''`
			} else {
				multilinevaluebuffer += value + "\n"
			}

			// When the last non-whitespace character on a line is a \, it will be
			// trimmed along with all whitespace (including newlines) up to the next
			// non-whitespace character or closing delimiter.
			if isLineEndingBackSlash || inMultilineArray {
				for j, c := range value {
					if unicode.IsSpace(c) {
						continue
					} else {
						value = value[j:]
						isLineEndingBackSlash = false
						break
					}
				}
				if isLineEndingBackSlash {
					value = ""
				}
			}
			if len(value) > 0 {
				if value[len(value)-1] == '\\' {
					isLineEndingBackSlash = true
					if value != "\\" {
						value = value[:len(value)-1]
					} else if value == "\\" {
						value = ""
					}
				}
			}
			multilinevalue += value

			// if value is array
			if inMultilineArray {
				multilinevalue += "\x20"
			}

			// input multilinevalue to parsedvalue
			if isMultilineEnd {
				parsedvalue = multilinevalue
			}
		}

		if !inMultiline {
			value = strings.TrimRight(value, "\x20")
		}

		// Write modified toml data to file
		if isMultilineEnd {
			switch k == strings.Trim(key, ` "`) && o == parsedvalue {
			case true:
				if isglobalkey == true {
					//fmt.Print(key, "\x20=\x20", v, "\x20", cline, "\n")
					writestring = key + "\x20=\x20" + v + "\x20" + cline
				} else if !isglobalkey && matchTable {
					// fmt.Print(key, "\x20=\x20", v, "\x20", cline, "\n")
					writestring = key + "\x20=\x20" + v + "\x20" + cline
				}
			case false:
				// fmt.Print(key, "\x20=\x20", multilinevaluebuffer, "\n")
				writestring = key + "\x20=\x20" + multilinevaluebuffer
			}
			isMultilineEnd = false
			inMultiline = false
			multilinevalue = ""
			multilinevaluebuffer = ""
			if i+1 < len(lines) {
				writestring += "\n"
			}
		} else if !inMultiline {
			switch k == strings.Trim(key, ` "`) && o == value {
			case true:
				if isglobalkey == true {
					// fmt.Print(key, "\x20=\x20", v, "\x20", cline, "\n")
					writestring = key + "\x20=\x20" + v + "\x20" + cline
				} else if !isglobalkey && matchTable {
					// fmt.Print(key, "\x20=\x20", v, "\x20", cline, "\n")
					writestring = key + "\x20=\x20" + v + "\x20" + cline
				}
				if i+1 < len(lines) {
					writestring += "\n"
				}
			case false:
				// fmt.Print(vline, cline, "\n")
				writestring = vline + cline
				if i+1 < len(lines) {
					writestring += "\n"
				}
			}
		}

		if writestring != "" {
			writebytes = append(writebytes, *(*[]byte)(unsafe.Pointer(&writestring))...)
		}
		writestring = ""

	}

	return writebytes
}
