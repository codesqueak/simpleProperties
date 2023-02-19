package simpleProperties

import (
	"bufio"
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func baseLoader(m map[string]string, path string) {
	println("Loading from: " + path)
	dir, filename := filepath.Split(path)
	fsys := os.DirFS(dir)
	file, err := fsys.Open(filename + ".yaml")
	if err == nil {
		loadYAML(file, m)
	}
	file, err = fsys.Open(filename + ".json")
	if err == nil {
		loadJSON(file, m)
	}
	file, err = fsys.Open(filename + ".properties")
	if err == nil {
		loadProperties(file, m)
	}
}

func loadOSEnvironment(m map[string]string) {
	println("-- environment --")
	for _, kv := range os.Environ() {
		parts := strings.Split(kv, "=")
		if len(parts) < 2 {
			log.Fatalf("Invalid environment variable string: %s", kv)
		}
		var key = parts[0]
		var value = parts[1]
		key = strings.Trim(key, " ")
		value = strings.Trim(value, " ")
		println(key, "-->", value)
		m[key] = value
	}
}

func loadProperties(file fs.File, m map[string]string) {
	println("-- properties --")
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := scanner.Text()
		strings.Trim(line, " ")
		if line == "" {
			continue
		}
		parts := strings.Split(line, "=")
		if len(parts) < 2 {
			log.Fatalf("Invalid property string: %s", line)
		}
		var key = parts[0]
		var value = parts[1]
		strings.Trim(key, " ")
		strings.Trim(value, " ")
		println(key, "-->", value)
		m[key] = value
	}
}

func loadJSON(file fs.File, m map[string]string) {
	println("-- json --")
	defer file.Close()
	byteValue, _ := io.ReadAll(file)
	var result map[string]interface{}
	err := json.Unmarshal([]byte(byteValue), &result)
	if err != nil {
		log.Fatalf("Invalid properties JSON file. Error %s", err)
	}
	extractKVMap(m, result, "")
}

func loadYAML(file fs.File, m map[string]string) {
	println("-- yaml --")
	defer file.Close()
	byteValue, _ := io.ReadAll(file)
	result := make(map[string]interface{})
	err := yaml.Unmarshal(byteValue, &result)
	if err != nil {
		log.Fatalf("Invalid properties YAML file. Error %s", err)
	}
	extractKVMap(m, result, "")
}

func extractKVMap(m map[string]string, json map[string]interface{}, prefix string) {
	var value string
	skip := false
	for key := range json {
		rawValue := json[key]
		name := prefix + key
		if rawValue != nil {
			switch valueType := rawValue.(type) {
			case int:
				value = strconv.Itoa(valueType)
			case string:
				value = valueType
			case float64:
				value = fmt.Sprintf("%g", valueType)
			case bool:
				value = strconv.FormatBool(valueType)
			case map[string]interface{}:
				skip = true
				extractKVMap(m, valueType, name+".")
			default:
				value = "???"
				println("!type", valueType)
			}
		} else {
			value = ""
		}
		if skip {
			skip = false
		} else {
			println(name, "-->", value)
			m[name] = value
		}
	}
}
