help:
	@egrep -h '\s##\s' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m  %-30s\033[0m %s\n", $$1, $$2}'
default: serve ## make serve

serve: server.c ## build and run server
	gcc server.c store.h -o server
	./server

client: client.c ## build and run client
	gcc client.c -o client
	./client

clean: ## remove built files
	rm -rf *.o *.out ./server ./client
