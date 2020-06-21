```bash
flv@genji ~ % docker inspect  $(docker ps --format "{{.ID}} {{.Names}}" |awk '/db/ {print $1}')|jq '.[0]|.NetworkSettings.Networks'
{
  "netbox_backend": {
    "IPAMConfig": null,
    "Links": null,
    "Aliases": [
      "db",
      "1a4d9d08d0aa"
    ],
    "NetworkID": "229889e6c51292a7ac83a5778c5863ec69a6ebb0b4d7e1e637d11883ade4e6b5",
    "EndpointID": "0cc465cb228f18f10aa21c165eef74ae7b5e88a0f83763fe61bd3d48b24f0ee8",
    "Gateway": "172.18.0.1",
    "IPAddress": "172.18.0.2",
    "IPPrefixLen": 16,
    "IPv6Gateway": "",
    "GlobalIPv6Address": "",
    "GlobalIPv6PrefixLen": 0,
    "MacAddress": "02:42:ac:12:00:02",
    "DriverOpts": null
  }
}

```

```bash
flv@genji ~ % docker exec -it netbox_netbox_1 /bin/bash                                              
[root@3ddd4ec6fec0 /]#                                                                                                                                                                                     
[root@3ddd4ec6fec0 /]# source /opt/netbox/venv/bin/activate                                          
(venv) [root@3ddd4ec6fec0 netbox]# python -V    
Python 3.6.8      
(venv) [root@3ddd4ec6fec0 /]# cd /opt/netbox/netbox/
(venv) [root@3ddd4ec6fec0 netbox]# 
(venv) [root@3ddd4ec6fec0 netbox]# python3 manage.py createsuperuser
/opt/netbox-2.8.6/venv/lib/python3.6/site-packages/cacheops/redis.py:21: RuntimeWarning: The cacheops cache is unreachable! Error: Error 99 connecting to localhost:6379. Cannot assign requested address.
  warnings.warn("The cacheops cache is unreachable! Error: %s" % e, RuntimeWarning)
Username (leave blank to use 'root'): admin
Email address: admin@blizzard.com
Password: 
Password (again): 
Superuser created successfully.

```

admin:admin_netbox