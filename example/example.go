package main

import (
	"fmt"

	bsonid "../../bsonid"
)

type TestStruct struct {
	Name string
	X    int
	Y    int
}

func main() {
	fmt.Println("sha3 test")
	fmt.Println("string 'test' SHA3 [", bsonid.Hash("test"), "]")
	fmt.Println("struct test struct SHA3 [", bsonid.Hash(TestStruct{
		Name: "test",
		X:    1,
		Y:    4,
	}), "]")

	fmt.Println("bsonid from seed (test struct) [", bsonid.NewFromSeed(TestStruct{
		Name: "test",
		X:    1,
		Y:    4,
	}), "]")
}
