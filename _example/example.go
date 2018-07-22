package main

import (
	"fmt"
	"time"
	"io/ioutil"
	"strconv"

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
	ConnMax float64 `toml:"connection_max"`
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
	file := "./_example/example.toml"
	if _, err := toml.DecodeFile(file, &config); err != nil {
		fmt.Println(err)
		return
	}

	// fmt.Printf("Title: %s\n", config.Title)
	// fmt.Printf("Owner: %s (%s, %s), Born: %s\n",
	// 	config.Owner.Name, config.Owner.Org, config.Owner.Bio,
	// 	config.Owner.DOB)
	fmt.Printf("Database: %s %v (Max conn. %d), Enabled? %v\n",
		config.DB.Server, config.DB.Ports, config.DB.ConnMax,
		config.DB.Enabled)
	// for serverName, server := range config.Servers {
	// 	fmt.Printf("Server: %s (%s, %s)\n", serverName, server.IP, server.DC)
	// }
	// fmt.Printf("Client data: %v\n", config.Clients.Data)
	// fmt.Printf("Client hosts: %v\n", config.Clients.Hosts)

	input, _ := ioutil.ReadFile(file)
	b, _ := tomlwriter.WriteValue("fizz", input, "database", "ports", config.DB.Ports)
    _ = ioutil.WriteFile(file, b, 0755)
	var a float64
	a, _ = strconv.ParseFloat("1e4", 32)
	fmt.Println(fmt.Sprintf("float test: %e", a))

	return

}
