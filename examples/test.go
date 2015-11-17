package main

import "fmt"

func main() {
    var channels map[*chan string]map[string]struct{}
    
    channels = make(map[int]map[string]struct{})
    channels[1] = make(map[string]struct{})
    channels[1]["hi"] = struct{}{}
    
    val1, has1 := channels[1]["hi"]
    val2, has2 := channels[1]["bye"]
    val3, has3 := channels[2]["hi"]
    
    fmt.Println(val1, has1)
    fmt.Println(val2, has2)
    fmt.Println(val3, has3)
}