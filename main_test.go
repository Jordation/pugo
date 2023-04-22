package main

import (
	"fmt"
	"testing"
)

func TestPickCaptains(T *testing.T) {
	for i := 0; i < 100; i++ {
		n1, n2 := GetCaptainIds()
		if n1 == n2 {
			T.Fail()
		}
		fmt.Println(n1, n2)
	}
}
