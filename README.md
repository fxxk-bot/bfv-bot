# BFV-BOT

> 战地五Q群机器人, 支持战绩查询、屏蔽查询、加群自动改名片、黑名单进服提醒...

发送`.help`查看完整功能菜单

## 群聊功能


### 绑定玩家

`bind=id`

![示例](/doc/3.png)

### 查询数据

`cx=id` 绑定后支持快捷查询(直接输入cx即可)

![示例](/doc/1.jpg)

### 屏蔽记录

`banlog=id` 绑定后支持快捷查询(直接输入banlog即可)

![示例](/doc/2.png)



### 开服信息

> 当群友发送消息: 开服了吗/查服/群组简短名称(比如miku), 机器人会进行回复



![示例](/doc/7.png)




### 自动修改群名片

> 1. 加群验证方式必须选择"需要回答问题并由管理员审核"
> 2. 机器人需要是管理员身份, 且设置为接受验证消息
> 3. 提供了错误的游戏id会在6个小时后触发第二次验证. 第三次在48个小时后, 第三次验证如果仍然无法确认的话, 则踢出 (可以通过私聊指令移除id检测: removecardcheck=qq)

![示例](/doc/4.png)

## 管理功能

> 私聊机器人触发命令

### 绑定GameId

`bindgameid=9428214840516`

> 绑定后, 机器人服务就知道查询哪个服的玩家列表
>
> 后面的这串数字, 可以在bfvrobot的服务器详情页面获取到
>
> 比如`https://www.bfvrobot.net/serverDetails?gameId=9428214840516&serverName=%5BBFV%20ROBOT%5D%5BMiku%5DYuki%27s%20Mixed&ownerId=1005839443554`
>
> url中9428214840516就是gameid

### 添加黑名单

`addblack=id`

> 顾名思义黑名单, 添加黑名单需要备注原因

### 移除黑名单

`removeblack=id`

### 移除id检测

`removecardcheck=qq`

> 参考自动修改群名片功能介绍

### 添加敏感词

`addsensitive=xxx`

> 发送的消息内含有指定内容消息就进行撤回并发送警告信息

### 移除敏感词

`removesensitive=xx`

### 开始检测黑名单

`op=start`

#### 服务器黑名单玩家进服提醒

> 需要添加了黑名单和绑定GameId之后, 并且启动黑名单检测, 才会接受到消息提醒

![示例](/doc/5.png)

#### 卡排队提醒

> 需要绑定GameId之后, 并且启动黑名单检测, 才会接受到消息提醒

![示例](/doc/6.png)

### 停止检测黑名单

`op=stop`

> 同时也会关闭卡排队提醒

## 配置文件

```yaml
server:
  # 机器人服务的端口, 后面配置napcat会用到
  port: 19997
  gin-mode: "release"
  # 战绩查询的背景图目录 图片长宽须是1220*728
  resource: "/bfv-bot/images"
  # 战绩查询的结果图目录
  output: "/bfv-bot/test"
  # 战绩查询所需的字体
  font: "/bfv-bot/bfv-font/HarmonyOS_Sans_SC_Medium.ttf"
  # 数据库类型 支持mysql/sqlite
  db-type: "mysql"

qq-bot:
  # napcat http服务地址
  address: http://192.168.93.130:3001
  # 机器人的qq
  qq: 123123
  # 加群欢迎信息
  welcome-msg: " 本服已接入机器人，如被踢请仔细阅读服务器限制或使用机器人自助查询。"
  # 管理员qq, 只有管理员能使用管理命令
  admin-qq: 123123
  # 发送黑名单/卡排队提醒的qq群
  admin-group: 123123
  # 启用机器人服务的群
  active-group:
    - 123123


bfv:
  # 群组唯一名称 比如miku... 这个配置是与<开服信息>搭配使用的, 机器人会使用这个唯一名称搜索服务器列表
  group-uni-name: "xxx"
  # 群组正式名称 这个配置可与<开服信息>搭配使用, 当群友发送的信息与该名称一致时, 则触发开服信息回复
  group-name: "xxx"
  # 卡排队阈值 当一边32人, 另一遍小于等于27人, 且有人在排队时, 触发卡排队提醒
  blocking-players: 27
  # 群组的服务器信息
  server:
      # 该服在群内的唯一标识
    - id: "100"
      # 服主pid. 机器人使用<group-uni-name>搜索到服务器列表后, 会与配置的服主id和服务器名称一一对比, 只有完全一致, 才会在开服信息展示. 避免同名服务器产生的干扰
      owner-id: 100854811xxxx
      # 服务器名称
      server-name: "[BFV ROBOT]xxx lv<100"


# 数据库配置 略
mysql:
  url: "localhost"
  port: "3306"
  config: "charset=utf8mb4&parseTime=True&loc=Local"
  db-name: "bfv_bot"
  username: "root"
  password: "123456"
  prefix: ""
  singular: false
  engine: ""
  max-idle-conns: 10
  max-open-conns: 100
  log-mode: error
  log-zap: true

sqlite:
  path: "/bfv-bot/"
  port: "3308"
  config: "charset=utf8mb4&parseTime=True&loc=Local"
  db-name: "bfv_bot"
  username: "root"
  password: "123456"
  max-idle-conns: 10
  max-open-conns: 100
  log-mode: "error"
  log-zap: true

ai:
  # 是否启用ai回复 当前暂不开放
  enable: false

# 日志配置
zap:
  # 日志等级
  level: debug
  prefix: ''
  format: console
  # 日志文件存放在哪个文件夹
  director: log
  encode-level: CapitalColorLevelEncoder
  stacktrace-key: stacktrace
  max-age: 2
  show-line: true
  log-in-console: true
```

