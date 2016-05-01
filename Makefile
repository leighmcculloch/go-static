test:
	go test -cover
	go vet

coverage:
	go test -coverprofile=c.out
	go tool cover -html=c.out
	rm c.out
