# gitdump

![workflow badge](https://github.com/yankeguo/gitdump/actions/workflows/go.yml/badge.svg) [![Go Reference](https://pkg.go.dev/badge/github.com/yankeguo/gitdump.svg)](https://pkg.go.dev/github.com/yankeguo/gitdump)

A tool for dumping hosted git repositories to local in batch.

## Features

- supported git hosting platforms
  - `github`
  - `gitee`
  - `gitea`
  - `coding`

## 中文使用说明

- [gitdump - 批量备份托管 Git 仓库](https://mp.weixin.qq.com/s/FLPoV5-9A-VuPit7QlFFYg)

## Usage

**Command**

```
./gitdump -conf config.yaml
```

**Configuration**

```yaml
dir: repos
concurrency: 3
accounts:
  - vendor: github
    # username, github username
    username: USERNAME
    # password, github personal token
    password: PERSONAL_TOKEN
  - vendor: gitee
    # username, gitee username
    username: USERNAME
    # password, gitee personal token
    password: PERSONAL_TOKEN
  - vendor: gitea
    # url, url of gitea instance
    url: https://your.gitea.com
    # username, gitea username
    username: USERNAME
    # password, gitea personal token
    password: PERSONAL_TOKEN
  - vendor: coding
    # url, url of coding instance
    url: https://your.coding.net
    # username, personal token username, displayed in token page, NOT YOUR CODING USERNAME
    username: TOKEN_USERNAME
    # password, personal token
    password: TOKEN
```

## Container Image

Check [GitHub Packages](https://github.com/yankeguo/gitdump/pkgs/container/gitdump) for available container images

Check [Dockerfile](Dockerfile) for details

By default, container image will execute `gitdump` every `6 hours`

All you need to do is to mount `/data` for data persistence, and put a `config.yaml` at `/data/config.yaml`

## Notification

Execution result will be delivered to environment variable `$NOTIFY_URL`, if given, by HTTP `POST`.

```
{"text": "MESSAGE..."}
```

## Credits

GUO YANKE, MIT License
