package ood

// package main

import (
	"fmt"
)

func main() {
	base := Base{"i am base"}
	rl := Real{base, "i am Real"}
	Bar(rl)
}

type Base struct {
	Name string
}

func (b *Base) Foo() {
	fmt.Println("i am now in Foo method...")
}

func (b *Base) Bar() {
	fmt.Println("i am now in Bar method...")
}

type Real struct {
	Base
	MyRealName string
}

func Bar(r Real) {
	r.Bar()
}
