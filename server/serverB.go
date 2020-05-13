package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"math"
	"io/ioutil"
)

var (
	port2 = "8083"
)

type Range struct {
    From float64
    To  float64
}

func main() {
	http.HandleFunc("/", HelloServer)
	http.ListenAndServe(":"+port2, nil)
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

     var resB float64=0.0
	 for i := p.From;i<p.To;i+=0.000000001 {
	 	resB+=2*math.Sqrt(1-i*i)*0.000000001
	 }
   
	 fmt.Println(" B "+string(data))
	 fmt.Fprintf(w,"%f", resB)
	//fmt.Fprintf(w, "Hello from port %s!", port2)
}
