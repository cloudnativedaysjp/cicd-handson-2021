# Chapter 5 CI pipeline by GitHub actions

## 5-1 GitHub Actionsを利用したContinuous Integration（CI）

GitHub Actionsは、GitHubに組み込まれたCI/CDシステムです。CI/CDに限らず、様々なイベントフックに対応してイベントをトリガーに自動処理を実行できます。

GitHubリポジトリの「.github/workflow/main.yml」ファイルに定義して利用します。

Chapter01～03で手動実行した、アプリケーションテスト、コンテナイメージビルド、コンテナイメージレジストリへのプッシュを「.github/workflow/main.yml」ファイルに定義して自動実行します。

各GitHubリポジトリのActions用のコンテナで実行される仕組みです。

以下は、ソースコードを変更後、GitHubリポジトリへの「$ git push」をトリガーに、コンテナイメージビルド、コンテナイメージレジストリ（GitHub Packages）へ「$ docker image push」が実行される定義です。

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

        # BuildKitによるコンテナイメージビルド
      - name: Build an image from Dockerfile
        run: |
          DOCKER_BUILDKIT=1 docker image build apps/ -t docker.pkg.github.com/${{ github.repository }}/go-image:${{ github.run_number }}

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
