package main

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/akiyosi/tomlwriter"
)

type tomlConfig struct {
	Title   string
	Owner   ownerInfo
	DB      database `toml:"database"`
	Servers map[string]server
	Clients clients
}

type ownerInfo struct {
	Name string
	Org  string `toml:"organization"`
	Bio  string
	DOB  time.Time
}

type database struct {
	Server  string
	Ports   []int
	ConnMax int `toml:"connection_max"`
	Enabled bool
}

type server struct {
	IP string
	DC string
}

type clients struct {
	Data  [][]interface{}
	Hosts []string
}

func main() {
	var config tomlConfig
	file := "./example.toml"
	if _, err := toml.DecodeFile(file, &config); err != nil {
		fmt.Println(err)
		return
	}

	b, _ := ioutil.ReadFile(file)

	// writing value of key in [owner] table:
	b, _ = tomlwriter.WriteValue(`"""`+"Learn Git and GitHub\n    without any code!"+`"""`, b, "owner", "organization", config.Owner.Org)

	// writing value in global key, tile is nil:
	b, _ = tomlwriter.WriteValue(`"writing string must be enclosed in double quote."`, b, nil, "title", config.Title)

	// wiriting date/time
	b, _ = tomlwriter.WriteValue("2018-07-24T00:00:00Z", b, "owner", "dob", config.Owner.DOB)

	// writing array
	b, _ = tomlwriter.WriteValue("[ 1081, 1082, 1083 ]", b, "database", "ports", config.DB.Ports)

	// writing interger
	b, _ = tomlwriter.WriteValue("9999", b, "database", "connection_max", config.DB.ConnMax)

	var i int
	for serverName, server := range config.Servers {
		i++
		b, _ = tomlwriter.WriteValue(`"192.168.122.`+`i"`, b, "servers."+serverName, "ip", server.IP)
	}

	file2 := "./example2.toml"
	_ = ioutil.WriteFile(file2, b, 0755)

	return
}
