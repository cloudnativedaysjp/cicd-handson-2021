# Chapter 9 CI/CD pipeline

## 9-1 アプリケーションの更新

アプリケーションを以下の手順で更新して、codeリポジトリに「$ git push」を実行します。

```
```

GitHub ActionsのCIの実行、GitHub Packagesへのコンテナイメージ保存、configリポジトリへのプルリクエストまでできることを確認します。

そして、configリポジトリのプルリクエストをマージします。

## 9-2 Argo CDによる同期の確認

マージ後、自動の場合はconfigリポジトリと同期するまで待ちます。
手動の場合はArgo CDのWebUIで「REFRESH」ボタンをクリックしてポーリング（手動）して、configリポジトリと同期させます。

同期後、WebUIでKubernetesクラスタに自動デプロイされること確認します。

デプロイ後、ブラウザでアプリケーションが更新されていることを確認します。

## 9-3 削除処理

学習終了後に各ツールの削除処理を実行してください。
削除処理のタイミングは、ご自身にお任せします。

この処理を実行するとこれまで作成してきたデータ等は削除されるので、
必要な場合はご自身でバックアップをお願いします。

### GitHub

GitHubに作成したものを削除します。

* codeとconfigリポジトリ

https://docs.github.com/ja/github/administering-a-repository/managing-repository-settings/deleting-a-repository

* Personal Access Token

以下URLにアクセスして、対象となるTokenリストにある「Delete」ボタンをクリックして削除してください。

https://github.com/settings/tokens

### minikube

Kubernetesクラスタを削除します。

```bash
$ minikube delete
```

### Others

事前準備でインストールした以下アプリケーションは、
不要であればアンインストールしてください。

* Git（Windows環境のみ）
* minikube
* Docker Desktop for Win/Mac
* Argo CD CLI

#### Windows

以下については、[コントロールパネル]-[プログラムと変更]から、対象のツールを選んでアンインストールしてください。

* Git（Windows環境のみ）
* minikube
* Docker Desktop for Win

以下については、任意のディレクトリに格納したexeファイルを削除してください。
通したパスも不要であれば削除してください。

パス設定参考サイト: https://www.atmarkit.co.jp/ait/articles/1805/11/news035.html

* Argo CD CLI

#### Mac

* minikube

Homebrewの場合

```bash
$ brew uninstall minikube
```

Binaryの場合

```bash
$ rm /usr/local/bin/minikube
```

* Docker Desktop for Mac

https://docs.docker.com/docker-for-mac/install/

* Argo CD CLI

```brew
$ brew uninstall argocd
```

#### GCP

* GKE

GCPダッシュボードの場合

対象となるKubernetesクラスタを選択して、削除してください。

gcloudコマンドの場合

```bash
$ gcloud container clusters delete cicd-cluster --zone <your-zone> --async
```

* Argo CD CLI

```bash
$ rm ~/argocd/bin/argocd
```
