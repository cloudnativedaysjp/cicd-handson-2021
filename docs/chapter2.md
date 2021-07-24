# Chapter 2

本章では、Docker/GitHub Package/Kubernetesの基本的な操作についておさらいします。  
ローカル環境でgo言語を使って、Docker imageを作成し、Kubernetesで起動してみましょう。  

## 2-1 ワークディレクトリを作成する

以下コマンドを実行して、ワークディレクトリを作成します。  
※ここでは、"C:\GoApp"というディレクトリを作成する例を紹介します。

```cmd
cd C:\
mkdir GoApp
```

## 2-2 Go言語でサンプルアプリを作成する

ワークディレクトリに、以下のファイルを作成します。  
※ここでは、"Hello Docker!!"と表示されるGoサンプルアプリのソースコード例を紹介します。  
ファイル名：`main.go`

```go
package main

import (
  "fmt"
  "log"
  "net/http"
)

func main(){
  http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
    log.Println("received request")
    fmt.Fprintf(w, "Hello Dcoker!!")
  })

  log.Println("start server")
  server := &http.Server{Addr: ":8080"}
  if err := server.ListenAndServe(); err !=nil {
     log.Println(err)
    }
  }
```

## 2-3 Dockerfileを作成する

ワークディレクトリに、以下のファイルを作成します。  
ファイル名：`Dockerfile`

```Dockerfile
FROM golang:latest

#ディレクトリ作成
WORKDIR /go/src/go-image
#ホストOSのmain.goをWORKDIRにコピー
COPY main.go .

#バイナリを生成
RUN go install -v .

#バイナリを実行
CMD ["go-image"]
```

## 2-4 Docker imageをビルドする

以下コマンドを実行して、DockerfileからDocker imageをビルドします。

```cmd
docker build -t docker-go .
```

## 2-5 Docker image一覧を確認する

以下コマンドを実行して、作成されたDocker imageを確認します。

```cmd
docker images
```

## 2-6 Dockerコンテナを起動する

以下コマンドを実行して、作成されたDocker imageからDockerコンテナを起動します。

```cmd
docker run --rm go-image
```

## 2-7 Dockerコンテナ状態を確認する

以下コマンドを実行して、起動されたDockerコンテナを動作確認します。

* 動作しているコンテナのみ表示する場合

```cmd
docker ps
```

* 動作・停止しているコンテナを表示する場合

```cmd
docker ps -a
```

## 2-8 GitHub Packagesへpushする

```cmd

```

## 2-9 Kubernetesでの起動（kubectl apply）

```cmd

```
