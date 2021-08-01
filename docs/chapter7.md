# Chapter 7 Secure container image

Kubernetesでは、コンテナイメージにOSレベルの脆弱性を含んでしまうことがあります。コンテナイメージのレベルで診断を行い、セキュアな環境を保持するためにCIで脆弱性診断を行います。

## 7-1 dockleによるコンテナイメージ診断

dockleは、CIS BenchmarkのDockerに関する項目、Dockerfileのベストプラクティスをベースにチェック可能な診断ツールです。
このdockleを定義に追加して実行してみましょう。

```
name: GitHub Actions CI

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

        # dockleによるイメージ診断
      - name: Run dockle
        uses: hands-lab/dockle-action@v1
        with:
          image: docker.pkg.github.com/${{ github.repository }}/gitops-go-app:${{ github.run_number }}

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

## 7-2 Trivyによるコンテナイメージ脆弱性診断

Trivyは、OSパッケージ情報、アプリケーションの依存関係などから脆弱性を検出します。定義に追加して実行してみましょう。

```
name: GitHub Actions CI

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

        # dockleによるイメージ診断
      - name: Run dockle
        uses: hands-lab/dockle-action@v1
        with:
          image: docker.pkg.github.com/${{ github.repository }}/gitops-go-app:${{ github.run_number }}

        # Trivyによるイメージスキャン
      - name: Run Trivy vulnerability scanner
        uses: aquasecurity/trivy-action@master
        with:
          image-ref: 'docker.pkg.github.com/${{ github.repository }}/gitops-go-app:${{ github.run_number }}'
          format: 'table'
          exit-code: '1'
          ignore-unfixed: true
          severity: 'CRITICAL,HIGH'

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
