package tomlwriter

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
	"unsafe"
)

func replaceLinebreak(str, code string) string {
	return strings.NewReplacer("\r\n", code, "\r", code, "\n", code).Replace(str)
}

func countAndReplaceSpaceRight(str string) (int, string) {
	var count int
	for z := 0; z < len(str); z++ {
		if str[len(str)-z-1] == '\n' {
			count++
		} else {
			break
		}
	}
	newstr := strings.TrimRight(str, "\n")

	return count, newstr
}

// parse value in toml file
func stringToTomlType(s string) interface{} {
	s = strings.Trim(s, " \t")

	if isTomlFloat(s) {
		num, _ := strconv.ParseFloat(s, 32)
		return fmt.Sprintf("%v", num)
	} else if isTomlInt(s) {
		num, _ := strconv.ParseInt(strings.Replace(s, "_", "", -1), 10, 32)
		return fmt.Sprintf("%v", num)
	} else if isTomlArray(s) {
		return strings.Replace(strings.Replace(s, "\x20", "", -1), ",", "", -1)
	}

	if len(s) > 1 {
		if s[0] == '"' && s[len(s)-1] == '"' {
			return s
		} else if s[0] == '\'' && s[len(s)-1] == '\'' {
			return `"` + s[1:len(s)-1] + `"`
		}
	}

	return s
}

func isTomlInt(s string) bool {
	for _, c := range s {
		if unicode.IsDigit(c) || c == '+' || c == '-' || c == '_' {
			continue
		} else {
			return false
		}
	}
	return true
}

func isTomlFloat(s string) bool {
	var dotCount int
	var eCount int
	for _, c := range s {
		if unicode.IsDigit(c) || (c == '.' && (dotCount == 0)) || (((c == 'e') || (c == 'E')) && eCount == 0) || c == '+' || c == '-' {
			if (c == 'e') || (c == 'E') {
				eCount++
			}
			if c == 'e' {
				dotCount++
			}
			continue

		} else {
			return false
		}
	}
	return true
}

func isTomlArray(s string) bool {
	if s[0] != '[' && s[len(s)-1] != ']' {
		return false
	}
	var brac int
	var kets int

	for _, c := range s {
		if c == '[' {
			brac++
			continue
		}
		if c == ']' {
			kets++
			continue
		}
	}
	if brac != kets {
		return false
	}

	return true
}

// v0.4.0 compare
func compareTomlValue(l, r string) bool {
	left := stringToTomlType(l)
	right := stringToTomlType(r)
	return left == right
}

