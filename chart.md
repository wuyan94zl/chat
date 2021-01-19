## 说明
该包实现了websocket聊天分组（组成员多少无限制，组个数无限制）服务  

### 要求
核心扩展 `gorm.io/gorm` ， `gorilla/websocket`   
使用gorm数据库迁移创建表ChatMessage  

### 使用
gin 框架下
```go
// 伪代码
router := gin.Default() // 获取路由实例
router.GET("/ws", func(c *gin.Context) {
    // DB为 gorm 数据库连接信息 类型（*gorm.DB）
    chart_server.Server(c.Writer, c.Request, DB, c.Query("channel_id"), 123)
})
```
js
```html
<script type="text/javascript">
var ws;
var channel_id = "456";
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
};
//连接关闭时触发
ws.onclose = function(evt) {
    console.log("关闭时触发")
};

// 发送消息
function sendMessage() {
    var obj = '{"channel_id":"' + channel_id + '","content":"发送的消息","type": "0"}'
    ws.send(obj)
}
```
