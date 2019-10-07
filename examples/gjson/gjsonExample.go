package main

import (
	"encoding/json"

	"github.com/tidwall/gjson"
)

// message represents a message from a neighbor to the router.
type message struct {
	Msg  interface{} `json:"msg"`
	Src  string      `json:"src"`
	Dst  string      `json:"dst"`
	Type string      `json:"type"`
}

func main() {
	// //Now an example using an external library called GJson.
	// //This might be useful for the assignment.

	// //Basic example
	// const json = `{"name":{"first":"Janet","last":"Prichard"},"age":47}`
	// value := gjson.Get(json, "name.last")
	// println(value.String())

	// //Basic example2
	// const json2 = `{"name": {"first": "Tom", "last": "Anderson"},
	// 	"age":37,
	// 	"children": ["Sara","Alex","Jack"],
	// 	"fav.movie": "Deer Hunter",
	// 	"friends": [
	// 	  {"first": "Dale", "last": "Murphy", "age": 44, "nets": ["ig", "fb", "tw"]},
	// 	  {"first": "Roger", "last": "Craig", "age": 68, "nets": ["fb", "tw"]},
	// 	  {"first": "Jane", "last": "Murphy", "age": 47, "nets": ["ig", "tw"]}
	// 	]
	//   }`

	// value2 := gjson.GetMany(json2, "name.last", "age", "children.1", "fav.movie")
	// for _, item := range value2 {
	// 	println(item.String())
	// }

	// //Example3, this time doing something functionally
	// //Don't really understand this right now :(
	// const json3 = `{
	// 	"programmers": [
	// 	  {
	// 		"firstName": "Janet",
	// 		"lastName": "McLaughlin",
	// 	  }, {
	// 		"firstName": "Elliotte",
	// 		"lastName": "Hunter",
	// 	  }, {
	// 		"firstName": "Jason",
	// 		"lastName": "Harold",
	// 	  }
	// 	]
	//   }`

	// value3 := gjson.Get(json, "programmers")
	// value3.ForEach(func(key, value gjson.Result) bool {
	// 	println(value.String())
	// 	return true // keep iterating
	// })

	// //Finally, there is a way to read straight from a byte array here.
	// //Which may be useful for you in the project.

	const json4 = `{
		"type": "update",
		"src": "192.168.0.2",
		"dst": "192.168.0.1",
		"msg": {
			"network": "192.168.1.0",
			"netmask": "255.255.255.0",
			"localpref": "100",
			"selfOrigin": "True",
			"ASPath": ["1"],
			"origin": "EGP"
		}
	}`

	temp, err := json.Marshal(json4)
	if err != nil {
		panic(err)
	}
	tempType := gjson.GetBytes(temp, "type").String()

	if tempType == "" {
		//It is always empty and I don't know why.
		println("its empty")
	}
	println(tempType)

}
