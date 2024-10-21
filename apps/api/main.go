package main

func main() {
	api := NewAPI()
	api.Serve(":3000")
}
