package utils

import (
	"fmt"
	"os"
)

func Assert(param interface{}) {
	if param == nil {
		os.Exit(2)
	}
	switch param.(type) {
	case string:
		if param == "" {
			os.Exit(2)
		}
	}
}

func AssertPanic(err error) {
	if err != nil {
		panic(err)
	}
}

func AssertPrint(err error) {
	if err != nil {
		fmt.Printf("Error found. Detail: \n%s\n", err.Error())
	}
}
