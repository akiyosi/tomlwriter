# tomlwriter

Tomlwriter will help it if you want to change specific values in the toml file while keeping comments, indentations, etc. 

## Installation

```
go get github.com/akiyosi/tomlwriter
```

## Usage
See `_example/example.go`. It make changes like the following:

#### Before

```
# This is a TOML document. Boom.

title = "TOML Example"

[owner]
name = "Tom Preston-Werner"
organization = "GitHub"
bio = "GitHub Cofounder & CEO\nLikes tater tots and beer."
dob = 1979-05-27T07:32:00Z # First class dates? Why not?

[database]
server = "192.168.1.1"
ports = [ 8001, 8001, 8002 ]
connection_max = 5000
enabled = true

[servers]

  # You can indent as you please. Tabs or spaces. TOML don't care.
  [servers.alpha]
  ip = "10.0.0.1"
  dc = "eqdc10"

  [servers.beta]
  ip = "10.0.0.2"
  dc = "eqdc10"

[clients]
data = [ ["gamma", "delta"], [1, 2] ] # just an update to make sure parsers support it

# Line breaks are OK when inside arrays
hosts = [
  "alpha",
  "omega"
]
```


#### After
```
# This is a TOML document. Boom.

title = "writing string must be enclosed in double quote."

[owner]
name = "Tom Preston-Werner"
organization = """Learn Git and GitHub
    without any code!"""
bio = "GitHub Cofounder & CEO\nLikes tater tots and beer."
dob = 2018-07-24T00:00:00Z # First class dates? Why not?

[database]
server = "192.168.1.1"
ports = [ 1081, 1082, 1083 ]
connection_max = 9999
enabled = true

[servers]

  # You can indent as you please. Tabs or spaces. TOML don't care.
  [servers.alpha]
  ip = "192.168.122.i"
  dc = "eqdc10"

  [servers.beta]
  ip = "192.168.122.i"
  dc = "eqdc10"

[clients]
data = [ ["gamma", "delta"], [1, 2] ] # just an update to make sure parsers support it

# Line breaks are OK when inside arrays
hosts = [
  "alpha",
  "omega"
]
```
