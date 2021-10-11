# Chapter Advance 3 CD pipeline by Flux

[Chapter 8](chapter8.md)ではArgo CDを使ったハンズオンを行いました。  
GitOpsを実現する方法はArgo CD以外にもあり、[Flux](https://github.com/fluxcd/flux2)も有力な選択肢の1つです。  
この章ではFluxを使ったハンズオンを行います。

Fluxは現在はv2が主流となっており、v1は[メンテナンスモード](https://github.com/fluxcd/flux/issues/3320)になっています。本章ではv2を使ったハンズオンを行います。

## 環境の初期化

[Chapter 8](chapter8.md)を行った方はArgo CD関連のリソースを削除してください。  
また、デプロイされている`goapp-deployment`も削除しておいてください。

```bash
# namespaceごとArgo CDを削除します
$ kubectl delete namespace argo

# goapp-deploymentの削除
$ kubectl delete deployment goapp-deployment
```

削除されているか確認します。

```bash
$ kubectl get all -n argo
No resources found in argo namespace.

$ kubectl get deploy
No resources found in default namespace.
```

defaultのネームスペースにある`dockerconfigjson-github-com`のSecretは削除せずに残しておいて下さい。  
もし削除してしまった、もしくはクラスターを再作成した場合は[Chapter 3](chapter3.md)の手順3-3-4より作成しておいてください。

## Flux cliのインストール

`flux`コマンドのインストールを行います。

```bash
# brew の場合
$ brew install fluxcd/tap/flux

# linux もしくは bash 環境
$ curl -s https://fluxcd.io/install.sh | sudo bash
```

インストールできたら実行できることを確認します。

```bash
$ flux version --client
flux: v0.18.1
```

また、checkを実行して問題ないことを確認します。

```bash
$ flux check --pre
► checking prerequisites
✔ Kubernetes 1.22.2 >=1.19.0-0
✔ prerequisites checks passed
```

## Fluxの初期設定

Fluxの初期設定(bootstrap)を行います。  
Fluxはこの初期設定でFluxのコンポーネントのマニフェストを精製し、GitHubへプッシュします。  
その後`flux-system`のnamespaceにコンポーネントをインストールします。  
今回は作業用に`flux`ブランチを使ってハンズオンを行っていきます。

```bash
# mainブランチからfluxブランチを作成します

# configディレクトリに移動します
$ cd cicd-handson-2021-config

# mainブランチにいるか確認します
$ git branch
* main

# fluxブランチを作成してプッシュします
$ git branch flux
$ git push origin flux
...
 * [new branch]      flux -> flux
```

現在配置されている`manifests/goapp.yaml`を一部修正します。

```bash
# configディレクトリで作業します
$ cd cicd-handson-2021-config

# fluxブランチに移動します
$ git checkout flux
Switched to branch 'flux'

# マニフェストを編集し、namespace情報を付与します
$ vi manifests/goapp.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: goapp-deployment
  namespace: default # <- 追加
spec:
...

# git diffで確認した後にコミットしてプッシュします
$ git diff
diff --git a/manifests/goapp.yaml b/manifests/goapp.yaml
index cd766f8..9b63c16 100644
--- a/manifests/goapp.yaml
+++ b/manifests/goapp.yaml
@@ -2,6 +2,7 @@ apiVersion: apps/v1
 kind: Deployment
 metadata:
   name: goapp-deployment
+  namespace: default
 spec:
   selector:

$ git add manifests/goapp.yaml
$ git commit -m "add namespace"
$ git push origin flux
```

`flux`ブランチの準備ができたのでbootstrap処理をしていきます。  

```bash
# ご自身のGitHubユーザー名とPersonal Access Tokenは適宜変更して下さい
$ flux bootstrap git \
  --url=https://github.com/<ご自身のGitHubユーザ名>/cicd-handson-2021-config \
  --branch=flux \
  --username=<ご自身のGitHubユーザ名> \
  --password=<ご自身のPersonal Access Token> \
  --token-auth=true \
  --path=manifests/
► cloning branch "flux" from Git repository "https://github.com/<ご自身のGitHubユーザー名>/cicd-handson-2021-config"
✔ cloned repository
► generating component manifests
✔ generated component manifests
...
...
✔ all components are healthy
```

「all components are healthy」が出てくれば成功です。

`flux`ブランチにfluxのコンポーネントのマニフェストが入っていることを確認しておきます。

```bash
# configディレクトリで作業します
$ cd cicd-handson-2021-config

# fluxブランチに移動します
$ git checkout flux

# fluxブランチをプルします
$ git pull origin flux

# manifestsディレクトリ内にflux-systemディレクトリが作成されていることを確認します
$ ls manifests/
flux-system     goapp.yaml
$ ls manifests/flux-system/
gotk-components.yaml	gotk-sync.yaml		kustomization.yaml
```

### TIPS
- 「✗ bootstrap failed with 1 health check failure(s)」と出てしまった場合は`goapp.yaml`内のmetadataにnamespaceの情報が書かれているか確認してください。

- なおこのbootstrap処理はGitHubにコンポーネントのマニフェストがプッシュされますが、マニフェストをプッシュせずにインストールのみを行いたい場合は「flux install」にてコンポーネントをインストールすることだけを実行することができます。


既に`goapp.yaml`が配置されているので、`goapp-deployment`がデプロイされていることを確認します。

```bash
$ kubectl get deployment
NAME               READY   UP-TO-DATE   AVAILABLE   AGE
goapp-deployment   1/1     1            1           10m
```

Fluxでは設定情報はカスタムリソースにて管理されます。  
bootstrapした情報がカスタムリソースとして登録されていることを確認します。

```bash
# GitHubの情報は「gitrepositories」に登録されます
$ kubectl get gitrepositories -n flux-system
NAME          URL                        READY   STATUS                                                            AGE
flux-system   <bootstrapしたrepository>   True    Fetched revision: flux/b3c3d640d4803ce5ac62c0a66dd7bab37dcbdfeb   16m
```

## リソースを変更し、fluxで自動的に更新されることを確認する

`goapp-deployment`のレプリカの数を変更し、Fluxによって自動的に反映されることを確認します。

```bash
# configディレクトリで作業します
$ cd cicd-handson-2021-config

# fluxブランチに移動します
$ git checkout flux

# fluxブランチを最新にします
$ git pull origin flux

# goapp.yamlを編集します
# "spec.replicas"を追加します
# 既に replicas がある方は数を変更してみてください
$ vi manifests/goapp.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: goapp-deployment
  namespace: default
spec:
  replicas: 3 # <-追加
  selector:
    matchLabels:
      app: goapp
...

# diffを確認した後にコミットしてプッシュします
$ git diff
diff --git a/manifests/goapp.yaml b/manifests/goapp.yaml
index b7d1bdf..08b909f 100644
--- a/manifests/goapp.yaml
+++ b/manifests/goapp.yaml
@@ -4,7 +4,7 @@ metadata:
   name: goapp-deployment
   namespace: default
 spec:
+  replicas: 3
   selector:
     matchLabels:
       app: goapp
...

$ git add manifests/goapp.yaml
$ git commit -m "add replicas"
$ git push origin flux
```

プッシュ後しばらくするとPodの数が変更されることを確認します。

```bash
$ kubectl get pod
NAME                                READY   STATUS    RESTARTS   AGE
goapp-deployment-76f55f57cb-5xds8   1/1     Running   0          50s
goapp-deployment-76f55f57cb-dd7qp   1/1     Running   0          100m
goapp-deployment-76f55f57cb-fwj5x   1/1     Running   0          50s
```

Fluxのコントローラーのログは「flux logs」で確認することができます。  
なにか挙動がおかしくなった時にはログを見ると良いでしょう。

```bash
# -f を付けることで表示し続けることが可能です
$ flux logs -f
2021-10-11T09:18:32.878Z info GitRepository/flux-system.flux-system - Reconciliation finished in 459.4463ms, next run in 1m0s
2021-10-11T09:19:33.257Z info GitRepository/flux-system.flux-system - Reconciliation finished in 445.6865ms, next run in 1m0s
...
```

Fluxはログにも出ている通りデフォルトでは1分毎に状態を確認します。  
手動で状態を反映させることも可能で、その場合は「flux reconcile」を実行することで手動で反映させることが可能です。

```bash
$ flux reconcile source git flux-system
► annotating GitRepository flux-system in flux-system namespace
✔ GitRepository annotated
◎ waiting for GitRepository reconciliation
✔ fetched revision flux/9bc6d15e54c1db8c66ed5df2fb2e325050511a7b
```

### TIPS

- 自動更新の間隔を設定するには？

  デフォルトでは更新間隔は1分ですが、自由に変更することが可能です。「flux bootstrap」や「flux create source」等を実行する際にオプションで「--interval」で指定することができます。  
  既存の設定を変更する場合は`gitrepository`リソース内の`spec.interval`の値を修正します。


## イメージが更新されたら自動でデプロイされるようにする

Fluxではイメージのタグが更新されたら自動でデプロイする仕組みがあります。(ImageOpsと呼ばれることもあります)  
Argo CDでも同様のことは[Argo CD Image Updater](https://github.com/argoproj-labs/argocd-image-updater)で実現することが可能ですが、まだ開発中のステータスとなっています。

イメージの更新をFluxで検知するためにはいくつかのFluxのコンポーネントが追加で必要になります。  
先程のbootstrap処理に「--components-extra」を追加してインストールします。

```bash
# flux bootstrapは冪等処理なので何度実行しても問題ありません
$ flux bootstrap git \
  --url=https://github.com/<ご自身のGitHubユーザ名>/cicd-handson-2021-config \
  --branch=flux \
  --username=<ご自身のGitHubユーザ名> \
  --password=<ご自身のPersonal Access Token> \
  --token-auth=true \
  --path=manifests/ \
  --components-extra=image-reflector-controller,image-automation-controller

# 新しくimage-automation-controllerおよびimage-reflector-controllerが起動されているのを確認します
$ kubectl get pod -n flux-system
NAME                                           READY   STATUS    RESTARTS   AGE
helm-controller-586f555bd9-88g42               1/1     Running   0          113m
image-automation-controller-79b8654998-779j6   1/1     Running   0          117s
image-reflector-controller-f9f78dd76-2jj7w     1/1     Running   0          117s
kustomize-controller-577c79958f-mr6q9          1/1     Running   0          113m
notification-controller-6f88f4dfc6-4h86m       1/1     Running   0          113m
source-controller-6bd9bd84db-xgg48             1/1     Running   0          113m
```

イメージの自動更新の定義を行います。  
イメージの自動更新の定義には下記の3つの定義が必要になります。
- image repository
- image policy
- image update

まず`flux-system`のnamespaceに`dockerconfigjson-github-com`のSecretをコピーします。(fluxのコントローラーが参照するため)

```bash
$ kubectl get secret dockerconfigjson-github-com -o yaml | sed 's/namespace: .*/namespace: flux-system/' | kubectl apply -f -
secret/dockerconfigjson-github-com created 
```

image repositoryの定義を行います。

```bash
$ flux create image repository go-image \
    --secret-ref dockerconfigjson-github-com \
    --image docker.pkg.github.com/<ご自身のGitHubユーザー名>/cicd-handson-2021-code/go-image
✚ generating ImageRepository
► applying ImageRepository
✔ ImageRepository created
◎ waiting for ImageRepository reconciliation
✔ ImageRepository reconciliation completed

# 正常に登録されているか確認します
# - READYが"True"になっていること
# - MESSAGEに"successful scan, found N tags"と表示されていること (例として8にしています)
$ flux get image repository
NAME    	READY	MESSAGE                      	LAST SCAN                	SUSPENDED
go-image	True 	successful scan, found 8 tags	2021-10-11T20:23:44+09:00	False
```

image policyの定義を行います。

```bash
# 今回のgo-imageのtagは数字なので「--select-numeric=asc」を指定します
# 「--filter-regex」を追加して数字のtagだけを参照するようにします
$ flux create image policy go-image --image-ref=go-image --select-numeric=asc --filter-regex='[0-9]+'
✚ generating ImagePolicy
► applying ImagePolicy
✔ ImageRepository created
◎ waiting for ImagePolicy reconciliation
✔ ImagePolicy reconciliation completed

# 最新のイメージtagが表示されているか確認します(例として最新tagを8にしています)
# 一部出力を省略しています
$ flux get image policy
NAME            READY   MESSAGE                                                 LATEST IMAGE
go-image        True    Latest image tag for '.../go-image' resolved to: 8      .../go-image:8
```

image updateの定義を行いますが、その前に`goapp.yaml`を編集します。

```bash
# configディレクトリで作業します
$ cd cicd-handson-2021-config

# fluxブランチに移動します
$ git checkout flux

# fluxブランチを最新にします
$ git pull origin flux

# goapp.yamlの「image:」に特定の文字列を付与します
# これはFluxが操作すべきimageの場所を特定するためです
# go-imageのtagの部分は最新よりも1つ前にしてください。
$ vi manifests/goapp.yaml
...
    spec:
      containers:
      - name: goapp
        image: docker.pkg.github.com/<ご自身のGitHubユーザー名>/cicd-handson-2021-code/go-image:<最新より1つ前のtag> # {"$imagepolicy": "flux-system:go-image"}
        ports:
...

# diffを確認した後にコミットしてプッシュします
$ git diff
diff --git a/manifests/goapp.yaml b/manifests/goapp.yaml
index b7d1bdf..255cac7 100644
--- a/manifests/goapp.yaml
+++ b/manifests/goapp.yaml
@@ -15,7 +15,7 @@ spec:
     spec:
       containers:
       - name: goapp
-        image: docker.pkg.github.com/<ご自身のGitHubユーザー名>/cicd-handson-2021-code/go-image:<最新tag>
+        image: docker.pkg.github.com/<ご自身のGitHubユーザー名>/cicd-handson-2021-code/go-image:<最新より1つ前のtag> # {"$imagepolicy": "flux-system:podinfo"}
         ports:
         - containerPort: 9090
       imagePullSecrets:

$ git add manifests/goapp.yaml
$ git commit -m "add image marker"
$ git push origin flux

# 最新の1つ前のイメージtagで動いているのを確認します
$ kubectl get deploy goapp-deployment -o jsonpath='{.spec.template.spec.containers[].image}{"\n"}'
docker.pkg.github.com/<ご自身のGitHubユーザー名>/cicd-handson-2021-code/go-image:<最新より1つ前のtag>
```

image updateの定義を行います。

```bash
$ flux create image update flux-system \
--git-repo-ref=flux-system \
--git-repo-path="manifests/" \
--checkout-branch=flux \
--push-branch=flux \
--author-name=fluxcdbot \
--author-email=fluxcdbot@users.noreply.github.com \
--commit-template="{{range .Updated.Images}}{{println .}}{{end}}"
✚ generating ImageUpdateAutomation
► applying ImageUpdateAutomation
✔ ImageRepository created
◎ waiting for ImageUpdateAutomation reconciliation
✔ ImageUpdateAutomation reconciliation completed
```

しばらくすると`goapp.yaml`の内容が最新のイメージtagに置き換えられます。  
その後その内容がクラスターへと反映されます。  

`goapp.yaml`の更新はFluxが行っています。gitのログを見てみましょう。

```bash
# configディレクトリで作業します
$ cd cicd-handson-2021-config

# fluxブランチに移動します
$ git checkout flux

# fluxブランチを最新にします
$ git pull origin flux

# gitのlogを確認します
$ git log --name-status HEAD^..HEAD
commit b9d17a4399d1e7fe4715fecff5851c1c7ced8178 (HEAD -> flux, origin/flux)
Author: fluxcdbot <fluxcdbot@users.noreply.github.com>
Date:   Mon Oct 11 12:10:34 2021 +0000

    docker.pkg.github.com/<ご自身のGitHubユーザー名>/cicd-handson-2021-code/go-image:<最新tag>

M       manifests/goapp.yaml
```

コミットメッセージやAuthorの部分は自由にカスタマイズすることが可能です。「flux create image update」時にオプションで設定してください。

codeリポジトリに対して何らかの更新を行い、新しいイメージを作成してみてください。  
ここまでの手順が上手く行っていれば、イメージを更新するだけでクラスターに自動的に反映される様子が確認できると思います。

## Fluxのアンインストール

ハンズオン終了後、Fluxの環境をクラスターから削除する場合は下記を実行します。

```bash
$ flux uninstall --namespace=flux-system
Are you sure you want to delete Flux and its custom resource definitions? [y/N] y
► deleting components in flux-system namespace
...
...
✔ Namespace/flux-system deleted
✔ uninstall finished
```

必要に応じてGitHub上の`flux`ブランチを削除します