#!/bin/bash

cd /opt/netbox/netbox/
source /opt/netbox/venv/bin/activate

#python3 manage.py syncdb --noinput
python3 manage.py migrate
python3 manage.py createsuperuser
python3 manage.py collectstatic --no-input

python3 manage.py runserver 0.0.0.0:8000  --insecure