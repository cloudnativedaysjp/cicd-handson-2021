# Chapter 2

本章では、Docker/GitHub Package/Kubernetesの基本的な操作についておさらいします。  
ローカル環境でgo言語を使って、Docker imageを作成し、Kubernetesで起動してみましょう。  

## 2-1 Dockerの基本操作をおさらいする

ここでは、Goアプリを使って、Dockerイメージのビルドから起動までの操作を確認します。

### 2-1-1 ワークディレクトリを作成する

#### コマンド実行
任意の場所にワークディレクトリを作成します。

```cmd
$ mkdir ~\goapp
$ cd ~\goapp
```

### 2-1-2 Go言語でサンプルアプリのソースコードを作成する

ワークディレクトリに、以下のファイルを作成します。  
※ここでは、"Hello Docker!!"と表示されるソースコード例を紹介します。  
- ファイル名：`main.go`

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
- ファイル名：`Dockerfile`

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
$ docker image build -t go-image:base .
```
※Docker v1.13 以降では、 旧`docker build`⇒新`docker image build`が推奨されています。

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
 => => naming to docker.io/library/go-image:base                                               0.0s
```

### 2-1-5 Docker image一覧を確認する

#### コマンド実行
作成されたDocker imageを確認します。

```
$ docker image ls
```
※Docker v1.13 以降では、 旧`docker images`⇒新`docker image ls`が推奨されています。

#### 実行結果

```
REPOSITORY   TAG       IMAGE ID       CREATED          SIZE
go-image     base      220026ab99c0   4 minutes ago   862MB
```

`go-image`が表示されていることを確認します。

### 2-1-6 Dockerコンテナを起動する

#### コマンド実行
作成されたDocker imageからDockerコンテナを起動します。  
- `--rm`：コンテナ終了時にコンテナ自動的に削除します。  
- `--name`：起動時のコンテナ名を指定します。  
- `-d`：コンテナをバックグラウンド実行します。  

```
$ docker container run --rm --name go-container -d go-image
```
※Docker v1.13 以降では、 旧`docker run`⇒新`docker container run`が推奨されています。

#### 実行結果

```
8d3a4f0c85e671284ec4cab7c973cbbf91318f54583059b15db4419f6903af56
```

コンテナIDが表示されていることを確認します。

### 2-1-7 Dockerコンテナ状態を確認する

#### コマンド実行
起動されたDockerコンテナを動作確認します。

```
$ docker container ls
```
※Docker v1.13 以降では、 旧`docker ps`⇒新`docker container ls`が推奨されています。

#### 実行結果

