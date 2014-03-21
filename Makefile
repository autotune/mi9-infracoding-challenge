all: mi9-infracoding-challenge

mi9-infracoding-challenge:
	go build

server: mi9-infracoding-challenge
	./mi9-infracoding-challenge

clean:
	rm -rf mi9-infracoding-challenge