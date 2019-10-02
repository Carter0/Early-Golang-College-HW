package main

import (
	"github.com/tidwall/gjson"
)

func main() {
	//Now an example using an external library called GJson.
	//This might be useful for the assignment.

	//Basic example
	const json = `{"name":{"first":"Janet","last":"Prichard"},"age":47}`
	value := gjson.Get(json, "name.last")
	println(value.String())

	//Basic example2
	const json2 = `{"name": {"first": "Tom", "last": "Anderson"},
		"age":37,
		"children": ["Sara","Alex","Jack"],
		"fav.movie": "Deer Hunter",
		"friends": [
		  {"first": "Dale", "last": "Murphy", "age": 44, "nets": ["ig", "fb", "tw"]},
		  {"first": "Roger", "last": "Craig", "age": 68, "nets": ["fb", "tw"]},
		  {"first": "Jane", "last": "Murphy", "age": 47, "nets": ["ig", "tw"]}
		]
	  }`

	value2 := gjson.GetMany(json2, "name.last", "age", "children.1", "fav.movie")
	for _, item := range value2 {
		println(item.String())
	}

	//Example3, this time doing something functionally
	//Don't really understand this right now :(
	const json3 = `{
		"programmers": [
		  {
			"firstName": "Janet", 
			"lastName": "McLaughlin", 
		  }, {
			"firstName": "Elliotte", 
			"lastName": "Hunter", 
		  }, {
			"firstName": "Jason", 
			"lastName": "Harold", 
		  }
		]
	  }`

	value3 := gjson.Get(json, "programmers")
	value3.ForEach(func(key, value gjson.Result) bool {
		println(value.String())
		return true // keep iterating
	})

	//Finally, there is a way to read straight from a byte array here.
	//Which may be useful for you in the project.

}
