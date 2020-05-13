package main

import(
    "math"
    "log"
    "html/template"
    "sync"
    "fmt"
    "net/http"
    "io/ioutil"
    "bytes"
    "encoding/json"
)

var doOnce sync.Once

type Resp struct {
        Res float64
        F float64
        Index int
}

type ServerDetails struct {
        Fr float64
        To float64
        Cal float64
}

type PiValue struct {
        SerA []ServerDetails
        SerB []ServerDetails
        SerC []ServerDetails
        Value float64
}
func round(num float64) int {
    return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
    output := math.Pow(10, float64(precision))
    return float64(round(num * output)) / output
}

func indexHTMLTemplateVariableHandler(w http.ResponseWriter, r *http.Request){

    tmpl := template.New("index.html")       //create a new template with some name
    tmpl, _= tmpl.ParseFiles("index.html") //parse some content and generate a template, which is an internal representation     
    
      var body []byte
 //    doOnce.Do(func(){
     url := "http://localhost:8080/"
    fmt.Println("URL:>", url)
    
    var a []ServerDetails
    var b []ServerDetails
    var c []ServerDetails
    ac:=0
    bc:=0
    cc:=0
    
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
    body, _= ioutil.ReadAll(resp.Body)
     if err != nil {
         print(err)
     }
   fmt.Println("response from load balancer : "+string(body))
 //   w.Write(body)

    var basket Resp
    err1 := json.Unmarshal(body, &basket)
    if err1 != nil {
        log.Println(err1)
    }
    fmt.Println(basket.F)
    idx:=basket.Index

    if idx == 0 {
        ii:=toFixed(i,3)
        dd:=toFixed(d,3)
        bb:=toFixed(basket.F,3)
        a1:=ServerDetails{ Fr: ii, To: dd, Cal: bb}
        a=append(a,a1)
        //a[ac]=a1
        ac=ac+1
    } else if idx == 1 {
        ii:=toFixed(i,3)
        dd:=toFixed(d,3)
        bb:=toFixed(basket.F,3)
        b1:=ServerDetails{ Fr: ii, To: dd, Cal: bb}
        b=append(b,b1)
        //b[bc]=b1
        bc=bc+1
    } else {
        ii:=toFixed(i,3)
        dd:=toFixed(d,3)
        bb:=toFixed(basket.F,3)
        c1:=ServerDetails{ Fr: ii, To: dd, Cal: bb}
        c=append(c,c1)        
        //c[cc]=c1
        cc=cc+1
    }

    }

   //  fmt.Println("response from load balancer : "+string(body))
 //    w.Write(body)
    s:=string(body)
   fmt.Println(s)

    var final Resp
    err1 := json.Unmarshal(body, &final)
    if err1 != nil {
        log.Println(err1)
    }
    fmt.Println(final.Res)
    d := PiValue{ SerA: a, SerB: b,SerC: c, Value: final.Res}
    fmt.Println(a)
    fmt.Println(b)
    fmt.Println(c)
    
    tmpl.Execute(w, d)

//})
}

func main() { 
   
    http.HandleFunc("/",indexHTMLTemplateVariableHandler)
    http.ListenAndServe(":8087", nil)
}


