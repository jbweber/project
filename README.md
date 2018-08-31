# Overview

[![Build Status](https://travis-ci.org/jbweber/project.svg?branch=master)](https://travis-ci.org/jbweber/project)

This is a quick demo project meant to mee the following criteria:

* Write a web application in your language of choice that returns the current date/time in JSON
* Write a simple test application that will query this “API” X times per second and record success/failure/TTLB (Time to last byte)
* Perform a blue-green deploy with the method/technology of your choosing while the test application is running and demonstrate there were no failed requests
  * Go from a single instance of v1 to a single instance of v2 gracefully

# Observations

* client using persistent connections means it stays pinned to old deployment in blue/green failover
* this could also happen if we used a DNS type failover but never re-resolve names
* need a way to force connection recycle
* does this happen on other types of ingress?
* when you do this failover is there a way to ensure that clients gracefully switch,
  because if you force it at the server side you could get a TCP connection reset and possible failed request.

`netstat -n | grep <address:port>` should observe ~N established connections and some other value of time_wait depending on configuration

# Opportunities

* clean refactoring of client code
* unit tests of client code
* instrument client and server code using prometheus
* dive into more ideas around how blue / green deployment is impacted by client implementation
* how does http/2 handle the blue / green deployment situation since it uses multiplexed TCP which has inherent keepalive characteristics
* learn more about k8s as this was my first usage in anger
* cleaner story around build + deployment of container images
* ci / cd job to handle
