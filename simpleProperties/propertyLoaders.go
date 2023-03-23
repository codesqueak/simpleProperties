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

const basePath = "resource/application"
const bootstrapPath = "resource/bootstrap"

var expression_matcher, _ = regexp.Compile("(\\$+\\{\\S+(:\\S+){0,1}})")
var name_matcher, _ = regexp.Compile("\\$\\{(\\S+?){1}(:\\S+?){0,1}}")

func BootPropertyLoader(path string) func(*Properties) {
	return func(p *Properties) {
		println("-- boot environment --")
		baseLoader(p, path)
		tempMap := p.bootKeyValueMap
		p.bootKeyValueMap = p.keyValueMap
		p.keyValueMap = tempMap
	}
}
func GlobalPropertyLoader(path string) func(*Properties) {
	return func(p *Properties) {
		println("-- global environment --")
		baseLoader(p, path)
	}
}

func ProfilePropertyLoader(path string) func(*Properties) {
	return func(p *Properties) {
		println("-- profile environment --")
		profileNames := strings.Split(p.keyValueMap["profile"], ",")
		if len(profileNames) > 0 {
			for _, name := range profileNames {
				if name != "" {
					name := strings.Trim(name, " ")
					baseLoader(p, path+"_"+name)
				}
			}
		}
	}
}

func LoadOSEnvironment() func(*Properties) {
	return func(p *Properties) {
		println("-- os environment --")
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

func baseLoader(p *Properties, path string) {
	println("Loading from: " + path)
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

func loadPropertiesFromFile(p *Properties, file fs.File) {
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
		setKV(p, parts[0], parts[1])
	}
}

func loadJSON(p *Properties, file fs.File) {
	println("-- json --")
	defer file.Close()
	byteValue, _ := io.ReadAll(file)
	var result map[string]interface{}
	err := json.Unmarshal([]byte(byteValue), &result)
	if err != nil {
		log.Fatalf("Invalid properties JSON file. Error %s", err)
	}
	extractKVMap(p, result, "")
}

func loadYAML(p *Properties, file fs.File) {
	println("-- yaml --")
	defer file.Close()
	byteValue, _ := io.ReadAll(file)
	result := make(map[string]interface{})
	err := yaml.Unmarshal(byteValue, &result)
	if err != nil {
		log.Fatalf("Invalid properties YAML file. Error %s", err)
	}
	extractKVMap(p, result, "")
}

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
				println("!type", valueType)
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

func setKV(p *Properties, key string, value string) {
	k := strings.Trim(key, " ")
	v := strings.Trim(value, " ")
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
		println("---")
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
