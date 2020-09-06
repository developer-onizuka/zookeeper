package main

import (
	"github.com/go-zookeeper/zk"
	"time"
	"fmt"
)

func main() {
	var rootpath string = "/myapp"
	c, _, err := zk.Connect([]string{"127.0.0.1:2181"}, time.Second)
	if err != nil {
		panic(err)
	}

	flags := int32(0)
	acl := zk.WorldACL(zk.PermAll)
	fmt.Printf("%+v\n", acl)
	c.Create(rootpath, []byte{0}, flags, acl)
	if err != nil {
		panic(err)
	}

	prefix := "/"
	childlock, err :=c.CreateProtectedEphemeralSequential(rootpath+prefix,[]byte{0}, acl)
	fmt.Printf("%v\n", childlock)
	length := len(childlock)
	myid := childlock[length-10:length]
	fmt.Printf("%v\n", myid)

	ticker := time.NewTicker(time.Second)
	for {
		var leader bool = true
		children, _, err := c.Children(rootpath)
		//fmt.Printf("%v\n",children)
		if err != nil {
			panic(err)
		}
		if children == nil {
			leader = false
		}
		for _, name := range children {
			if myid > name[len(name)-10:len(name)] {
				leader = false
				break
			}
		}
		if leader == true {
			fmt.Printf("I am leader\n")
		} else if leader == false {
			fmt.Printf("I am follower\n")
		}
		<-ticker.C
	}
}

