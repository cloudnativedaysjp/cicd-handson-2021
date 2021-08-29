# Chapter 3

本章では、ローカル環境にてGo言語で書かれたアプリを使用し、Docker/GitHub Packages/Kubernetesの一連の流れを実践しながら、基本操作についておさらいしていきます。  

- 目次
  - [3-1 Dockerの基本操作をおさらいする](https://github.com/cloudnativedaysjp/cicd-handson-2021/blob/main/docs/chapter3.md#3-1-1-%E3%82%AB%E3%83%AC%E3%83%B3%E3%83%88%E3%83%87%E3%82%A3%E3%83%AC%E3%82%AF%E3%83%88%E3%83%AA%E3%82%92%E5%A4%89%E6%9B%B4%E3%81%99%E3%82%8B)
    - build-ship-runの流れを確認します。
  - [3-2 GitHub Packagesの基本操作をおさらいする](https://github.com/cloudnativedaysjp/cicd-handson-2021/blob/main/docs/chapter3.md#3-2-github-packages%E3%81%AE%E5%9F%BA%E6%9C%AC%E6%93%8D%E4%BD%9C%E3%82%92%E3%81%8A%E3%81%95%E3%82%89%E3%81%84%E3%81%99%E3%82%8B)
    - トークン生成、DockerイメージPushの流れを確認します。
  - [3-3 Kubernetesの基本操作をおさらいする](https://github.com/cloudnativedaysjp/cicd-handson-2021/blob/main/docs/chapter3.md#3-3-kubernetes%E3%81%AE%E5%9F%BA%E6%9C%AC%E6%93%8D%E4%BD%9C%E3%82%92%E3%81%8A%E3%81%95%E3%82%89%E3%81%84%E3%81%99%E3%82%8B)
    - DockerイメージPull、ポッド起動の流れを確認します。

## 3-1 Dockerの基本操作をおさらいする

ここでは、Goアプリを使って、Dockerイメージのビルドから起動までの操作を確認します。

### 3-1-1 カレントディレクトリを変更する

#### コマンド実行
[リポジトリの作成](https://github.com/cloudnativedaysjp/cicd-handson-2021/blob/main/docs/chapter1.md#%E3%83%AA%E3%83%9D%E3%82%B8%E3%83%88%E3%83%AA%E3%81%AE%E4%BD%9C%E6%88%90)で、GitHubからcloneしたローカル環境のリポジトリへカレントディレクトリを変更します。

```cmd
$ cd cicd-handson-2021-code/apps
```

### 3-1-2 Go言語アプリソースコードを確認する

Dockerイメージのビルドに必要な、以下ディレクトリが存在することを確認します。 

```
cicd-handson-2021-code/apps/server
```

### 3-1-3 Dockerfileを作成する

Dockerfileは、Dockerイメージのビルド時に、事前に実施しておきたい操作をコードとして記述したファイルです。  
主に、OSやミドルウェア、コマンド実行、デーモン実行、環境変数などの設定を記述することが可能で、  
Dockerfileを使ってイメージをビルドすることで、その設定をコードとして柔軟に管理することができます。  

`cicd-handson-2021-code/apps`配下に、以下のDockerfileを作成します。  
※ここでは、コンテナ内ポート9090で公開するサーバ側アプリをビルドするためのDockerfileを作成します。  
- ファイル名：`Dockerfile`

```Dockerfile
# ベースイメージ指定
FROM golang:latest

# ワークディレクトリを指定
WORKDIR /app

# ホストOSのapps内全てをWORKDIRにコピー
COPY . ./

# ビルド
RUN go build -o ./server-run ./server

# コンテナのポートを9090で公開
EXPOSE 9090

# アプリ実行
CMD [ "./server-run" ]
```

`cicd-handson-2021-code/apps`配下に、Dockerfileが配置されていることを確認します。

### 3-1-4 DockerfileをリポジトリへPushする

ここでは、作成したDockerfileを`cicd-handson-2021-code`リポジトリへPushします。

#### コマンド実行
gitコマンド初回実行の場合は、任意のメールアドレス、ユーザ名を設定します。

```git
$ git config user.email "you@example.com"
$ git config user.name "Your Name"
```

#### コマンド実行
`Dockerfile`をインデックス(コミット対象)に追加します。

```git
$ git add Dockerfile
```

#### コマンド実行
インデックスにある`Dockerfile`をコミットします。
- `-m`：コメントメッセージを設定します。

```git
$ git commit -m "Add Dockerfile"
```

#### 実行結果

```
[main da54a41] Add Dockerfile
 1 file changed, 17 insertions(+)
 create mode 100644 apps/Dockerfile
```
 
`1 file changed`、`apps/Dockerfile`と変更が表示されていることを確認します。

#### コマンド実行
コミット内容を、リモートリポジトリ`origin`上の`main`ブランチへ反映します。

```git
$ git push origin main
```

#### 実行結果

```
Enumerating objects: 5, done.
Counting objects: 100% (5/5), done.
Delta compression using up to 16 threads
Compressing objects: 100% (3/3), done.
Writing objects: 100% (4/4), 561 bytes | 561.00 KiB/s, done.
Total 4 (delta 0), reused 0 (delta 0), pack-reused 0
To https://github.com/<GITHUB_USER>/cicd-handson-2021-code.git
   71c5d70..a80279a  main -> main
```
※`<GITHUB_USER>`は、GitHubユーザ名に置き換わっている状態。

`cicd-handson-2021-code`リポジトリの`main`ブランチへPushされていることを確認します。

### 3-1-5 Docker imageをビルドする

#### コマンド実行
作成したDockerfileを使用して、Dockerイメージをビルドします。  
- `-t`："名前:タグ"形式で名前とオプションのタグを指定します。

```bash
$ docker image build -t go-image:base .
```
※Docker v1.13 以降では、 旧`docker build`⇒新`docker image build`コマンドが推奨されています。

#### 実行結果

```
[+] Building 20.2s (10/10) FINISHED
 => [internal] load build definition from Dockerfile                                           0.0s
 => => transferring dockerfile: 365B                                                           0.0s
 => [internal] load .dockerignore                                                              0.0s
 => => transferring context: 2B                                                                0.0s
 => [internal] load metadata for docker.io/library/golang:1.16                                 2.9s
 => [auth] library/golang:pull token for registry-1.docker.io                                  0.0s
 => [1/4] FROM docker.io/library/golang:1.16@sha256:87cbbe43ece5024f0745be543c81ae6bf7b88291  15.8s
 => => resolve docker.io/library/golang:1.16@sha256:87cbbe43ece5024f0745be543c81ae6bf7b88291a  0.0s
 => => sha256:0c6e622a0ff6a2c83bb5b6f0f939367cf40083754b34e48f17cfcc73b0 129.04MB / 129.04MB  11.9s
 => => sha256:54406b2e8bb95003bcec911562d5606af1a17de7a2fcc7e8d7258fcb6e9a2fe 1.80kB / 1.80kB  0.0s
 => => sha256:7669e289697491b9ef28ca9dd0958a6f7ff9ddc0d50cefe1a58acc74b37ddd9b 155B / 155B     0.3s
 => => sha256:87cbbe43ece5024f0745be543c81ae6bf7b88291a8bc2b4429a43b7236254ec 2.36kB / 2.36kB  0.0s
 => => sha256:019c7b2e3cb8185c3c454f24679a7c83add6d31427e319d8bee35b12707f9a3 6.99kB / 6.99kB  0.0s
 => => extracting sha256:0c6e622a0ff6a2c83bb5b6f0f939367cf40083754b34e48f17cfcc73b05ad99c      3.6s
 => => extracting sha256:7669e289697491b9ef28ca9dd0958a6f7ff9ddc0d50cefe1a58acc74b37ddd9b      0.0s
 => [internal] load build context                                                              0.1s
 => => transferring context: 7.67MB                                                            0.1s
 => [2/4] WORKDIR /app                                                                         0.1s
 => [3/4] COPY . ./                                                                            0.1s
 => [4/4] RUN go build -o ./server-run ./server                                                1.1s
 => exporting to image                                                                         0.1s
 => => exporting layers                                                                        0.1s
 => => writing image sha256:d7aa4942052cd829b7556ea6dd31615b968b6c2ea946edb534b77207fbeaa175   0.0s
 => => naming to docker.io/library/go-image:base                                               0.0s
```

### 3-1-6 Docker image一覧を確認する

#### コマンド実行
作成されたDocker imageを確認します。

```bash
$ docker image ls
```
※Docker v1.13 以降では、 旧`docker images`⇒新`docker image ls`コマンドが推奨されています。

#### 実行結果

```
REPOSITORY   TAG       IMAGE ID       CREATED          SIZE
go-image     base      220026ab99c0   4 minutes ago    938MB
```

`go-image`が表示されていることを確認します。

### 3-1-7 Dockerコンテナを起動する

#### コマンド実行
作成されたDocker imageからDockerコンテナを起動します。  
- `--rm`：コンテナ終了時にコンテナ自動的に削除します。  
- `--name`：起動時のコンテナ名を指定します。  
- `-d`：コンテナをバックグラウンド実行します。  

```bash
$ docker container run --rm --name go-container -d go-image:base
```
※Docker v1.13 以降では、 旧`docker run`⇒新`docker container run`コマンドが推奨されています。

#### 実行結果

```
8d3a4f0c85e671284ec4cab7c973cbbf91318f54583059b15db4419f6903af56
```

コンテナIDが表示されていることを確認します。

### 3-1-8 Dockerコンテナ状態を確認する

#### コマンド実行
起動されたDockerコンテナを動作確認します。

```bash
$ docker container ls
```
※Docker v1.13 以降では、 旧`docker ps`⇒新`docker container ls`コマンドが推奨されています。

#### 実行結果

```
CONTAINER ID   IMAGE           COMMAND          CREATED          STATUS          PORTS           NAMES
8d3a4f0c85e6   go-image:base   "./server-run"   2 minutes ago    Up 2 minutes    9090/tcp        go-container
```

`PORTS`に`9090/tcp`と表示されていることを確認します。

### 3-1-9 Goアプリのレスポンスを確認する

#### コマンド実行
GoアプリへGETリクエストを送信し、レスポンスを確認します。

```cmd
$ curl http://localhost:9090/health
```

#### 実行結果

```
curl: (7) Failed to connect to localhost port 9090: Connection refused
```

`Connection refused`と表示され、9090ポートへの接続に失敗することを確認します。  

### Tips ポートへの接続に失敗する理由を理解する  
ここで、少し考えてみましょう。なぜ、接続に失敗するのでしょうか？  
Goアプリにて公開設定したポート9090は、コンテナ内に限定されたコンテナポートであることを理解する必要があります。つまり、コンテナ内からコンテナポート9090へのアクセスは可能ですが、コンテナ外部(ローカルPC)からのコンテナポートへのアクセスは不可能であることを意味します。この場合には、Dockerのポートフォワーディング機能を利用し、ホストマシンのポートをコンテナポートに紐付け、コンテナ外との通信を可能にすることができます。以下で手順を確認します。

### 3-1-10 Dockerコンテナを停止する

#### コマンド実行
ポート設定が出来ていないコンテナを停止します。

```bash
$ docker container stop go-container
```
※Docker v1.13 以降では、 旧`docker stop`⇒新`docker container stop`コマンドが推奨されています。  
また、このコマンドでは、`docker container stop <"NAME" or "CONTAINER ID">`のように、コンテナ名以外に、コンテナIDを指定してコンテナを停止させることも出来ます。

#### 実行結果

```
go-container
```

### 3-1-11 Dockerコンテナ状態を確認する

#### コマンド実行
[Dockerコンテナを起動する](https://github.com/cloudnativedaysjp/cicd-handson-2021/blob/main/docs/chapter3.md#3-1-7-docker%E3%82%B3%E3%83%B3%E3%83%86%E3%83%8A%E3%82%92%E8%B5%B7%E5%8B%95%E3%81%99%E3%82%8B)で、`--rm`オプションを指定して起動されたDockerコンテナが、コンテナ停止と共に削除されていることを確認します。

```bash
$ docker container ls
```

#### 実行結果

```
CONTAINER ID   IMAGE      COMMAND                  CREATED          STATUS          PORTS           NAMES
```

`go-container`がコンテナ一覧に存在していないことを確認します。

### 3-1-12 ポートフォワーディング設定をしてDockerコンテナを起動する

#### コマンド実行
作成されたDocker imageからDockerコンテナを起動します。  
※ここでは、9090で公開したコンテナ内ポートを、9091でコンテナ外へ公開するための紐付け設定を行います。  
- `-p`：{コンテナ外部側ポート}:{コンテナ内部側ポート}の書式で記述可能です。

```bash
$ docker container run --rm --name go-container -d -p 9091:9090 go-image:base
```
※Docker v1.13 以降では、 旧`docker run`⇒新`docker container run`コマンドが推奨されています。

#### 実行結果

```
d94c925240845c03b2f2dff0d43aea9d9b7c2f86184309e84b2cb4e93ff97c0a
```

コンテナIDが表示されていることを確認します。

### 3-1-13 Dockerコンテナ状態を確認する

#### コマンド実行
起動されたDockerコンテナを動作確認します。

```bash
$ docker container ls
```

#### 実行結果

```
CONTAINER ID   IMAGE           COMMAND                  CREATED          STATUS          PORTS                                       NAMES
d94c92524084   go-image:base   "go run /work/main.go"   15 seconds ago   Up 14 seconds   0.0.0.0:9091->9090/tcp, :::9091->9090/tcp   go-container
```

`PORTS`にポートフォワーディング設定がされていることを確認します。

### 3-1-14 Goアプリのレスポンスを確認する

#### コマンド実行
9090ポート(コンテナ内部)⇒9091ポート(コンテナ外部)で公開されたGoアプリへGETリクエストを送信し、レスポンスを確認します。

```cmd
$ curl http://localhost:9091/health
```

または、任意のブラウザで http://localhost:9091/health へアクセスしてください。

#### 実行結果

```
{"status":"Healthy"}
```

`{"status":"Healthy"}`がレスポンスされることを確認します。

### 3-1-15 Dockerコンテナを停止する

#### コマンド実行

```bash
$ docker container stop go-container
```

#### 実行結果

```
go-container
```

### 3-1-16 Dockerコンテナ状態を確認する

#### コマンド実行
[ポートフォワーディング設定をしてDockerコンテナを起動する](https://github.com/cloudnativedaysjp/cicd-handson-2021/blob/main/docs/chapter3.md#3-1-12-%E3%83%9D%E3%83%BC%E3%83%88%E3%83%95%E3%82%A9%E3%83%AF%E3%83%BC%E3%83%87%E3%82%A3%E3%83%B3%E3%82%B0%E8%A8%AD%E5%AE%9A%E3%82%92%E3%81%97%E3%81%A6docker%E3%82%B3%E3%83%B3%E3%83%86%E3%83%8A%E3%82%92%E8%B5%B7%E5%8B%95%E3%81%99%E3%82%8B)で、`--rm`オプションを指定して起動されたDockerコンテナが、コンテナ停止と共に削除されていることを確認します。

```bash
$ docker container ls
```

#### 実行結果

```
CONTAINER ID   IMAGE      COMMAND                  CREATED          STATUS          PORTS           NAMES
```

`go-container`がコンテナ一覧に存在していないことを確認します。

## 3-2 GitHub Packagesの基本操作をおさらいする

ここでは、ローカル環境で作成したDockerイメージをGitHub Packagesへpushします。

### 3-2-1 DockerでGitHub Packagesの認証を行う

#### コマンド実行
[GitHubのPersonal Access Tokenの取得](https://github.com/cloudnativedaysjp/cicd-handson-2021/blob/main/docs/chapter0.md#github%E3%81%AEpersonal-access-token%E3%81%AE%E5%8F%96%E5%BE%97)で作成し、ローカルディレクトリへ保存したPersonal Access Token(以下PAT)を使用して、`docker login`コマンドで、GitHub Packagesの認証を受けることができます。また、ローカルファイルからトークンを読み取るDockerの`--password-stdin`フラグを使うことをおすすめします。

```cmd
$ cd "token.txtを保存した任意のローカルディレクトリ"
$ cat token.txt | docker login https://docker.pkg.github.com -u <GITHUB_USER> --password-stdin
```

#### 実行結果

```
Login Succeeded
```

### 3-2-2 Dockerタグを付与する

ここでは、[Docker imageをビルドする](https://github.com/cloudnativedaysjp/cicd-handson-2021/blob/main/docs/chapter3.md#3-1-5-docker-image%E3%82%92%E3%83%93%E3%83%AB%E3%83%89%E3%81%99%E3%82%8B)で作成したDockerイメージにタグ付けを行います。<GITHUB_USER>をリポジトリを所有するGitHubユーザ名で、REPOSITORYをプロジェクトを含むリポジトリの名前で、IMAGE_NAMEをパッケージもしくはイメージの名前で、VERSIONをビルドの時点のパッケージバージョンで置き換えてください。

#### コマンド実行

```bash
$ docker image tag go-image:base docker.pkg.github.com/<GITHUB_USER>/cicd-handson-2021-code/go-image:base
```
※Docker v1.13 以降では、 旧`docker tag`⇒新`docker image tag`コマンドが推奨されています。  
`<GITHUB_USER>`をGitHubユーザ名に置き換えてコマンドを実行します。

### 3-2-3 Docker image一覧を確認する

#### コマンド実行
新しくタグ付けされたDocker imageを確認します。

```bash
$ docker image ls
```

#### 実行結果

```bash
REPOSITORY                                                            TAG       IMAGE ID       CREATED       SIZE
go-image                                                              base      220026ab99c0   4 hours ago   862MB
docker.pkg.github.com/<GITHUB_USER>/cicd-handson-2021-code/go-image   base      220026ab99c0   4 hours ago   862MB
```
※`<GITHUB_USER>`は、GitHubユーザ名に置き換わっている状態。

`docker.pkg.github.com/<GITHUB_USER>/cicd-handson-2021-code/go-image`が表示されていることを確認します。

### 3-2-4 GitHub PackagesへDockerイメージをpushする

#### コマンド実行

```bash
$ docker image push docker.pkg.github.com/<GITHUB_USER>/cicd-handson-2021-code/go-image:base
```
※Docker v1.13 以降では、 旧`docker push`⇒新`docker image push`コマンドが推奨されています。  
`<GITHUB_USER>`をGitHubユーザ名に置き換えてコマンドを実行します。

#### 実行結果

```bash
The push refers to repository [docker.pkg.github.com/<GITHUB_USER>/cicd-handson-2021-code/go-image]
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
※`<GITHUB_USER>`は、GitHubユーザ名に置き換わっている状態。  
端末によっては、すべて`Pushed`になるまでしばらくかかります。

実行結果が、全て`Pushed`になっていることを確認し、GitHub画面のPackagesに`go-image`が公開されていることを確認します。
GitHub Packagesの画面に出てくるまで、数分かかることがあります。

![image](https://user-images.githubusercontent.com/45567889/130482375-96c65eb2-429d-453d-a311-a00390e24c94.png)


## 3-3 Kubernetesの基本操作をおさらいする

ここでは、Minikubeを使って、Kubernetesクラスタの起動からオブジェクトの作成までの操作を確認します。

### 3-3-1 Minikubeを起動する

[Minikube](https://github.com/cloudnativedaysjp/cicd-handson-2021/blob/main/docs/chapter0.md#minikube)を参考に、`minikube start`、`kubectl get nodes`、`kubectl get pods`コマンドを実行し、Minikubeを起動します。

### 3-3-2 ローカル環境のDockerイメージを削除する

#### コマンド実行
Github PackagesにPushしたDockerイメージを使って、ポッドを作成するためのマニフェストを記述する前に、[Docker imageをビルドする](https://github.com/cloudnativedaysjp/cicd-handson-2021/blob/main/docs/chapter3.md#3-1-5-docker-image%E3%82%92%E3%83%93%E3%83%AB%E3%83%89%E3%81%99%E3%82%8B)と[Dockerタグを付与する](https://github.com/cloudnativedaysjp/cicd-handson-2021/blob/main/docs/chapter3.md#3-2-3-docker%E3%82%BF%E3%82%B0%E3%82%92%E4%BB%98%E4%B8%8E%E3%81%99%E3%82%8B)で作成した2つのローカル環境のDockerイメージを削除しておきます。

```bash
$ docker image rm go-image:base
$ docker image rm docker.pkg.github.com/<GITHUB_USER>/cicd-handson-2021-code/go-image:base
```
※`<GITHUB_USER>`をGitHubユーザ名に置き換えてコマンドを実行します。

#### 実行結果

```bash
Untagged: go-image:base
Untagged: docker.pkg.github.com/<GITHUB_USER>/cicd-handson-2021-code/go-image:base
Untagged: docker.pkg.github.com/<GITHUB_USER>/cicd-handson-2021-code/go-image@sha256:11dd65371181d74b33c84b18f4f6ba87537cdbab7c884ef12ee6429865c0f640
Deleted: sha256:f0e3dbfbe4c7cf3d65a39f33d9cb5bb049ca063b1b40dc7a55ef055d983c0b91
```
※`<GITHUB_USER>`は、GitHubユーザ名に置き換わっている状態。

### 3-3-3 Docker image一覧を確認する

#### コマンド実行
ローカルにgo-imageのDocker imageが存在しないことを確認します。

```bash
$ docker image ls
```

#### 実行結果

```
REPOSITORY                    TAG       IMAGE ID       CREATED       SIZE
```

`go-image`と`docker.pkg.github.com/<GITHUB_USER>/cicd-handson-2021-code/go-image`が一覧に存在しないことを確認します。  
※`<GITHUB_USER>`は、GitHubユーザ名に置き換わっている状態。

### 3-3-4 Dockerコンテナレジストリ認証用のクレデンシャル(Secret)を作成する

#### コマンド実行
- `--save-config`：作成した現在の設定をannotationに保存します。
- `--docker-server`：Dockerレジストリサーバを指定します。
- `--docker-username`：GitHub登録時のユーザ名を指定します。
- `--docker-password`：[GitHubのPersonal Access Tokenの取得](https://github.com/cloudnativedaysjp/cicd-handson-2021/blob/main/docs/chapter0.md#github%E3%81%AEpersonal-access-token%E3%81%AE%E5%8F%96%E5%BE%97)で作成したGitHubのPAT(Personal Access Token)を指定します。
- `--docker-email`：GitHub登録時のメールアドレスを指定します。

```bash
$ kubectl create secret docker-registry --save-config dockerconfigjson-github-com --docker-server=docker.pkg.github.com --docker-username=<GITHUB_USER> --docker-password=<PERSONAL_ACCESS_TOKEN> --docker-email=<GITHUB_EMAIL>
```
※ghcr.io への読み書きについて、[Working with the Container registry](https://docs.github.com/en/packages/working-with-a-github-packages-registry/working-with-the-container-registry#authenticating-to-the-container-registry)によると、`--docker-password`には、GitHubのPATを指定する必要があります。間違えてGitHub登録時のパスワードを入力しないよう注意が必要です。

#### 実行結果

```
secret/dockerconfigjson-github-com created
```

### 3-3-5 カレントディレクトリを変更する

#### コマンド実行
[リポジトリの作成](https://github.com/cloudnativedaysjp/cicd-handson-2021/blob/main/docs/chapter1.md#%E3%83%AA%E3%83%9D%E3%82%B8%E3%83%88%E3%83%AA%E3%81%AE%E4%BD%9C%E6%88%90)で、GitHubからcloneしたローカル環境のリポジトリへカレントディレクトリを変更します。

```cmd
$ cd cicd-handson-2021-config/manifests
```

### 3-3-6 マニフェストファイルを作成する

Kubernetesでは、作成するポッドのリソース構成をマニフェストファイルにコードで記述することができます。

`cicd-handson-2021-config/manifests`配下に、以下のマニフェストファイルを作成します。  
- ファイル名：`goapp.yaml`

```yaml
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
        image: docker.pkg.github.com/<GITHUB_USER>/cicd-handson-2021-code/go-image:base
        ports:
        - containerPort: 9090
      imagePullSecrets:
      - name: dockerconfigjson-github-com
```
※`imagePullSecrets`に[Dockerコンテナレジストリ認証用のクレデンシャル(Secret)を作成する](https://github.com/cloudnativedaysjp/cicd-handson-2021/blob/main/docs/chapter3.md#3-3-4-docker%E3%82%B3%E3%83%B3%E3%83%86%E3%83%8A%E3%83%AC%E3%82%B8%E3%82%B9%E3%83%88%E3%83%AA%E8%AA%8D%E8%A8%BC%E7%94%A8%E3%81%AE%E3%82%AF%E3%83%AC%E3%83%87%E3%83%B3%E3%82%B7%E3%83%A3%E3%83%ABsecret%E3%82%92%E4%BD%9C%E6%88%90%E3%81%99%E3%82%8B)で作成したクレデンシャル(Secret)保存名`dockerconfigjson-github-com`を指定し忘れないよう注意が必要です。  
`<GITHUB_USER>`は、GitHubユーザ名に置き換わっている状態。

### 3-3-7 マニフェストファイルをリポジトリへPushする

ここでは、作成したマニフェストファイルを`cicd-handson-2021-config`リポジトリへPushします。

#### コマンド実行
`goapp.yaml`をインデックス(コミット対象)に追加します。

```git
$ git add goapp.yaml
```

#### コマンド実行
インデックスにある`goapp.yaml`をコミットします。

```git
$ git commit -m "Add goapp.yaml"
```

#### 実行結果

```
[main 28d94de] Add goapp.yaml
 1 file changed, 17 insertions(+)
 create mode 100644 manifests/goapp.yaml
```
 
`1 file changed`、`manifests/goapp.yaml`と変更が表示されていることを確認します。

#### コマンド実行
コミット内容を、リモートリポジトリ`origin`上の`main`ブランチへ反映します。

```git
$ git push origin main
```

#### 実行結果

```
Enumerating objects: 5, done.
Counting objects: 100% (5/5), done.
Delta compression using up to 16 threads
Compressing objects: 100% (3/3), done.
Writing objects: 100% (4/4), 561 bytes | 561.00 KiB/s, done.
Total 4 (delta 0), reused 0 (delta 0), pack-reused 0
To https://github.com/<GITHUB_USER>/cicd-handson-2021-config.git
   71c5d70..a80279a  main -> main
```
※`<GITHUB_USER>`は、GitHubユーザ名に置き換わっている状態。

`cicd-handson-2021-config`リポジトリの`main`ブランチへPushされていることを確認します。

### 3-3-8 ポッドを作成する

#### コマンド実行
マニフェストファイルからポッドを作成します。
- `-f`：ファイル名を指定します。

```cmd
$ kubectl apply -f goapp.yaml
```

#### 実行結果

```
deployment.apps/goapp-deployment created
```

`goapp-deployment`の作成結果が表示されていることを確認します。

### 3-3-9 ポッド一覧を確認する

#### コマンド実行
- `-o wide`：より詳細なリストを表示します。

```bash
$ kubectl get deploy,pods -o wide
```

#### 実行結果

```bash
NAME                                    READY   UP-TO-DATE   AVAILABLE   AGE     CONTAINERS   IMAGES                                                                     SELECTOR
deployment.apps/goapp-deployment        1/1     1            1           8m56s   goapp        docker.pkg.github.com/<GITHUB_USER>/cicd-handson-2021-code/go-image:base   app=goapp

NAME                                    READY   STATUS       RESTARTS    AGE     IP           NODE       NOMINATED NODE   READINESS GATES
pod/goapp-deployment-6c85ff5cc-6pc89    1/1     Running      0           8m56s   172.17.0.2   minikube   <none>           <none>
```
※`<GITHUB_USER>`は、GitHubユーザ名に置き換わっている状態。

`STATUS`が`Running`になっていることを確認します。  

### 3-3-10 ポートフォアーディング設定を行う。

ローカルのポートを任意のPodのポートに転送するための設定を行います。ここではローカルの9092番ポートをdeploymentの9090番ポートに転送します。

#### コマンド実行

```bash
$ kubectl port-forward deployment.apps/goapp-deployment 9092:9090
```

#### 実行結果

```
Forwarding from 127.0.0.1:9092 -> 9090
Forwarding from [::1]:9092 -> 9090
```

`9092 -> 9090`の転送設定が表示されていることを確認します。

### 3-3-11 アプリへ接続確認する

#### コマンド実行
任意のブラウザにて、以下のURLを実行し、疎通できていることを確認します。

```
http://localhost:9092/health
```

#### 実行結果

```
{"status":"Healthy"}
```

`{"status":"Healthy"}`が表示されていることを確認します。

### 3-3-12 ポッドを削除する

#### コマンド実行
マニフェストファイルから作成したポッドを削除します。
- `-f`：ファイル名を指定します。

```bash
$ kubectl delete -f goapp.yaml
```

#### 実行結果

```
deployment.apps "goapp-deployment" deleted
```

`goapp-deployment`の削除結果が表示されていることを確認します。
