default: serve

serve: server.c
	gcc server.c -o server
	./server

client: client.c
	gcc client.c -o client
	./client

clean:
	rm -rf *.o *.out ./server ./client
