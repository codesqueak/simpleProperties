package main

import (
	"fmt"
	"github.com/codesqueak/simpleProperties/pkg/simpleProperties"
)

// just some code to run various test scenarios
func main() {
	println("-- start --")

	properties := simpleProperties.DefaultProperties()

	properties.Load()

	println("\n----------------------------------------------------")
	println("-- properties --")
	for _, k := range properties.GetKeys() {
		println(k, "-->", properties.GetProperty(k))
	}

	println("\n----------------------------------------------------")
	println("-- boot properties --")
	for _, k := range properties.GetBootKeys() {
		println(k, "-->", properties.GetBootProperty(k))
	}

	println("\n----------------------------------------------------")
	println("-- eval properties --")
	for _, k := range properties.GetEvalKeys() {
		println(k, "-->", properties.GetEvalProperty(k))
	}

	println("\n----------------------------------------------------")
	println("-- eval parts --")
	for _, k := range properties.GetEvalKeys() {
		l := properties.GetExprProperty(k)
		fmt.Println(k, "-->")
		element := l.Front()
		for element != nil {
			fmt.Println("   ...> ", element.Value)
			element = element.Next()
		}

	}

	println("-- end --")

	s1 := "abc123"
	s2 := s1

	s2 = s2 + "xyz"
	println(s1)
	println(s2)

}
