# wdir

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/52funny/wdir)
![GitHub](https://img.shields.io/github/license/52funny/wdir)

[English](https://github.com/52funny/wdir/blob/master/README.md) | 简体中文

基于 Go 的目录索引

## 预览

![1](https://raw.githubusercontent.com/52funny/wdir/master/pics/1.png)
![2](https://raw.githubusercontent.com/52funny/wdir/master/pics/2.png)
![3](https://raw.githubusercontent.com/52funny/wdir/master/pics/3.png)

## 用法

```sh
git clone https://github.com/52funny/wdir
cd wdir
go build
./wdir
```

在 wdir 文件夹中有默认的配置文件`config.yaml`。你可以用`-c`命令来指定配置文件。
`./wdir -c config`

## 配置

在 wdir 目录下有`config.yaml`, 编码为`UTF-8`

配置模版

```yaml
config:
  port: 8080
  template: compress
  path: /Users/52funny
  logpath: log
```

| 配置项          | 配置说明 |
| --------------- | -------- |
| config.port     | 端口     |
| config.template | 模版目录 |
| config.path     | 索引目录 |
| config.logpath  | 日志目录 |

## 第三方库

- [fasthttp](https://github.com/valyala/fasthttp)
- [viper](https://github.com/spf13/viper)
- [filetype](https://github.com/h2non/filetype)
- [bulma](https://github.com/jgthms/bulma)
- [Fuse.js](https://github.com/krisk/Fuse)
