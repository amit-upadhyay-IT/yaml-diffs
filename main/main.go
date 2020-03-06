package main

import (
	"enumutil/utils/io"
	"fmt"
	"log"
	"strings"
)

func main() {

	file1Path := "files/config1.yml"
	file2Path := "files/config2.yml"

	lines1, _ := io.ReadFile(file1Path, false)
	lines2, _ := io.ReadFile(file2Path, false)

	result1 := GetMapFromYAML(lines1, "Config1")
	GetMapFromYAML(lines2, "Config2")
	fmt.Println(result1)

}

func GetMapFromYAML(lines []string, rootName string) map[string]interface{} {
	// pre-process lines to remove comments, empty lines, etc
	fileLines := make([]string, 0)
	for _, line := range lines {
		if !isDummyLine(line) {
			fileLines = append(fileLines, line)
		}
	}
	// forming-up the map after pre-processing
	var result map[string]interface{} = make(map[string]interface{})
	FormUpMapFromYAML(fileLines, result, rootName)
	return result
}

// a yml file or json file are structured as a tree, so just
// doing a dfs kind of traversal to formup the map
func FormUpMapFromYAML(lines []string, result map[string]interface{}, parentKey string) {

	for i := 0; i < len(lines); i++ {
		line := lines[i]
		if isDummyLine(line) {
			continue
		}
		// we can identify the child node based on the space indentation or RHS value after colon, currently using the later
		keyVal := strings.Split(line, ":")
		keyVal = []string{keyVal[0], strings.Join(keyVal[1:], ":")}  // handling case where line can have more than one colon
		if isChildNodeFound(line) {
			nextOffset := getOffset(lines[i:])
			data, ok := result[parentKey].(map[string]interface{})
			if ok {
				FormUpMapFromYAML(lines[i+1 : i+nextOffset], data, strings.TrimSpace(keyVal[0]))
				i = i + nextOffset - 1 // subtracting one coz, the loop will increment the counter anyway.
			} else {
				// since, no map could be found with parentKey, so here I am adding a new map to the parentKey where we will initialize their child nodes
				// Note that creating a key with parentKey here is necessary coz, we want to keep the references of the child nodes in the root node
				result[parentKey] = make(map[string]interface{})
				data, _ = result[parentKey].(map[string]interface{})
				FormUpMapFromYAML(lines[i+1 : i+nextOffset], data, strings.TrimSpace(keyVal[0]))
				i = i + nextOffset - 1
			}
		} else {
			// set up things in parent node
			key := strings.TrimSpace(keyVal[0])
			val := strings.TrimSpace(keyVal[1])
			node := make(map[string]interface{})
			node[key] = val
			// get existing map elements and append them to node, otherwise we will loose up the already added nodes of same level (same height)
			data, ok := result[parentKey].(map[string]interface{})
			if ok {
				for k, v := range data {
					node[k] = v
				}
			}
			result[parentKey] = node
		}
	}
}

func countPrefixSpaces(word string) int {
	count := 0
	for i := 0; i < len(word); i++ {
		if word[i] == ' ' {
			count++
		} else {
			return count
		}
	}
	return count
}

// a dummy line can be a comment, empty lines, or lines beginning with space and followed by comment
// trim beginning white spaces and if after trimming the first character is `#` then its a dummy line
func isDummyLine(line string) bool {
	lin := strings.TrimSpace(line)
	return strings.HasPrefix(lin, "#") || len(lin) == 0 || !strings.Contains(lin, ":")
}

func isChildNodeFound(line string) bool {
	keyVal := strings.Split(line, ":")
	if len(keyVal) > 1 {
		if keyVal[1] == "" {
			return true
		} else {
			if v := strings.TrimSpace(keyVal[1]); strings.HasPrefix(v, "&") {  // TODO: store this reference somewhere, later use this to form same level nodes
				return true
			} else {
				return false // this will be the leaf node
			}
		}
	}
	// in any case length should be greater than 1, so this code point should not reach
	log.Fatal("this shouldn't have happened, something is terribly wrong!")
	return false
}

func getOffset(lines []string) int {
	prefixSpacesCount := countPrefixSpaces(lines[0])
	index := -1
	for i := 1; i < len(lines); i++ {
		spcCnt := countPrefixSpaces(lines[i])
		if spcCnt == prefixSpacesCount {
			index = i
			break
		}
	}
	if index == -1 {
		return len(lines)
	}
	return index
}
