package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUnitLoad(t *testing.T) {
	_, err := unitLoad("fs_copy")
	assert.NoError(t, err)

	_, err = unitLoad("-")
	assert.Error(t, err)
}

func TestLoadAllUnits(t *testing.T) {
	if len(loadAllUnits()) == 0 {
		t.Fail()
	}
}
