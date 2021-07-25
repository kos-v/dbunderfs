build: clean
	go build -o dbfs main.go

clean:
	rm -f ./dbfs
