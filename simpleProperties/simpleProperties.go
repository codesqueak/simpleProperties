package simpleProperties

import "strings"

var keyValueMap = make(map[string]string, 128)
var basePath = "resource/application"
var bootstrapPath = "resource/boostrap"

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
//note: If mixed properties, JSON and YAML files are present, all will be read, but .yaml overridden by .json override by .properties

func init() {
	println("-- init --")
	baseLoader(keyValueMap, bootstrapPath)
}

func loadEnvironment() {
	println("-- loadEnvironment --")

	baseLoader(keyValueMap, basePath)
	profileNames := strings.Split(keyValueMap["profile"], ",")
	if len(profileNames) > 0 {
		for _, name := range profileNames {
			if name != "" {
				name := strings.Trim(name, " ")
				baseLoader(keyValueMap, basePath+"_"+name)
			}
		}
	}
	loadOSEnvironment(keyValueMap)
}

type propertyLoader interface {
	baseLoader(f string) map[string]string
	loader() map[string]string
}
