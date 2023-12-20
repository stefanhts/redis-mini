#include <errno.h>
#include <netinet/in.h>
#include <netinet/ip.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/socket.h>
#include <unistd.h>

#include "store.h"

#define PORT 3000
#define BUFF_SIZE 1048

int handleConnection(int cfd) {
  char r_buff[BUFF_SIZE];
  char *str;
  char *respond = "Hello ";
  while (1) {
    ssize_t b_recv = recv(cfd, r_buff, sizeof(r_buff), 0);
    if (b_recv <= 0) {
      perror("no message recieved");
      close(cfd);
      return 1;
    }
    str = malloc(strlen(r_buff) + 1);
    strcpy(str, r_buff);
    if (b_recv >= BUFF_SIZE) {
      r_buff[BUFF_SIZE - 1] = '\0';
    } else {
      r_buff[b_recv] = '\0';
    }
    if (!strcmp(r_buff, "PING")) {
      send(cfd, "PONG", strlen("PONG"), 0);
    } else {
      char *token, *tofree;
      while((token = strsep(&str, " ")) != NULL) {
        send(cfd, token, strlen(token), 0);
        send(cfd, ", ", strlen(", "), 0);
      }
    }
    printf("got: %s", r_buff);
  }

  return 0;
}

int main() {
  initStore(4);

  // disable output buffering
  setbuf(stdout, NULL);

  struct sockaddr_in server_info = {0};
  struct sockaddr_in client_info = {0};
  server_info.sin_family = AF_INET;
  server_info.sin_port = htons(PORT);
  server_info.sin_addr.s_addr = ntohl(0);

  socklen_t server_info_len = sizeof(server_info);
  socklen_t client_info_len = sizeof(client_info);

  // create and configure socket
  int sfd = socket(server_info.sin_family, SOCK_STREAM, 0);
  if (sfd < 0) {
    perror("error creating socket");
    exit(1);
  }

  // bind socket
  if (bind(sfd, (struct sockaddr *)&server_info, server_info_len)) {
    perror("error binding socket");
    exit(1);
  }

  // listen on socket
  if (listen(sfd, 5)) {
    perror("error listening on socket");
    exit(1);
  }

  while (1) {
    // accept connection on socket
    int cfd = accept(sfd, (struct sockaddr *)&client_info, &client_info_len);
    if (cfd <= 0) {
      perror("error accepting client connection");
      exit(1);
    }
    handleConnection(cfd);
  }

  return 0;
}
