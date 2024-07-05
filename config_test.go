// -*- tab-width: 2 -*-

package config

import (
	"strings"
	"testing"
)

var defaultConfig = `#
numConfig=4
numList=1,2,3
float= 6.4 // ok
boolConfig=true
stringConfig=Chris Lane
testSlashSlash = 3 // number should be 3 not 0
# test comment
`

var (
	envVarKeyPrefix = "TINYCONFIG_OVERRIDE_"
	envVarKey       = "sigma_3"
	envVarValue     = "3.2"
)

// later some kind of testing of the actual file names and so on.

func TestSortTime(t *testing.T) {
	config := Config{}
	strReader := strings.NewReader(defaultConfig)

	err := addConfigFromReader(strReader, &config)
	if err != nil {
		t.Log("Couldn't parse default config", err)
		t.Fail()
	}

	if config["testSlashSlash"].IntVal != 3 {
		t.Log("Expected 3 got", config["numConfig"].IntVal)
		t.Fail()
	}

	if config["numConfig"].IntVal != 4 {
		t.Log("Expected 4 got", config["numConfig"].IntVal)
		t.Fail()
	}

	if !config["boolConfig"].BoolVal {
		t.Log("Expected true got", config["boolConfig"].BoolVal)
		t.Fail()
	}

	if config["float"].Float64Val != 6.4 { // is == test ok?
		t.Log("Expected 6.4 got", config["float"].Float64Val)
		t.Fail()
	}

	if config["stringConfig"].StrVal != "Chris Lane" {
		t.Log("Expected 3 got", config["stringConfig"].StrVal)
		t.Fail()
	}

	if len(strings.Split(config["numList"].StrVal, ",")) != 3 {
		t.Log("Expected 3 1,2,3 got", strings.Split(config["numList"].StrVal, ","))
		t.Fail()
	}

	if config[envVarKey].Float64Val != 0 {
		t.Log("Env var override beforeenv set",
			config[envVarKey].Float64Val,
			"should be zero")
		t.Fail()
	}

	t.Setenv(envVarKeyPrefix+envVarKey, envVarValue)
	overrideConfigFromEnv(&config)

	if config[envVarKey].Float64Val != 3.2 {
		t.Log("Env var override beforeenv set",
			config[envVarKey].Float64Val,
			"should be 3.2")
		t.Fail()
	}
}
