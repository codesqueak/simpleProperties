package simpleProperties

import (
	"bufio"
	"container/list"
	"encoding/json"
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

const (
	basePath        = "resources/application"
	bootstrapPath   = "resources/bootstrap"
	expressionRegex = "(\\$+\\{\\S+(:\\S+){0,1}})"
	nameRegex       = "\\$\\{(\\S+?){1}(:\\S+?){0,1}}"
)

var expression_matcher, _ = regexp.Compile(expressionRegex) // find all ${} expressions
var name_matcher, _ = regexp.Compile(nameRegex)             // find name & default values ${abc:xyz) -> abc, :xyz

// BootPropertyLoader load properties from the boostrap file(s)
func BootPropertyLoader(path string) func(*Properties) error {
	return func(p *Properties) error {
		if err := baseLoader(p, path); err != nil {
			return err
		}
		tempMap := p.bootKeyValueMap
		p.bootKeyValueMap = p.keyValueMap
		p.keyValueMap = tempMap
		return nil
	}
}

// GlobalPropertyLoader load properties from the application property file(s)
func GlobalPropertyLoader(path string) func(*Properties) error {
	return func(p *Properties) error {
		return baseLoader(p, path)
	}
}

// ProfilePropertyLoader load properties from the application_<profile> property file(s)
func ProfilePropertyLoader(path string) func(*Properties) error {
	return func(p *Properties) error {
		profileNames := strings.Split(p.keyValueMap["profile"], ",")
		if len(profileNames) > 0 {
			for _, profileName := range profileNames {
				if profileName != "" {
					name := strings.TrimSpace(profileName)
					if err := baseLoader(p, path+"_"+name); err != nil {
						return err
					}
				}
			}
		}
		return nil
	}
}

// LoadOSEnvironment load properties from the O/S environment
func LoadOSEnvironment() func(*Properties) error {
	return func(p *Properties) error {
		for _, kv := range os.Environ() {
			parts := strings.Split(kv, "=")
			if len(parts) < 2 {
				continue
			}
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			setKV(p, key, value)
		}
		return nil
	}
}

// LoadCLIParameters loads -key=value CLI parameters
func LoadCLIParameters() func(*Properties) error {
	return func(p *Properties) error {
		if len(os.Args) > 1 { // ignore run param
			var args = os.Args[1:]
			for _, argString := range args {
				arg := []rune(argString)
				if len(arg) >= 3 { // smallest is -k=
					if arg[0] == '-' {
						arg = arg[1:]
						before, after, found := strings.Cut(string(arg), "=")
						if found && len(before) > 0 { // allow blank values
							setKV(p, before, after)
						}
					}
				}
			}
		}
		return nil
	}
}

//
// utilities
//

// load properties from the file specified in the path.  Look for .yaml, .json and .properties files with the
// load order being .yaml least to .properties highest
func baseLoader(p *Properties, path string) error {
	dir, filename := filepath.Split(path)
	fsys := os.DirFS(dir)
	file, err := fsys.Open(filename + ".yaml")
	if err == nil {
		if err = loadYAML(p, file); err != nil {
			return err
		}
	}
	file, err = fsys.Open(filename + ".json")
	if err == nil {
		if err = loadJSON(p, file); err != nil {
			return err
		}
	}
	file, err = fsys.Open(filename + ".properties")
	if err == nil {
		if err = loadPropertiesFromFile(p, file); err != nil {
			return err
		}
	}
	return nil
}

// load properties from the specified .properties file
func loadPropertiesFromFile(p *Properties, file fs.File) error {
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "#") { // ignore comments
			continue
		}
		parts := strings.Split(line, "=")
		if len(parts) < 2 {
			return errors.New("Invalid property string: " + line)
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		if key == "" { // blank key, odd !
			continue
		}
		setKV(p, key, value)
	}
	return nil
}

// load properties from the specified .json file
func loadJSON(p *Properties, file fs.File) error {
	defer file.Close()
	byteValue, _ := io.ReadAll(file)
	var result map[string]interface{}
	err := json.Unmarshal([]byte(byteValue), &result)
	if err != nil {
		return errors.New("Invalid properties JSON file. Error: " + err.Error())
	}
	extractKVMap(p, result, "")
	return nil
}

// load properties from the specified .yaml file
func loadYAML(p *Properties, file fs.File) error {
	defer file.Close()
	byteValue, _ := io.ReadAll(file)
	result := make(map[string]interface{})
	err := yaml.Unmarshal(byteValue, &result)
	if err != nil {
		return errors.New("Invalid properties YAML file. Error: " + err.Error())
	}
	extractKVMap(p, result, "")
	return nil
}

// recursively work through a map of key -> value, and convert each found value into a string.
// if a value is a structure then the kv pairs are extracted by a recursive call, with each key
// being prefixed with key of the origination kv pair
//
// for example:
//
//	{ "level1": {
//	             "level2": "my value"
//	            }
//	}
//
// then the recursive prefix would be level1, giving a full property key of level1.level2
func extractKVMap(p *Properties, json map[string]interface{}, prefix string) {
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
				extractKVMap(p, valueType, name+".")
			default:
				value = "???"
			}
		} else {
			value = ""
		}
		if skip {
			skip = false
		} else {
			setKV(p, name, value)
		}
	}
}

// put a kev pair into the property map. leading / trailing white space is removed
func setKV(p *Properties, key string, value string) {
	k := strings.TrimSpace(key)
	v := strings.TrimSpace(value)
	if k != "" {
		if containsExpression(value) {
			// value with evaluation fields
			p.evalKeyValueMap[k] = v
			p.evalExprMap[k] = extractExpressions(value)

		} else {
			// simple value
			p.keyValueMap[k] = v
		}
	}
}

// extract all expression in the rhs property
func extractExpressions(value string) *list.List {
	expr := value
	var parts = expression_matcher.FindAllString(expr, -1)
	fmt.Println(parts)
	l := list.New()
	for _, v := range parts {
		name := name_matcher.FindStringSubmatch(v) // gives ${} then name, then :default if it exists
		defaultValue := name[2]
		if defaultValue != "" {
			defaultValue = defaultValue[1:]
		}
		l.PushBack(&exprParts{name[0], name[1], defaultValue})
	}
	return l
}

func containsExpression(s string) bool {
	r := expression_matcher.FindAllString(s, -1)
	return len(r) > 0
}

type exprParts struct {
	full         string // with  ${abc:xyz}, this is ${abc:xyz}
	name         string // with  ${abc:xyz}, this is abc
	defaultValue string // with  ${abc:xyz}, this is xyz
}
