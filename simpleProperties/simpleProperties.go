package simpleProperties

var internalProperties = &Properties{nil, make(map[string]string, 32), nil}

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
// note: If mixed properties, JSON and YAML files are present, all will be read, but .yaml overridden by .json override by .properties

func init() {
	println("-- load bootstrap --")

	operations := []func(p *Properties){}
	// loaders
	operations = append(operations, GlobalPropertyLoader(basePath))
	operations = append(operations, ProfilePropertyLoader(basePath))
	operations = append(operations, LoadOSEnvironment())
	//
	// evaluators
	operations = append(operations, DefaultEvaluator(""))
	//
	// boot properties
	f := GlobalPropertyLoader(bootstrapPath) // loads into keyValueMap. Will need to move to bootKeyValueMap
	f(internalProperties)
	internalProperties = &Properties{internalProperties.keyValueMap, make(map[string]string, 128), operations}
}

func NewProperties() *Properties {
	return &Properties{internalProperties.bootKeyValueMap, internalProperties.keyValueMap, internalProperties.operations}
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

type Properties struct {
	bootKeyValueMap map[string]string
	keyValueMap     map[string]string
	operations      []func(p *Properties)
}
