# go-sms
SMS Library for telemessage service - Golang

for more info :
http://developer.telemessage.com/REST/

# Example
```golang

package main

import (
	"github.com/NaySoftware/go-sms"
	"fmt"
)


func main() {

	client := telemsg.NewClient("username","password")
	client.NewMsg("+123456789", "Hello World")
	status, err := client.Send()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(status.ResultDescription)
	}

}

```
