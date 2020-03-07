package main

import (
	"../helper"
	"encoding/json"
	"fmt"
	"github.com/amit-upadhyay-it/goutils/io"
)

func main() {

	file1Path := "files/config3.yml"
	file2Path := "files/config2.yml"

	lines1, _ := io.ReadFile(file1Path, false)
	lines2, _ := io.ReadFile(file2Path, false)

	result1 := helper.GetMapFromYAML(lines1, "Config1")
	result2 := helper.GetMapFromYAML(lines2, "Config2")
	r1, _ := json.Marshal(result1["Config1"])
	r2, _ := json.Marshal(result2["Config2"])
	fmt.Println(string(r1))
	fmt.Println("--------------------------------------------------------------------")
	fmt.Println(string(r2))
}
