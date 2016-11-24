package settings

import (
	"encoding/json"
	"io/ioutil"
)

type Settings struct {
	Root  string `json: "root"`
	Redis redis  `json: "redis"`
	Paths []path `json: "paths"`
}

type redis struct {
	DB   int    `json: db`
	Addr string `json: "addr"`
	Pass string `json: "password"`
}

type path struct {
	Name string `json: "name"`
	Path string `json: "path"`
}

func Load(filename string) (*Settings, error) {

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var settings Settings
	err = json.Unmarshal(data, &settings)
	if err != nil {
		return nil, err
	}

	return &settings, nil
}
