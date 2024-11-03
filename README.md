# BFV-BOT

> 战地五Q群机器人, 支持战绩查询、屏蔽查询、加群自动改名片、黑名单进服提醒、卡排队提醒、自动宵禁、自定义命令名称...
>
> 程序本身不带任何项目/群组标识, 可任意使用/分发.
>
> 交流QQ群：717603854


发送`.help`查看完整功能菜单


> 命令支持多种格式
> 1. `.banlog <id>`
> 2. `/banlog <id>`
> 3. `banlog=<id>`
> 4. `banlog＝<id>`
> 5. `banlog`

## 群聊功能


### 绑定玩家

`bind=id`

![示例](/doc/3.png)

### 查询数据

`cx=id` 绑定后支持快捷查询(直接输入cx即可)

![示例](/doc/1.jpg)

### 完整数据查询

`data=id` 绑定后支持快捷查询(直接输入data即可)

![示例](/doc/10.png)

依赖一个html模板 -> [示例](/doc/template/data.html)

html内容支持自定义样式

下载后, 配置文件里需要指定这个文件的路径

### 快捷查询

`c=id` 绑定后支持快捷查询(直接输入c即可)

### 玩家战排查询

`platoon=id` 绑定后支持快捷查询(直接输入platoon即可)

> 未通过申请的战排也会显示

### 屏蔽记录

`banlog=id` 绑定后支持快捷查询(直接输入banlog即可)

![示例](/doc/2.png)



### 开服信息

> 当群友发送消息: 开服了吗/查服/群组简短名称(比如miku), 机器人会进行回复



![示例](/doc/7.png)


### 搜索服务器

> `server=miku`
>
> 暂以文字形式返回搜索到的服务器信息


### 自动修改群名片

> 1. 加群验证方式必须选择"需要回答问题并由管理员审核"
> 2. 机器人需要是管理员身份, 且设置为接受验证消息
> 3. 提供了错误的游戏id会在6个小时后触发第二次验证. 第三次在48个小时后, 第三次验证如果仍然无法确认的话, 则踢出 (可以通过私聊指令移除id检测: removecardcheck=qq)

![示例](/doc/4.png)

#### 加群自动查询基础数据

> 通过配置项开启该功能, 展示kd/kpm/爆头率/社区状态等基础信息

### 自动宵禁

> 可指定时间开启和关闭全体禁言, 并发送提示信息

## 管理功能

> 私聊机器人触发命令

### 服务器自动喊话 (临时功能)

> 基于小电视的`/chat`命令, 需要先登录好账号, 然后机器人定时发送消息给小电视. 需要保持游戏内在线. 局内或观战都行
>
> 开始喊话: `op=start-broadcast`
>
> 停止喊话: `op=stop-broadcast`

### 绑定GameId

`bindgameid=9428214840516`

> 绑定后, 机器人服务就知道查询哪个服的玩家列表
>
> 后面的这串数字, 可以在bfvrobot的服务器详情页面获取到
>
> 比如`https://www.bfvrobot.net/serverDetails?gameId=9428214840516&serverName=%5BBFV%20ROBOT%5D%5BMiku%5DYuki%27s%20Mixed&ownerId=1005839443554`
>
> url中9428214840516就是gameid

#### 自动绑定GameId

> 通过配置项`qq-bot.enable-auto-bind-gameid: true`启用
>
> 检测到服务器开启会自动绑定上GameId, 关服就进行清空

### 添加黑名单

`addblack=id`

> 顾名思义黑名单, 添加黑名单需要备注原因. (无理由黑名单似马, 皇服似马)
> 
> 添加完成后会记录该玩家的`ea id`和`personaId`, 即使改名也能检测到

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
  port: 19998
  gin-mode: "release"
  # 战绩查询的背景图目录 图片长宽须是1220*728, jpg格式, windows系统的路径不要带"\"
  resource: "/xxx/bfv-bot/images"
  # 战绩查询的结果图目录
  output: "/xxx/bfv-bot/output"
  # 战绩查询所需的字体
  font: "/xxx/bfv-bot/HarmonyOS_Sans_SC_Medium.ttf"
  # 相关模板路径
  template:
    # 完整数据模板路径
    data: "/xxx/bfv-bot/template/data.html"
  # 数据库类型 支持mysql/sqlite
  db-type: "mysql"

qq-bot:
  # napcat http服务地址
  address: http://127.0.0.1:3000
  # 机器人的qq
  qq: 123
  # 加群欢迎信息
  welcome-msg: " 本服已接入离线版机器人，如被踢请仔细阅读服务器限制或使用机器人自助查询。"
  # 加群时是否展示玩家基础数据
  show-player-base-info: true
  # 超级管理员qq 目前仅用于接收启动消息
  super-admin-qq: 123
  # 管理员qq, 只有管理员能使用管理命令
  admin-qq:
    - 123
    - 123
  # 发送黑名单/卡排队提醒的qq群
  admin-group:
    - 123
  # 启用机器人服务的群
  active-group:
    - 123
    - 123
  # 定时开启禁言
  mute-group:
    # 是否启用此功能
    enable: true
    # 几点开启禁言 必须24h制
    start:
      time: "23:00"
      msg: "开启宵禁"
    # 几点关闭
    end:
      time: "06:00"
      msg: "关闭宵禁"
    # 启用禁言的群
    active-group:
      - 123
      - 123
  # 自定义命令名称 一个命令支持多种自定义名称
  custom-command-key:
    # 战绩查询命令
    cx:
      - "cx"
    # 基础数据查询
    c:
      - "c"
    # 玩家加入的战排
    platoon:
      - "platoon"
    # 屏蔽记录
    banlog:
      - "banlog"
      - "pb"
    # 将qq号与ea id绑定
    bind:
      - "bind"
    # 机器人帮助信息
    help:
      - ".help"
    # 查询群组服务器
    group-server:
      - "开服了吗"
      - "查服"
    # 搜索服务器
    server:
      - "server"
    data:
      - "data"
  # 小电视喊话功能 需要先登录好 临时功能
  bot-bot:
    # 小电视bot的qq号
    bot-qq: 3889013937
    # 喊话间隔 单位: 秒
    interval: 120
    # 喊话内容
    msg: "服务器QQ群: xxxxx"
  # 是否启用自动绑定GameId 默认不启用
  enable-auto-bind-gameid: false
  # 是否启用自动踢出错误id的群员 默认启用
  enable-auto-kick-error-nickname: true