对应SQL文件.

```sql
-- ----------------------------
-- Table structure for bind
-- ----------------------------
DROP TABLE IF EXISTS `bind`;
CREATE TABLE `bind`  (
  `qq` bigint(0) NOT NULL,
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `pid` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  PRIMARY KEY (`qq`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for blacklist
-- ----------------------------
DROP TABLE IF EXISTS `blacklist`;
CREATE TABLE `blacklist`  (
  `id` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL,
  `reason` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for card_check
-- ----------------------------
DROP TABLE IF EXISTS `card_check`;
CREATE TABLE `card_check`  (
  `qq` bigint(0) NOT NULL,
  `group_id` bigint(0) NOT NULL,
  `fail_cnt` bigint(0) NOT NULL,
  `next_check_time` bigint(0) NOT NULL,
  PRIMARY KEY (`qq`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for ignorelist
-- ----------------------------
DROP TABLE IF EXISTS `ignorelist`;
CREATE TABLE `ignorelist`  (
  `id` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for sensitive
-- ----------------------------
DROP TABLE IF EXISTS `sensitive`;
CREATE TABLE `sensitive`  (
  `id` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = Dynamic;

SET FOREIGN_KEY_CHECKS = 1;
```



## napcat配置

### 一键脚本

https://napneko.github.io/zh-CN/guide/getting-started#%E4%B8%80%E9%94%AE%E6%92%B8%E7%8C%AB

不要选择docker安装, 配置和升级比较麻烦

### json配置

安装完启动完后, 可以在这个地方找到配置

`/opt/QQ/resources/app/app_launcher/napcat/config`

修改onebot11_[机器人qq号].json 配置

1. http.enable=true
2. http.port 端口必须与bfv-bot的配置一致 (比如napcat中http.port的配置是3001, napcat的http.host一般情况可以不填. 那么bot的qq-bot.address的配置就得是: http://<napcat的ip>:3001)
3. http.enablePost=true
4. http.postUrls改为bfv-bot的访问路径 (比如bfv-bot的server.port为19997, 那么postUrl就需要是http://<bfvbot的ip>:19997/api/event/post)

```json
{
    "http": {
        "enable": true,
        "host": "",
        "port": 3001,
        "secret": "",
        "enableHeart": false,
        "enablePost": true,
        "postUrls": ["http://192.168.93.1:19997/api/event/post"] 
    },
    "ws": {
        "enable": false,
        "host": "",
        "port": 3001
    },
    "reverseWs": {
        "enable": false,
        "urls": []
    },
    "GroupLocalTime": {
        "Record": false,
        "RecordList": []
    },
    "debug": false,
    "heartInterval": 30000,
    "messagePostFormat": "array",
    "enableLocalFile2Url": true,
    "musicSignUrl": "",
    "reportSelfMessage": false,
    "token": ""
}
```
改完后 重启`napcat`

正常的话就能看见控制台日志

```
Log":true,"fileLogLevel":"debug","consoleLogLevel":"info"}
2024-10-07 01:47:56 [WARN] () | [Native] Error: Native Not Init
2024-10-07 01:47:56 [INFO] () | [Notice] [OneBot11]
    HTTP服务 已启动, :3001
    HTTP上报服务 已启动, 上报地址: http://192.168.93.1:19997/api/event/post
    WebSocket服务 未启动, :3001
    WebSocket反向服务 未启动, 反向地址:
2024-10-07 01:47:56 [INFO] () | [OneBot] [HTTP Server Adapter] Start On Port 3001
```



## 程序下载

https://github.com/fxxk-bot/bfv-bot/releases

配置文件路径是绝对路径

`./bfv-bot /bfv/config.yaml`



### 启动成功

控制台

![示例](/doc/8.png)

qq消息

如果程序启动, qq没有接收到信息, 说明napcat地址的配置有问题

![示例](/doc/9.png)



优化代码中, 准备开源
