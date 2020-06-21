FROM   nginx:1.19.0-alpine

COPY ./nginx.conf /etc/nginx/nginx.conf
COPY ./60-netbox-server.conf /etc/nginx/sites-enabled/
COPY ./80-upstream-netbox.conf  /etc/nginx/conf.d/