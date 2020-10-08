FROM   nginx:1.19.0-alpine

WORKDIR /

COPY ../../netbox-docker/docker/fake/nginx-selfsigned.crt /etc/ssl/certs/nginx-selfsigned.crt
COPY ../../netbox-docker/docker/fake/nginx-selfsigned.key /etc/ssl/private/nginx-selfsigned.key
COPY ../../netbox-docker/docker/fake/dhparam.pem /etc/ssl/certs/dhparam.pem
COPY ./nginx.conf /etc/nginx/nginx.conf
#COPY ./60-netbox-server.conf /etc/nginx/sites-enabled/
COPY ./60-netbox-server.conf /etc/nginx/conf.d/
COPY ./80-upstream-netbox.conf  /etc/nginx/conf.d/