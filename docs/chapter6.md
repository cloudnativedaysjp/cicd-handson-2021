# Chapter 6 Application testing

Chapter 5 で作成した `main.yml` にアプリケーションテストのstepを追加します。

## 6-1 アプリケーションテストの追加

最初に追加するアプリケーションテストは、失敗します。  
その段階でCIはエラーとなり後続するコンテナイメージビルド、コンテナレジストリへのプッシュ処理は実行されません。  
実際に実行して確認してみましょう。  
今回のアプリケーションはmakeコマンドでテストすることができるようになっています。  

```yaml
name: GitHub Actions CI

# mainブランチへの「git push」をトリガー
on:
  push:
    branches: [ main ]

jobs:
  build:
    name: GitOps Workflow
    runs-on: ubuntu-latest
    steps:
      - name: Checkout code
        uses: actions/checkout@v2

        # アプリケーションテスト
      - name: Application test
        run: |
          cd apps
          make run-test

        # BuildKitによるコンテナイメージビルド
      - name: Build an image from Dockerfile
        run: |
          DOCKER_BUILDKIT=1 docker image build . -f apps/Dockerfile -t docker.pkg.github.com/${{ github.repository }}/go-image:${{ github.run_number }}

        # コンテナイメージをGitHub Packagesに「docker image push」
      - name: GitHub Packages login
        uses: docker/login-action@v1
        with:
          registry: docker.pkg.github.com
          username: ${{ secrets.USERNAME }}
          password: ${{ secrets.PERSONAL_ACCESS_TOKEN }}

      - name: Push image to GitHub Packages
        run: docker image push docker.pkg.github.com/${{ github.repository }}/go-image:${{ github.run_number }}
```

`main.yaml` の修正をしたらリポジトリにプッシュします

```git
$ git add .github/workflows/main.yaml
$ git commit -m "add test step"
$ git push origin main
```

## 6-2 アプリケーションの修正

アプリケーションテストが通るように修正してみましょう。  
問題のあるコードは`apps/server/landscape.go` です。  
100 行目を修正します。  
`proj.Twitter` となっているところを `proj.Crunchbase` に修正して保存します。

```go
# エディタ等で修正します
$ vi apps/server/landscape.go
...
...
 92                                                 go func(i int, proj SubItem) {
 93                                                         defer wg.Done()
 94                                                         list[i] = Project{
 95                                                                 Name:        proj.Name,
 96                                                                 Description: proj.Description,
 97                                                                 HomepageUrl: proj.HomepageUrl,
 98                                                                 Project:     getProject(proj.Project, proj.Crunchbase, ml),
 99                                                                 RepoUrl:     proj.RepoUrl,
100                                                                 Crunchbase:  proj.Crunchbase, #ここを直す
101                                                                 StarCount:   getStarCount(proj.RepoUrl),
102                                                         }
103                                                 }(i, proj)
```

`landscape.go` を修正したらテストが通るようになっているか確認します。

```bash
$ cd apps
$ make run-test

# make コマンドが無い場合はこちらで確認してください
$ go test ./server
```

「ok」が出力されたら修正成功です。

```
ok  	github.com/cloudnativedaysjp/cicd-handson-2021/cicd-landscape/server
```

修正したコードをリポジトリにプッシュします。

```git
$ git add apps/server/landscape.go
$ git commit -m "fix test"
$ git push origin main
```

問題なくCIが進み、GitHub Packagesにイメージが格納されていれば成功です。
