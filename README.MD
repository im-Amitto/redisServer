# Redis Functionalities using Golang

## How to run?

```
  ./main
```

## How to compile?

### Prerequisite
- Go

### Compile

```
  go build main.go
```

## Commands
```
 SET key value [expiration EX seconds|PX milliseconds]
 GET key
 DEL key [key ...]
 TTL key
 EXPIRE key seconds
 ZADD key score member [score member ...]
 ZRANGE key start stop [WITHSCORES]
 ZRANK key member
 
```

## Why GoLang?

Golang is has fast-paving compilation and execution. Unlike C++/JAVA/Python it supports Multithreading And Concurrency in much better way. While Go may not provide you with a wide range of pre-built data structure but that enables you to build new data structures and get the better understanding.
Redis is one of the fastest cache/database available in the market thus it's replica should also be developed using one of the fastest language hence go is one of the perfect choice.

## Why not java/c++

Java and C++ both are written years back with different mindset and Multithreading And Concurrency were not hot topics back then thus they were not developed with a mindset of supporting such features. Community from both languages have provided updates that support Multithreading And Concurrency but they are not that efficient when compare to GoLan and the reason is simple C++/Java has some code which can't be touched for one issue or other but Go was written recently and with a focus on providing simplest language with maximum benefits.

## Futher Improvements

- Currently the application use goroutine(thread) to backup the data after a specific time period as well as when program exist, we can replace it with a channel or worker logic.
- Golang does not come with set and ordered_set and it does not have any popular library like c++ have stl so i have written the data structures myself and they do have scope for improvement.
- New features can be added

## Data Structures

- A Map was used for implementing the GET/SET/EXPIRE/TTL/DELETE commands . Through this data structure, we are able to support the mentioned operations in an average runtime complexity of O(1) which is constant time.
- A set like data structure was developed to keep track of all the keys present in O(1)
- For the ordered set functionality, an 'sorted set' data structure which was implemented over struct has been used.
- Complexity of the functions for data structures is mentioned at the top of functions in their respective files


## Multithreading Support

Yes, the implementation offers multithreading support. 
I haven't created multiple client because i did not find the need to implement it to prove that it can support multithreading.
It can handle multiple clients connected to the cache in parallel. Race conditions are handled through the use of channels.
Currently their are a total of 4 workers are placed, out of which 3 are responsible for read or single operation related functionalities(e.g. GET, ZADD, TTL...). The 4th one is responsible for handle write related operation thus the workers uses the first in first out to complete the query so their won't be a case where the data will be affected.
Also using channels can be proved beneficial when their is high system load and you wanna prioritize a specific type of operation. 

## Demo
### Click to play
[![Watch the video](https://i.ibb.co/30BfspP/Screenshot-from-2020-04-21-07-18-36.png)](https://www.loom.com/share/6b1acc97401149e498d19dc1678637f6)