```
CONTAINER ID   IMAGE      COMMAND                  CREATED          STATUS          PORTS           NAMES
8d3a4f0c85e6   go-image   "go run /work/main.go"   2 minutes ago    Up 2 minutes                    go-container
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

### Tips ここで、少し考えてみましょう。なぜ、接続に失敗するのでしょうか？  
Goアプリにて公開設定した8080ポートは、コンテナ内に限定されたコンテナポートであることを理解する必要があります。つまり、コンテナ内からコンテナポート(8080)へのアクセスは可能ですが、コンテナ外部(ローカルPC)からのコンテナポートへのアクセスは不可能であることを意味します。この場合には、Dockerのポートフォワーディング機能を利用し、ホストマシンのポートをコンテナポートに紐付け、コンテナ外との通信を可能にすることができます。以下で手順を確認します。

### 2-1-9 Dockerコンテナを停止する

#### コマンド実行
ポート設定が出来ていないコンテナを停止します。

```
$ docker container stop go-container
```
※Docker v1.13 以降では、 旧`docker stop`⇒新`docker container stop`が推奨されています。  
また、このコマンドでは、`docker container stop <"CONTAINER ID" or "NAME">`のように、コンテナID、または、コンテナ名を指定してコンテナを停止させることが出来ます。

#### 実行結果

```
go-container
```

### 2-1-10 Dockerコンテナ状態を確認する

#### コマンド実行
2-1-6で、`--rm`オプションを指定して起動されたDockerコンテナが、コンテナ停止と共に削除されていることを確認します。

```
$ docker container ls
```

#### 実行結果

```
CONTAINER ID   IMAGE      COMMAND                  CREATED          STATUS          PORTS           NAMES
```

`go-container`がコンテナ一覧に存在していないことを確認します。

### 2-1-11 ポートフォワーディング設定をしてDockerコンテナを起動する

#### コマンド実行
作成されたDocker imageからDockerコンテナを起動します。  
- `-p`：{コンテナ外部側ポート}:{コンテナ内部側ポート}の書式で記述することが可能。

```
$ docker container run --rm --name go-container -d -p 9000:8080 go-image
```
※Docker v1.13 以降では、 旧`docker run`⇒新`docker container run`が推奨されています。

#### 実行結果

```
d94c925240845c03b2f2dff0d43aea9d9b7c2f86184309e84b2cb4e93ff97c0a
```

コンテナIDが表示されていることを確認します。

### 2-1-12 Dockerコンテナ状態を確認する

#### コマンド実行
起動されたDockerコンテナを動作確認します。

```
$ docker container ls
```

#### 実行結果

```
CONTAINER ID   IMAGE      COMMAND                  CREATED          STATUS          PORTS                                       NAMES
d94c92524084   go-image   "go run /work/main.go"   15 seconds ago   Up 14 seconds   0.0.0.0:9000->8080/tcp, :::9000->8080/tcp   go-container
```

`PORTS`にポートフォワーディング設定がされていることを確認します。

### 2-1-13 Goアプリのレスポンスを確認する

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

### 2-1-14 Dockerコンテナを停止する

#### コマンド実行
ポート設定が出来ていないコンテナを停止します。

```
$ docker container stop go-container
```

#### 実行結果

```
go-container
```

### 2-1-15 Dockerコンテナ状態を確認する

#### コマンド実行
2-1-11で、`--rm`オプションを指定して起動されたDockerコンテナが、コンテナ停止と共に削除されていることを確認します。

```
$ docker container ls
```

#### 実行結果

```
CONTAINER ID   IMAGE      COMMAND                  CREATED          STATUS          PORTS           NAMES
```

`go-container`がコンテナ一覧に存在していないことを確認します。

## 2-2 GitHub Packagesの基本操作をおさらいする

ここでは、ローカル環境で作成したDockerイメージをGitHub Packagesへpushします。

### 2-2-1 トークン情報を作成・保存する

[GitHub Docs : Creating a personal access token](https://docs.github.com/en/github/authenticating-to-github/keeping-your-account-and-data-secure/creating-a-personal-access-token)のドキュメントに従い、Dockerログイン時に使用する以下の権限を付与したPAT(Personal Access Token)情報を取得します。  
- write:packages(Upload packages to github package registry)  
- read:packages(Download packages from github package registry)

![image](https://user-images.githubusercontent.com/45567889/129010278-33158cc5-55d9-47ff-bcc3-ce9c07f2fa2c.png)

![image](https://user-images.githubusercontent.com/45567889/128994241-87aefb3a-d670-455f-9001-115c2f52fa7f.png)

生成されたPATをファイルに保存します。  
- ファイル名：`token.txt`

```
"生成されたPATのプレーンテキスト"
```

### 2-2-2 DockerでGitHub Packagesの認証を行う

#### コマンド実行
`docker login`コマンドを使い、DockerでGitHub Packagesの認証を受けることができます。クレデンシャルをセキュアに保つ貯めに、個人アクセストークンは自分のコンピュータのローカルファイルに保存し、ローカルファイルからトークンを読み取るDockerの`--password-stdin`フラグを使うことをおすすめします。

```
$ cat token.txt | docker login https://docker.pkg.github.com -u USERNAME --password-stdin
```

#### 実行結果

```
Login Succeeded
```

### 2-2-3 GitHubのリポジトリを作成する

ここでは、Dockerイメージをpushするためのリポジトリを以下の名前で作成します。  
- リポジトリ名：`cicd-handson`

![image](https://user-images.githubusercontent.com/45567889/129009968-dd180fa6-9363-47ac-834b-2bbcff7b3be8.png)

```
OWNER/cicd-handson
```
※`OWNER`は、オーナー名に置き換わっている状態。

上記のようなリポジトリが作成されていることを確認します。

### 2-2-4 Dockerタグを付与する

ここでは、2-1-4で作成したDockerイメージにタグ付けを行います。OWNERをリポジトリを所有するユーザもしくはOrganizationアカウントの名前で、REPOSITORYをプロジェクトを含むリポジトリの名前で、IMAGE_NAMEをパッケージもしくはイメージの名前で、VERSIONをビルドの時点のパッケージバージョンで置き換えてください。

#### コマンド実行

```
$ docker image tag go-image:base docker.pkg.github.com/<OWNER>/cicd-handson/go-image:base
```
※Docker v1.13 以降では、 旧`docker tag`⇒新`docker image tag`が推奨されています。  
`OWNER`をオーナー名に置き換えてコマンドを実行します。

#### 実行結果

特に表示されません。

### 2-2-5 Docker image一覧を確認する

#### コマンド実行
新しくタグ付けされたDocker imageを確認します。

```
$ docker image ls
```

#### 実行結果

```
REPOSITORY                                                  TAG       IMAGE ID       CREATED       SIZE
go-image                                                    base      220026ab99c0   4 hours ago   862MB
docker.pkg.github.com/naka-teruhisa/cicd-handson/go-image   base      220026ab99c0   4 hours ago   862MB
```

`docker.pkg.github.com/naka-teruhisa/cicd-handson/go-image`が表示されていることを確認します。

### 2-2-1 GitHub PackagesへDockerイメージをpushする

#### コマンド実行

```
$ docker image push docker.pkg.github.com/<OWNER>/cicd-handson/go-image:base
```
※Docker v1.13 以降では、 旧`docker push`⇒新`docker image push`が推奨されています。  
`OWNER`をオーナー名に置き換えてコマンドを実行します。

#### 実行結果

```
The push refers to repository [docker.pkg.github.com/OWNER/cicd-handson/go-image]
3b29f7317fbb: Pushed
23255aaac099: Pushed
903b16bd3e46: Pushed
2f3906bb26c9: Pushed
d1c59e37fbfc: Pushed
ad83f0aa5c0a: Pushed
5a9a65095453: Pushed
4b0edb23340c: Pushed
afa3e488a0ee: Pushed
base: digest: sha256:11dd65371181d74b33c84b18f4f6ba87537cdbab7c884ef12ee6429865c0f640 size: 2209
```
※`OWNER`は、オーナー名に置き換わっている状態。

全て`Pushed`になっていることを確認します。

![image](https://user-images.githubusercontent.com/45567889/129005604-570c3ed0-f298-4421-b57f-ea2616e47da3.png)

GitHub画面のPackagesに`go-image`が公開されていることを確認します。

## 2-3 Kubernetesの基本操作をおさらいする

ここでは、Minikubeを使って、Kubernetesクラスタの起動からオブジェクトの作成までの操作を確認します。

### 2-3-1 Minikubeを起動する

Chapter00の[Minikube](https://github.com/cloudnativedaysjp/cicd-handson-2021/blob/main/docs/chapter0.md#minikube)を参考に、`minikube start`、`kubectl get nodes`、`kubectl get pods`コマンドを実行し、Minikubeが起動されていることを確認します。

### 2-3-2 ローカル環境のDockerイメージを削除する

#### コマンド実行
Github PackagesにPushしたDockerイメージを使って、ポッドを作成するためのマニフェストを記述する前に、2-1-4と2-2-4で作成した2つのローカル環境のDockerイメージを削除しておきます。

```
$ docker image rm go-image:base
$ docker image rm docker.pkg.github.com/OWNER/cicd-handson/go-image:base
```
※`OWNER`をオーナー名に置き換えてコマンドを実行します。

#### 実行結果

```
Untagged: go-image:base
Untagged: docker.pkg.github.com/OWNER/cicd-handson/go-image:base
Untagged: docker.pkg.github.com/OWNER/cicd-handson/go-image@sha256:11dd65371181d74b33c84b18f4f6ba87537cdbab7c884ef12ee6429865c0f640
```
※`OWNER`は、オーナー名に置き換わっている状態。

### 2-3-3 Docker image一覧を確認する

#### コマンド実行
新しくタグ付けされたDocker imageを確認します。

```
$ docker image ls
```

#### 実行結果

```
REPOSITORY                    TAG       IMAGE ID       CREATED       SIZE
```

`go-image`と`docker.pkg.github.com/OWNER/cicd-handson/go-image`が一覧に存在しないことを確認します。

### 2-3-4 Dockerコンテナレジストリ認証用のクレデンシャル(Secret)を作成する

#### コマンド実行

```
$ kubectl create secret docker-registry --save-config dockerconfigjson-github-com \
   --docker-server=docker.pkg.github.com \
   --docker-username=<USERNAME> \
   --docker-password=<PASSWORD> \
   --docker-email=<EMAIL>
