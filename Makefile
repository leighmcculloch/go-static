test:
	go test -v -cover

coverage:
	go test -coverprofile=c.out
	go tool cover -html=c.out
	rm c.out
