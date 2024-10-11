### 出现问题的功能

例如

> 查询功能, cx=xxxx, 失败

### log文件夹内的日志

例如

error.log

```
2024/10/07 - 02:27:35.290 ERROR 请求失败
```


info.log

```
2024/10/07 - 02:27:35.290	[34mINFO[0m	bfv-bot/main.go:56	服务成功启动在 Port: 19997
2024/10/07 - 02:27:47.042	[34mINFO[0m	initialize/gin.go:32	路由注册完成
2024/10/07 - 02:27:47.519	[34mINFO[0m	initialize/cron.go:43	定时任务初始化完成
2024/10/07 - 02:27:47.520	[34mINFO[0m	initialize/sensitive.go:18	加载到4条敏感词
2024/10/07 - 02:27:47.707	[34mINFO[0m	bfv-bot/main.go:56	服务成功启动在 Port: 19997
```

debug.log

```
2024/10/07 - 02:27:09.338	[35mDEBUG[0m	initialize/ai.go:18	支持的xx
2024/10/07 - 02:27:09.351	[35mDEBUG[0m	initialize/ai.go:19	xxxx-3.5-4K-0205, xxxx Speed
```



### 版本信息

* 服务启动后发送的Git Branch/Commit信息

例如

> Git Branch: master
> 
> Git Commit: 0ecb2b4e5ec725fac2540e9c0c8a8e17ed79c974
> 
> Build Time: 2024-10-07T00:43:12
