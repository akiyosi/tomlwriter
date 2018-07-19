//package tomlwriter
package main

import (
  "fmt"
  "io/ioutil"
  "strings"
  "unicode"
)

// WriteValue takes new value v, bytes b, table t, key k(, old value o),
//  replace old value o with new value v 
func WriteValue(v interface{}, b []byte, t interface{}, k interface{}, o interface{}) {
//func WriteValue(b []byte) {
  var matchTable bool
  var inMultiline bool
  var inMultilineString bool
  var inMultilineLiteral bool
  var inMultilineArray bool
  var isLineEndingBackSlash bool
  var key, value, multilinevalue, parsedvalue string

  lines := strings.Split(string(b), "\n")
  
  for i := 0; i < len(lines); i++ {
    fmt.Println("------------ ", multilinevalue)
	line := lines[i]
	
    // Comment
    // A hash symbol marks the rest of the line as a comment.
	var vline, cline string
	if strings.Contains(line, "#") {
	  vline = line[:strings.Index(line, "#")]
	  cline = "#" + line[strings.Index(line, "#")+1:]
	} else {
	  vline = line
      cline = ""
	}

	// if Table t is nil
	if t == "" {
	} else {

	  // Inline Table
	  // Inline tables are enclosed in curly braces { and }
	  if strings.Contains(vline, "{") && strings.Contains(vline, "}") && !strings.Contains(vline, "="){
		table := strings.Replace(strings.Split(string(vline), "=")[0], " ", "", -1)
		if t != table {
          matchTable = false
		} else { matchTable = true }
	  }

      // Array of Tables
	  if strings.Contains(vline, "[[") && strings.Contains(vline, "]]") && !strings.Contains(vline, "="){
		table := strings.Replace(strings.Split(string(strings.Split(string(vline), "]]")[0]), "[[")[1], " ", "", -1)
		if t != table {
          matchTable = false
		} else {
		  matchTable = true
		}
	  }

	  // Table
	  // They appear in square brackets on a line by themselves.
	  if strings.Contains(vline, "[") && strings.Contains(vline, "]") && !strings.Contains(vline, "="){
		table := strings.Replace(strings.Split(string(strings.Split(string(vline), "]")[0]), "[")[1], " ", "", -1)
		if t != table {
          matchTable = false
		} else { matchTable = true }
	  }

	}

	// if !strings.Contains(line, "=")
    if strings.Contains(line, "=") {
      key = strings.Replace(strings.Split(string(vline), "=")[0], " ", "", -1)
	  value = strings.Split(string(vline), "=")[1]
      if unicode.IsSpace([]rune(value)[0]) {
          value = value[1:]
      }
    } else {
	  value = vline
    }

    if len(value) > 2 {
      // multi line string
      if strings.Contains(value[:3], `"""`) && !inMultilineString {
	    inMultiline = true
	    inMultilineString = true
        value = value[3:]
	  }
      if strings.Contains(value[len(value)-3:], `"""`) && inMultilineString {
	    inMultiline = false
	    inMultilineString = false
        value = value[:len(value)-3]

        parsedvalue = multilinevalue + value
        multilinevalue = ""
	  }

      // multi line literal
      if strings.Contains(value[:3], `'''`) && !inMultilineLiteral {
	    inMultiline = true
	    inMultilineLiteral = true
        value = value[3:]
	  }
      if strings.Contains(value[len(value)-3:], `'''`) && inMultilineLiteral {
	    inMultiline = false
	    inMultilineLiteral = false
        value = value[:len(value)-3]

        parsedvalue = multilinevalue + value
        multilinevalue = ""
	  }
    }

    // multi line array
	if strings.Contains(value, `[`) && !strings.Contains(value, `]`) && !inMultilineArray {
	  inMultiline = true
	  inMultilineArray = true
	}
	if !strings.Contains(value, `[`) && strings.Contains(value, `]`) && inMultilineArray {
	  inMultiline = false
	  inMultilineArray = false

      parsedvalue = multilinevalue + value
      multilinevalue = ""
	}

    if !inMultiline {
      parsedvalue = value
    } else {

      // When the last non-whitespace character on a line is a \, it will be 
      // trimmed along with all whitespace (including newlines) up to the next
      // non-whitespace character or closing delimiter. 
      if isLineEndingBackSlash {
          for j, c := range value {
              if unicode.IsSpace(c) {
                  continue
              } else {
                  value = value[j:]
                  isLineEndingBackSlash = false
                  break
              }
          }
          if isLineEndingBackSlash { value = "" }
      }
      if value[len(value)-1] == '\\' {
          isLineEndingBackSlash = true
      }
      fmt.Println(value)
      multilinevalue += value
    }

    if matchTable && k == key && o == parsedvalue {
        fmt.Println(key, v, cline, " !match")
    } else {
        fmt.Print(vline, cline, "\n")
    }

  }
}

// test
func main() {
  // file := "/Users/akiyoshi/go/src/github.com/BurntSushi/toml/_examples/example.toml"
  input, _ := ioutil.ReadFile(file)

  WriteValue("hoge", input, "database", "ports", "[ 8001, 8001, 8002 ]") 

}
