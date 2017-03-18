package core

import (
	"io/ioutil"
	"encoding/json"
)

var unitsDir = "./data/units/"

type Unit struct {
	ID string
	Name string
	Version string
	Creator string
	Command string
	Setting map[string]string
	Inputs map[string]string
	Outputs map[string]string
}

// load unit json file
func unitLoad(id string) (Unit, error) {
	var unit Unit

	file, err := ioutil.ReadFile(unitsDir + id + "/unit.json")

	if err != nil {
		return unit, err
	}

	err = json.Unmarshal(file, &unit)
	if err != nil {
		return unit, err
	}

	return unit, nil
}