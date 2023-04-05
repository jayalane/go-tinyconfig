// -*- tab-width: 2 -*-

// Package config implements a simple key-value config system that reads from disk with passed in defaults.
package config

import (
	"bufio"
	"io"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
)

// StringOrInt is a value with .StrVal, .IntVal and .BoolVal methods.
// Parsing happens at file read time not use time
type StringOrInt struct {
	StrVal     string
	IntVal     int
	BoolVal    bool
	Float64Val float64
}

// Config is a string-key to string/int/boolean value map.
type Config map[string]StringOrInt

// ReadConfig takes a default config and overrides it from ./config.txt file.
// if filename is "", then don't read a file
func ReadConfig(filename string, defaultConfig string) (Config, error) {

	config := Config{}

	strReader := strings.NewReader(defaultConfig)
	err := addConfigFromReader(strReader, &config)
	if err != nil {
		return nil, err
	}
	err = readConfigFromFile(filename, &config)
	if err != nil {
		return config, err
	}
	envVarStruct, ok := config["configEnvVar"]
	if !ok {
		return config, err
	}

	envFileName := getEnvVarName(envVarStruct.StrVal, filename)
	err = readConfigFromFile(envFileName, &config)

	return config, err
}

func getEnvVarName(envToken string, filename string) string {
	if !strings.Contains(filename, ".") {
		return filename + "_" + envToken
	}
	segments := strings.SplitN(filename, ".", 2)
	return segments[0] + "_" + envToken + "." + segments[1]
}

func readConfigFromFile(filename string, config *Config) error {

	if len(filename) == 0 {
		log.Println("No config file specified, using default")
		return nil
	}
	binaryFilename, err := os.Executable()
	if err != nil {
		panic(err)
	}
	filePath := path.Join(path.Dir(binaryFilename), filename)

	file, err := os.Open(filePath)
	// non-existing is not fatal

	if err != nil {
		if os.IsNotExist(err) {
			log.Println(filePath, "does not exist, using defaults")
			return nil
		}
		log.Println("Warning: can't open config file, using defaults,", filename, filePath, err.Error())
		return err
	}
	log.Println("Using config file", filePath)
	defer file.Close()
	fileReader := bufio.NewReader(file)
	err = addConfigFromReader(fileReader, config)
	if err != nil {
		log.Println("Warning: can't use config file, using defaults,", filename, filePath, err.Error())
		return err
	}
	return nil
}

// addConfigFromReader merges the parsed config from reader into Config
func addConfigFromReader(reader io.Reader, config *Config) error {

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		if err := scanner.Err(); err != nil {
			if err == io.EOF {
				break
			}
			log.Println("Error reading config", err)
			return err
		}
		if len(line) > 0 && line[:1] == "#" { // TODO  space space #
			continue
		}

		// check if the line has = sign
		// and process the line. Ignore the rest.
		slashSlash := strings.Index(line, "//")
		var l string
		if slashSlash >= 0 {
			l = line[:slashSlash]
		} else {
			l = line
		}
		if equal := strings.Index(l, "="); equal >= 0 {
			if key := strings.TrimSpace(l[:equal]); len(key) > 0 {
				value := ""
				if len(l) > equal {
					value = strings.TrimSpace(l[equal+1:])
				}
				// assign the config map

				//bool
				bool := true
				if value == "true" {
					bool = true
					value = "1"
				} else if value == "false" {
					bool = false
					value = "0"
				}
				// int
				num, err := strconv.Atoi(value)
				if err != nil {
					num = 0
				}

				//float64
				f64, err := strconv.ParseFloat(value, 64)
				if err != nil {
					f64 = 0.0
				}
				(*config)[key] = StringOrInt{value, num, bool, f64}
				log.Println("Setting config", key, "to", value)

			}
		}
	}
	return nil
}
