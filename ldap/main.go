package main

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-ldap/ldap/v3"
	"io/ioutil"
)

type Config struct {
	Username     string
	Password     string
	LdapUrl      string
	LdapPort     string
	BindDN       string
	BindPass     string
	UidAttribute string
	UserBaseDN   string
	Attributes   []string
}

var config *Config

func init() {
	cfg, err := readConfig("config.json")
	if err != nil {
		panic(err)
	}
	config = cfg
}

func main() {
	userDn, err := testLdap()
	if err != nil {
		panic(err)
	}
	println("nice to see you ldap; " + userDn)
}

func readConfig(path string) (*Config, error) {
	var config Config
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(b, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

func getLdapConn() (*ldap.Conn, error) {
	url := fmt.Sprintf("%s:%s", config.LdapUrl, config.LdapPort)
	conn, err := ldap.DialURL(url)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	// Reconnect with TLS
	err = conn.StartTLS(&tls.Config{InsecureSkipVerify: true})
	if err != nil {
		return nil, err
	}

	err = conn.Bind(config.BindDN, config.BindPass)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func testLdap() (string, error) {
	conn, err := getLdapConn()
	if err != nil {
		return "", err
	}

	filter := fmt.Sprintf("(%s=%s)", config.UidAttribute, config.Username)
	searchRequest := ldap.NewSearchRequest(
		config.UserBaseDN,
		ldap.ScopeWholeSubtree, ldap.DerefInSearching, 0, 0, false,
		filter,
		config.Attributes,
		nil,
	)

	sr, err := conn.Search(searchRequest)
	if err != nil {
		return "", err
	}
	if len(sr.Entries) != 1 {
		return "", errors.New("user does not exist or too many users")
	}
	userdn := sr.Entries[0].DN
	err = conn.Bind(userdn, config.Password)
	if err != nil {
		return "", err
	}
	if err = conn.Bind(config.BindDN, config.BindPass); err != nil {
		return "", err
	}
	return userdn, nil
}
