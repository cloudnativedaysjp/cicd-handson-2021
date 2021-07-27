# Chapter 2

本章では、Docker/GitHub Package/Kubernetesの基本的な操作についておさらいします。  
ローカル環境でgo言語を使って、Docker imageを作成し、Kubernetesで起動してみましょう。  

## 2-1 Dockerの基本操作をおさらいする

ここでは、Goアプリを使って、Dockerイメージのビルドから起動までの操作を確認します。

### 2-1-1 ワークディレクトリを作成する

#### コマンド実行
任意の場所にワークディレクトリを作成します。

```cmd
$ mkdir ~\GoApp
$ cd ~\GoApp
```

### 2-1-2 Go言語でサンプルアプリのソースコードを作成する

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

### 2-1-3 Dockerfileを作成する

Dockerfileは、Dockerイメージのビルド時に、事前に実施しておきたい操作をコードとして記述したファイルです。  
主に、OSやミドルウェア、コマンド実行、デーモン実行、環境変数などの設定を記述することが可能で、  
Dockerfileを使ってイメージをビルドすることで、その設定をコードとして柔軟に管理することができます。  

ワークディレクトリに、以下のDockerfileを作成します。  
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

### 2-1-4 Docker imageをビルドする

#### コマンド実行
DockerfileからDocker imageをビルドします。

```
$ docker build -t go-image:latest .
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

### 2-1-5 Docker image一覧を確認する

#### コマンド実行
作成されたDocker imageを確認します。

```
$ docker images
```

#### 実行結果

```
REPOSITORY   TAG       IMAGE ID       CREATED          SIZE
go-image     latest    220026ab99c0   4 minutes ago   862MB
```

### 2-1-6 Dockerコンテナを起動する

#### コマンド実行
作成されたDocker imageからDockerコンテナを起動します。  
`-d`：コンテナをバックグラウンド実行します。

```
$ docker run -d go-image
```

#### 実行結果

```
8d3a4f0c85e671284ec4cab7c973cbbf91318f54583059b15db4419f6903af56
```

コンテナIDが表示されていることを確認します。

### 2-1-7 Dockerコンテナ状態を確認する

#### コマンド実行
起動されたDockerコンテナを動作確認します。

```
$ docker ps
```

#### 実行結果

```
CONTAINER ID   IMAGE      COMMAND                  CREATED          STATUS          PORTS           NAMES
8d3a4f0c85e6   go-image   "go run /work/main.go"   2 minutes ago    Up 2 minutes                    inspiring_sinoussi
```

`PORTS`に何も割り当たっていないことを確認します。

### 2-1-8 Goアプリのレスポンスを確認する

#### コマンド実行
GoアプリへGETリクエストを送信し、レスポンスを確認します。

```cmd
$ curl http://localhost:8080
```

#### 実行結果

```
curl: (7) Failed to connect to localhost port 8080: Connection refused
```

`Connection refused`と表示され、8080ポートへの接続に失敗することを確認します。  

#### ここで、少し考えてみましょう。なぜ、接続に失敗するのでしょうか？  
Goアプリにて公開設定した8080ポートは、コンテナ内に限定されたコンテナポートであることを理解する必要があります。つまり、コンテナ内からコンテナポート(8080)へのアクセスは可能ですが、コンテナ外部(ローカルPC)からのコンテナポートへのアクセスは不可能であることを意味します。この場合には、Dockerのポートフォワーディング機能を利用し、ホストマシンのポートをコンテナポートに紐付け、コンテナ外との通信を可能にすることができます。以下で手順を確認します。

### 2-1-9 ポートフォワーディング設定をしてDockerコンテナを起動する

#### コマンド実行
作成されたDocker imageからDockerコンテナを起動します。
`-p`：{コンテナ外部側ポート}:{コンテナ内部側ポート}の書式で記述することが可能。

```
$ docker run -d -p 9000:8080 go-image
```

#### 実行結果

```
d94c925240845c03b2f2dff0d43aea9d9b7c2f86184309e84b2cb4e93ff97c0a
```

コンテナIDが表示されていることを確認します。

### 2-1-10 Dockerコンテナ状態を確認する

#### コマンド実行
起動されたDockerコンテナを動作確認します。

```
$ docker ps
```

#### 実行結果

```
CONTAINER ID   IMAGE      COMMAND                  CREATED          STATUS          PORTS                                       NAMES
d94c92524084   go-image   "go run /work/main.go"   15 seconds ago   Up 14 seconds   0.0.0.0:9000->8080/tcp, :::9000->8080/tcp   elegant_germain
```

`PORTS`にポートフォワーディング設定がされていることを確認します。

### 2-1-11 Goアプリのレスポンスを確認する

#### コマンド実行
8080ポート(コンテナ外部)⇔9000ポート(コンテナ内部)で公開されたGoアプリへGETリクエストを送信し、レスポンスを確認します。

```cmd
$ curl http://localhost:9000
```

#### 実行結果

```
Hello Dcoker!!
```

9000ポートへの接続に成功し、`Hello Dcoker!!`と表示されることを確認します。

## 2-2 GitHub Packagesの基本操作をおさらいする

★見せ方を再確認

ここでは、GitHub PackagesへDockerイメージをpushします。

### 2-2-1

```

