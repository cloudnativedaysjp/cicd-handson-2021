# Chapter 2

本章では、Docker/GitHub Package/Kubernetesの基本的な操作についておさらいします。  
ローカル環境でgo言語を使って、Docker imageを作成し、Kubernetesで起動してみましょう。  

## 2-1 ワークディレクトリを作成する

#### コマンド実行
任意の場所にワークディレクトリを作成します。

```cmd
mkdir ~\GoApp
cd ~\GoApp
```

## 2-2 Go言語でサンプルアプリのソースコードを作成する

ワークディレクトリに、以下のファイルを作成します。  
※ここでは、"Hello Docker!!"と表示されるソースコード例を紹介します。  
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
#ベースイメージ指定
FROM golang:latest

#ディレクトリ作成
RUN mkdir /work
#ホストOSのmain.goをWORKDIRにコピー
COPY main.go /work

#Goアプリ実行
CMD ["go", "run", "/work/main.go"]
```

## 2-4 Docker imageをビルドする

#### コマンド実行
DockerfileからDocker imageをビルドします。

```cmd
docker build -t go-image:latest .
```

#### 実行結果

```
[+] Building 1.1s (8/8) FINISHED
 => [internal] load build definition from Dockerfile                                           0.0s
 => => transferring dockerfile: 32B                                                            0.0s
 => [internal] load .dockerignore                                                              0.0s
 => => transferring context: 2B                                                                0.0s
 => [internal] load metadata for docker.io/library/golang:latest                               0.9s
 => [internal] load build context                                                              0.0s
 => => transferring context: 29B                                                               0.0s
 => [1/3] FROM docker.io/library/golang:latest@sha256:4544ae57fc735d7e415603d194d9fb09589b8ad  0.0s
 => CACHED [2/3] RUN mkdir /work                                                               0.0s
 => CACHED [3/3] COPY main.go /work                                                            0.0s
 => exporting to image                                                                         0.0s
 => => exporting layers                                                                        0.0s
 => => writing image sha256:220026ab99c08d4d2592f66f728f078d1c48a8e4d1e14d77630e1497df058642   0.0s
 => => naming to docker.io/library/go-image:latest                                             0.0s
```

## 2-5 Docker image一覧を確認する

#### コマンド実行
作成されたDocker imageを確認します。

```cmd
docker images
```

#### 実行結果

```
REPOSITORY   TAG       IMAGE ID       CREATED          SIZE
go-image     latest    220026ab99c0   4 minutes ago   862MB
```

## 2-6 Dockerコンテナを起動する

#### コマンド実行
作成されたDocker imageからDockerコンテナを起動します。  
`-d`：コンテナをバックグラウンド実行します。

```cmd
docker run -d go-image
```

#### 実行結果

```
8d3a4f0c85e671284ec4cab7c973cbbf91318f54583059b15db4419f6903af56
```

## 2-7 Dockerコンテナ状態を確認する

#### コマンド実行
起動されたDockerコンテナを動作確認します。

```cmd
docker ps
```

#### 実行結果

```
CONTAINER ID   IMAGE      COMMAND                  CREATED          STATUS          PORTS           NAMES
8d3a4f0c85e6   go-image   "go run /work/main.go"   2 minutes ago    Up 2 minutes                    inspiring_sinoussi
```

`PORTS`に何も割り当たっていないことを確認します。

## 2-8 Goアプリのレスポンスを確認する

#### コマンド実行
GoアプリへGETリクエストを送信し、レスポンスを確認します。

```
curl http://localhost:8080
```

#### 実行結果

```
curl: (7) Failed to connect to localhost port 8080: Connection refused
```

`Connection refused`と表示され、8080ポートへの接続に失敗することを確認します。  

#### ここで、少し考えてみましょう。なぜ、接続に失敗するのでしょうか？  
Goアプリにて公開設定した8080ポートは、コンテナ内に限定されたコンテナポートであることを理解する必要があります。つまり、コンテナ内からアクセスする際は、コンテナポート(8080)を使用できますが、コンテナ外部(ローカルPC)からアクセスする際は、使用することができません。この場合には、ポートフォワーディング機能を利用し、ホストマシンのポートをコンテナポートに紐付け、コンテナ外との通信を可能にすることができます。以下で手順を確認します。

## 2-9 ポートフォワーディング設定をしてDockerコンテナを起動する

#### コマンド実行
作成されたDocker imageからDockerコンテナを起動します。
`-p`：{コンテナ外部側ポート}:{コンテナ内部側ポート}の書式で記述することが可能。

```cmd
docker run -d -p 9000:8080 go-image
```

#### 実行結果

```
d94c925240845c03b2f2dff0d43aea9d9b7c2f86184309e84b2cb4e93ff97c0a
```

## 2-10 Dockerコンテナ状態を確認する

#### コマンド実行
起動されたDockerコンテナを動作確認します。

```cmd
docker ps
```

#### 実行結果

```
CONTAINER ID   IMAGE      COMMAND                  CREATED          STATUS          PORTS                                       NAMES
d94c92524084   go-image   "go run /work/main.go"   15 seconds ago   Up 14 seconds   0.0.0.0:9000->8080/tcp, :::9000->8080/tcp   elegant_germain
```

`PORTS`にポートフォワーディング設定がされていることを確認します。

## 2-11 Goアプリのレスポンスを確認する

#### コマンド実行
8080ポート(コンテナ外部)⇔9000ポート(コンテナ内部)で公開されたGoアプリへGETリクエストを送信し、レスポンスを確認します。

```
curl http://localhost:9000
```

#### 実行結果

```
Hello Dcoker!!
```

9000ポートへの接続に成功し、`Hello Dcoker!!`と表示されることを確認します。

## 2-12 GitHub Packagesへpushする

```cmd

```

## 2-13 Kubernetesでの起動（kubectl apply）

```cmd

```
