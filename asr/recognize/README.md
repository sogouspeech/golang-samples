# 搜狗知音平台Go语言示例

## 认证

* 参考文档 [快速指南](https://docs.zhiyin.sogou.com/docs/asr/quickstart/scli) 获取鉴权所需appid、token
* 运行下述命令设置鉴权信息
```bash
export SOGOU_SPEECH_ENDPOINT=$(scli show address)
export SOGOU_SPEECH_APPID=$(scli show appid)
export SOGOU_SPEECH_TOKEN=$(scli show token)
```

## 运行示例文件

安装本示例：

```bash
go get -u github.com/sogouspeech/golang-samples/asr/recognize
```

进行语音识别：

```bash
recognize /path/to/pcm/file
```

