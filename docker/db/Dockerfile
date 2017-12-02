FROM percona:5.7
MAINTAINER hkheradm@gmail.com

RUN apt-get update && apt-get install -y wget libdbi-perl libdbd-mysql-perl
RUN wget -O /usr/bin/pt-heartbeat https://percona.com/get/pt-heartbeat && chmod +x /usr/bin/pt-heartbeat
