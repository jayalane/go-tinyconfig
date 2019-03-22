// -*- tab-width: 2 -*-

package config

import "testing"
import "strings"

var defaultConfig = `#
numConfig=4
numList=1,2,3
boolConfig=true
stringConfig=Chris Lane
# test comment
`

func TestSortTime(t *testing.T) {
	config := Config{}

	strReader := strings.NewReader(defaultConfig)
	err := addConfigFromReader(strReader, &config)
	if err != nil {
		t.Log("Couldn't parse default config", err)
		t.Fail()
	}
	if 4 != config["numConfig"].IntVal {
		t.Log("Expected 3 got", config["numConfig"].IntVal)
		t.Fail()
	}
	if !config["boolConfig"].BoolVal {
		t.Log("Expected true got", config["boolConfig"].BoolVal)
		t.Fail()
	}
	if "Chris Lane" != config["stringConfig"].StrVal {
		t.Log("Expected 3 got", config["stringConfig"].StrVal)
		t.Fail()
	}
	if 3 != len(strings.Split(config["numList"].StrVal, ",")) {
		t.Log("Expected 3 1,2,3 got", strings.Split(config["numList"].StrVal, ","))
		t.Fail()
	}
}
