# gitdump

A tool for dumping hosted git repositories in batch.

## Features

* supported hosted git repositories
  * `github`
  * `gitee`
  * `gitea`
  * `coding`

## Usage

**Command**

```
./gitdump -conf config.yaml
```

**Configuration**

```yaml
dir: repos
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

## Upstream

https://git.guoyk.net/go-guoyk/gimir

Due to various reasons, codebase is detached from upstream.

## Credits

Guo Y.K., MIT License