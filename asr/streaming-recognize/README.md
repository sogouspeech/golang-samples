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

假设您已成功安装了Go开发环境，但本示例额外依赖于 [PortAudio](http://www.portaudio.com/) 抓取麦克风音频。

在Linux下，您可以使用包管理工具安装PortAudio，例如 apt：

```bash
apt install portaudio19-dev
```

在MacOS下，您可以使用brew工具安装PortAudio以及pkg-config：

```bash
brew install portaudio
brew install pkg-config
```

在Windows下，为了让cgo能够工作，您需要首先安装gcc。
这里演示通过 [scoop](https://scoop.sh/) 安装gcc，假设您已安装scoop，进入Powershell:

```bash
# 假设当前在 D:/somepath/portaudio 目录

# 安装后可能需要您注销再登录才能直接通过命令行调用gcc
scoop install wget gcc tar pkg-config

# 直接下载msys2环境下编好的portaudio库
wget http://repo.msys2.org/mingw/x86_64/mingw-w64-x86_64-portaudio-190600_20161030-3-any.pkg.tar.xz
tar -xf mingw-w64-x86_64-portaudio-190600_20161030-3-any.pkg.tar.xz

# 进入解压后的目录，将 D:/somepath/portaudio/mingw64/lib/pkgconfig/portaudio-2.0.pc 文件中
# prefix=/mingw64 改为 prefix=D:/somepath/portaudio/mingw64
# 注意要使用 '/' 而不是 '\'

# 之后，要导出环境变量供之后使用：
$env:Path += ";D:\somepath\portaudio\mingw64\bin"
$env:PKG_CONFIG_PATH = "D:\Develop\portaudio\mingw64\lib\pkgconfig"
```


安装本示例：

```bash
go get -u github.com/sogouspeech/golang-samples/asr/streaming-recognize
```

进行语音识别：

```bash
streaming-recognize
```
