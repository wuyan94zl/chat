## 说明
基于websocket实现的聊天功能  

### 要求
核心扩展 `gorm.io/gorm` ， `gorilla/websocket`     

### 使用

使用 `gorm` 连接 `mysql` 数据库 `DB` 执行迁移表 `DB.AutoMigrate(chart.ChatMessage{})`

gin 框架下
```go
package main
// 伪代码
import (
    "github.com/wuyan94zl/chart"
)
router := gin.Default()
router.GET("/ws", func(c *gin.Context) {
    // DB为 gorm 数据库连接信息 类型（*gorm.DB）
    // c.Query("channel_id") 为管道ID 多个以 `,` 隔开,相当于当前用户能接收到多个管道的信息
    // 123 为用户ID
    chart.Server(c.Writer, c.Request, c.Query("channel_id"), 123, DB)
})
router.Run(":8888")
```
js
```html
<script type="text/javascript">
var ws;
var channel_id = "123,456,789";
var port = '8888'
ws = new WebSocket("ws://localhost:" + port + "/ws?channel_id=" + channel_id);
//连接打开时触发
ws.onopen = function(evt) {
    console.log("连接打开时触发")
};
//接收到消息时触发
ws.onmessage = function(evt) {
    let data = $.parseJSON(evt.data)
    console.log(data)
    // data.channel_id 为管道id
};
//连接关闭时触发
ws.onclose = function(evt) {
    console.log("关闭时触发")
};

// 发送消息
function sendMessage() {
    // 像管道 123 发送消息
    let obj = '{"channel_id": "123","content": "发送的消息","type": "0"}'
    ws.send(obj)
    // 像管道 456 发送消息
    obj = '{"channel_id": "456","content": "发送的消息","type": "0"}'
    ws.send(obj)
}
```
