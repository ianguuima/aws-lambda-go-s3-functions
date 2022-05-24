GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o main main.go
7z a main.zip main
rm main