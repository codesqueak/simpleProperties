package main

import (
	"propertySource/simpleProperties"
	_ "propertySource/simpleProperties"
)

func main() {
	println("-- start --")

	properties := simpleProperties.NewProperties()

	properties.Load()

	println("-- properties --")
	for _, k := range properties.GetKeys() {
		println(k, "-!->", properties.GetProperty(k))
	}

	println("-- boot properties --")
	for _, k := range properties.GetBootKeys() {
		println(k, "-*->", properties.GetBootProperty(k))
	}

	println("-- end --")

}
