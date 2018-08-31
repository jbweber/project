# Observations

* client using persistent connections means it stays pinned to old deployment in blue/green failover
* this could also happen if we used a DNS type failover but never re-resolve names
* need a way to force connection recycle
* does this happen on other types of ingress?
* when you do this failover is there a way to ensure that clients gracefully switch,
  because if you force it at the server side you could get a TCP connection reset and possible failed request.


`netstat -n | grep <address:port>` should observe ~N established connections and some other value of time_wait depending on configuration
