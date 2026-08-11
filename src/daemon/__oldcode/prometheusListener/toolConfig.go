package main

import (
	"encoding/json"
	"fmt"
	cerr "github.com/jeanfrancoisgratton/customError/v3"
	hf "github.com/jeanfrancoisgratton/helperFunctions/v3"
	"os"
	"path/filepath"
	"strings"
)

// loadEnvironmentFile: load the JSON config file
func loadEnvironmentFile() (Config_s, *cerr.CustomError) {
	var payload Config_s

	rcFile := filepath.Join("/etc", "prometheusSDlistener", "prometheusSDlistener.json")
	_, err := os.Stat(rcFile)
	// We need to create the environment file if it does not exist
	if os.IsNotExist(err) {
		f := fmt.Sprintf("Configuration file %s not found", rcFile)
		panic(f)
	}

	jFile, err := os.ReadFile(rcFile)
	if err != nil {
		return Config_s{}, &cerr.CustomError{Title: "Error reading the file", Message: err.Error()}
	}
	err = json.Unmarshal(jFile, &payload)
	if err != nil {
		return Config_s{}, &cerr.CustomError{Title: "Error unmarshalling JSON", Message: err.Error()}
	} else {
		return payload, nil
	}
}

// SaveEnvironmentFile: save the config in a json file
func (cs Config_s) SaveEnvironmentFile() *cerr.CustomError {
	jStream, err := json.MarshalIndent(cs, "", "  ")
	if err != nil {
		return &cerr.CustomError{Title: err.Error(), Fatality: cerr.Fatal}
	}
	rcFile := filepath.Join("/etc", "prometheusSDlistener", "prometheusSDlistener.json")
	if err = os.WriteFile(rcFile, jStream, 0644); err != nil {
		return &cerr.CustomError{Title: "Unable to write JSON file", Message: err.Error(), Fatality: cerr.Fatal}
	}

	return nil
}

// setup: prompts user for configuration values
func setup() *cerr.CustomError {
	cfg := Config_s{}

	cp := hf.GetStringValFromPrompt("Enter the absolute path to your SSL certificate, without the extension: ")
	if !strings.HasPrefix(cp, "/") {
		return &cerr.CustomError{Title: "Invalid path", Message: "The path must be absolute"}
	}
	cp = strings.TrimSuffix(cp, ".crt")
	cp = strings.TrimSuffix(cp, ".key")
	cfg.Cert = cp + ".crt"
	cfg.Key = cp + ".key"
	cfg.Port = uint(hf.GetIntValFromPrompt("Enter the port the listener should listen on: "))
	if cfg.Port < 1025 || cfg.Port > 65535 {
		return &cerr.CustomError{Title: "Invalid port number", Message: "The port number must be between 1025 and 65535", Fatality: cerr.Fatal}
	}
	cfg.TargetDir = hf.GetStringValFromPrompt("Enter the path where the hostnames should be added/removed: ")
	if err := os.MkdirAll(cfg.TargetDir, os.ModePerm); err != nil {
		return &cerr.CustomError{Title: "Unable to create directory", Message: err.Error()}
	}
	return cfg.SaveEnvironmentFile()
}
