// +build !windows

package main

import (
	"fmt"
	"internal/server/args"
)

func main() {


	// Separate args parsing, as in the future we may need to decouple the server from this args
	path, port := args.ParseArgs(); // The program will exit if the arguments are not correct, so we always get correct values here


	fmt.Printf("Server will start on port: %d\n", port)
	fmt.Printf("Serving directory %v\n", path)

//	fmt.Println( server.Hello() )
}



