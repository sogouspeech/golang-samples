# 语音识别服务Go语言使用示例

本目录包含 [语音识别API](https://zhiyin.sogou.com/product/recognition/) 的Go语言使用示例

## 示例

### recognize

`recognize` 命令发送音频(不超过 ~1分钟)到搜狗知音平台并打印出文本转写结果。

`streaming-recognize` 命令实时识别流式语音(时间限制见[配额和限制](https://docs.zhiyin.sogou.com/docs/asr/resources/quotas))到搜狗知音平台并打印出文本转写结果。
