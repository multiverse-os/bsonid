package main

import (
	"fmt"

	id "github.com/multiverse-os/levelup/id"
)

type TestStruct struct {
	Name string
	X    int
	Y    int
}

func main() {
	fmt.Println("id library")
	fmt.Println("==============================================================")
	fmt.Println("string 'test' Hash(id).String() [", id.Hash("test").String(), "]")
	fmt.Println("string 'test' Hash(id).Bytes() [", id.Hash("test").Bytes(), "]")
	fmt.Println("==============================================================")
	fmt.Println("string 'test' HashAs(Sha3, id).String() [", id.HashAs(id.SHA3, "test").String(), "]")
	fmt.Println("string 'test' HashAs(Sha3, id).Bytes() [", id.HashAs(id.SHA3, "test").Bytes(), "]")
	fmt.Println("==============================================================")
	fmt.Println("struct test struct HashAs(XXHash32, id) [", id.HashAs(id.XXH32, TestStruct{Name: "test", X: 1, Y: 4}), "]")
	fmt.Println("struct test struct HashAs(SHA3, id) [", id.HashAs(id.XXH32, TestStruct{Name: "test", X: 1, Y: 4}), "]")
	fmt.Println("==============================================================")
	fmt.Println("string 'test' Hash(id.New().String()).String() [", id.Hash(id.New().String()).String(), "]")
	fmt.Println("string 'test' Hash(id.New().Bytes()).Bytes() [", id.Hash(id.New().Bytes()).Bytes(), "]")
	fmt.Println("==============================================================")
	fmt.Println("bsonid from seed (test struct) [", id.NewFromSeed(TestStruct{Name: "test", X: 1, Y: 4}), "]")
	fmt.Println("string 'test' String() [", id.New().String(), "]")
	fmt.Println("string 'test' Bytes() [", id.New().Bytes(), "]")
	fmt.Println("string 'test' Uint32 [", id.New().UInt32(), "]")
}
