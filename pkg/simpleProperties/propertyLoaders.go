package simpleProperties

import (
	"bufio"
	"container/list"
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

const basePath = "resources/application"
const bootstrapPath = "resources/bootstrap"

var expression_matcher, _ = regexp.Compile("(\\$+\\{\\S+(:\\S+){0,1}})") // find all ${} expressions
var name_matcher, _ = regexp.Compile("\\$\\{(\\S+?){1}(:\\S+?){0,1}}")   // find name & default values ${abc:xyz) -> abc, :xyz

// BootPropertyLoader load properties from the the boostrap file(s)
func BootPropertyLoader(path string) func(*Properties) {
	return func(p *Properties) {
		log.Println("-- boot environment --")
		baseLoader(p, path)
		tempMap := p.bootKeyValueMap
		p.bootKeyValueMap = p.keyValueMap
		p.keyValueMap = tempMap
	}
}

// GlobalPropertyLoader load properties from the the application property file(s)
func GlobalPropertyLoader(path string) func(*Properties) {
	return func(p *Properties) {
		log.Println("-- global environment --")
		baseLoader(p, path)
	}
}

// ProfilePropertyLoader load properties from the the application_<profile> property file(s)
func ProfilePropertyLoader(path string) func(*Properties) {
	return func(p *Properties) {
		log.Println("-- profile environment --")
		profileNames := strings.Split(p.keyValueMap["profile"], ",")
		if len(profileNames) > 0 {
			for _, profileName := range profileNames {
				if profileName != "" {
					name := strings.Trim(profileName, " \t")
					baseLoader(p, path+"_"+name)
				}
			}
		}
	}
}

// LoadOSEnvironment load properties from the the O/S environment
func LoadOSEnvironment() func(*Properties) {
	return func(p *Properties) {
		log.Println("-- os environment --")
		for _, kv := range os.Environ() {
			parts := strings.Split(kv, "=")
			if len(parts) < 2 {
				log.Fatalf("Invalid environment variable string: %s", kv)
			}
			setKV(p, parts[0], parts[1])
		}
	}
}

//
// utilities
//

// load properties from the file specified in the path.  Look for .yaml, .json and .properties files with the
// load order being .yaml least to .properties most
func baseLoader(p *Properties, path string) {
	log.Println("Loading from: " + path)
	dir, filename := filepath.Split(path)
	fsys := os.DirFS(dir)
	file, err := fsys.Open(filename + ".yaml")
	if err == nil {
		loadYAML(p, file)
	}
	file, err = fsys.Open(filename + ".json")
	if err == nil {
		loadJSON(p, file)
	}
	file, err = fsys.Open(filename + ".properties")
	if err == nil {
		loadPropertiesFromFile(p, file)
	}
}

// load properties from the specified .properties file
func loadPropertiesFromFile(p *Properties, file fs.File) {
	log.Println("-- properties --")
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := scanner.Text()
		strings.Trim(line, " \t")
		if line == "" {
			continue
		}
		parts := strings.Split(line, "=")
		if len(parts) < 2 {
			log.Fatalf("Invalid property string: %s", line)
		}
		setKV(p, parts[0], parts[1])
	}
}

// load properties from the specified .json file
func loadJSON(p *Properties, file fs.File) {
	log.Println("-- json --")
	defer file.Close()
	byteValue, _ := io.ReadAll(file)
	var result map[string]interface{}
	err := json.Unmarshal([]byte(byteValue), &result)
	if err != nil {
		log.Fatalf("Invalid properties JSON file. Error %s", err)
	}
	extractKVMap(p, result, "")
}

// load properties from the specified .yaml file
func loadYAML(p *Properties, file fs.File) {
	log.Println("-- yaml --")
	defer file.Close()
	byteValue, _ := io.ReadAll(file)
	result := make(map[string]interface{})
	err := yaml.Unmarshal(byteValue, &result)
	if err != nil {
		log.Fatalf("Invalid properties YAML file. Error %s", err)
	}
	extractKVMap(p, result, "")
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
				log.Println("!type", valueType)
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
	k := strings.Trim(key, " \t")
	v := strings.Trim(value, " \t")
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
		fmt.Println("name=", name)
		defaultValue := name[2]
		if defaultValue != "" {
			defaultValue = defaultValue[1:]
		}
		l.PushBack(&exprParts{name[0], name[1], defaultValue})
		log.Println("---")
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