```

## 2-3 Kubernetesの基本操作をおさらいする

ここでは、Minikubeを使って、Kubernetesクラスタの起動からオブジェクトの作成までの操作を確認します。

### 2-3-1 Minikubeを起動する

#### コマンド実行
Docker driverを利用してMinikubeを起動し、立ち上げたDockerコンテナ上にKubernetesクラスタを構築します。

```
$ minikube start --driver=docker
```
#### 実行結果

```
* Microsoft Windows 10 Pro 10.0.19043 Build 19043 上の minikube v1.22.0
* 設定を元に、 docker ドライバを使用します
* コントロールプレーンのノード minikube を minikube 上で起動しています
* イメージを Pull しています...
* Kubernetes v1.21.2 のダウンロードの準備をしています
    > preloaded-images-k8s-v11-v1...: 502.14 MiB / 502.14 MiB  100.00% 13.77 Mi
    > gcr.io/k8s-minikube/kicbase...: 361.09 MiB / 361.09 MiB  100.00% 8.21 MiB
* docker container (CPUs=2, Memory=1986MB) を作成しています...
* Docker 20.10.7 で Kubernetes v1.21.2 を準備しています...
  - 証明書と鍵を作成しています...
  - Control Plane を起動しています...
  - RBAC のルールを設定中です...
* Kubernetes コンポーネントを検証しています...
  - イメージ gcr.io/k8s-minikube/storage-provisioner:v5 を使用しています
* 有効なアドオン: storage-provisioner, default-storageclass
* 完了しました！ kubectl が「"minikube"」クラスタと「"default"」ネームスペースを使用するよう構成されました
```
`完了しました！`と表示されていることを確認します。

### 2-3-2 Minikubeの状態を確認する

#### コマンド実行

```
$ minikube status
```

#### 実行結果

```
type: Control Plane
host: Running
kubelet: Running
apiserver: Running
kubeconfig: Configured
```

### 2-3-3 ノード一覧を確認する

#### コマンド実行
Minikubeによって作成されたKubernetesクラスタのノードを確認します。

```
$ kubectl get nodes
```
#### 実行結果

```
NAME       STATUS   ROLES                  AGE   VERSION
minikube   Ready    control-plane,master   25s   v1.21.2
```

`minikube`というノードが表示されていることを確認します。

### 2-3-4 ポッド一覧を確認する

#### コマンド実行
`-A`：全NameSpaceの結果を取得
`-o wide`：各Podの実行ホストIPを表示

```
$ kubectl get pods -A -o wide
```

#### 実行結果

```
NAMESPACE     NAME                               READY   STATUS    RESTARTS   AGE   IP             NODE       NOMINATED NODE   READINESS GATES
kube-system   coredns-558bd4d5db-s8r7s           1/1     Running   0          33s   172.17.0.2     minikube   <none>           <none>
kube-system   etcd-minikube                      1/1     Running   0          48s   192.168.49.2   minikube   <none>           <none>
kube-system   kube-apiserver-minikube            1/1     Running   0          48s   192.168.49.2   minikube   <none>           <none>
kube-system   kube-controller-manager-minikube   1/1     Running   0          43s   192.168.49.2   minikube   <none>           <none>
kube-system   kube-proxy-tqzc2                   1/1     Running   0          33s   192.168.49.2   minikube   <none>           <none>
kube-system   kube-scheduler-minikube            1/1     Running   0          43s   192.168.49.2   minikube   <none>           <none>
kube-system   storage-provisioner                1/1     Running   1          44s   192.168.49.2   minikube   <none>           <none>
```

`STATUS`が`Running`になっていることを確認します。

### 2-3-5 マニフェストファイルを作成する

Kubernetesでは、作成するポッドのリソース構成をマニフェストファイルにコードで記述することができます。

ワークディレクトリに、以下のマニフェストファイルを作成します。  
ファイル名：`nginx.yaml`

```
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
spec:
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx:latest
        ports:
        - containerPort: 80
```

### 2-3-6 ポッドを作成する

#### コマンド実行
マニフェストファイルからポッドを作成します。

```
$ kubectl apply -f nginx.yaml
```

#### 実行結果

```
deployment.apps/nginx-deployment created
```

`created`と表示されていることを確認します。

### 2-3-7 ポッド一覧を確認する

#### コマンド実行
`-o wide`：各Podの実行ホストIPを表示

```
$ kubectl get pods -o wide
```

#### 実行結果

```
NAME                               READY   STATUS    RESTARTS   AGE   IP           NODE       NOMINATED NODE   READINESS GATES
nginx-deployment-585449566-8tkwb   1/1     Running   0          4s    172.17.0.3   minikube   <none>           <none>
```

`STATUS`が`Running`になっていることを確認します。

