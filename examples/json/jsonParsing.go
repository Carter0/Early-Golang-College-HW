package main

import (
	"encoding/json"
	"log"
	"os"
)

//Message is the testing data for json.
type Message struct {
	Name string
	Body string
	Time int64
}

func main() {

	//Encoding the json message.
	m := Message{"Alice", "Hello", 1294706395881547000}
	b, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}

	//Decoding the json message.
	var m2 Message
	err = json.Unmarshal(b, &m2)
	if err != nil {
		panic(err)
	}

	//println(m.Name)

	//An example of using json.NewDecoder to decode a stream of json.
	dec := json.NewDecoder(os.Stdin)
	enc := json.NewEncoder(os.Stdout)
	for {
		//Remember an empty interface is just a data type without an methods
		//The map function says that we have a hashmap of string keys for the value interface
		//The jsonelementname must be a string, and its value can be anything. so make map[string] = anything
		var v map[string]interface{}
		if err := dec.Decode(&v); err != nil {
			log.Println(err)
			return
		}
		//Look through the whole json message
		//and if the fieldName != "Name" delete the key:value pair
		for k := range v {
			if k != "Name" {
				delete(v, k)
			}
		}
		//encode the map back into a json object and print it out.
		if err := enc.Encode(&v); err != nil {
			log.Println(err)
		}
	}
}
