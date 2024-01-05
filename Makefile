generate:
	go generate ./...

genmock-usecase:
	 mockery.exe --dir=./usecase --name UsecaseInterface --filename mocks/mock_UsecaseInterface.go

test:
	go test ./...