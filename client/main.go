package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	fmt.Print("Your name?: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	username := scanner.Text()

	con, err := net.Dial("tcp", ":3500")
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected!")
	defer con.Close()

	con.Write([]byte(username))

	go func() {
		for {
			reply := make([]byte, 512)
			n, err := con.Read(reply)
			if err == nil {
				str := fmt.Sprintf("\r%s", string(reply)[0:n])
				fmt.Println(str)
				fmt.Printf("[%s]: >> ", username)
			}
		}
	}()
	for {
		fmt.Printf("[%s]: >> ", username)
		scanner.Scan()
		text := scanner.Text()
		_, err = con.Write([]byte(text))
		if err != nil {
			panic(err)
		}
	}
}