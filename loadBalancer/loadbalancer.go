package main

import (
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
	 "fmt"
	 "io/ioutil"
	 "strconv"
	//"encoding/json"
	"bytes"
)

// Server is used to describe the individual servers
// that are connected to the Load Balancer
type Server struct {
	Route        string
	Alive        bool
	ReverseProxy *httputil.ReverseProxy
}

// ServerList is all the servers the Load Balancer
// has access to. The index of the server accessed
// most recently is stored in ServerList.Latest
type ServerList struct {     
	Servers []Server
	Latest  int
}

var res float64=0.0


// isAlive checks if a server is available by sending
// a TCP request to server.Route and checking if it
// successfully responds back
func (server *Server) isAlive() bool {
	timeout := time.Duration(1 * time.Second)

	log.Println("Started Health Check For:", server.Route)
	_, err := net.DialTimeout("tcp", server.Route, timeout)
	if err != nil {
		log.Println(server.Route, "Is Dead")
		log.Println("Health Check Error:", err)
		server.Alive = false
		return false
	}

	log.Println(server.Route, "Is Alive")
	server.Alive = true
	return true
}

// init is used to create the ServerList by taking in
// a slice of routes that need to be connected to
// the server and convert them to the Server
// struct format and store them all
// in ServerList.Servers slice
func (serverList *ServerList) init(serverRoutes []string) {
	log.Println("Creating Server List For Routes:", serverRoutes)

	for _, serverRoute := range serverRoutes {
		var localServer Server

		localServer.Route = serverRoute
		localServer.Alive = localServer.isAlive()

		origin, _ := url.Parse("http://" + serverRoute)
		director := func(req *http.Request) {
			req.Header.Add("X-Forwarded-Host", req.Host)
			req.Header.Add("X-Origin-Host", origin.Host)
			req.URL.Scheme = "http"
			req.URL.Host = origin.Host
		}
		localServer.ReverseProxy = &httputil.ReverseProxy{Director: director}

		log.Println("Server", localServer, "Added To Server List")
		serverList.Servers = append(serverList.Servers, localServer)
	}

	serverList.Latest = -1
	log.Println("Successfully Created Server List:", serverList)

}

// nextServer facilitates the round robin selection
// of each server by getting back to the first
// server after the last server is passed
func (serverList *ServerList) nextServer() int {
	return (serverList.Latest + 1) % len(serverList.Servers)
}

// loadBalance takes in the request and based on Round Robin method
// assigns it to a particular server in ServerList.Servers. If no
// servers are present it responds with a http.StatusServiceUnavailable
// status back to the client and if there are servers present it then
// checks if the server is alive and then only routes the request to it,
// otherwise it loops through the entire ServerList.Servers to find
// another alive server until it gets back to the first server
// it tried accessing and then responds with a
// http.StatusServiceUnavailable status
// back to the client
func (serverList *ServerList) loadBalance(w http.ResponseWriter, r *http.Request) {

	if len(serverList.Servers) > 0 {
		serverCount := 0
		for index := serverList.nextServer(); serverCount < len(serverList.Servers); index = serverList.nextServer() {
			if serverList.Servers[index].isAlive() {

				    body, err := ioutil.ReadAll(r.Body)
				    if err != nil {
				        http.Error(w, err.Error(), http.StatusInternalServerError)
				        return
				    }

				    // you can reassign the body if you need to parse it as multipart
				    r.Body = ioutil.NopCloser(bytes.NewReader(body))

				    // create a new url from the raw RequestURI sent by the client
				    url:="http://"+serverList.Servers[index].Route
				    proxyReq, err := http.NewRequest(r.Method, url, bytes.NewReader(body))

				    // We may want to filter some headers, otherwise we could just use a shallow copy
				    // proxyReq.Header = req.Header
				    proxyReq.Header = r.Header
				    client := &http.Client{}
				    resp, err := client.Do(proxyReq)
				    if err != nil {
				        http.Error(w, err.Error(), http.StatusBadGateway)
				        return
				    }
				    defer resp.Body.Close()
				    resBody, _ := ioutil.ReadAll(resp.Body)
				    f,_:=strconv.ParseFloat(string(resBody),8)
				    res+=f
    				fmt.Println("response Body: ", f)
					fmt.Println("Ans after adding : ", res)
					pi:=fmt.Sprintf("%f",res)
					fmt.Println(pi)
					b1:=[]byte("Value of Pi : "+pi)
					fmt.Println(b1)
					w.Write(b1)
					

				log.Println("Routing Request",string(body), "To", serverList.Servers[index].Route)
				//serverList.Servers[index].ReverseProxy.ServeHTTP(w, r)

				serverList.Latest = index
				log.Println("Updated Latest Server To:", serverList.Latest)

				return
			}
			serverCount++
			serverList.Latest = serverList.nextServer()
		}
	}
	log.Println("No Servers Available")
	http.Error(w, "No Servers Available", http.StatusServiceUnavailable)
}

// We can either import this as a package or use initialize
// the ServerList by providing a list of server routes to
// connect to and then create a server for the Load Balancer
func main() {
	var serverList ServerList
	loadBalancerPort := "8080"

	serverRoutes := []string{
		"localhost:8081",
		"localhost:8083",
		"localhost:8085",
	}

	serverList.init(serverRoutes)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		serverList.loadBalance(w, r)
	})

	http.ListenAndServe(":"+loadBalancerPort, nil)
}
