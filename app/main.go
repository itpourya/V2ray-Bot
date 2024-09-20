package main

func main() {
	app := GenarateAPP()

	go app.Start()

	ch := make(chan struct{})
	<-ch
}
