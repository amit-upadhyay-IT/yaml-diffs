package helper

import (
	"log"
	"strings"
)

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
	RefactorMap(result, rootName)
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

		keyVal := strings.Split(line, ":")
		if len(keyVal) > 1 { // i.e. it doesn't have prefix "- "
			keyVal = []string{keyVal[0], strings.Join(keyVal[1:], ":")}  // handling case where line can have more than one colon
		}
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
			val := ""
			if len(keyVal) > 1 {
				val = strings.TrimSpace(keyVal[1])
			}
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
	return strings.HasPrefix(lin, "#") || len(lin) == 0 || (!strings.HasPrefix(lin, "- ") && !strings.Contains(lin, ":"))  // order maters here :P
}

// tells if the current node has a child node or not
// we can identify the child node based on the space indentation or RHS value after colon, currently using the later
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
	} else if len(keyVal) == 1 {  // then ideally it should start with `- `
		if v := strings.TrimSpace(keyVal[0]); strings.HasPrefix(v, "- ") {  // since its an array symbol and it doesn't contain `:` so the node itself is a child
			return false
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

// iterate over map of maps,
// when encountered that key is beginning with `- `, collect the keys in a list and change the reference of the parentKey to the collected list
/*
algorithm:
* In the map each key can contain a value as a map or as a primitive
* if key contains a map:
	* if key has the prefix `- `:
		* read the value map of the above key
		* if the value in the value map are not null/empty:
			* create an object of that key-value pair and append it to a list (eg: named listValues)
		* else if value in the value map is empty then:
			* simply append the keyName in the above list (eg: named listValues)
	* else if key doesn't have prefix:
		* do nothing
* else if key contains a value:
	* do nothing
 */
func RefactorMap(result map[string]interface{}, parentKey string) {
	data, ok := result[parentKey].(map[string]interface{})
	if ok {  // map is found
		listConcatenationRequired := false
		keysInData := getKeysFromMap(data)
		if keysInData != nil && len(keysInData) > 0 && strings.HasPrefix(keysInData[0], "- ") {
			assertionOnKeySymantic(keysInData)  // an assertion to verify if yaml file is having correct syntax, if one key has prefix `- ` then all the keys under that node should have that prefix
			listConcatenationRequired = true
		}
		listOfMapOrPrimitives := make([]interface{}, 0)
		if listConcatenationRequired {
			// create a separate map with each key value pair in data map
			// append it to a list
			for key, val := range data {
				if strings.HasPrefix(key, "- ") {
					key = strings.TrimLeft(key, "- ")
				}
				if val == nil || val == "" {
					listOfMapOrPrimitives = append(listOfMapOrPrimitives, key)
				} else {
					mapToAppend := make(map[string]interface{})
					mapToAppend[key] = val
					listOfMapOrPrimitives = append(listOfMapOrPrimitives, mapToAppend)
				}
			}
			result[parentKey] = listOfMapOrPrimitives
			for _, item := range listOfMapOrPrimitives {
				if node, ok := item.(map[string]interface{}); ok {
					RefactorMap(node, getOneKeyNameFromMap(node))
				}
			}
			// run a loop for the values present in result[parentKey] and call method recursively on each item in listOfMapOrPrimitives
		} else {
			// iterate through the map and recursively call for each node
			for key, _ := range data {
				RefactorMap(data, key)
			}
		}
		//RefactorMap(data, parentKey)
	} else {
		// child node is found

	}
}

func assertionOnKeySymantic(data []string) {
	for _, item := range data {
		if !strings.HasPrefix(item, "- ") {
			log.Fatal("The key in the yaml file seem to have problem, not all key are array")
		}
	}
}

func getKeysFromMap(dic map[string]interface{}) []string {
	res := make([]string, 0)
	for k, _ := range dic {
		res = append(res, k)
	}
	return res
}

func getOneKeyNameFromMap(dic map[string]interface{}) string {
	res := ""
	for key, _ := range dic {
		res = key
		break
	}
	return res
}
