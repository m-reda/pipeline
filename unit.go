package main

import (
	"encoding/json"
	"io/ioutil"
)

var unitsDir = "./data/units/"

// Unit type represents a unit information
type Unit struct {
	ID      string
	Name    string
	Version string
	Creator string
	Group   string
	Command string
	Setting map[string]map[string]string
	Inputs  map[string]string
	Outputs map[string]string
}

// load unit json file
func unitLoad(id string) (Unit, error) {
	var unit Unit

	// read the unit file
	file, err := ioutil.ReadFile(unitsDir + id + "/unit.json")
	if err != nil {
		return unit, err
	}

	// decode the unit json
	err = json.Unmarshal(file, &unit)
	if err != nil {
		return unit, err
	}

	return unit, nil
}

func loadAllUnits() []Unit {
	var units []Unit

	// load all the units directories
	files, _ := ioutil.ReadDir(unitsDir)

	for _, f := range files {
		// ignore bin directory
		if f.Name() == "bin" {
			continue
		}

		// load the unit and add it to the units list
		if unit, err := unitLoad(f.Name()); err == nil {
			units = append(units, unit)
		}
	}

	return units
}
