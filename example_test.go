package plumb

import (
	"fmt"
)

func ExamplePort_Recv() {
	port, err := Open("edit", 0)
	if err != nil {
		fmt.Println(err)
		return
	}

	msg, err := port.Recv()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("%s", msg)
}

func ExamplePort_Send() {
	port, err := Open("send", 1)
	if err != nil {
		fmt.Println(err)
		return
	}

	port.Send(&Msg{Dst: "edit", Wdir: "/tmp", Data: []byte("/etc/passwd:9")})
}
