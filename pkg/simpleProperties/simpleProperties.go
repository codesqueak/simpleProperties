package simpleProperties

import "container/list"

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
	println("-- load bootstrap --")

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
	// boot properties
	f := BootPropertyLoader(bootstrapPath)
	f(internalProperties)
}

func NewProperties() *Properties {
	return &Properties{internalProperties.bootKeyValueMap,
		internalProperties.keyValueMap,
		internalProperties.evalKeyValueMap,
		internalProperties.evalExprMap,
		internalProperties.operations}
}

func (p *Properties) Load() {
	for _, f := range p.operations {
		f(p)
	}
}

func (p *Properties) GetBootProperty(key string) string {
	if key == "" {
		return ""
	} else {
		return p.bootKeyValueMap[key]
	}
}

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

type Properties struct {
	bootKeyValueMap map[string]string
	keyValueMap     map[string]string
	evalKeyValueMap map[string]string
	evalExprMap     map[string]*list.List
	operations      []func(p *Properties)
}
