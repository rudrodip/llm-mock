package main

func main() {
	api := NewAPIServer(":8080")
	api.Run()
}
