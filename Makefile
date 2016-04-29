test:
	go test -v -cover

coverage:
	go test -coverprofile=c.out
	go tool cover -html=c.out
	rm c.out

coverage-update-readme:
	COVER="$$(go test -cover | grep coverage | sed -E "s/.*coverage: ([0-9]+)\.[0-9]%.* of statements/\1/")" \
	perl -i -pe's/(coverage)-(\d+%)-([a-z]+)\./\1-'"$$COVER"'%-\3./' README.md
