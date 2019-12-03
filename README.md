# frp-api

## Usage

```shell script
# frp directory include frpc and frpc.ini
FRP=/opt/frp PORT=5600 go run server.go
```

```shell script
curl -X POST http://localhost:5600/api/frpc -H "Content-Type:application/json" -d '{"ip": "202.204.62.62", "port":"8000", "name": "test"}'
```

then the `/opt/frp/frpc.ini` will be changed to

```yaml
...

[test]
local_ip = 202.204.62.62
local_port = 8000
remote_port = 10001 # randomly generate
```
