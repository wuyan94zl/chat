## 说明
基于websocket实现的聊天服务
### 提供 4 个回调接口供业务端自定义实现
```go
// 定义业务处理结构体，实现下面接口
type data struct {
}
// SendMessage ws 发送消息时回调，可以在这里处理消息记录
func (d *data) SendMessage(msg chart.Message) {
}
// LoginServer 登录成功后回调
func (d *data) LoginServer(uid uint64) {
}
// 退出登陆回调
func (d *data) LogoutServer(uid uint64) {
}
// 错误消息回调
func (d *data) ErrorLogServer(err error) {
}
```
### 提供3个内部接口实现聊天扩展
```go
    // uid 用户 加入 channel_one，channel_two 频道的监听
    chart.JoinChannelIds(uid, "channel_one", "channel_two")
	// uid 用户 退出 channel_one，channel_two 频道的监听 
	chart.UnJoinChannelIds(uid, "channel_one")
	// 业务端 以uid用户主动发送channel消息
	chart.SendMessageToChannelIds(uid, "test", "channel_one", "channel_two")
```

### 依赖
依赖扩展 `gorilla/websocket`     

### 使用
`go get github.com/wuyan94zl/chart`

### 完整示例
```go
package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/wuyan94zl/chart"
)

// 定义业务处理结构体，实现下面接口
type data struct {
}

// SendMessage ws 发送消息时回调，可以在这里处理消息记录
func (d *data) SendMessage(msg chart.Message) {
	fmt.Println("send message callback ", msg.ChannelId, msg.Content, msg.Type, msg.SendTime, msg.UserId)
}

// LoginServer 登录成功后回调
func (d *data) LoginServer(uid uint64) {
	fmt.Println("login callback ", uid)
	// 业务逻辑查询出用户应该监听 "channel_one", "channel_two" 2个 channel
	// 监听 "channel_one", "channel_two" 频道
	chart.JoinChannelIds(uid, "channel_one", "channel_two")
	// 退出监听 "channel_one" 频道 监听
	chart.UnJoinChannelIds(uid, "channel_one")

	// 外部业务端发送channel消息
	fmt.Println(chart.SendMessageToChannelIds(uid, "test", "channel_one", "channel_two"))
}
func (d *data) LogoutServer(uid uint64) {
	// 退出登陆回调
	fmt.Println("logout ", uid)
}
func (d *data) ErrorLogServer(err error) {
	// 错误消息回调
	fmt.Println("err: ", err)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		uid, _ := strconv.Atoi(r.FormValue("id"))
		chart.NewServer(w, r, uint64(uid), &data{})
	})
	err := http.ListenAndServe(":8888", mux)
	if err != nil {
		return
	}
}

```
html
```html
<button onclick="sendMessage()">send</button>
<script type="text/javascript">
    ws = new WebSocket("ws://localhost:8888/ws?id=888");
    //连接打开时触发
    ws.onopen = function (evt) {
        console.log("连接成功")
    };
    //接收到消息时触发
    ws.onmessage = function (evt) {
        console.log("收到消息：",evt.data)
        // data.channel_id 为管道id
    };
    //连接关闭时触发
    ws.onclose = function (evt) {
        console.log("关闭时触发")
    };
    // 发送消息
    function sendMessage() {
        // 向频道 channel_one 发送消息
        let obj = '{"channel_id": "channel_one","content": "发送的消息","type": "3"}'
        ws.send(obj)
        // 向频道 channel_two 发送消息
        obj = '{"channel_id": "channel_two","content": "发送的消息","type": "3"}'
        ws.send(obj)
    }
</script>
```
**运行 `go run main.go`**  
浏览器打开html文件即可
