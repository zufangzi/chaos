package main

import (
	"fmt"
	"reflect"
)

func main() {
	// test_closure()
	// test_pointer()
	var inta int
	inta = 10
	test_ref(&inta)
	fmt.Println(inta)

	//
	// aMap := make(map[string]string)
	// aMap["hello"] = "kitty"
	// fmt.Println(aMap["hello"])
}

func test_ref(a *int) {
	(*a)++
	fmt.Println(reflect.ValueOf(a))
}

func test_helloworld() {
	fmt.Println("hello world golang!")
}

func test_slice() {
	mySlice := make([]int, 5, 10)
	for _, v := range mySlice {
		fmt.Println(v)
	}
}

func test_closure() {
	j := 5
	printer := func(k int) func() {
		i := 2
		return func() {
			fmt.Printf("i, j, k: %d, %d, %d\n", i, j, k*j*i)
		}
	}(1)

	printer()
	j *= 2
	printer()
	fmt.Println("down")
}

func test_pointer() {

	y, err := test_pointer_to_be_invoked()
	fmt.Println(y)
	fmt.Println(err.Op)
}

type QuickError struct {
	Op   string
	Path string
}

func test_pointer_to_be_invoked() (y int, err QuickError) {
	return 3, QuickError{"Operation", "/usr/local/bin"}
}
