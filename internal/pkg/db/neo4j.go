package neo4j

import (
	"fmt"

	"github.com/neo4j/neo4j-go-driver/neo4j"
)

type DBConfig struct {
	URI      string
	UserName string
	Pwd      string
}

var config DBConfig
var driver neo4j.Driver

func LoadConfig(uri, username, password string) {
	config = DBConfig{
		URI:      uri,
		UserName: username,
		Pwd:      password,
	}

	d, err := neo4j.NewDriver(uri, neo4j.BasicAuth(username, password, ""))
	if err != nil {
		panic(fmt.Errorf("LoadConfig error: %v", err))
	}
	driver = d
}

func GetConnection() (session neo4j.Session, err error) {
	s, err := driver.NewSession(neo4j.SessionConfig{})
	return s, err
}
