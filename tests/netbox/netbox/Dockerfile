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
    tar xzf ./netbox-v2.4.7.tar.gz -C /opt && \
    cd /opt/ && \
    ln -s netbox-2.4.7/ netbox

RUN groupadd --system netbox && \
    adduser --system --gid netbox netbox && \
    chown --recursive netbox /opt/netbox/netbox/media/

RUN curl -s http://pki.blizzard.net/pki/blizzardroot.pem -o /etc/pki/ca-trust/source/anchors/blizzardroot.crt
RUN curl -s http://pki.blizzard.net/pki/BlizzardCA01.pem -o /etc/pki/ca-trust/source/anchors/BlizzardCA01.crt
RUN update-ca-trust extract

#RUN pip3 install -i https://mirrors.aliyun.com/pypi/simple/ --trusted-host mirrors.aliyun.com  virtualenv && \
#    python3 -m virtualenv --python=/usr/bin/python3 /opt/netbox/venv
#RUN /opt/netbox/venv/bin/pip3 install -i https://mirrors.aliyun.com/pypi/simple/ --trusted-host mirrors.aliyun.com -r /opt/netbox/requirements.txt

RUN pip3 install -i http://pip.battle.net:6081/api/pypi/pip-virtual/simple --trusted-host pip.battle.net virtualenv && \
    python3 -m virtualenv --python=/usr/bin/python3 /opt/netbox/venv
RUN /opt/netbox/venv/bin/pip install -i http://pip.battle.net:6081/api/pypi/pip-virtual/simple --trusted-host pip.battle.net -r /opt/netbox/requirements.txt


RUN  cp ./configuration.py /opt/netbox/netbox/netbox/

EXPOSE 8000

ENTRYPOINT [ "/bin/bash", "/init.sh" ]
