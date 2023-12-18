#include <netinet/in.h>
#include <netinet/ip.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/socket.h>
#include <unistd.h>

int main() {
  int fd = socket(AF_INET, SOCK_STREAM, 0);
  if (0 >= fd) {
    perror("socket");
    return -1;
  }

  struct sockaddr_in client_info = {0};
  client_info.sin_family = AF_INET;
  client_info.sin_port = htons(3000);
  client_info.sin_addr.s_addr = ntohl(INADDR_LOOPBACK);

  int rv =
      connect(fd, (const struct sockaddr *)&client_info, sizeof(client_info));

  char buff[1024];
  while (1) {
    printf("redis >");
    scanf("%1023s", buff);
    if (!strcmp(buff, ".exit")) {
      exit(0);
    }

    write(fd, buff, strlen(buff));

    char rbuf[64] = {};

    int n = read(fd, rbuf, sizeof(rbuf) - 1);
    rbuf[n + 1] = '\0';

    printf("%s\n", rbuf);
  }
  close(fd);
}
