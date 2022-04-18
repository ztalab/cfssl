# Golang cilog 日志组件通用包

## Tips

执行 `export GOSUMDB=off` 关闭代码检查   

## V2

**V2 版本目前仅支持输出到 LogProxy**

### Features
1. level 改造符合公司日志标准
2. 兼容 Zap 所有方法
3. 额外添加 `package.Info()` 系列方法
4. 支持动态自定义 Field
5. 支持启用 Caller 记录调用函数
6. Error 级别记录调用 Stack

默认情况下日志会同时输出到 stdout, 以及 Redis.    

### Usage

在目录 [v2/example](./v2/example) 查看使用示例.   

#### Notice
HookConfig 设置为 `nil`, 则不输出到 Redis.       
本地测试环境下请设置为 `nil`.       

## V1

基于 logrus 封装，日志输出到 Redis.      

在目录 [example](./example) 查看使用示例.   