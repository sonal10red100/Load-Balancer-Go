package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"math"
	"io/ioutil"
)

var (
	port3 = "8085"
)

type Range struct {
    From float64
    To  float64
}

func main() {
	http.HandleFunc("/", HelloServer)
	http.ListenAndServe(":"+port3, nil)
}

func HelloServer(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	 if err!=nil{
	 	panic(err)
	 }

	 var p Range
	 json.Unmarshal([]byte(string(data)),&p)
	 
	 fmt.Printf("from : %f",p.From)
	 fmt.Println()
	 fmt.Printf("to : %f",p.To)

     var resC float64=0.0
	 for i := p.From;i<p.To;i+=0.000000001 {
	 	resC+=2*math.Sqrt(1-i*i)*0.000000001
	 }
   
	 fmt.Println(" C "+string(data))
	 fmt.Fprintf(w,"%f", resC)
	//fmt.Fprintf(w, "Hello from port %s!", port3)
}
