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