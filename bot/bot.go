package bot

import (
    "crypto/sha1"
    "encoding/xml"
    "fmt"
    "io"
    "reflect"
    "sort"
    "strings"
    // "time"
)

const (
    //Input your Token
    Token = "your_token_string_ever_guessed"
)

// type---------------

type WeiBot struct {
    token     string
    Signature string
    Timestamp string
    Nonce     string
    Echostr   string
    XMLBody   []byte
    Message   ReceivedMsg
    Replay    ReplayMsg
    // Timestamp time.Duration
}

type ReceivedMsg struct {
    XMLName      xml.Name `xml:xml`
    ToUserName   string   `xml:"ToUserName"`
    FromUserName string   `xml:"FromUserName"`
    CreateTime   string   `xml:"CreateTime"`
    MsgType      string   `xml:"MsgType"`
    MsgId        string   `xml:"MsgId"`

    //------text----------
    Content string `xml:"Content,omitempty"`

    //-------media--------
    PicUrl       string `xml:",omitempty"`
    MediaId      string `xml:"MediaId,omitempty"`
    Format       string `xml:"Format,omitempty"`
    ThumbMediaId string `xml:"ThumbMediaId,omitempty"`
    Location_X   string `xml:"Location_X,omitempty"`
    Location_Y   string `xml:"Location_Y,omitempty"`
    Scale        string `xml:"Scale,omitempty"`
    Label        string `xml:"Label,omitempty"`
    Title        string `xml:"Title,omitempty"`
    Description  string `xml:"Description,omitempty"`
    Url          string `xml:"Url,omitempty"`

    //------event---------
    Event string `xml:"Event,omitempty"`
}

type ReplayMsg struct {
    XMLName xml.Name `xml:"xml"`
    ReceivedMsg
}

func NewWeiBot() *WeiBot {
    return &WeiBot{token: Token}
}

//Register args
func (wbot *WeiBot) OnConnect(signature, timestamp, nonce, echostr string) {
    // fmt.Println("on_connect: ", signature, timestamp, nonce, echostr)
    wbot.Signature = signature
    wbot.Timestamp = timestamp
    wbot.Nonce = nonce
    wbot.Echostr = echostr

}

func (w *WeiBot) Validate() (is_true bool, echstr string) {
    li := []string{w.token, w.Timestamp, w.Nonce}
    sort.Strings(li)
    str_li := strings.Join(li, "")

    _Signature := rtn_sha1(str_li)
    // fmt.Println(_Signature)

    if _Signature == w.Signature {
        return true, w.Echostr
    } else {
        return false, ""
    }
}

func (w *WeiBot) OnMessage(xmlbody []byte) {
    // w.XMLBody = xmlbody
    w.parse_msg(xmlbody)
}

func (w *WeiBot) parse_msg(xmlbody []byte) {
    // err := xml.Unmarshal(w.XMLBody, &w.Message)
    err := xml.Unmarshal(xmlbody, &w.Message)
    if err != nil {
        fmt.Printf("error: %v\n", err)
        return
    }
}

func (w *WeiBot) ReplayMessage() ReplayMsg {
    var replay ReplayMsg
    replay.ToUserName = w.Message.FromUserName
    replay.FromUserName = w.Message.ToUserName
    replay.CreateTime = w.Message.CreateTime
    replay.MsgId = w.Message.MsgId
    return replay
}

func (r *ReplayMsg) ToString() string {
    output, err := xml.MarshalIndent(r, "  ", "    ")
    if err != nil {
        fmt.Printf("error: %v\n", err)
    }
    return string(output)
}

func (r *ReplayMsg) Text(content string) {
    r.MsgType = "text"
    r.Content = content
}

func (r *ReplayMsg) Img() {
    r.MsgType = "image"
}

func (r *ReplayMsg) Voice() {
    r.MsgType = "voice"
}

func (r *ReplayMsg) Video() {
    r.MsgType = "video"
}

func (r *ReplayMsg) Location() {
    r.MsgType = "location"
}

func (r *ReplayMsg) Link() {
    r.MsgType = "link"
}

// func (r *ReplayMsg) ReplayEvent() {
//  r.MsgType = "event"
// }

//对字符串进行sha1加密
func rtn_sha1(s string) string {
    h := sha1.New()
    io.WriteString(h, s)
    return fmt.Sprintf("%x", h.Sum(nil))
}

//遍历struct 打印字段名和值
func traver_struct(v interface{}) {
    value := reflect.ValueOf(v)
    for i := 0; i < value.NumField(); i++ {
        val := value.Field(i)
        fmt.Printf("Field %s: %v\n", val.Type().Name(), val.String())
    }
}

//struct convert to map
func struct_to_map(s interface{}) {
    m := make(map[string]interface{})
    value := reflect.ValueOf(s)
    for i := 0; i < value.NumField(); i++ {
        m[value.Type().Field(i).Name] = value.Field(i)
    }
    fmt.Println(m)
}
