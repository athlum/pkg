package utils

import (
	"fmt"
	"testing"
)

func Test_Logger(t *testing.T) {
	c := "/usr/bin/bash -c \"sleep 10000\""
	fmt.Printf("%v, %#v\n", c, SplitCommand(c))
	c = "/usr/bin/bash -c 'sleep 10000'"
	fmt.Printf("%v, %#v\n", c, SplitCommand(c))
	c = "docker run --env=\"asd=asd\" -e \"asd1=asd1\" test bash"
	fmt.Printf("%v, %#v\n", c, SplitCommand(c))
	c = "python"
	fmt.Printf("%v, %#v\n", c, SplitCommand(c))
}

func Test_LocalIP(t *testing.T) {
	fmt.Println(GetLocalIP())
	fmt.Println(GetLocalIPCached())
}
