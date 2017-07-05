# monexec
Light Supervisor for process on Go (with optional Consul autoregistration)

# installation

```
go get -v -u github.com/reddec/monexec/...
```


# modes

* `forever` - restart always
* `critical` - run once, on error stop and kill others 
* `restart` - restart N times on error, if not successed - die alone
* `oneshot` - run once, die alone on error