```

#### 実行結果

```
secret/dockerconfigjson-github-com created
```

### 2-3-6 マニフェストファイルを作成する

Kubernetesでは、作成するポッドのリソース構成をマニフェストファイルにコードで記述することができます。

ワークディレクトリに、以下のマニフェストファイルを作成します。  
- ファイル名：`goapp.yaml`

```
apiVersion: apps/v1
kind: Deployment
metadata:
  name: goapp-deployment
spec:
  selector:
    matchLabels:
      app: goapp
  template:
    metadata:
      labels:
        app: goapp
    spec:
      containers:
      - name: goapp
        image: docker.pkg.github.com/naka-teruhisa/cicd-handson/go-image:base
        ports:
        - containerPort: 80
```

### 2-3-7 ポッドを作成する

#### コマンド実行
マニフェストファイルからポッドを作成します。

```
$ kubectl apply -f goapp.yaml
```

#### 実行結果

```
deployment.apps/goapp-deployment created
```

### 2-3-8 ポッド一覧を確認する

#### コマンド実行
`-o wide`：各Podの実行ホストIPを表示

```
$ kubectl get pods -o wide
```

#### 実行結果

```
NAME                               READY   STATUS    RESTARTS   AGE   IP           NODE       NOMINATED NODE   READINESS GATES
goapp-deployment-6c85ff5cc-6pc89   1/1     Running   0          4s    172.17.0.3   minikube   <none>           <none>
```

`STATUS`が`Running`になっていることを確認します。

### 2-3-9 ポッドを削除する

#### コマンド実行
マニフェストファイルから作成したポッドを削除します。

```
$ kubectl delete -f goapp.yaml
```

#### 実行結果

```
deployment.apps "goapp-deployment" deleted
```
