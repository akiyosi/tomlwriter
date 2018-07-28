package tomlwriter

import (
	"reflect"
	"testing"
    "time"
)

type args struct {
	newvalue interface{}
	b        []byte
	table    interface{}
	keyname  interface{}
	oldvalue interface{}
}

type test struct {
	name string
	args args
	want1 []byte
	want2 int
}


func Test_replaceLinebreak(t *testing.T) {
	type args struct {
		str  string
		code string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := replaceLinebreak(tt.args.str, tt.args.code); got != tt.want {
				t.Errorf("replaceLinebreak() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_countAndReplaceSpaceRight(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name  string
		args  args
		want  int
		want1 string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := countAndReplaceNrRight(tt.args.str)
			if got != tt.want {
				t.Errorf("countAndReplaceNrRight() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("countAndReplaceNrRight() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_stringToTomlType(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want interface{}
	}{
		// TODO: Add test cases.
		{" 1", args{`"moji"`}, `"moji"`, },
		{" 2", args{`3e+9`}, "3000000000.000000", },
		{" 3", args{`+1.23e+2`}, "123.000000", },
		{" 4", args{`-1.23e-2`}, "-0.012300", },
		{" 5", args{`1_230`}, "1230.000000", },
		{" 6", args{`-39`}, "-39.000000", },
		{" 7", args{`+3_9_8_7_`}, "3987.000000", },
		{" 8", args{`true`}, "true", },
		{" 9", args{`false`}, "false", },
		{" 10", args{`[ 1, 2 ]`}, "[12]", },
		{" 11", args{`[  [1,2,3],  [4, 5, 6] ]`}, "[[123][456]]", },
		{" 12", args{`[  [ "hoge","fuga"],  ["fizz","buzz"]]`}, `[["hoge""fuga"]["fizz""buzz"]]`, },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := stringToTomlType(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("stringToTomlType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isTomlDateTime(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
        {" 1", args{`2017-07-24T15:00:00Z`}, true, },
        {" 2", args{`2017-07-24T15:00:00-09:00`}, true, },
        {" 3", args{`2017-07-24T15:00:00.123456-09:00`}, true, },
        {" 4", args{`2017-07-24 15:00:00Z`}, true, },
        {" 5", args{`2017-07-24`}, true, },
        {" 6", args{`15:00:00`}, true, },
        {" 7", args{`15:00:00.123456-09:00`}, true, },
        {" 8", args{`2017-07-24T15:00:00X`}, false, },
        {" 9", args{`XXXX-07-24T15:00:00X`}, false, },
        {" 10", args{`"2017-07-24T15:00:00Z"`}, false, },
        {" 11", args{`2017/07/24`}, false, },
        {" 12", args{`15-00-00`}, false, },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := isTomlDateTime(tt.args.s)
            if got != tt.want {
				t.Errorf("isTomlDateTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isTomlInt(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
		{" 1", args{`+99`}, true, },
		{" 2", args{`42`}, true, },
		{" 3", args{`0`}, true, },
		{" 4", args{`-17`}, true, },
		{" 5", args{`1_234_5`}, true, },
		{" 6", args{`0_4`}, true, },
		{" 7", args{`abc`}, false, },
		{" 7", args{`1e8`}, false, },
		{" 7", args{`'hoge'`}, false, },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isTomlInt(tt.args.s); got != tt.want {
				t.Errorf("isTomlInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isTomlFloat(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
		{" 1", args{`+1.0`}, true, },
		{" 2", args{`3.1415`}, true, },
		{" 3", args{`-0.1`}, true, },
		{" 4", args{`5e+22`}, true, },
		{" 5", args{`1e6`}, true, },
		{" 6", args{`-2E-2`}, true, },
		{" 7", args{`abc`}, false, },
		{" 8", args{`123abc`}, false, },
		{" 9", args{`"123"`}, false, },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isTomlFloat(tt.args.s); got != tt.want {
				t.Errorf("isTomlFloat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isTomlArray(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
		{" 1", args{`[  [ "hoge","fuga"],  ["fizz","buzz"]]`}, true, },
		{" 2", args{`[   "hoge","fuga",  ["fizz","buzz"]]`}, false, },
		{" 3", args{`[ 1 2 3]`}, true, },
		{" 3", args{`"[ 1 2 3]"`}, false, },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isTomlArray(tt.args.s); got != tt.want {
				t.Errorf("isTomlArray() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_compareTomlValue(t *testing.T) {
	type args struct {
		l string
		r string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := compareTomlValue(tt.args.l, tt.args.r); got != tt.want {
				t.Errorf("compareTomlValue() = %v, want %v", got, tt.want)
			}
		})
	}
}



func TestWriteValue(t *testing.T) {
	tests := newTestData()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got1, got2 := WriteValue(tt.args.newvalue, tt.args.b, tt.args.table, tt.args.keyname, tt.args.oldvalue)
			//got, got1 := WriteValue(tt.args.newvalue, tt.args.b, tt.args.table, tt.args.keyname, tt.args.oldvalue)
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("\nWriteValue() ---------------------- : \n%v\n\nWant ------------------------------ : \n%v",string(got1), string(tt.want1))
				t.Errorf("\nWriteValue() ---------------------- : \n%v\n\nWant ------------------------------ : \n%v", got1, tt.want1)
			}
			if got2 != tt.want2 {
				t.Errorf("\nWriteValue() got2 = %v, want2 %v", got2, tt.want2)
			}
		})
	}
}

func newTestData() []*test {

	// Data start
	data := []*test{

		//// Toml global Key/Value Pair | XXX SUCCESS XXX
		{"Test global Key-Value Pair | Success Case:",
			args{`"fizz"`, // Value to write
				[]byte( // Input Bytes

`# This is a TOML document.

title = "TOML Example"

[owner]
name = "Tom Preston-Werner"
dob = 1979-05-27T07:32:00-08:00 # First class dates`),

				nil,               // Table
				"title",           // Key
				"TOML Example"}, // Old value
			[]byte( // Expected Bytes

`# This is a TOML document.

title = "fizz"

[owner]
name = "Tom Preston-Werner"
dob = 1979-05-27T07:32:00-08:00 # First class dates`), 3,
		}, // -----------------------------------------------------

		//// Toml global Key/Value Pair | XXX FAIL XXX
		{"Test global Key-Value Pair | Fail Case :",
			args{`"fizz"`, // Value to write
				[]byte( // Input Bytes

`# This is a TOML document.

title = "TOML Example"

[owner]
name = "Tom Preston-Werner"
dob = 1979-05-27T07:32:00-08:00 # First class dates`),

				"nil",             // Table
				"title",           // Key
				"TOML Example"}, // Old value
			[]byte( // Expected Bytes

`# This is a TOML document.

title = "TOML Example"

[owner]
name = "Tom Preston-Werner"
dob = 1979-05-27T07:32:00-08:00 # First class dates`), 0,
		}, // -----------------------------------------------------

		//// Toml Key/Value Pair in Table
		{"Test Key-Value Pair in Table: ",
			args{`"fizz"`, // Value to write
				[]byte( // Input Bytes

`# This is a TOML document.

title = "TOML Example"

[note]
# this is note
"name" = note
180720 = "Add note."
180721 = "Fix note."

[owner]
name = "Tom Preston-Werner"
dob = 1979-05-27T07:32:00-08:00 # First class dates`),

				"owner",                 // Table
				"name",                  // Key
				"Tom Preston-Werner"}, // Old value
			[]byte( // Expected Bytes

`# This is a TOML document.

title = "TOML Example"

[note]
# this is note
"name" = note
180720 = "Add note."
180721 = "Fix note."

[owner]
name = "fizz"
dob = 1979-05-27T07:32:00-08:00 # First class dates`), 12,
		}, // -----------------------------------------------------

		//// Toml Dotted keys
		{"Test Dotted keys: ",
			args{"false", // Value to write
				[]byte( // Input Bytes

`name = "Orange"
physical.color = "orange"
physical.shape = "round"
site."google.com" = true`),

				nil,                 // Table
				`site."google.com"`, // Key
				true},             // Old value
			[]byte( // Expected Bytes

`name = "Orange"
physical.color = "orange"
physical.shape = "round"
site."google.com" = false`), 4,

		}, // -----------------------------------------------------

		//// Toml Muliline string
		{"Test Multiline string: ",
			args{`"""` + "\n       fizz\\\n       buzz" + `"""`, // Value to write
				[]byte( // Input Bytes

`str1 = "The quick brown fox jumps over the lazy dog."

str2 = """
The quick brown \


  fox jumps over \
    the lazy dog."""

str3 = """\
       The quick brown \
       fox jumps over \
       the lazy dog.\
       """`),

				nil,    // Table
				"str2", // Key
				"The quick brown fox jumps over the lazy dog."}, // Old value
			[]byte( // Expected Bytes

`str1 = "The quick brown fox jumps over the lazy dog."

str2 = """
       fizz\
       buzz"""

str3 = """\
       The quick brown \
       fox jumps over \
       the lazy dog.\
       """`), 3,

		}, // -----------------------------------------------------

		//// Toml Muliline string 2
		{"Test Multiline string 2: ",
			args{`"""` + "\n       fizz\\\n       buzz" + `"""`, // Value to write
				[]byte( // Input Bytes

`str1 = "The quick brown fox jumps over the lazy dog."

str3 = """
       The quick brown # this contained
       fox jumps over 
       the lazy dog.
       """`),

				nil,    // Table
				"str1", // Key
				"The quick brown fox jumps over the lazy dog."}, // Old value
			[]byte( // Expected Bytes

`str1 = """
       fizz\
       buzz"""

str3 = """
       The quick brown # this contained
       fox jumps over 
       the lazy dog.
       """`), 1,

		}, // -----------------------------------------------------

		//// Toml Muliline literal
		{"Test Muliline literal: ",
			args{`'''` + "\n       fizz\\\n       buzz" + `'''`, // Value to write
				[]byte( // Input Bytes

`str1 = "cat."

str4 = '''
       The quick brown 

       fox jumps over 
       the lazy dog.
       '''`),

				nil,    // Table
				"str4", // Key
				`
       The quick brown 

       fox jumps over 
       the lazy dog.
       `}, // Old value
			[]byte( // Expected Bytes

`str1 = "cat."

str4 = '''
       fizz\
       buzz'''`), 3,

		}, // -----------------------------------------------------


		//// Toml Literal string
		{"Toml Literal string: ",
			args{`'<\s*\t*\d\+>'`, // Value to write
				[]byte( // Input Bytes

`# What you see is what you get.
winpath  = 'C:\Users\nodejs\templates'
winpath2 = '\\ServerX\admin$\system32\'
quoted   = 'Tom "Dubs" Preston-Werner'
regex    = '<\i\c*\s*>'`),

				nil,             // Table
				"regex",         // Key
				`<\i\c*\s*>`}, // Old value
			[]byte( // Expected Bytes

`# What you see is what you get.
winpath  = 'C:\Users\nodejs\templates'
winpath2 = '\\ServerX\admin$\system32\'
quoted   = 'Tom "Dubs" Preston-Werner'
regex = '<\s*\t*\d\+>'`), 5,

		}, // -----------------------------------------------------

		//// Toml Integer
		{"Toml Integer: ",
			args{1024, // Value to write
				[]byte( // Input Bytes

`int1 = +99
int2 = 42
int3 = 0
int4 = -17
int5 = 1_000
int6 = 5_349_221
int7 = 1_2_3_4_5     # VALID but discouraged`),

				nil,          // Table
				"int7",       // Key
				12345 }, // Old value
			[]byte( // Expected Bytes

`int1 = +99
int2 = 42
int3 = 0
int4 = -17
int5 = 1_000
int6 = 5_349_221
int7 = 1024 # VALID but discouraged`), 7,
		}, // -----------------------------------------------------

		//// Toml Float
		{"Toml Float : ",
			args{8e+99, // Value to write
				[]byte( // Input Bytes

`# exponent
flt4 = 5e+22
flt5 = 1e6
flt6 = -2E-2`),

				nil,    // Table
				"flt5", // Key
				1e6 }, // Old value
			[]byte( // Expected Bytes

`# exponent
flt4 = 5e+22
flt5 = 8e+99
flt6 = -2E-2`), 3,

		}, // -----------------------------------------------------

		//// Toml Date
		{"Toml Date : ",
			args{"2009-05-27T07:32:00Z", // Value to write
				[]byte( // Input Bytes

`[owner]
name = "Tom Preston-Werner"
organization = "GitHub"
bio = "GitHub Cofounder & CEO\nLikes tater tots and beer."
dob = 1979-05-27T07:32:00Z # First class dates? Why not?`),

				"owner",    // Table
				"dob", // Key
				time.Date(1979, 5 ,27 ,7 ,32 ,0 , 000000, time.UTC) }, // Old value
			[]byte( // Expected Bytes

`[owner]
name = "Tom Preston-Werner"
organization = "GitHub"
bio = "GitHub Cofounder & CEO\nLikes tater tots and beer."
dob = 2009-05-27T07:32:00Z # First class dates? Why not?`), 5,

		}, // -----------------------------------------------------


		//// Toml Array
		{"Toml Array: ",
			args{`"fizz"`, // Value to write
				[]byte( // Input Bytes

`[database]
server = "192.168.1.1"
ports = [ 8001, 8001, 8002 ]
connection_max = 1e5
enabled = true`),

				"database",            // Table
				"ports",         // Key
				[]int{8001, 8001, 8002} }, // Old value
			[]byte( // Expected Bytes

`[database]
server = "192.168.1.1"
ports = "fizz"
connection_max = 1e5
enabled = true`), 3,

		}, // -----------------------------------------------------


		//// Toml Array
		{"Toml Array: ",
			args{"[ 99.2, 23.4 ]", // Value to write
				[]byte( // Input Bytes

`arr5 = [ [ 1, 2 ], ["a", "b", "c"] ]
arr7 = [
  1, 2, 3
]
arr8 = [
  1,
  2, # this is ok
]`),

				nil,            // Table
				"arr7",         // Key
				[]int{1,2,3} }, // Old value
			[]byte( // Expected Bytes

`arr5 = [ [ 1, 2 ], ["a", "b", "c"] ]
arr7 = [ 99.2, 23.4 ]
arr8 = [
  1,
  2, # this is ok
]`), 2,

		}, // -----------------------------------------------------

		//// Toml Table
		{"Toml Table: ",
			args{`"fizz"`, // Value to write
				[]byte( // Input Bytes

`[dog."tater.man"]
type.name = "pug"`),

				nil,         // Table
				"type.name", // Key
				"pug"},    // Old value
			[]byte( // Expected Bytes

`[dog."tater.man"]
type.name = "fizz"`), 2,

		}, // -----------------------------------------------------

		//// Toml Inline Table
		{"Toml Inline Table: ",
			args{`"fizz"`, // Value to write
				[]byte( // Input Bytes

`key = value
name = { first = "Tom", last = "Preston-Werner" }
point = { x = 1, y = 2 }
animal = { type.name = "pug" }`),

				"name",   // Table
				"first",  // Key
				"Tom"}, // Old value
			[]byte( // Expected Bytes

`key = value
name = { first = "fizz", last = "Preston-Werner" }
point = { x = 1, y = 2 }
animal = { type.name = "pug" }`), 2,

		}, // -----------------------------------------------------

		//// Toml Array Table
		{"Toml Array Table: ",
			args{`"fizz"`, // Value to write
				[]byte( // Input Bytes

`[[products]]
name = "Hammer"
sku = 738594937

[[products]]

[[products]]
name = "Nail"
sku = 284758393
color = "gray"`),

				"[products]", // Table
				"name",       // Key
				"Nail"},    // Old value
			[]byte( // Expected Bytes

`[[products]]
name = "Hammer"
sku = 738594937

[[products]]

[[products]]
name = "fizz"
sku = 284758393
color = "gray"`), 8,

		}, // -----------------------------------------------------

		//// Toml write new Key-Value
		{"Toml : Write new global Key-Value: ",
			args{`"fizz"`, // Value to write
				[]byte( // Input Bytes

`key = value

name = { first = "Tom", last = "Preston-Werner" }
[animal]
type.name = "pug"`),

				"",       // Table
				"newKey", // Key
				""},      // Old value
			[]byte( // Expected Bytes

`key = value

name = { first = "Tom", last = "Preston-Werner" }
newKey = "fizz"
[animal]
type.name = "pug"`), 4,

		}, // -----------------------------------------------------

		//// Toml write new Key-Value
		{"Toml : Write new global Key-Value 2: ",
			args{`"fizz"`, // Value to write
				[]byte( // Input Bytes

`key = value

name = { first = "Tom", last = "Preston-Werner" }


[animal]
type.name = "pug"`),

				"",       // Table
				"newKey", // Key
				""},      // Old value
			[]byte( // Expected Bytes

`key = value

name = { first = "Tom", last = "Preston-Werner" }
newKey = "fizz"


[animal]
type.name = "pug"`), 4,

		}, // -----------------------------------------------------

		//// Toml write new Key-Value in Table
		{"Toml : Write new Key-Value in Table : ",
			args{`"fizz"`, // Value to write
				[]byte( // Input Bytes

`name = { first = "Tom", last = "Preston-Werner" }
[animal]
type.name = "pug"

[bird]
type.name = "Sparrows"
`),

				"bird",   // Table
				"newKey", // Key
				""},      // Old value
			[]byte( // Expected Bytes

`name = { first = "Tom", last = "Preston-Werner" }
[animal]
type.name = "pug"

[bird]
type.name = "Sparrows"
newKey = "fizz"
`), 7,
		}, // -----------------------------------------------------

		//// Toml write new Key-Value in Array Table
		{"Toml : Write new Key-Value in Array table: ",
			args{`"fizz"`, // Value to write
				[]byte( // Input Bytes

`name = { first = "Tom", last = "Preston-Werner" }
[animal]
type.name = "pug"

[[birds]]
type.name = "Sparrows"
[[birds]]
type.name = "Songbirds"
[[birds]]
type.name = "Ducks"
`),

				"[birds]",   // Table
				"type.name", // Key
				""},         // Old value
			[]byte( // Expected Bytes

`name = { first = "Tom", last = "Preston-Werner" }
[animal]
type.name = "pug"

[[birds]]
type.name = "Sparrows"

[[birds]]
type.name = "fizz"

[[birds]]
type.name = "Songbirds"
[[birds]]
type.name = "Ducks"
`), 7,
		}, // -----------------------------------------------------

		//// Toml write new Key-MultilineValue in Array Table
		{"Toml : Write new Key-MultilineValue in Array table: ",
			args{`"""` + "\n  fizz\n  fizz\n" + `"""`, // Value to write
				[]byte( // Input Bytes

`[[birds]]
type.name = """
  Songbirds
  Songbirds
"""
[[birds]]
type.name = """
  Ducks
  Ducks
"""`),

				"[birds]",   // Table
				"type.name", // Key
				""},         // Old value
			[]byte( // Expected Bytes

`[[birds]]
type.name = """
  Songbirds
  Songbirds
"""

[[birds]]
type.name = """
  fizz
  fizz
"""

[[birds]]
type.name = """
  Ducks
  Ducks
"""`), 6,
		}, // -----------------------------------------------------



		//// Toml write too complex toml file
		{"Toml : Write too complex toml file : ",
			args{`'hogefuga'`, // Value to write
				[]byte( // Input Bytes

// --------------------------------------------------------
`[[plugins]]
repo = 'Shougo/dein.vim'

#[[plugins]]
#repo = 'ryanoasis/vim-devicons'
# 
#[[plugins]]
#repo = 'scrooloose/nerdtree'
#hook_add = '''
#  " appeylance
#  " vim-devicons
#  let g:webdevicons_conceal_nerdtree_brackets = 1
#  let g:WebDevIconsNerdTreeAfterGlyphPadding = ' '
#  " dir-icons
#  let g:WebDevIconsUnicodeDecorateFolderNodes = 1
#  let g:DevIconsEnableFoldersOpenClose = 1
#  let g:WebDevIconsUnicodeDecorateFolderNodesDefaultSymbol = ''
#  let g:DevIconsDefaultFolderOpenSymbol = ''
#  " file-icons
#  let g:WebDevIconsUnicodeDecorateFileNodesExtensionSymbols = {}
#  let g:WebDevIconsUnicodeDecorateFileNodesExtensionSymbols['html'] = ''
#  let g:WebDevIconsUnicodeDecorateFileNodesExtensionSymbols['css'] = ''
#  let g:WebDevIconsUnicodeDecorateFileNodesExtensionSymbols['md'] = ''
#  let g:WebDevIconsUnicodeDecorateFileNodesExtensionSymbols['txt'] = ''
#  " open nerdtree
#  nnoremap <silent> <Space>d   :call NERDTreeCWD()<CR>
#'''
# let g:webdevicons_conceal_nerdtree_brackets = 
# let g:webdevicons_conceal_nerdtree_brackets = 
# let g:webdevicons_conceal_nerdtree_brackets = 
# --- gonvim need it
[[plugins]]
repo = 'equalsraf/neovim-gui-shim'
#hook_source = '''
#  let g:deoplete#enable_at_startup = 1
#  let g:deoplete#max_list = 30
#'''
# " nnoremap <Space>fb :GonvimFuzzyBLines<CR>
# let g:webdevicons_conceal_nerdtree_brackets = 
# let g:webdevicons_conceal_nerdtree_brackets = 
# let g:webdevicons_conceal_nerdtree_brackets = 
# let g:webdevicons_conceal_nerdtree_brackets = 
hoge
`),
// --------------------------------------------------------

				"[plugins]",   // Table
				"repo", // Key
				nil},      // Old value
			[]byte( // Expected Bytes 


// --------------------------------------------------------
`[[plugins]]
repo = 'Shougo/dein.vim'

#[[plugins]]
#repo = 'ryanoasis/vim-devicons'
# 
#[[plugins]]
#repo = 'scrooloose/nerdtree'
#hook_add = '''
#  " appeylance
#  " vim-devicons
#  let g:webdevicons_conceal_nerdtree_brackets = 1
#  let g:WebDevIconsNerdTreeAfterGlyphPadding = ' '
#  " dir-icons
#  let g:WebDevIconsUnicodeDecorateFolderNodes = 1
#  let g:DevIconsEnableFoldersOpenClose = 1
#  let g:WebDevIconsUnicodeDecorateFolderNodesDefaultSymbol = ''
#  let g:DevIconsDefaultFolderOpenSymbol = ''
#  " file-icons
#  let g:WebDevIconsUnicodeDecorateFileNodesExtensionSymbols = {}
#  let g:WebDevIconsUnicodeDecorateFileNodesExtensionSymbols['html'] = ''
#  let g:WebDevIconsUnicodeDecorateFileNodesExtensionSymbols['css'] = ''
#  let g:WebDevIconsUnicodeDecorateFileNodesExtensionSymbols['md'] = ''
#  let g:WebDevIconsUnicodeDecorateFileNodesExtensionSymbols['txt'] = ''
#  " open nerdtree
#  nnoremap <silent> <Space>d   :call NERDTreeCWD()<CR>
#'''
# let g:webdevicons_conceal_nerdtree_brackets = 
# let g:webdevicons_conceal_nerdtree_brackets = 
# let g:webdevicons_conceal_nerdtree_brackets = 
# --- gonvim need it

[[plugins]]
repo = 'hogefuga'

[[plugins]]
repo = 'equalsraf/neovim-gui-shim'
#hook_source = '''
#  let g:deoplete#enable_at_startup = 1
#  let g:deoplete#max_list = 30
#'''
# " nnoremap <Space>fb :GonvimFuzzyBLines<CR>
# let g:webdevicons_conceal_nerdtree_brackets = 
# let g:webdevicons_conceal_nerdtree_brackets = 
# let g:webdevicons_conceal_nerdtree_brackets = 
# let g:webdevicons_conceal_nerdtree_brackets = 
hoge
`), 32,
		}, // -----------------------------------------------------



		// **** End
	}
	return data
}
