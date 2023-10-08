# Family power monitor server 家庭用电监控服务端

用于接受监控模块的MQTT消息，并推送到[message-pusher-server](https://github.com/songquanpeng/message-pusher)，进而推送到用户企业微信（包括微信插件通知）

# Feature

- 接收电量MQTT消息，定时推送消息（企业微信bot）
- 部署参数可配置

# System structure

- [自行组装的统计模块](https://post.smzdm.com/p/aqxqv867/)
- [MQTT Server](https://hub.docker.com/_/eclipse-mosquitto)
- Go Monitor（current project）
- [Message Pusher](https://github.com/songquanpeng/message-pusher)

![power-monitor.jpg](doc%2Fpower-monitor.jpg)

![monitor.jpeg](doc%2Fmonitor.jpeg)

![WechatIMG696.jpeg](doc%2FWechatIMG696.jpeg)

# Quick Start

docker部署：

    docker run -d chajiuqqq/power-monitor-server

环境变量：

| ENV       | Default     | Description                    |
|-----------|-------------|--------------------------------|
| BROKER_IP | 127.0.0.1   | mqtt server ip                 |
| PORT      | 1883        | mqtt server port               |
| PROFILE   | dev         | dev/prod/test                  |
| CRON      | 0 0 1 * * * | 定时通知表达式（UTC时间），示例是每天早上北京时间9点推送 |



# TODO

v1.1

- [Y] 持久化电量MQTT消息
- 提供bot接收、处理消息的API
  - 获取当日统计、获取当前自然月统计
  - 重置全部、重置今天、重置昨天

v1.2

- 微信bot提供预充值统计模式
  - 电量余额维护（充钱操作累加操作、显示余额、显示充钱日志）
  - 余额阈值报警
- 用电历史统计（自然月、充值周期）