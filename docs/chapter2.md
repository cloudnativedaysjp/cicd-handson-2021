# Chapter 2 Application

ハンズオンで使用するデモアプリについて説明します。

## Simple CICD Landscape!

本アプリは [Cloud Native Landscape](https://landscape.cncf.io/) の中から、CICDに関するプロジェクトの一覧を表示します。

![app](../apps/screenshot.png)

アプリ仕様

- CNCF Relation によるプロジェクトの分類
- クリックでプロジェクトの GitHub（ない場合はホームページ）へ移動
- GitHub のスター数をAPI経由で取得して表示
  - 起動中はメモリ上にキャッシュされる
- Description が設定されている場合はマウスオーバーで表示

特徴

- シングルバイナリでWebページとサーバをホスト
- Go言語で開発（WebAssembly、サーバ実装 ＋ テストコード）

## テスト・動作確認

`apps` ディレクトリに移動し、作業環境に合わせたコマンドを実行してください。

```bash
# apps ディレクトリに移動
cd apps
```

### macOS, Linux

```bash
# テスト実行
go test ./server

# サーバ起動
go run ./server
```

### Windows

```pwsh
# テスト実行
go test .\server\

# サーバ起動
go run .\server\
```

現時点では、テストの実行結果は失敗となります。

サーバを起動後、次のいずれかの方法で動作確認を行ってください。

### コマンドラインで確認

ヘルスチェック用のエンドポイントにアクセスし、レスポンスを確認します。

```bash
curl http://127.0.0.1:9090/health

# json 形式のレスポンス
{
  "status": "Healthy"
}
```

### Webブラウザで確認

ブラウザで <http://127.0.0.1:9090> にアクセスし、Webページが表示されることを確認します。

ヘルスチェック用のエンドポイント (<http://127.0.0.1:9090/health>)、API 用のエンドポイント (<http://127.0.0.1:9090/landscape>) にもアクセス可能です。

## ※注意点

GitHub スター数の取得に利用している GitHub API には 60回/h の利用制限があり、これを超えるとスター数を取得できなくなります。<br/>
そのためアプリの再起動を数回行うとスター数が表示されなくなりますが、アプリの動作自体に影響はありません。

参考：<https://docs.github.com/en/rest/overview/resources-in-the-rest-api#rate-limiting>
