package constants

func ErrorConstants(command int) string {
	if command == 0 {
		return "Invalid Input"
	} else if command == 1 {
		return "Invalid parameters - SET key value [expiration EX seconds|PX milliseconds]"
	} else if command == 2 {
		return "Invalid parameters - GET key"
	} else if command == 3 {
		return "Invalid parameters - DEL key [key ...]"
	} else if command == 4 {
		return "Invalid parameters - TTL key"
	} else if command == 5 {
		return "Invalid parameters - EXPIRE key seconds"
	} else if command == 6 {
		return "Invalid parameters - ZADD key score member [score member ...]"
	} else if command == 7 {
		return "Invalid parameters - ZRANGE key start stop [WITHSCORES]"
	} else if command == 8 {
		return "Invalid parameters - ZRANK key member"
	} else if command == 9 {
		return "WRONGTYPE Operation against a key holding the wrong kind of value"
	}
	return "Not Identified"
}
