package main

import (
	"fmt"
	"github.com/amit-upadhyay-it/goutils/io"
	"../helper"
)

func main() {

	file1Path := "files/config3.yml"
	file2Path := "files/config2.yml"

	lines1, _ := io.ReadFile(file1Path, false)
	lines2, _ := io.ReadFile(file2Path, false)

	result1 := helper.GetMapFromYAML(lines1, "Config1")
	helper.GetMapFromYAML(lines2, "Config2")
	fmt.Println(result1)
}

