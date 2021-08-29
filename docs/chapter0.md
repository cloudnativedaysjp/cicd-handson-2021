# Chapter 0

このハンズオンでは、下記の事前準備が必要です。
それぞれについて、手順を紹介します。

* GitHubアカウント
* Git（Windows環境のみ）
* Docker
* Kubernetesクラスタ（ローカルのKubernetesでも可）
* Argo CD CLI

## GitHubアカウント

今回の演習では、主にGitHubを利用して行います。
そのため、事前にGitHubアカウントを作成してください。
https://github.com/signup


## GitHubのPersonal Access Tokenの取得

[GitHub Docs : Creating a personal access token](https://docs.github.com/en/github/authenticating-to-github/keeping-your-account-and-data-secure/creating-a-personal-access-token)のドキュメントに従い、Dockerログイン時に使用する以下の権限を付与したPersonal Access Token(以下PAT)情報を取得します。

- workflow (Update GitHub Action workflows)
- write:packages (Upload packages to GitHub Package Registry)
- delete:packages (Delete packages from GitHub Package Registry)
- admin:org (Full control of orgs and teams, read and write org projects)

![image](https://user-images.githubusercontent.com/45567889/129031847-9778cd34-5642-4d9f-bf3d-06a9b1b32089.png)

![image](https://user-images.githubusercontent.com/45567889/128994241-87aefb3a-d670-455f-9001-115c2f52fa7f.png)

任意のローカルディレクトリに、以下のPATファイルを作成します。ハンズオン当日に利用するため、なくさないようにしてください。
また、このPATを利用することで自身のユーザと同等の権限を持つことができてしまうため、流出しないように注意してください。
※git cloneした「cicd-handson-2021-code」「cicd-handson-2021-config」ディレクトリには、置かないで下さい。
- ファイル名：`token.txt`

```
"生成されたPATのプレーンテキスト"
```

## Git

Windows環境では、デフォルトでGitコマンドを実行できないため、Git for Windowsをインストールしてください。
以下サイトからダウンロードできます。

https://gitforwindows.org/


インストール後、コマンドプロンプトまたはGit Bashターミナルで下記のように「git version」コマンドを実行できれば問題ありません。

```git
$ git version
git version 2.32.0.windows.2
```

※この作業は、Windows環境のみです。

## Go

今回はGo言語で書かれたアプリケーションを用いてハンズオンを行っていきます。
そのため、ビルドするためにGoのランタイムが必要なgoバイナリのインストールを行います。
下記の手順を参考に「1.16」系のバージョンを導入してください。

* インストール手順
	* https://golang.org/doc/install
* バージョン1.16のインストーラのダウンロード元
	* https://golang.org/dl/

インストール後、コマンドプロンプトまたはGit Bashターミナルで下記のように「go version」コマンドを実行し、1.16.7が表示できれば問題ありません。
すべてのバージョンで動作確認は行っていませんが、その他のバージョンでも動作すると思います。

```bash
$ go version
go version go1.16.7 darwin/amd64
```


## Docker

手元の環境でDockerが利用できるように、下記の手順に従って利用しているマシンにDockerをインストールしてください。
macOS・Windows・Linuxなど、様々なOSに対応しています。

https://docs.docker.com/get-docker/

dockerコンテナを確認する「docker ps」コマンドが正常に終了することを確認してください。
エラーが出力されていなければ、下記のように1つもコンテナが表示されていなくても問題ありません。

```bash
$ docker container ls
CONTAINER ID   IMAGE     COMMAND   CREATED   STATUS    PORTS     NAMES
```

## Kubernetesクラスタ

今回のハンズオンを試す環境は、インターネットへの外向き通信が許可されているKubernetesクラスタであれば試すことができます。そのため、Minikube・kind(Kubernetes in Docker)・microk8sなどのローカル上でクラスタが立ち上がる様々なソフトウェアも利用できます。

今回はローカルKubernetesとしてMinikubeを利用した手順と、クラウドのマネージドKubernetesサービスのGKE（Google Kubernetes Engine）を紹介します。

### Minikube

Minikubeはデフォルトで"type: LoadBalancer"のServiceが利用できたり、様々な機能を有効化することもできるため、軽く試すにはうってつけです。
ローカルマシンに 2 CPU / 8GB Memory以上の余剰リソースがある場合はローカルKubernetesを利用することが可能です。

利用を開始するには、次の手順に従って各OSごとに適切なバイナリをインストールしてください。
https://minikube.sigs.k8s.io/docs/start/

今回はHomebrewでインストールした下記のバージョンのMinikubeを利用します。

```bash
$ minikube version
minikube version: v1.15.1
commit: 23f40a012abb52eff365ff99a709501a61ac5876
```

バイナリをインストール後は、「minikube start」コマンドを利用してKubernetesクラスタを立ち上げます。minikubeでは、裏側で仮想化技術を利用してクラスタを立ち上げます。KVM・VirtualVox・Hyper-V・Hyperkit・Dockerなど、様々な仮想化ドライバーを選択することができます。今回は、Docker driverを利用してKubernetesクラスタを立ち上げます。
Docker driverでは、Dockerコンテナを立ち上げ、そのコンテナを1つのコンピュータと見立てて、Kubernetesノードとして利用します。そのため、Kubernetesクラスタ上でコンテナAを起動すると、そのコンテナAはホストマシン上で起動しているコンテナの上で起動されている状態になります（nested container）。

```bash
# Minikubeを利用してKubernetesクラスタの起動
$ minikube start --driver=docker
😄  Darwin 10.15.7 上の minikube v1.15.1

✨  設定を元に、 docker ドライバを使用します
👍  コントロールプレーンのノード minikube を minikube 上で起動しています
🚜  Pulling base image ...
💾  Kubernetes v1.19.4 のダウンロードの準備をしています
    > preloaded-images-k8s-v6-v1.19.4-docker-overlay2-amd64.tar.lz4: 486.35 MiB
🔥  docker container (CPUs=2, Memory=8100MB) を作成しています...
🐳  Docker 19.03.13 で Kubernetes v1.19.4 を準備しています...
🔎  Kubernetes コンポーネントを検証しています...
🌟  有効なアドオン: default-storageclass, storage-provisioner
🏄  Done! kubectl is now configured to use "minikube" cluster and "default" namespace by default
```

正常にKubernetesクラスタの起動ができたら、Kubernetesクラスタの情報を確認してみます。下記のようにminikubeのノードがReadyになっており、各PodがRunningになっていることを確認します。

```bash
# Minikubeで起動したKubernetesクラスタのノード一覧を表示
$ kubectl get nodes
NAME       STATUS   ROLES    AGE     VERSION
minikube   Ready    master   2m14s   v1.19.4

# Kubernetesクラスタ上で起動しているPodの一覧を表示
$ kubectl get pods -A -owide
NAMESPACE     NAME                               READY   STATUS    RESTARTS   AGE     IP             NODE       NOMINATED NODE   READINESS GATES
kube-system   coredns-f9fd979d6-nqnc9            1/1     Running   0          2m18s   172.17.0.2     minikube   <none>           <none>
kube-system   etcd-minikube                      1/1     Running   0          2m21s   192.168.49.2   minikube   <none>           <none>
kube-system   kube-apiserver-minikube            1/1     Running   0          2m21s   192.168.49.2   minikube   <none>           <none>
kube-system   kube-controller-manager-minikube   1/1     Running   0          2m21s   192.168.49.2   minikube   <none>           <none>
kube-system   kube-proxy-6zcvh                   1/1     Running   0          2m18s   192.168.49.2   minikube   <none>           <none>
kube-system   kube-scheduler-minikube            1/1     Running   0          2m21s   192.168.49.2   minikube   <none>           <none>
kube-system   storage-provisioner                1/1     Running   2          2m22s   192.168.49.2   minikube   <none>           <none>
```

このハンズオンを終了後、クラスタを削除するには下記のコマンドを実行してください。

```bash
# minikubeクラスタの削除
$ minikube stop
✋  ノード "minikube" を停止しています...
🛑  SSH 経由で「minikube」の電源をオフにしています...
🛑  1台のノードが停止しました。
```

### GKE (Google Kubernetes Engine)

GKEは、Kubernetesの元となったBorgを開発したGoogleが提供している、マネージドKubernetesサービスです。Kubernetesのコントロールプレーンやワーカーノード、様々なシステムコンポーネントをGoogleが管理してくれます。
リモート環境にKubernetesクラスタが起動するため、ローカルマシンに潤沢なリソースが存在しない場合でも問題ありません。

GCPアカウントをまだ作成していない場合、下記の手順に従って作成してください。
GCPを利用したことがない方はアカウント発行時に$300の無料枠があるため、この無料枠を使ってクラスタを立ち上げることも可能です。
https://cloud.google.com/free/?hl=ja

また、GCPを操作するためのCLIツール「gcloud」を下記の手順に沿ってインストールしてください。
https://cloud.google.com/sdk/docs/install?hl=ja

アカウントを作成し、CLIインストール後は、下記の手順に従ってKubernetesクラスタを作成します。
（参考：https://cloud.google.com/kubernetes-engine/docs/quickstart）

CLIを利用できるように、まずはログイン処理を行います。Webブラウザが起動するため、そこでGoogleアカウントの認証情報を入力してください。

```bash
# GCPへのログイン
$ gcloud auth login
```

次に、gcloudコマンドを利用してGKEクラスタを起動します。
初期状態では複数リージョンにまたがったクラスタが起動されるため、1ノードを指定しても各リージョンに1ノードずつ起動し、3ノードのクラスタが作成されます。

```bash
# GKEを利用してKubernetesクラスタを作成
$ gcloud container clusters create cicd-cluster --num-nodes=1
WARNING: Starting in January 2021, clusters will use the Regular release channel by default when `--cluster-version`, `--release-channel`, `--no-enable-autoupgrade`, and `--no-enable-autorepair` flags are not specified.
WARNING: Currently VPC-native is not the default mode during cluster creation. In the future, this will become the default mode and can be disabled using `--no-enable-ip-alias` flag. Use `--[no-]enable-ip-alias` flag to suppress this warning.
WARNING: Starting with version 1.18, clusters will have shielded GKE nodes by default.
WARNING: Your Pod address range (`--cluster-ipv4-cidr`) can accommodate at most 1008 node(s).
WARNING: Starting with version 1.19, newly created clusters and node-pools will have COS_CONTAINERD as the default node image when no image type is specified.
Creating cluster cicd-cluster in asia-northeast1... Cluster is being health-checked (master is healthy)...done.
Created [https://container.googleapis.com/v1/projects/cyberagent-001/zones/asia-northeast1/clusters/cicd-cluster].
To inspect the contents of your cluster, go to: https://console.cloud.google.com/kubernetes/workload_/gcloud/asia-northeast1/cicd-cluster?project=cyberagent-001
kubeconfig entry generated for cicd-cluster.

NAME          LOCATION         MASTER_VERSION   MASTER_IP    MACHINE_TYPE  NODE_VERSION     NUM_NODES  STATUS
cicd-cluster  asia-northeast1  1.19.9-gke.1900  34.85.4.209  e2-medium     1.19.9-gke.1900  3          RUNNING
```

正常にKubernetesクラスタの起動ができたら、Kubernetesクラスタの情報を確認してみます。下記のようにGKEのノードが3台ともReadyになっていることを確認します。

```bash
# GKEで起動したKubernetesクラスタのノード一覧を表示
$ kubectl get nodes -o wide
NAME                                          STATUS   ROLES    AGE    VERSION            INTERNAL-IP   EXTERNAL-IP      OS-IMAGE                             KERNEL-VERSION   CONTAINER-RUNTIME
gke-cicd-cluster-default-pool-1028475f-bt3m   Ready    <none>   113s   v1.19.9-gke.1900   10.240.0.21   xx.xxx.xxx.xxx   Container-Optimized OS from Google   5.4.89+          containerd://1.4.3
gke-cicd-cluster-default-pool-3347e48f-db7f   Ready    <none>   113s   v1.19.9-gke.1900   10.240.0.22   xx.xxx.xxx.xxx   Container-Optimized OS from Google   5.4.89+          containerd://1.4.3
gke-cicd-cluster-default-pool-701af3a5-8f0v   Ready    <none>   112s   v1.19.9-gke.1900   10.240.0.19   xx.xxx.xxx.xxx   Container-Optimized OS from Google   5.4.89+          containerd://1.4.3
```

このハンズオンを終了後、クラスタを削除するには下記のコマンドを実行してください。

```bash
# GKEクラスタの削除
$ gcloud container clusters delete cicd-cluster
The following clusters will be deleted.
 - [cicd-cluster] in [asia-northeast1]

Do you want to continue (Y/n)?  y

Deleting cluster cicd-cluster...done.
Deleted [https://container.googleapis.com/v1/projects/PROJECT/zones/asia-northeast1/clusters/cicd-cluster].
```

## Argo CD CLI

Argo CD専用のCLIをインストールします。

### Windows

以下をダウンロードして、任意のディレクトリにこのexeファイルを格納して、パスを通します。

https://github.com/argoproj/argo-cd/releases/download/v2.0.5/argocd-windows-amd64.exe

パス設定参考サイト:
https://www.atmarkit.co.jp/ait/articles/1805/11/news035.html

### Mac

```bash
$ brew install argocd
```

### Linux（Cloud Shell）

Cloud Shellから接続を抜けても継続してargocdコマンドを実行できるように、argocd/binフォルダを作成して配置します。

```bash
$ mkdir -p ~/argocd/bin
$ sudo curl -sSL -o ~/argocd/bin/argocd https://github.com/argoproj/argo-cd/releases/download/v2.0.5/argocd-linux-amd64
```

```bash
# 実行権付与
$ sudo chmod +x ~/argocd/bin/argocd
```

```bash
# パス設定
$ export PATH="~/argocd/bin:${PATH}"
$ echo PATH="\"~/argocd/bin:\${PATH}\"" >> ~/.bashrc
```

これにて事前準備は完了です。お疲れさまでした。
ハンズオン当日はスライドによる説明を行った後に、[Chapter1](chapter1.md)から進めていきます！
お楽しみに！
