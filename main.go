// SPDX-License-Identifier: BSD-3-Clause

// Package main implements the core logic
package main

import (
	"fmt"

	"github.com/9elements/bmc-test-go/pkg/testcase"
)

const bmcTestGoVersion = "v0.0.0"

func printResult(name string, result testcase.Result) {
	fmt.Printf("TEST [%s][%s]:", name, result.Name)

	if result.Error != nil {
		fmt.Printf(" Error: %v", result.Error)
	} else {
		fmt.Printf(" PASS")
	}

	fmt.Println()
}

func printResults(test testcase.Test) {
	for _, r := range test.Results {
		printResult(test.Name, r)
	}
}

func main() {
	fmt.Println(bmcTestGoVersion)

	result1 := testcase.Result{Name: "step 1", Error: nil}
	result2 := testcase.Result{Name: "step 2", Error: nil}

	test := testcase.Test{
		Name:        "nop test",
		Description: "showcasing how to use the structs",
		Function: func(_ testcase.Config) []testcase.Result {
			return []testcase.Result{result1, result2}
		},
		Results: []testcase.Result{result1, result2},
	}

	printResults(test)
}
