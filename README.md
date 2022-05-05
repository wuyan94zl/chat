## 说明
基于websocket实现的聊天服务
### 提供5个回调接口供业务端自定义实现
```go
// 定义业务处理结构体，实现下面接口
type data struct {
}
// SendMessage ws 发送消息时回调，可以在这里处理消息记录
func (d *data) SendMessage(msg chart.Message) {
}
// 离线消息延时队列发送处理
func (d *data) DelaySendMessage(channelId string, msg Message, uids []uint64){
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
    chat.JoinChannelIds(uid, "channel_one", "channel_two")
	// uid 用户 退出 channel_one，channel_two 频道的监听 
	chat.UnJoinChannelIds(uid, "channel_one")
	// 业务端 发送消息到channel_id消息
	chat.SendMessageToChannelIds(uid, "test", "channel_one", "channel_two")
    // 业务端 发送消息给单个用户，uid和toUid可以相同
    chat.SendMessageToUid(uid, toUId uint64, msg string, tp uint8)
```

### 依赖
依赖扩展 `gorilla/websocket`     

### 使用
`go get github.com/wuyan94zl/chart`

### 完整示例
