# TikTink

前后端分离开发的短视频应用app的后端

V.0.1.0
更新：
1. Logic, DAO层引入自定义上下文
2. 基于自定义 context 实现请求处理链路追踪
3. logger 日志打印模块化封装处理

在controller层中不应该有自定义行为，需要定制化的
操作应该封装好后在controller层调用

如果需要调用下游提供的服务，应传入空白context，调用结束后会
将所有链路信息记录在context中。
