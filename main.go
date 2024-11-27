package main

func main() {
	// Start server
	server := NewAPIServer(":8080")
	server.Run()
}
