package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

var config *Config

type Config struct {
	Id1     string
	Id2     string
	Token   string
	Channel string
}

type SlackMessage struct {
	Channel  string `json:"channel"`
	Text     string `json:"text"`
	Username string `json:"username"`
}

func init() {
	cfg, err := readConfig("config.json")
	if err != nil {
		panic(err)
	}
	config = cfg
}

func main() {
	msg := &SlackMessage{
		Channel:  config.Channel,
		Text:     "Hello slack",
		Username: "Gopher",
	}
	if err := send(msg); err != nil {
		panic(err)
	}
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

func send(msg *SlackMessage) error {
	req, err := createRequest(msg)
	if err != nil {
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("[%d-error] %s\n", resp.StatusCode, resp.Status)
		return errors.New(resp.Status)
	}
	return nil
}

func createRequest(msg *SlackMessage) (*http.Request, error) {
	b, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}
	slackUrl := fmt.Sprintf("https://hooks.slack.com/services/%s/%s/%s", config.Id1, config.Id2, config.Token)

	buff := bytes.NewBuffer(b)
	req, err := http.NewRequest("POST", slackUrl, buff)
	req.Header.Add("Content-type", "application/json")

	return req, nil
}
