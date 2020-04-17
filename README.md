# Round Robin Load Balancer
A round robin load balancer implemented in Go that calculates the value of Pi by distributing the calculation.

A list of routes to servers is provided that generates a load balancer using round robin algorithm to route requests to each server.

Usage:
1. go run all the number of servers from /server
2. go run loadBalancer/loadbalancer.go
3. go run main1.go

-loops through all available servers using round robin methodology and then routes the request to the appropriate server.
-performs health checks on each server to check whether ther server is alive or not : In the case a server goes down, the load balancer routes the request to the next available server.
  
