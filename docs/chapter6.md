# Chapter 6 Application testing

Chapter 5 で作成した `main.yml` にアプリケーションテストのstepを追加します。

## 6-1 アプリケーションテストの追加

最初に追加するアプリケーションテストは、失敗します。
その段階でCIはエラーとなり後続するコンテナイメージビルド、コンテナレジストリへのプッシュ処理は実行されません。
実際に実行して確認してみましょう。

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
        #アプリケーションテストが失敗する内容#

        # BuildKitによるコンテナイメージビルド
      - name: Build an image from Dockerfile
        run: |
          DOCKER_BUILDKIT=1 docker image build . -f app/Dockerfile -t docker.pkg.github.com/${{ github.repository }}/gitops-go-app:${{ github.run_number }}

        # コンテナイメージをGitHub Packagesに「docker image push」
      - name: GitHub Packages login
        uses: docker/login-action@v1
        with:
          registry: docker.pkg.github.com
          username: ${{ secrets.USERNAME }}
          password: ${{ secrets.PERSONAL_ACCESS_TOKEN }}

      - name: Push image to GitHub Packages
        run: docker image push docker.pkg.github.com/${{ github.repository }}/gitops-go-app:${{ github.run_number }}
```

アプリケーションテストが通るように修正してみましょう。

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
        #アプリケーションテストが成功する内容#

        # BuildKitによるコンテナイメージビルド
      - name: Build an image from Dockerfile
        run: |
          DOCKER_BUILDKIT=1 docker image build . -f app/Dockerfile -t docker.pkg.github.com/${{ github.repository }}/gitops-go-app:${{ github.run_number }}

        # コンテナイメージをGitHub Packagesに「docker image push」
      - name: GitHub Packages login
        uses: docker/login-action@v1
        with:
          registry: docker.pkg.github.com
          username: ${{ secrets.USERNAME }}
          password: ${{ secrets.PERSONAL_ACCESS_TOKEN }}

      - name: Push image to GitHub Packages
        run: docker image push docker.pkg.github.com/${{ github.repository }}/gitops-go-app:${{ github.run_number }}
```

問題なくテストが通り、GitHub Packagesにイメージが格納されていれば成功です。
