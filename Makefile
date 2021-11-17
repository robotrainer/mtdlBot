GR=go run
GB=go build


all: main

main: initmongo.go initbot.go command.go notification.go main.go
	$(GB) -o build/mtdlBot

run: initmongo.go initbot.go command.go notification.go main.go
	$(GR) initmongo.go initbot.go command.go notification.go main.go