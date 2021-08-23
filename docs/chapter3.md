# Chapter 3

本章では、ローカル環境でGo言語で書かれたアプリを使って、Docker/GitHub Packages/Kubernetesの一連の流れを実践しながら、基本操作についておさらいしていきます。  

- [3-1 Dockerの基本操作をおさらいする](https://github.com/cloudnativedaysjp/cicd-handson-2021/blob/main/docs/chapter3.md#3-1-1-%E3%82%AB%E3%83%AC%E3%83%B3%E3%83%88%E3%83%87%E3%82%A3%E3%83%AC%E3%82%AF%E3%83%88%E3%83%AA%E3%82%92%E5%A4%89%E6%9B%B4%E3%81%99%E3%82%8B)：build-ship-runの流れを確認します。
- [3-2 GitHub Packagesの基本操作をおさらいする](https://github.com/cloudnativedaysjp/cicd-handson-2021/blob/main/docs/chapter3.md#3-2-github-packages%E3%81%AE%E5%9F%BA%E6%9C%AC%E6%93%8D%E4%BD%9C%E3%82%92%E3%81%8A%E3%81%95%E3%82%89%E3%81%84%E3%81%99%E3%82%8B)：トークン生成、DockerイメージPushの流れを確認します。
- [3-3 Kubernetesの基本操作をおさらいする](https://github.com/cloudnativedaysjp/cicd-handson-2021/blob/main/docs/chapter3.md#3-3-kubernetes%E3%81%AE%E5%9F%BA%E6%9C%AC%E6%93%8D%E4%BD%9C%E3%82%92%E3%81%8A%E3%81%95%E3%82%89%E3%81%84%E3%81%99%E3%82%8B)：DockerイメージPull、ポッド起動の流れを確認します。

## 3-1 Dockerの基本操作をおさらいする

ここでは、Goアプリを使って、Dockerイメージのビルドから起動までの操作を確認します。

### 3-1-1 カレントディレクトリを変更する

#### コマンド実行
[リポジトリの作成](https://github.com/cloudnativedaysjp/cicd-handson-2021/blob/main/docs/chapter1.md#%E3%83%AA%E3%83%9D%E3%82%B8%E3%83%88%E3%83%AA%E3%81%AE%E4%BD%9C%E6%88%90)で、GitHubからcloneしたローカル環境のリポジトリへカレントディレクトリを変更します。

```cmd
$ cd ~\cicd-handson-2021-code\apps
```

### 3-1-2 Go言語アプリソースコードを確認する

Dockerイメージのビルドに必要な、以下ディレクトリが存在することを確認します。 

```
~\cicd-handson-2021-code\apps\server
```

### 3-1-3 Dockerfileを作成する

Dockerfileは、Dockerイメージのビルド時に、事前に実施しておきたい操作をコードとして記述したファイルです。  
主に、OSやミドルウェア、コマンド実行、デーモン実行、環境変数などの設定を記述することが可能で、  
Dockerfileを使ってイメージをビルドすることで、その設定をコードとして柔軟に管理することができます。  

ワークディレクトリに、以下のDockerfileを作成します。  
※ここでは、コンテナ内ポート9090で公開するサーバ側アプリをビルドするためのDockerfileを作成します。  
- ファイル名：`Dockerfile`

```Dockerfile
# ベースイメージ指定
FROM golang:latest

# ワークディレクトリを指定
WORKDIR /app

# ホストOSのapp内全てをWORKDIRにコピー
COPY . ./

# ビルド
RUN go build -o ./server-run ./server

# コンテナのポートを9090で公開
EXPOSE 9090

# アプリ実行
CMD [ "./server-run" ]
```

### 3-1-4 Docker imageをビルドする

#### コマンド実行
DockerfileからDocker imageをビルドします。

```bash
$ docker image build -t go-image:base .
```
※Docker v1.13 以降では、 旧`docker build`⇒新`docker image build`コマンドが推奨されています。

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

### 3-1-5 Docker image一覧を確認する

#### コマンド実行
作成されたDocker imageを確認します。

```bash
$ docker image ls
```
※Docker v1.13 以降では、 旧`docker images`⇒新`docker image ls`コマンドが推奨されています。

#### 実行結果

```
REPOSITORY   TAG       IMAGE ID       CREATED          SIZE
go-image     base      220026ab99c0   4 minutes ago    862MB
```

`go-image`が表示されていることを確認します。

### 3-1-6 Dockerコンテナを起動する

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

### 3-1-7 Dockerコンテナ状態を確認する

#### コマンド実行
起動されたDockerコンテナを動作確認します。

```bash
$ docker container ls
```
※Docker v1.13 以降では、 旧`docker ps`⇒新`docker container ls`コマンドが推奨されています。

#### 実行結果

```
CONTAINER ID   IMAGE           COMMAND                  CREATED          STATUS          PORTS           NAMES
8d3a4f0c85e6   go-image:base   "go run /work/main.go"   2 minutes ago    Up 2 minutes                    go-container
```

`PORTS`に何も割り当たっていないことを確認します。

### 3-1-8 Goアプリのレスポンスを確認する

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

### Tips ここで、少し考えてみましょう。なぜ、接続に失敗するのでしょうか？  
Goアプリにて公開設定した9090ポートは、コンテナ内に限定されたコンテナポートであることを理解する必要があります。つまり、コンテナ内からコンテナポート(9090)へのアクセスは可能ですが、コンテナ外部(ローカルPC)からのコンテナポートへのアクセスは不可能であることを意味します。この場合には、Dockerのポートフォワーディング機能を利用し、ホストマシンのポートをコンテナポートに紐付け、コンテナ外との通信を可能にすることができます。以下で手順を確認します。

### 3-1-9 Dockerコンテナを停止する

#### コマンド実行
ポート設定が出来ていないコンテナを停止します。

```bash
$ docker container stop go-container
```
※Docker v1.13 以降では、 旧`docker stop`⇒新`docker container stop`コマンドが推奨されています。  
また、このコマンドでは、`docker container stop <"CONTAINER ID" or "NAME">`のように、コンテナID、または、コンテナ名を指定してコンテナを停止させることが出来ます。

#### 実行結果

```
go-container
```

### 3-1-10 Dockerコンテナ状態を確認する

#### コマンド実行
[3-1-6](https://github.com/cloudnativedaysjp/cicd-handson-2021/blob/main/docs/chapter3.md#3-1-6-docker%E3%82%B3%E3%83%B3%E3%83%86%E3%83%8A%E3%82%92%E8%B5%B7%E5%8B%95%E3%81%99%E3%82%8B)で、`--rm`オプションを指定して起動されたDockerコンテナが、コンテナ停止と共に削除されていることを確認します。

```bash
$ docker container ls
```

#### 実行結果

```
CONTAINER ID   IMAGE      COMMAND                  CREATED          STATUS          PORTS           NAMES
```

`go-container`がコンテナ一覧に存在していないことを確認します。

### 3-1-11 ポートフォワーディング設定をしてDockerコンテナを起動する

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

### 3-1-12 Dockerコンテナ状態を確認する

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

### 3-1-13 Goアプリのレスポンスを確認する

#### コマンド実行
9090ポート(コンテナ内部)⇒9091ポート(コンテナ外部)で公開されたGoアプリへGETリクエストを送信し、レスポンスを確認します。

```cmd
$ curl http://localhost:9091/health
```

#### 実行結果

```
StatusCode        : 200
StatusDescription : OK
Content           : {"status":"Healthy"}

RawContent        : HTTP/1.1 200 OK
                    Content-Length: 21
                    Content-Type: text/plain; charset=utf-8
                    Date: Sat, 21 Aug 2021 07:05:58 GMT

                    {"status":"Healthy"}

Forms             : {}
Headers           : {[Content-Length, 21], [Content-Type, text/plain; charset=utf-8], [Date, Sat, 21 Aug 2021 07:05:58 GMT]}
Images            : {}
InputFields       : {}
Links             : {}
ParsedHtml        : System.__ComObject
RawContentLength  : 21
```

9091ポートへの接続に成功し、`200`レスポンスが返却されることを確認します。

### 3-1-14 Dockerコンテナを停止する

#### コマンド実行
ポート設定が出来ていないコンテナを停止します。

```bash
$ docker container stop go-container
```

#### 実行結果

```
go-container
```

### 3-1-15 Dockerコンテナ状態を確認する

#### コマンド実行
[3-1-11](https://github.com/cloudnativedaysjp/cicd-handson-2021/blob/main/docs/chapter3.md#3-1-11-%E3%83%9D%E3%83%BC%E3%83%88%E3%83%95%E3%82%A9%E3%83%AF%E3%83%BC%E3%83%87%E3%82%A3%E3%83%B3%E3%82%B0%E8%A8%AD%E5%AE%9A%E3%82%92%E3%81%97%E3%81%A6docker%E3%82%B3%E3%83%B3%E3%83%86%E3%83%8A%E3%82%92%E8%B5%B7%E5%8B%95%E3%81%99%E3%82%8B)で、`--rm`オプションを指定して起動されたDockerコンテナが、コンテナ停止と共に削除されていることを確認します。

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

### 3-2-1 トークン情報を作成・保存する

[GitHub Docs : Creating a personal access token](https://docs.github.com/en/github/authenticating-to-github/keeping-your-account-and-data-secure/creating-a-personal-access-token)のドキュメントに従い、Dockerログイン時に使用する以下の権限を付与したPersonal Access Token(以下PAT)情報を取得します。  
- write:packages(Upload packages to github package registry)  
- read:packages(Download packages from github package registry)

![image](https://user-images.githubusercontent.com/45567889/129031847-9778cd34-5642-4d9f-bf3d-06a9b1b32089.png)

![image](https://user-images.githubusercontent.com/45567889/128994241-87aefb3a-d670-455f-9001-115c2f52fa7f.png)

任意のローカルディレクトリに、以下のPATファイルを作成します。  
※git cloneした「cicd-handson-2021」ディレクトリには、置かないで下さい。
- ファイル名：`token.txt`

```
"生成されたPATのプレーンテキスト"
```

### 3-2-2 DockerでGitHub Packagesの認証を行う

#### コマンド実行
`docker login`コマンドを使い、DockerでGitHub Packagesの認証を受けることができます。クレデンシャルをセキュアに保つ貯めに、個人アクセストークンは自分のコンピュータのローカルファイルに保存し、ローカルファイルからトークンを読み取るDockerの`--password-stdin`フラグを使うことをおすすめします。

```cmd
$ cd ~\token.txt
$ cat token.txt | docker login https://docker.pkg.github.com -u USERNAME --password-stdin
```

#### 実行結果

```
Login Succeeded
```

### 3-2-3 Dockerタグを付与する

ここでは、[3-1-4](https://github.com/cloudnativedaysjp/cicd-handson-2021/blob/main/docs/chapter3.md#3-1-4-docker-image%E3%82%92%E3%83%93%E3%83%AB%E3%83%89%E3%81%99%E3%82%8B)で作成したDockerイメージにタグ付けを行います。<GITHUB_USER>をリポジトリを所有するGitHubユーザ名で、REPOSITORYをプロジェクトを含むリポジトリの名前で、IMAGE_NAMEをパッケージもしくはイメージの名前で、VERSIONをビルドの時点のパッケージバージョンで置き換えてください。

#### コマンド実行

```bash
$ docker image tag go-image:base docker.pkg.github.com/<GITHUB_USER>/cicd-handson-2021-code/go-image:base
```
※Docker v1.13 以降では、 旧`docker tag`⇒新`docker image tag`コマンドが推奨されています。  
`<GITHUB_USER>`をGitHubユーザ名に置き換えてコマンドを実行します。

### 3-2-4 Docker image一覧を確認する

#### コマンド実行
新しくタグ付けされたDocker imageを確認します。

```bash
$ docker image ls
```

#### 実行結果

```
REPOSITORY                                                            TAG       IMAGE ID       CREATED       SIZE
go-image                                                              base      220026ab99c0   4 hours ago   862MB
docker.pkg.github.com/<GITHUB_USER>/cicd-handson-2021-code/go-image         base      220026ab99c0   4 hours ago   862MB
```

`docker.pkg.github.com/<GITHUB_USER>/cicd-handson-2021-code/go-image`が表示されていることを確認します。

### 3-2-5 GitHub PackagesへDockerイメージをpushする

#### コマンド実行

```bash
$ docker image push docker.pkg.github.com/<GITHUB_USER>/cicd-handson-2021-code/go-image:base
```
※Docker v1.13 以降では、 旧`docker push`⇒新`docker image push`コマンドが推奨されています。  
`<GITHUB_USER>`をGitHubユーザ名に置き換えてコマンドを実行します。

#### 実行結果

```
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

![image](https://user-images.githubusercontent.com/45567889/130482375-96c65eb2-429d-453d-a311-a00390e24c94.png)

実行結果が、全て`Pushed`になっていることを確認し、GitHub画面のPackagesに`go-image`が公開されていることを確認します。

## 3-3 Kubernetesの基本操作をおさらいする

ここでは、Minikubeを使って、Kubernetesクラスタの起動からオブジェクトの作成までの操作を確認します。

### 3-3-1 Minikubeを起動する

[Minikube](https://github.com/cloudnativedaysjp/cicd-handson-2021/blob/main/docs/chapter0.md#minikube)を参考に、`minikube start`、`kubectl get nodes`、`kubectl get pods`コマンドを実行し、Minikubeを起動します。

### 3-3-2 ローカル環境のDockerイメージを削除する

#### コマンド実行
Github PackagesにPushしたDockerイメージを使って、ポッドを作成するためのマニフェストを記述する前に、[3-1-4](https://github.com/cloudnativedaysjp/cicd-handson-2021/blob/main/docs/chapter3.md#3-1-4-docker-image%E3%82%92%E3%83%93%E3%83%AB%E3%83%89%E3%81%99%E3%82%8B)と[3-2-4](https://github.com/cloudnativedaysjp/cicd-handson-2021/blob/main/docs/chapter3.md#3-2-4-docker%E3%82%BF%E3%82%B0%E3%82%92%E4%BB%98%E4%B8%8E%E3%81%99%E3%82%8B)で作成した2つのローカル環境のDockerイメージを削除しておきます。

```bash
$ docker image rm go-image:base
$ docker image rm docker.pkg.github.com/<GITHUB_USER>/cicd-handson-2021-code/go-image:base
```
※`<GITHUB_USER>`をGitHubユーザ名に置き換えてコマンドを実行します。

#### 実行結果

```
Untagged: go-image:base
Untagged: docker.pkg.github.com/<GITHUB_USER>/cicd-handson-2021-code/go-image:base
Untagged: docker.pkg.github.com/<GITHUB_USER>/cicd-handson-2021-code/go-image@sha256:11dd65371181d74b33c84b18f4f6ba87537cdbab7c884ef12ee6429865c0f640
```
※`<GITHUB_USER>`は、GitHubユーザ名に置き換わっている状態。

### 3-3-3 Docker image一覧を確認する

#### コマンド実行
新しくタグ付けされたDocker imageを確認します。

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
- `--docker-password`：[3-2-1](https://github.com/cloudnativedaysjp/cicd-handson-2021/blob/main/docs/chapter3.md#3-2-1-%E3%83%88%E3%83%BC%E3%82%AF%E3%83%B3%E6%83%85%E5%A0%B1%E3%82%92%E4%BD%9C%E6%88%90%E4%BF%9D%E5%AD%98%E3%81%99%E3%82%8B)で作成したGitHubのPAT(Personal Access Token)を指定します。
- `--docker-email`：GitHub登録時のメールアドレスを指定します。

```bash
$ kubectl create secret docker-registry --save-config dockerconfigjson-github-com \
   --docker-server=docker.pkg.github.com \
   --docker-username=<GITHUB_USER> \
   --docker-password=<PERSONAL_ACCESS_TOKEN> \
   --docker-email=<GITHUB_EMAIL>
```
または、
```bash
$ kubectl create secret docker-registry --save-config dockerconfigjson-github-com --docker-server=docker.pkg.github.com --docker-username=<DOCKER_USER> --docker-password=<PERSONAL_ACCESS_TOKEN> --docker-email=<DOCKER_EMAIL>
```
※ghcr.io への読み書きについて、[Working with the Container registry](https://docs.github.com/en/packages/working-with-a-github-packages-registry/working-with-the-container-registry#authenticating-to-the-container-registry)によると、`--docker-password`には、GitHubのPATを指定する必要があります。間違えてGitHub登録時のパスワードを入力しないよう注意が必要です。

#### 実行結果

```
secret/dockerconfigjson-github-com created
```

### 3-3-5 マニフェストファイルを作成する

Kubernetesでは、作成するポッドのリソース構成をマニフェストファイルにコードで記述することができます。

ワークディレクトリに、以下のマニフェストファイルを作成します。  
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
※`imagePullSecrets`に[3-3-4](https://github.com/cloudnativedaysjp/cicd-handson-2021/blob/main/docs/chapter3.md#3-3-4-docker%E3%82%B3%E3%83%B3%E3%83%86%E3%83%8A%E3%83%AC%E3%82%B8%E3%82%B9%E3%83%88%E3%83%AA%E8%AA%8D%E8%A8%BC%E7%94%A8%E3%81%AE%E3%82%AF%E3%83%AC%E3%83%87%E3%83%B3%E3%82%B7%E3%83%A3%E3%83%ABsecret%E3%82%92%E4%BD%9C%E6%88%90%E3%81%99%E3%82%8B)で作成したクレデンシャル(Secret)保存名`dockerconfigjson-github-com`を指定し忘れないよう注意が必要です。  
`<GITHUB_USER>`は、GitHubユーザ名に置き換わっている状態。

### 3-3-6 ポッドを作成する

#### コマンド実行
マニフェストファイルからポッドを作成します。
- `-f`：ファイル名を指定します。

```cmd
$ cd "goapp.yamlを保存したローカルディレクトリ"
$ kubectl apply -f goapp.yaml
```

#### 実行結果

```
deployment.apps/goapp-deployment created
```

### 3-3-7 ポッド一覧を確認する

#### コマンド実行
- `-o wide`：各PodのIPアドレスを表示します。

```bash
$ kubectl get deploy,pods -o wide
```

#### 実行結果

```
NAME                                    READY   UP-TO-DATE   AVAILABLE   AGE     CONTAINERS   IMAGES                                                                     SELECTOR
deployment.apps/goapp-deployment        1/1     1            1           8m56s   goapp        docker.pkg.github.com/<GITHUB_USER>/cicd-handson-2021-code/go-image:base   app=goapp

NAME                                    READY   STATUS       RESTARTS    AGE     IP           NODE       NOMINATED NODE   READINESS GATES
pod/goapp-deployment-6c85ff5cc-6pc89    1/1     Running      0           8m56s   172.17.0.2   minikube   <none>           <none>
```
`STATUS`が`Running`になっていることを確認します。  
※`<GITHUB_USER>`は、GitHubユーザ名に置き換わっている状態。

### 3-3-8 ポートフォアーディング設定を行う。

ローカルのポートを任意のPodのポートに転送するための設定を行います。ここではローカルの9090番ポートをdeploymentの9092番ポートに転送します。

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

### 3-3-9 アプリへ接続確認する

#### コマンド実行
任意のブラウザにて、以下のURLを実行し、疎通できていることを確認します。

```bash
http://localhost:9091/health
```

#### 実行結果

```
"status":"Healthy"
```
`"status":"Healthy"`が表示されていることを確認します。

### 3-3-10 ポッドを削除する

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
