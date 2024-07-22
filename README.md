# wdir

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/52funny/wdir)
![GitHub](https://img.shields.io/github/license/52funny/wdir)

English | [简体中文](https://github.com/52funny/wdir/blob/master/README_zhCN.md)

Directory indexing system based on golang

## Preview

![1](https://raw.githubusercontent.com/52funny/wdir/master/pics/1.png)
![2](https://raw.githubusercontent.com/52funny/wdir/master/pics/2.png)
![3](https://raw.githubusercontent.com/52funny/wdir/master/pics/3.png)

## Usage

```sh
git clone https://github.com/52funny/wdir
cd wdir
go build
./wdir
```

There is a default configuration file `config.yaml` in the wdir folder. you can specify the configuration file with the `-c` command.
`./webd -c config`

### Docker

```sh
docker run -d --name wdir -p 9194:8080 -v /Users/52funny:/mnt  52funny/wdir
```

8080 is the port inside the container, 9194 is the port you want to map locally. /mnt is the directory to be indexed inside the container, /Users/52funny is the directory mapped locally to the container.

## Configuration

There is `config.yaml` in the wdir directory, encoded as `UTF-8`

Configuration template

```yaml
config:
  port: 8080
  path: /Users/52funny
  log_path: log
  show_hidden_files: false
```

| Configuration item       | Configuration instructions |
| ------------------------ | -------------------------- |
| config.port              | Server port                |
| config.path              | Index catalog              |
| config.log_path          | Log catalog                |
| config.show_hidden_files | Show Hidden Files          |

## Third Party Library

- [fasthttp](https://github.com/valyala/fasthttp)
- [viper](https://github.com/spf13/viper)
- [filetype](https://github.com/h2non/filetype)
- [bulma](https://github.com/jgthms/bulma)
- [Fuse.js](https://github.com/krisk/Fuse)
