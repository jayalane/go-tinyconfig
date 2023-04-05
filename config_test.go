// -*- tab-width: 2 -*-

package config

import "testing"
import "strings"

var defaultConfig = `#
numConfig=4
numList=1,2,3
float= 6.4 // ok
boolConfig=true
stringConfig=Chris Lane
testSlashSlash = 3 // number should be 3 not 0
# test comment
`

// TODO some kind of testing of the actual file names and so on.  

func TestSortTime(t *testing.T) {
	config := Config{}

	strReader := strings.NewReader(defaultConfig)
	err := addConfigFromReader(strReader, &config)
	if err != nil {
		t.Log("Couldn't parse default config", err)
		t.Fail()
	}
	if 3 != config["testSlashSlash"].IntVal {
		t.Log("Expected 3 got", config["numConfig"].IntVal)
		t.Fail()
	}
	if 4 != config["numConfig"].IntVal {
		t.Log("Expected 4 got", config["numConfig"].IntVal)
		t.Fail()
	}
	if !config["boolConfig"].BoolVal {
		t.Log("Expected true got", config["boolConfig"].BoolVal)
		t.Fail()
	}
	if 6.4 != config["float"].Float64Val { // is == test ok?
		t.Log("Expected 6.4 got", config["float"].Float64Val)
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
