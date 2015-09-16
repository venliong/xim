package main

import (
	"fmt"

	"github.com/liuhengloveyou/nodenet"
)

func init() {
	nodenet.SetWorker(LOGIC_TEMPORARY, tempGroupWorker)
}

func tempGroupWorker(data interface{}) (result interface{}, err error) {
	fmt.Println(data)

	return data, nil
}
