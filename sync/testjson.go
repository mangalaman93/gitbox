package main

import (
	"encoding/json"
	"fmt"
)

type Test_struct1 struct {
	Id   int
	Name string
}

type Test_struct2 struct {
	Id int
}

type Response1 struct {
	Page   int
	Fruits []string
}

func main() {
	ts1 := &Test_struct1{Id: 2, Name: "arbit"}
	j1, err := json.Marshal(ts1)
	fmt.Println(err)
	fmt.Println(string(j1))
	fmt.Println(j1)

	ts2 := &Test_struct2{}
	json.Unmarshal(j1, &ts2)
	fmt.Println(ts2.Id)
}
