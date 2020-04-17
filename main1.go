package main

import(
    "fmt"
    "net/http"
    "io/ioutil"
    "bytes"
    "encoding/json"
)


func main(){
    url := "http://localhost:8080/"
    fmt.Println("URL:>", url)
    for i := -1.0; i <= 1.0; i+=0.05 {
    d:=i+0.05
    reqBody, err := json.Marshal(map[string]float64{
    "from": i,
    "to": d,
    })
    if err != nil {
        print(err)
    }
    resp, err := http.Post("http://localhost:8080/",
        "application/json", bytes.NewBuffer(reqBody))
    if err != nil {
        print(err)
    }
    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        print(err)
    }
    fmt.Println("response from load balancer : "+string(body))
}

}