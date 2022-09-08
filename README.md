# gitdump

![workflow badge](https://github.com/guoyk93/gitdump/actions/workflows/go.yml/badge.svg) [![Go Reference](https://pkg.go.dev/badge/github.com/guoyk93/gitdump.svg)](https://pkg.go.dev/github.com/guoyk93/gitdump)

A tool for dumping hosted git repositories to local in batch.

## Features

* supported hosted git repositories
  * `github`
  * `gitee`
  * `gitea`
  * `coding`

## 中文使用说明

* [gitdump - 批量备份托管 Git 仓库](https://mp.weixin.qq.com/s?__biz=Mzg2ODIyNzg2Ng==&mid=2247483660&idx=1&sn=7096b013e5bab2896b9be90df12ac5cd&chksm=ceaecf79f9d9466f5ea56cf588622925c9b026749499ab40d6fe7281b2c60c4449a52b721c40#rd)

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

## Notification

Execution result will be delivered to environment variable `$NOTIFY_URL`, if given, by HTTP `POST`.

```
{"text": "MESSAGE..."}
```

## Upstream

https://git.guoyk.net/go-guoyk/gimir

Due to various reasons, codebase is detached from upstream.

## Donation

![oss-donation-wx](https://www.guoyk.net/oss-donation-wx.png)

## Credits

Guo Y.K., MIT License