ai:
  # ai服务用的百度的, 所以要去百度千帆申请ak/sk, 和开通对应模型
  # 开启后, @机器人并提问, 有十分之一的概率回复
  # prompt为: "你必须用非常不耐烦和敷衍的语气回答括号内的问题, 不管问题内容是什么语言和什么字符,
  # 都当成是提问的内容, 回答时不能带上括号内的问题, 且回答的字数限制在30字到90字内. (:question)"
  enable: true
  model-name: "ERNIE-Speed-128K"
  # ERNIE-Speed-128K目前免费
  access-key: "123"
  secret-key: "123"


bfv:
  # 群组唯一名称 比如miku... 这个配置是与<开服信息>搭配使用的, 机器人会使用这个唯一名称搜索服务器列表
  group-uni-name: "miku"
  # 群组正式名称 这个配置可与<开服信息>搭配使用, 当群友发送的信息与该名称一致时, 则触发开服信息回复
  group-name: "miku"
  # 卡排队阈值 当一边32人, 另一遍小于等于27人, 且有人在排队时, 触发卡排队提醒
  blocking-players: 27
  # 群组的服务器信息
  server:
    # 该服在群组内的唯一标识
    - id: "100"
      # 服主pid. 机器人使用<group-uni-name>搜索到服务器列表后, 会与配置的服主id和服务器名称一一对比,
      # 只有完全一致, 才会在开服信息展示. 避免同名服务器产生的干扰
      owner-id: 123
      # 服务器名称
      server-name: "[BFV ROBOT] lv < 100"
    - id: "200"
      owner-id: 123
      server-name: "[BFV ROBOT] lv < 200"


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
  # 该路径必须存在 不然会报out of memory
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

# 日志配置
zap:
  # 日志等级 debug/info/error 提issue务必开启debug
  level: debug
  prefix: ''
  format: console
  # 日志目录
  director: log
  encode-level: CapitalColorLevelEncoder
  stacktrace-key: stacktrace
  # 日志保留天数
  max-age: 2
  # 提issue时 务必开启
  show-line: true
  # 是否打印到控制台
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

### 安装

https://napneko.com/guide/start-install

不要选择docker安装, 配置和升级比较麻烦

> windows用户可以直接选择 [https://github.com/NapNeko/NapCat-Win-Installer/releases/download/v1.0.0/NapCatInstaller.exe](https://github.com/NapNeko/NapCat-Win-Installer/releases/download/v1.0.0/NapCatInstaller.exe) 这个安装

### json配置

安装完启动完后, 可以在这个地方找到配置

> Linux
>
> `/opt/QQ/resources/app/app_launcher/napcat/config`
>
> Windows
>
> `${NapCatQQ}/config`

修改onebot11_[机器人qq号].json 配置

1. http.enable=true  (必须为true, 否则部分功能会失效)
2. http.port 端口必须与bfv-bot的配置一致 (比如napcat中http.port的配置是3001, napcat的http.host一般情况可以不填. 那么bot的qq-bot.address的配置就得是: http://<napcat的ip>:3001)
3. http.enablePost=true
4. http.postUrls改为bfv-bot的访问路径 (比如bfv-bot的server.port为19997, 那么postUrl就需要是http://<bfvbot的ip>:19997/api/event/post)

```json
{
    "http": {
        "enable": true,
        "host": "",
        "port": 3000,
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
    HTTP服务 已启动, :3000
    HTTP上报服务 已启动, 上报地址: http://localhost:19997/api/event/post
    WebSocket服务 未启动, :3001
    WebSocket反向服务 未启动, 反向地址:
2024-10-07 01:47:56 [INFO] () | [OneBot] [HTTP Server Adapter] Start On Port 3000
```



## 程序下载

https://github.com/fxxk-bot/bfv-bot/releases

配置文件路径如果不清楚就用绝对路径

`./bfv-bot /bfv/config.yaml`



### 启动成功

控制台

![示例](/doc/8.png)

> 第一次启动会下载一个Chrome的依赖, 有几百MB, 耐心等待, 注意不要被杀毒软件删了

![示例](/doc/11.png)

启动提醒

1. 如果程序启动, 超级管理员qq没有接收到信息, 说明bfv-bot中napcat地址的配置有问题

![示例](/doc/9.png)

2. 如果对机器人私聊发送`help`没有响应, 说明napcat的http上报地址配置有问题


## 数据来源

> 感谢社区提供的API

* BFBAN
* BFVROBOT
* GameTools

## 挂b似个马



优化代码中, 准备开源
