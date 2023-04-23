package simpleProperties

import "container/list"
import "log"

var internalProperties = &Properties{
	make(map[string]string, 32),
	make(map[string]string, 32),
	make(map[string]string, 32),
	make(map[string]*list.List, 32),
	nil}

// load simpleProperties with this precedence
//
// 1. boot properties
// 2. system simpleProperties
// 3. default simpleProperties files
// 4. profile simpleProperties files
// 5. profile properties ( x n)
// 6. command line arguments
// 7. evaluate references
//
// note: If mixed properties, JSON and YAML files are present, all will be read, but .yaml overridden by .json overridden by .properties

func init() {
	log.Println("-- load bootstrap --")

	operations := []func(p *Properties){}
	// loaders
	operations = append(operations, GlobalPropertyLoader(basePath))
	operations = append(operations, ProfilePropertyLoader(basePath))
	//	operations = append(operations, LoadOSEnvironment())
	//
	// evaluators
	operations = append(operations, BasicEvaluator())
	// operations = append(operations, DefaultEvaluator())
	//
	internalProperties.operations = operations
	// boot properties
	f := BootPropertyLoader(bootstrapPath)
	f(internalProperties)
}

// DefaultProperties create a default properties structure. This will contain the bootstrap properties and default operations
// to load the properties via the default operations, call the Load method
func DefaultProperties() *Properties {
	return &Properties{copyKV(internalProperties.bootKeyValueMap),
		make(map[string]string, 32),
		make(map[string]string, 32),
		make(map[string]*list.List, 32),
		copyOps(internalProperties.operations),
	}
}

// EmptyProperties create a blank properties structure with no data or operations
func EmptyProperties() *Properties {
	return &Properties{make(map[string]string, 32),
		make(map[string]string, 32),
		make(map[string]string, 32),
		make(map[string]*list.List, 32),
		nil,
	}
}

// Load execute the list operations for property loading
func (p *Properties) Load() {
	for _, f := range p.operations {
		f(p)
	}
}

// GetBootProperty get a bootstrap property (if it exists)
func (p *Properties) GetBootProperty(key string) string {
	if key == "" {
		return ""
	} else {
		return p.bootKeyValueMap[key]
	}
}

// GetProperty get a global property (if it exists). will fall back to boostrap properties if not
// held in the global property map
func (p *Properties) GetProperty(key string) string {
	if key == "" {
		return ""
	} else {
		if p.keyValueMap != nil {
			v := p.keyValueMap[key]
			if v == "" {
				return p.bootKeyValueMap[key]
			} else {
				return v
			}
		} else {
			return p.bootKeyValueMap[key]
		}
	}
}

// GetEvalProperty get a value from the map of evaluated properties
func (p *Properties) GetEvalProperty(key string) string {
	if key == "" {
		return ""
	} else {
		if p.evalKeyValueMap != nil {
			return p.evalKeyValueMap[key]
		}
		return ""
	}
}

// GetExprProperty get the list of expression data associated with a value key
func (p *Properties) GetExprProperty(key string) *list.List {
	if key == "" {
		return nil
	} else {
		if p.evalExprMap != nil {
			return p.evalExprMap[key]
		}
		return nil
	}
}

func (p *Properties) GetBootKeys() []string {
	keys := []string{}
	for k := range p.bootKeyValueMap {
		keys = append(keys, k)
	}
	return keys
}

func (p *Properties) GetKeys() []string {
	keys := []string{}
	for k := range p.keyValueMap {
		keys = append(keys, k)
	}
	return keys
}

func (p *Properties) GetEvalKeys() []string {
	keys := []string{}
	for k := range p.evalKeyValueMap {
		keys = append(keys, k)
	}
	return keys
}

func copyKV(in map[string]string) map[string]string {
	c := make(map[string]string)
	for k, v := range in {
		c[k] = v
	}
	return c
}

func copyOps(in []func(p *Properties)) []func(p *Properties) {
	c := make([]func(p *Properties), len(in))
	for i, v := range in {
		c[i] = v
	}
	return c
}

type Properties struct {
	bootKeyValueMap map[string]string
	keyValueMap     map[string]string
	evalKeyValueMap map[string]string
	evalExprMap     map[string]*list.List
	operations      []func(p *Properties)
}
