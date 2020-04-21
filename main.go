package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/im-Amitto/redisServer/constants"
	"github.com/im-Amitto/redisServer/redis"
)

func handleOperation(command []string) string {
	operation := strings.ToLower(command[0])
	commandLength := len(command)

	switch operation {
	case "set":
		if commandLength == 3 || commandLength == 5 {
			return redis.SET(command)
		}
		return constants.ErrorConstants(1)

	case "get":
		if commandLength == 2 {
			return redis.GET(command)
		}
		return constants.ErrorConstants(2)

	case "ttl":
		if commandLength == 2 {
			return redis.TTL(command)
		}
		return constants.ErrorConstants(4)

	case "del":
		if commandLength >= 2 {
			var sliced []string = command[1:commandLength]
			return redis.DEL(sliced)
		}
		return constants.ErrorConstants(3)

	case "expire":
		if commandLength == 3 {
			i1, err := strconv.Atoi(command[2])
			if err == nil {
				return redis.Expire(command[1], i1, "seconds")
			} else {
				return constants.ErrorConstants(5)
			}
		}
		return constants.ErrorConstants(5)

	case "zadd":
		if commandLength >= 4 && commandLength%2 == 0 {
			var sliced []string = command[2:commandLength]
			return redis.ZADD(command[1], sliced)
		}
		return constants.ErrorConstants(6)

	case "zrange":
		if commandLength == 4 || commandLength == 5 {
			var sliced []string = command[2:commandLength]
			return redis.ZRANGE(command[1], sliced)
		}
		return constants.ErrorConstants(7)

	case "zrank":
		if commandLength == 3 {

			return redis.ZRANK(command[1], command[2])
		}
		return constants.ErrorConstants(8)

	default:
		return constants.ErrorConstants(0)
	}
	return constants.ErrorConstants(0)
}

func worker(c1 chan []string) {
	go func() {
		for {
			select {
			case data := <-c1:
				fmt.Println(handleOperation(data))
			}
		}
	}()
}

func workerW(c1 chan []string) {
	go func() {
		for {
			select {
			case data := <-c1:
				fmt.Println(handleOperation(data))
			}
		}
	}()
}

func main() {
	redis.Restore()
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Redis server")
	fmt.Println("---------------------")
	fmt.Println("Available Commands")
	fmt.Println("1. SET key value [expiration EX seconds|PX milliseconds]")
	fmt.Println("2. GET key")
	fmt.Println("3. DEL key [key ...]")
	fmt.Println("4. TTL key ")
	fmt.Println("5. EXPIRE key seconds")
	fmt.Println("6. ZADD key score member [score member ...]")
	fmt.Println("7. ZRANGE key start stop [WITHSCORES]")
	fmt.Println("8. ZRANK key member")

	readChannel := make(chan []string)
	writeChannel := make(chan []string)
	worker(readChannel)
	worker(readChannel)
	worker(readChannel)
	workerW(writeChannel)

	go func() {
		for {
			time.Sleep(time.Second * 20)
			redis.BackUp()
		}
	}()

	fmt.Println("Redis server started")
	for {
		text, _ := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)
		if strings.Compare("quit", text) == 0 {
			redis.BackUp()
			fmt.Println("bye")
			break
		}
		data := strings.Split(text, " ")
		operation := strings.ToLower(data[0])
		if operation == "get" || operation == "ttl" || operation == "zrank" || operation == "zrange" {
			readChannel <- data
		} else {
			writeChannel <- data
		}
	}

}
