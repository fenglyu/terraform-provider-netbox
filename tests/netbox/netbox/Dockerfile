ARG version="2.8.6"
ARG piprepo="https://pypi.tuna.tsinghua.edu.cn/simple"

FROM centos:centos8

COPY ./ /

RUN dnf clean all && \
    dnf makecache && \
    dnf install -y \
    curl gcc make \
    python3 python3-devel python3-pip \
    python3-setuptools libxml2-devel \
    libxslt-devel libffi-devel \
    openssl-devel redhat-rpm-config && \
    dnf clean all && \
   # curl -o netbox.tar.gz https://github.com/netbox-community/netbox/archive/v2.8.6tar.gz && \
    tar xzf ./v2.8.6.tar.gz -C /opt && \
    cd /opt/ && \
    ln -s netbox-2.8.6/ netbox

RUN groupadd --system netbox && \
    adduser --system --gid netbox netbox && \
    chown --recursive netbox /opt/netbox/netbox/media/

RUN pip3 install -i https://pypi.tuna.tsinghua.edu.cn/simple virtualenv && \
    python3 -m virtualenv --python=/usr/bin/python3 /opt/netbox/venv

RUN /opt/netbox/venv/bin/pip install -i https://pypi.tuna.tsinghua.edu.cn/simple -r /opt/netbox/requirements.txt

#RUN cd netbox/netbox/ && \
#    cp configuration.example.py configuration.py && \
#    sed -i "s###g" configuration.py
RUN  cp ./configuration.py /opt/netbox/netbox/netbox/

#RUN cd /opt/netbox/netbox/ && \
#    /opt/netbox/venv/bin/python3 manage.py migrate && \
#    /opt/netbox/venv/bin/python3 manage.py createsuperuser && \
#    /opt/netbox/venv/bin/python3 manage.py collectstatic --no-input
#RUN chown a+x /init.sh
EXPOSE 8000

ENTRYPOINT [ "/bin/bash", "/init.sh" ]