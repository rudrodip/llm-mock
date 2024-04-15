package main

import "fmt"

func main() {
	fmt.Println(
		`
	_    _    __  __     __  __         _   
	| |  | |  |  \/  |___|  \/  |___  __| |__
	| |__| |__| |\/| |___| |\/| / _ \/ _| / /
	|____|____|_|  |_|   |_|  |_\___/\__|_\_\
	`)
	api := NewAPIServer(":6969")
	api.Run()
}
