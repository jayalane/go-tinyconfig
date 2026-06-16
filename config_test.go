// -*- tab-width: 2 -*-

package config

import (
	"os"
	"path"
	"strings"
	"testing"
)

var defaultConfig = `#
numConfig=4
numList=1,2,3
float= 6.4 // ok
boolConfig=true
boolConfig2 = 1
falseBoolConfig = false
falseBoolConfig2 = 0  # testing midline hash
stringConfig=Chris Lane
testSlashSlash = 3 // number should be 3 not 0
sigma_3 = 4.0
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

	if !config["boolConfig2"].BoolVal {
		t.Log("Expected true got", config["boolConfig2"].BoolVal)
		t.Fail()
	}

	if config["falseBoolConfig"].BoolVal {
		t.Log("Expected false got", config["falseBoolConfig"].BoolVal)
		t.Fail()
	}

	if config["falseBoolConfig2"].BoolVal {
		t.Log("Expected false got", config["falseBoolConfig2"].BoolVal)
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

	if config[envVarKey].Float64Val != 4.0 {
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

// TestPerBinaryConfigFallback verifies that when config.txt is absent,
// ReadConfig falls back to "${binary_name}_config.txt" sitting next to
// the executable.
func TestPerBinaryConfigFallback(t *testing.T) {
	binaryFilename, err := os.Executable()
	if err != nil {
		t.Fatal("can't find executable", err)
	}

	dir := path.Dir(binaryFilename)

	// Skip if a real config.txt happens to live beside the test binary,
	// since that would be used instead of the fallback.
	if _, err := os.Stat(path.Join(dir, "config.txt")); err == nil {
		t.Skip("config.txt exists beside test binary; skipping fallback test")
	}

	fallback := path.Join(dir, path.Base(binaryFilename)+"_config.txt")

	err = os.WriteFile(fallback, []byte("numConfig=99\n"), 0o600)
	if err != nil {
		t.Fatal("can't write fallback config", err)
	}

	defer os.Remove(fallback)

	config, err := ReadConfig("config.txt", defaultConfig)
	if err != nil {
		t.Fatal("ReadConfig failed", err)
	}

	if config["numConfig"].IntVal != 99 {
		t.Log("Expected fallback override 99 got", config["numConfig"].IntVal)
		t.Fail()
	}
}
