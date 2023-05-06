// -*- tab-width: 2 -*-

// Package config implements a simple key-value config system that reads
// from disk with passed in defaults.

// I used this file to test GPT-4 but it was not totally successful.

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

// Config is a string-key to string/int/boolean value map.
type Config map[string]StringOrInt

// StringOrInt is a value with .StrVal, .IntVal and .BoolVal methods.
// Parsing happens at file read time not use time
type StringOrInt struct {
	StrVal     string
	IntVal     int
	BoolVal    bool
	Float64Val float64
}

func getEnvVarFilename(envToken string, filename string) string {
	if !strings.Contains(filename, ".") {
		return filename + "_" + envToken
	}
	segments := strings.SplitN(filename, ".", 2)
	return segments[0] + "_" + envToken + "." + segments[1]
}

// ReadConfig takes a default config and overrides it from
// ./config.txt file.  if filename is "", then don't read a file; if
// config.txt has a "configEnvVar" setting, it will check the
// environment for a key equal to the vlaue of that setting; if
// present, it will use the value of that environment variable to make
// another file name, e.g. config_PROD.txt and open that file and add
// it to the config.
func ReadConfig(filename string, defaultConfig string) (Config, error) {

	config := Config{}

	strReader := strings.NewReader(defaultConfig)
	err := addConfigFromReader(strReader, &config)
	if err != nil {
		return nil, err
	}

	if len(filename) == 0 {
		log.Println("No config file specified, using default")
		return config, nil
	}

	err = readConfigFile(filename, &config)
	if err != nil {
		log.Println("Warning: can't use config file, using defaults,",
			filename, err.Error())
		return config, err
	}

	// first featyure: Check for the "configEnvVar" key in the config
	envVarKey, ok := config["configEnvVar"]
	if !ok {
		return config, nil
	}
	log.Println("Found configEnvVar", envVarKey.StrVal)
	envVarName := envVarKey.StrVal
	envVarValue := os.Getenv(envVarName)
	log.Println("Found configEnvVar value", envVarValue)
	if envVarValue != "" {
		envfilename := getEnvVarFilename(envVarValue, filename)
		err = readConfigFile(envfilename, &config)
		if err != nil {
			log.Println("Warning: can't use second config file, using defaults,",
				envfilename, err.Error())
			return config, err
		}
	}

	// second feature
	overrideConfigFromEnv(&config)

	return config, nil
}

func readConfigFile(filename string, config *Config) error {
	binaryFilename, err := os.Executable()
	if err != nil {
		panic(err)
	}

	filePath := path.Join(path.Dir(binaryFilename), filename)

	file, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			log.Println(filePath, "does not exist, using defaults")
			return nil
		}
		log.Println("Warning: can't open config file, using defaults,", filename, filePath, err.Error())
		return err
	}

	log.Println("Reading config file", filePath)
	defer file.Close()

	fileReader := bufio.NewReader(file)

	err = addConfigFromReader(fileReader, config)
	if err != nil {
		log.Println("Warning: can't use config file, using defaults,", filename, filePath, err.Error())

		return err
	}

	return nil
}

// overrideConfigFromEnv makes a little config
// file from the env and returns a reader to it.
func overrideConfigFromEnv(config *Config) {
	var res strings.Builder

	e := os.Environ()

	for _, v := range e {
		if strings.HasPrefix(v, "TINYCONFIG_OVERRIDE_") {
			configK := v[len("TINYCONFIG_OVERRIDE_"):]

			log.Println("Env override Setting config", configK)
			res.WriteString(configK)
		}
	}
	strReader := strings.NewReader(res.String())

	_ = addConfigFromReader(strReader, config)
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
