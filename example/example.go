package main

import (
    "fmt"
    "gwebot/bot"
    "io/ioutil"
    "log"
    "net/http"
)

func ReciveMsgHandler(w http.ResponseWriter, r *http.Request) {
    // var weibot bot.WeiBot
    var weibot bot.WeiBot = *bot.NewWeiBot()

    r.ParseForm()
    signature := r.FormValue("signature")
    timestamp := r.FormValue("timestamp")
    nonce := r.FormValue("nonce")
    echostr := r.FormValue("echostr")
    weibot.OnConnect(signature, timestamp, nonce, echostr)

    if vali, _echostr := weibot.Validate(); vali {

        if r.Method == "GET" {
            // fmt.Println("Do something with GET")
            fmt.Println(_echostr)
            fmt.Fprintf(w, weibot.Echostr)

        } else if r.Method == "POST" {
            xmlbody, err := ioutil.ReadAll(r.Body)
            if err != nil {
                fmt.Printf("error: %v\n", err)
            } else {
                weibot.OnMessage(xmlbody)
                replay := weibot.ReplayMessage()
                text_content := "Content for text type test!"
                replay.Text(text_content)
                fmt.Fprintf(w, replay.ToString())
            }
        }

    } else {
        fmt.Println("Invalidate!")
    }

}

func main() {
    http.HandleFunc("/gWebot", ReciveMsgHandler)
    err := http.ListenAndServe("127.0.0.1:8000", nil)
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}
