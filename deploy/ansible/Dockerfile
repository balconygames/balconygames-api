FROM cytopia/ansible

RUN apk add openssh-client

RUN eval `ssh-agent`

ADD ./certs/ssh/id_rsa /root/.ssh/
ADD ./certs/ssh/id_rsa.pub /root/.ssh/
ADD ./certs/ssh/known_hosts /root/.ssh/

RUN chmod 700 /root/.ssh/
RUN chmod 644 /root/.ssh/id_rsa.pub
RUN chmod 600 /root/.ssh/id_rsa

ADD ./scripts/common.sh /app/
ADD ./scripts/deploy.sh /app/
ADD ./scripts/provision.sh /app/

RUN chmod +x /app/common.sh
RUN chmod +x /app/deploy.sh
RUN chmod +x /app/provision.sh
