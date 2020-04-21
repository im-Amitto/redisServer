package redis

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/im-Amitto/redisServer/constants"

	"github.com/im-Amitto/redisServer/sortedset"

	"github.com/im-Amitto/redisServer/set"
)

var getSetMap = make(map[string]string)
var ttlMap = make(map[string]int)
var keysExist = set.ItemSet{}
var orderedSetMap = make(map[string]*sortedset.SortedSet)

func isset(arr []string, index int) bool {
	return (len(arr) > index)
}

func Expire(key string, timer int, timerType string) string {
	flag := false
	if value, ok := ttlMap[key]; ok {
		if value <= 0 {
			flag = true
		}
	} else {
		flag = true
	}
	ttlMap[key] = timer
	if flag {
		go func() {
			for k1, typ := key, timerType; true; {
				if typ == "seconds" {
					time.Sleep(time.Second * 1)
				} else {
					time.Sleep(time.Millisecond * 1)
				}
				if value, ok := ttlMap[k1]; ok {
					if value == 1 {
						DEL([]string{k1})
						break
					} else if value == -1 {
						break
					} else {
						ttlMap[key] = value - 1
					}
				} else {
					break
				}
			}
		}()
	}
	return "ok"
}

//SET handle redis set operation
func SET(command []string) string {
	DEL([]string{command[1]})
	flag := true
	if isset(command, 3) {
		if command[3] == "EX" || command[3] == "ex" {
			i1, err := strconv.Atoi(command[4])
			if err == nil {
				flag = false
				Expire(command[1], i1, "seconds")
			} else {
				return "Invalid command"
			}
		} else if command[3] == "PX" || command[3] == "px" {
			i1, err := strconv.Atoi(command[4])
			if err == nil {
				flag = false
				Expire(command[1], i1, "milliseconds")
			} else {
				return "Invalid command"
			}
		} else {
			return "Invalid command"
		}
	}
	keysExist.Add(command[1])
	getSetMap[command[1]] = command[2]
	if flag {
		ttlMap[command[1]] = -1
	}
	return "ok"
}

//ZADD handle redis zadd operation
func ZADD(key string, values []string) string {
	var tempSet *sortedset.SortedSet
	if keysExist.Has(key) {
		if _, ok := getSetMap[key]; ok {
			DEL([]string{key})
			tempSet = sortedset.New()
		} else {
			tempSet = orderedSetMap[key]
		}
	} else {
		tempSet = sortedset.New()
	}

	for i := 0; i < len(values); i = i + 2 {
		score, _ := strconv.Atoi(values[i])
		sc := sortedset.SCORE(score)
		tempSet.AddOrUpdate(values[i+1], sc, 0)
	}
	orderedSetMap[key] = tempSet
	keysExist.Add(key)
	ttlMap[key] = -1
	return "ok"
}

//ZRANGE handle redis zadd operation
func ZRANGE(key string, command []string) string {
	defaultMessage := "(nil)"
	start, _ := strconv.Atoi(command[0])
	stop, _ := strconv.Atoi(command[1])
	withscores := "false"
	if isset(command, 2) {
		withscores = command[2]
	}

	if keysExist.Has(key) {
		defaultMessage = constants.ErrorConstants(9)
	}
	if value, ok := orderedSetMap[key]; ok {
		output := ""
		tempArray := value.GetByRankRange(start, stop, false)
		if strings.ToLower(withscores) == "withscores" {
			for _, s := range tempArray {
				output = output + "(1): " + s.Key() + " (2): " + strconv.Itoa(int(s.Score())) + "\n"
			}
		} else {
			for _, s := range tempArray {
				output = output + s.Key() + "\n"
			}
		}
		return output + "ok"
	}
	return defaultMessage
}

//ZRANK handle redis zadd operation
func ZRANK(key string, member string) string {
	defaultMessage := "(nil)"
	if keysExist.Has(key) {
		defaultMessage = constants.ErrorConstants(9)
	}
	if value, ok := orderedSetMap[key]; ok {
		return strconv.Itoa(value.FindRank(member))
	}
	return defaultMessage
}

//GET handle redis get operation
func GET(command []string) string {
	defaultMessage := "(nil)"
	if keysExist.Has(command[1]) {
		defaultMessage = constants.ErrorConstants(9)
	}
	if value, ok := getSetMap[command[1]]; ok {
		return value
	}
	return defaultMessage
}

//DEL handle redis DEL operation
func DEL(keys []string) string {
	flag := "0"
	for _, key := range keys {
		if keysExist.Has(key) {
			flag = "1"
			keysExist.Delete(key)
			delete(ttlMap, key)
			if _, ok := getSetMap[key]; ok {
				delete(getSetMap, key)
			} else if _, ok := orderedSetMap[key]; ok {
				delete(orderedSetMap, key)
			}
		}
	}
	return flag
}

//TTL handle redis TTL operation
func TTL(command []string) string {
	if value, ok := ttlMap[command[1]]; ok {
		return strconv.Itoa(value)
	}
	return "-2"
}

func BackUp() {
	f, err := os.Create("./backup/backup.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	line1 := ""
	for k, v := range getSetMap {
		line1 = line1 + k + ":" + v + ","
	}
	f.WriteString(line1)
	f.WriteString("\n")
	line2 := ""
	for k, v := range orderedSetMap {
		line2 = line2 + k + ":"
		tempElements := v.GetByRankRange(0, -1, false)
		for _, s := range tempElements {
			line2 = line2 + strconv.Itoa(int(s.Score())) + "$" + s.Key() + "||"
		}
		line2 = line2 + ","
	}
	f.WriteString(line2)
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
}

func Restore() {
	file, err := os.Open("./backup/backup.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	map1 := scanner.Text()
	scanner.Scan()
	map2 := scanner.Text()
	map1Data := strings.Split(map1, ",")
	for _, m := range map1Data {
		if m == "" {
			break
		}
		splitted := strings.Split(m, ":")
		key := splitted[0]
		value := splitted[1]
		command := []string{"set", key, value}
		SET(command)
	}
	map2Data := strings.Split(map2, ",")

	for _, m := range map2Data {
		if m == "" {
			break
		}
		temp := []string{}
		splitted := strings.Split(m, ":")
		key := splitted[0]
		values := strings.Split(splitted[1], "||")
		for _, setV := range values {
			if setV != "" {
				splitted := strings.Split(setV, "$")
				temp = append(temp, splitted[0])
				temp = append(temp, splitted[1])
			}
		}
		ZADD(key, temp)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
