# gogate
Report server/service status. This is not just another metrics system, it is
just a tool to reflect your current infrastructure facts without effort

# status

work in progress

# run
 
```
go install gogate
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

loop forever to  report status of the server

- / redirect to master
- /self: self status
- put /identifier forward to master
- /health show health status json
