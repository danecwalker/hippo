bin/hippo: cmd/main.go
	go build -o bin/hippo cmd/main.go

.PHONY: bin/hippo