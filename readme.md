# gogate
Report server/service status. This is not just another metrics system, it is
just a tool to reflect your current infrastructure facts without effort

# when to use it?

- When you have a horriable infrastructure, you may find this useful. 
- When you have a good infrastructure, and a CMDB to control the `desired status`, you can still have it to serve as an view to `current status`.
 
# status

work in progress

# run
 
```
go install github.com/woosley/gogate
gogate -h
```

sample

```
[root@localhost gogate]# gogate --is-master --key ip
http server started on [::]:1234

[root@localhost gogate]# http http://localhost:1234/self
HTTP/1.1 200 OK
Content-Length: 264
Content-Type: application/json; charset=UTF-8
Date: Mon, 28 Aug 2017 08:36:32 GMT

{
    "Hostname": "localhost",
    "Interfaces": [
,,,,,


[root@localhost gogate]# http http://localhost:1234/
HTTP/1.1 200 OK
Content-Length: 278
Content-Type: application/json; charset=UTF-8
Date: Mon, 28 Aug 2017 08:36:36 GMT

{
    "10.0.2.15": {
        "Hostname": "localhost",

```
# design

loop forever to report status of the server

- GET /: Redirect to master on node; Return all nodes on master
- GET /self: current node status
- POST /: Not allowed on node; Create new record on master
- GET /health: show health status json
