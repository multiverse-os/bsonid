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
	fmt.Println("This is the real reason to provide hashing of objects, so we")
	fmt.Println("can seed an id from the checksum hash of any given object")
	fmt.Println("string 'test' NewFromSeed(id.Hash(\"test\")).Bytes() [", id.NewFromSeed(id.Hash("test")).Bytes(), "]")
	fmt.Println("string 'test' NewFromSeed(id.Hash(\"test\")).String() [", id.NewFromSeed(id.Hash("test")).String(), "]")
	fmt.Println("==============================================================")
	fmt.Println("id.NewShort() creates a short ID where we remove some of the elements, dropping the size from 12 to 8")
	fmt.Println("string 'test' String() [", id.NewShort().String(), "]")
	fmt.Println("string 'test' Bytes() [", id.NewShort().Bytes(), "]")
	fmt.Println("string 'test' Uint32 [", id.NewShort().UInt32(), "]")
	fmt.Println("==============================================================")
	fmt.Println("Sometimes we will need IDs that are derived from something and")
	fmt.Println("will always reproduce the same ID given an input like a name of")
	fmt.Println("a collection for example. Then we can use that to get the uint32")
	fmt.Println("id from the name and use it to do a lookup without doing any string")
	fmt.Println("comparison:")
	fmt.Println(" id:Hash(\"collectionName\").UInt32():", id.Hash("collectionName").UInt32())
	fmt.Println(" id:Hash(\"collectionName\").Bytes():", id.Hash("collectionName").Bytes())
	fmt.Println(" id:Hash(\"collectionName\").String():", id.Hash("collectionName").String())
	fmt.Println(" id.HashAs(id.SHA3, \"collectionName\").UInt32():", id.HashAs(id.SHA3, "collectionName").UInt32())
	fmt.Println("==============================================================")
	fmt.Println("bsonid from seed (test struct) [", id.NewFromSeed(TestStruct{Name: "test", X: 1, Y: 4}), "]")
	fmt.Println("string 'test' String() [", id.New().String(), "]")
	fmt.Println("string 'test' Bytes() [", id.New().Bytes(), "]")
	fmt.Println("string 'test' Uint32 [", id.New().UInt32(), "]")
}
