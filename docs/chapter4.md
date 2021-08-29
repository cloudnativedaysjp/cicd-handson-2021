# Chapter 4 Optimize Container Image

コンテナイメージを最適化するため、Chapter 3 で作成した以下の `Dockerfile` を編集します。

```text:./apps/Dockerfile
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

コンテナイメージの最適化にはいくつかのポイントがあります。

* コンテナイメージのビルド時間を短くすること
* コンテナイメージのサイズを小さくすること
* コンテナイメージのセキュリティを高めること
* など

Chapter 4では、`コンテナイメージのサイズを小さくする` ように最適化していきましょう。

# Dockerfile と コンテナイメージを確認

まずは先程のイメージを再度ビルドして、コンテナイメージのサイズを確認します。

```bash
$ cd ./cicd-handson-2021-code/apps

$ docker image build -t go-image:base .
[+] Building 1.9s (10/10) FINISHED
...
=> naming to docker.io/library/go-image:base    0.0s

$ docker image ls
REPOSITORY                    TAG       IMAGE ID       CREATED          SIZE
go-image                      base      e9e77e06562e   10 seconds ago   959MB
```

SIZE欄にあるように、コンテナイメージサイズは `959MB` であることが分かります。

現在の Dockerfile では、`golang:latest` をPullし、go buildを行うことでアプリケーションのビルドを行っています。つまり、この Dockerfile を使用して、コンテナイメージをビルドする場合は、**アプリケーションのビルドに必要な goライブラリを含んだ状態でコンテナイメージがビルドされる** ことになり、アプリケーション実行時には必要の無いコンポーネントが含まれています。このため、コンテナイメージサイズが肥大化します。

# Dockerfile の編集
## マルチステージビルド と distrolessイメージの活用
以下のように Dockerfile を編集します。

```text:./apps/Dockerfile
# [編集] ベースイメージをbuilderイメージとして指定
FROM golang:1.16 as builder

# ワークディレクトリを指定
WORKDIR /app

# ホストOSのapps内全てをWORKDIRにコピー
COPY . ./

# ビルド
RUN go build -o ./server-run ./server

# [追加] 軽量なdistrolessイメージを指定
FROM gcr.io/distroless/base

# コンテナのポートを9090で公開
EXPOSE 9090

# [追加] golang:1.16でビルドしたアプリケーションをコピー
COPY --from=builder app/server-run /.
COPY web /web

# アプリ実行
CMD [ "./server-run" ]
```

通常Dockerfileを記述する際には、 `FROM` 行をひとつ記述し、任意のベースイメージを指定します。今回のケースだと `FROM golang:latest`です。ソースコードをコンテナ内でビルドする場合は必要なライブラリをイメージ内に含める必要があるため、最終的に作成するコンテナイメージのサイズが肥大化してしまいます。この問題を、マルチステージビルドを用いることで解決することができます。

マルチステージビルドを行う際は、Dockerfile 内に `FROM` を複数行記述します。最後に記述された `FROM` ステージが最終的なコンテナイメージとなります。

今回の手順では、最初に記述する `FROM golang:1.16 as builder` でアプリケーションをビルドし、最後に記述する `FROM gcr.io/distroless/base` のコンテナ内にコピーします。このようにすることで、最終的なコンテナイメージには、golangのライブラリなどを含まずにアプリケーションのみを残すことができます。

>[補足1]
>
>`golang:latest`のように`latest`タグを使用すると、タイミング次第では想定外のイメージをプルすることになりトラブルとなる可能性があります。このため、基本的には `golang:1.16` のように特定のバージョンを指定します。
>
>[補足2]
>
>`distrolessイメージ` は、アプリケーションの実行に特化したコンテナイメージであり、パッケージマネージャ、シェル、や不要なプロセスを含まないなど、必要最低限のコンポーネントで構成されており、Googleが公開しています。イメージサイズを縮小するだけでなくセキュリティ面の脆弱性の軽減も期待できます。
>
>\[GitHub\]: [GoogleContainerTools/distroless](https://github.com/GoogleContainerTools/distroless)
>

## コンテナイメージのビルド
ここでは、コンテナイメージtagを `distroless` にしてビルドします。

```bash
$ docker image build -t go-image:distroless .
```

## コンテナイメージのサイズの再確認

```bash
$ docker image ls
REPOSITORY                    TAG       IMAGE ID       CREATED              SIZE
go-image                      distroless    7677fd6819ba   7 seconds ago        27.1MB
go-image                      base          e9e77e06562e   About a minute ago   959MB
```

SIZE欄を比較すると、圧倒的なサイズ差があることが分かります。このように必要最低限のコンポーネントのみを取り込んだ状態でビルドすることで、サイズの縮小、ビルドスピードの向上につながり、コンテナイメージを最適化できます。

最後にリポジトリにDockerfileをPushしておきましょう。

```bash
$ git add Dockerfile
$ git commit -m "Fix Dockerfile"
$ git push origin main
```

# おまけ: scratchイメージの利用
Dockerfileの編集
```text:./apps/Dockerfile
# ベースイメージをbuilderイメージとして指定
FROM golang:1.16 as builder

# ワークディレクトリを指定
WORKDIR /app

# ホストOSのapps内全てをWORKDIRにコピー
COPY . ./

# [編集] ビルド
RUN CGO_ENABLED=0 go build -o ./server-run ./server

# [編集] scratchイメージを指定
FROM scratch

# コンテナのポートを9090で公開
EXPOSE 9090

# [追加] golang:1.16でビルドしたアプリケーションをコピー
COPY --from=builder app/server-run /.
COPY web /web

# アプリ実行
CMD [ "./server-run" ]
```

>[補足]
>
>C言語のコードを使用していないため、go build時に `CGO_ENABLED=0` で `cgo` を無効化することで静的リンクにしています。

コンテナイメージのビルド
```bash
$ docker image build -t go-image:scratch .
```

コンテナイメージのサイズ確認
```bash
$ docker image ls

REPOSITORY                    TAG       IMAGE ID       CREATED         SIZE
go-image                      scratch      5b871daf4d8f   13 minutes ago   7.85MB
go-image                      distroless   7677fd6819ba   2 hours ago      27.1MB
go-image                      base         e9e77e06562e   2 hours ago      959MB
```

コンテナ起動
```bash
$ docker container run --name go-container -d -p 9091:9090 go-image:scratch
```

動作確認
```bash
$ curl http://localhost:9091/health

{"status":"Healthy"}
```

# Chapter 4 まとめ
Chapter 4 では、`マルチステージビルド` と `distroless イメージの利用` によってコンテナイメージのサイズを小さくする最適化を行いました。ただし、デバッグや何らかの目的に合わせて、必要なツールをコンテナイメージに取り込んでおくことも、障害対応の観点では重要です。サービスの要件に合わせてコンテナイメージを作成するようにしましょう。
