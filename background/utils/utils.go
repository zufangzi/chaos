package utils

import (
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
