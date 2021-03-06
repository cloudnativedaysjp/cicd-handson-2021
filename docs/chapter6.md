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
          DOCKER_BUILDKIT=1 docker image build apps/ -t docker.pkg.github.com/${{ github.repository }}/go-image:${{ github.run_number }}

        # イメージをプッシュする為の「docker login」
      - name: GitHub Packages login
        uses: docker/login-action@v1
        with:
          registry: docker.pkg.github.com
          username: ${{ secrets.USERNAME }}
          password: ${{ secrets.PERSONAL_ACCESS_TOKEN }}

        # コンテナイメージをGitHub Packagesに「docker image push」
      - name: Push image to GitHub Packages
        run: docker image push docker.pkg.github.com/${{ github.repository }}/go-image:${{ github.run_number }}
```

```diff
diff --git a/.github/workflows/main.yml b/.github/workflows/main.yml
index e09bb49..8e9f670 100644
--- a/.github/workflows/main.yml
+++ b/.github/workflows/main.yml
@@ -13,6 +13,12 @@ jobs:
       - name: Checkout code
         uses: actions/checkout@v2

+        # アプリケーションテスト
+      - name: Application test
+        run: |
+          cd apps
+          make run-test
+
         # BuildKitによるコンテナイメージビルド
       - name: Build an image from Dockerfile
         run: |
```

`main.yml` の修正をしたらリポジトリにプッシュします

```git
$ git add .github/workflows/main.yml
$ git commit -m "add test step"
$ git push origin main
```

GitHub Actionsのページを確認すると失敗していることが確認できます。

![テスト追加後に失敗したGitHub Actionsの確認](images/chapter6/chapter06-001.png)

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

GitHub Actionsのページを確認すると成功していることが確認できます。
問題なくCIが進み、GitHub Packagesに新しいイメージが格納されていれば成功です。

![テスト修正後に成功したGitHub Actionsの確認](images/chapter6/chapter06-002.png)
![テスト修正後に成功したGitHub Actionsによってプッシュされたイメージ](images/chapter6/chapter06-003.png)

---
[Chapter 7 Secure container image](chapter7.md)へ