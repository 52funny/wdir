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

### Docker

```sh
docker run -d --name wdir -p 9194:8080 -v /Users/52funny:/mnt  52funny/wdir
```

8080 是容器内部的端口，9194 是你要映射本地的端口。/mnt 是容器内要索引的目录, /Users/52funny 是本地映射到容器内的目录。

## 配置

在 wdir 目录下有`config.yaml`, 编码为`UTF-8`

配置模版

```yaml
config:
  port: 8080
  path: /Users/52funny
  log_path: log
  show_hidden_files: false
```

| 配置项                   | 配置说明     |
| ------------------------ | ------------ |
| config.port              | 端口         |
| config.path              | 索引目录     |
| config.log_path          | 日志目录     |
| config.show_hidden_files | 显示隐藏目录 |

## 第三方库

- [fasthttp](https://github.com/valyala/fasthttp)
- [viper](https://github.com/spf13/viper)
- [filetype](https://github.com/h2non/filetype)
- [bulma](https://github.com/jgthms/bulma)
- [Fuse.js](https://github.com/krisk/Fuse)