// WriteValue takes new value, file path, table name, key name, old value,
// return bytes replaced old value with new value, and it's line number
func WriteValue(newvalue interface{}, b []byte, table interface{}, keyname interface{}, oldvalue interface{}) ([]byte, int) {

	if newvalue == nil || newvalue == "" {
		return b, 0
	}
	if table == nil {
		table = ""
	}
	if keyname == nil || keyname == "" {
		return b, 0
	}
	if oldvalue == nil {
		oldvalue = ""
	}

	v := fmt.Sprintf("%v", newvalue)
	t := fmt.Sprintf("%v", table)
	k := fmt.Sprintf("%v", keyname)
	//o := fmt.Sprintf("%v", oldvalue)
	// o := valueToString(oldvalue)
	// o := oldvalue

	var o string
	switch oldvalue.(type) {
	case bool, int, float64:
		o = fmt.Sprintf("%v", oldvalue)
	case string:
		if oldvalue == "" {
			o = ""
		} else {
			o = fmt.Sprintf(`"%v"`, oldvalue)
		}
	default:
		n := fmt.Sprintf("%v", oldvalue)
		if n[0] == '[' && n[len(n)-1] == ']' {
			o = strings.Replace(fmt.Sprintf("%v", oldvalue), "\x20", ",\x20", -1)
		} else {
			o = fmt.Sprintf(`"%v"`, oldvalue)
		}
	}

	var matchTable bool
	var matchArrayTable bool
	var matchKeyInArrayTable bool
	var doneWriteNewValue bool
	var inMultiline bool
	var inMultilineString bool
	var inMultilineLiteral bool
	var inMultilineArray bool
	var isLineEndingBackSlash bool
	var isMultilineEnd bool
	var isglobalkey bool
	var isInlineTableMatch bool
	isglobalkey = true
	var key, value, multilinevaluebuffer, multilinevalue, parsedvalue string
	var writestring string

	var writebytes []byte
	var writeLinenumber int

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

		// Write new entry
		if !inMultiline {
			if (strings.Contains(vline, "[") && strings.Contains(vline, "]") && !strings.Contains(vline, "=")) || i == len(lines)-1 {

				// if old value is nil and writeLinenumber == 0, matchTable is true, then,
				// key and new value write in the previous table as new entry
				if writeLinenumber == 0 && doneWriteNewValue == false {
					if (t != "" && o == "" && matchTable == true) || (t == "" && o == "" && isglobalkey == true) {
						co, trimedRightString := countAndReplaceSpaceRight(string(writebytes))
						writebytes = *(*[]byte)(unsafe.Pointer(&trimedRightString))
						if matchArrayTable == true && matchKeyInArrayTable {
							writestring += "\n[" + t + "]"
						}
						writestring += "\n" + k + "\x20=\x20" + v
						for a := 0; a < co; a++ {
							writestring += "\n"
						}
						doneWriteNewValue = true
					}
				}
			}
		}

		// Toml Inline Table
		// Inline tables are enclosed in curly braces { and }
		if strings.Contains(vline, "{") && strings.Contains(vline, "}") && strings.Contains(vline, "=") && !inMultiline {
			if strings.Contains(v, "\n") {
				return b, 0
			}

			table := strings.Replace(strings.Split(string(vline), "=")[0], "\x20", "", 1)
			if t != "" && t == table {
				inlinetable := strings.Trim(strings.SplitAfterN(string(vline), "=", 2)[0], "\x20")
				inlinetablevalue := strings.Trim(strings.SplitAfterN(string(vline), "=", 2)[1], " {}")
				inlines := strings.Split(inlinetablevalue, ",")
				inlinestring := inlinetable + "\x20" + "{ "
				for a, inline := range inlines {
					key = strings.Split(string(inline), "=")[0]
					value = strings.Split(string(inline), "=")[1]

					// Not support array of inline table
					if !strings.Contains(value, "[") && strings.Contains(inlinetablevalue, "[") {
						return b, 0
					}

					// switch strings.Trim(k, `"`) == strings.Trim(key, ` "`) && o == strings.Trim(value, "\x20") {
					switch strings.Trim(k, `"`) == strings.Trim(key, ` "`) && compareTomlValue(o, strings.Trim(value, "\x20")) {
					case true:
						isInlineTableMatch = true
						inlinestring += key + "=\x20" + v
						writeLinenumber = i + 1
					case false:
						inlinestring += inline
					}
					if a < len(inlines)-1 {
						inlinestring += ","
					}

				}
				inlinestring += " }"

				if isInlineTableMatch == false {
					writestring = vline + cline
					if i+1 < len(lines) {
						writestring += "\n"
					}
				} else {
					writestring = inlinestring
				}
				if cline != "" {
					writestring += "\x20" + cline
				}
				if i+1 < len(lines) {
					writestring += "\n"
				}
				isInlineTableMatch = false

				writebytes = append(writebytes, *(*[]byte)(unsafe.Pointer(&writestring))...)
				writestring = ""

				continue
			}
		}

		// Toml Table | Array Table
		// They appear in square brackets on a line by themselves.
		if strings.Contains(vline, "[") && strings.Contains(vline, "]") && !strings.Contains(vline, "=") && !inMultiline {
			table := strings.Replace(strings.Replace(strings.Replace(string(vline), "[", "", 1), "]", "", 1), "\x20", "", -1)
			if strings.Contains(vline, "[[") && strings.Contains(vline, "]]") {
				matchKeyInArrayTable = false
				matchArrayTable = true
			} else {
				matchArrayTable = false
			}
			if t != "" && t != table {
				matchTable = false
			} else {
				matchTable = true
			}
			isglobalkey = false
		}

		// if t == "" && !isglobalkey && o == "" {
		//   return b, 0
		// }

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
				multilinevalue = `"`
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
				multilinevalue = `'`
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
				if inMultilineArray {
					multilinevaluebuffer += value + cline
				} else {
					multilinevaluebuffer += value
				}
				if !isMultilineEnd {
					multilinevaluebuffer += "\n"
				}
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
				if !strings.Contains(string(vline), "=") && strings.Contains(string(vline), `"""`) {
					parsedvalue = multilinevalue + `"`
				} else if !strings.Contains(string(vline), "=") && strings.Contains(string(vline), `'''`) {
					parsedvalue = multilinevalue + `'`
				} else {
					parsedvalue = multilinevalue
				}
			}
		}

		if !inMultiline {
			value = strings.TrimRight(value, "\x20")
		}
		if matchTable && strings.Trim(k, `"`) == strings.Trim(key, ` "`) {
			matchKeyInArrayTable = true
		}

		// writestring += "\n" + "o :::: " + fmt.Sprintf("%v", stringToTomlType(o))
		// writestring += "\n" + "parsedvalue :::: |" + fmt.Sprintf("%v", stringToTomlType(parsedvalue)) + "|"
		// writestring += "\n" + fmt.Sprintf("+++ %v", compareTomlValue(o, parsedvalue)) + "\n\n"

		// Write modified toml data
		if isMultilineEnd {
			//switch strings.Trim(k, `"`) == strings.Trim(key, ` "`) && o == parsedvalue && o != "" {
			switch strings.Trim(k, `"`) == strings.Trim(key, ` "`) && compareTomlValue(o, parsedvalue) && o != "" {
			case true:
				if isglobalkey == true && t == "" {
					//fmt.Print(key, "\x20=\x20", v, "\x20", cline, "\n")
					writestring += key + "\x20=\x20" + v
					if cline != "" {
						writestring += "\x20" + cline
					}
					writeLinenumber = i + 1
				} else if !isglobalkey && matchTable {
					// fmt.Print(key, "\x20=\x20", v, "\x20", cline, "\n")
					writestring += key + "\x20=\x20" + v
					if cline != "" {
						writestring += "\x20" + cline
					}
					writeLinenumber = i + 1
				} else {
					writestring += key + "\x20=\x20" + multilinevaluebuffer
				}
			case false:
				// fmt.Print(key, "\x20=\x20", multilinevaluebuffer, "\n")
				writestring += key + "\x20=\x20" + multilinevaluebuffer
			}
			isMultilineEnd = false
			inMultiline = false
			multilinevalue = `"`
			multilinevaluebuffer = ""
			if i+1 < len(lines) {
				writestring += "\n"
			}
		} else if !inMultiline {

			//switch strings.Trim(k, `"`) == strings.Trim(key, ` "`) && o == value && o != "" {
			switch strings.Trim(k, `"`) == strings.Trim(key, ` "`) && compareTomlValue(o, parsedvalue) && o != "" {
			case true:
				if isglobalkey == true && t == "" {
					// fmt.Print(key, "\x20=\x20", v, "\x20", cline, "\n")
					writestring += key + "\x20=\x20" + v
					if cline != "" {
						writestring += "\x20" + cline
					}
					writeLinenumber = i + 1
				} else if !isglobalkey && matchTable {
					// fmt.Print(key, "\x20=\x20", v, "\x20", cline, "\n")
					writestring += key + "\x20=\x20" + v
					if cline != "" {
						writestring += "\x20" + cline
					}
					writeLinenumber = i + 1
				} else {
					writestring += vline + cline
				}
				if i+1 < len(lines) {
					writestring += "\n"
				}
			case false:
				// fmt.Print(vline, cline, "\n")
				writestring += vline + cline
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

	return writebytes, writeLinenumber
}
