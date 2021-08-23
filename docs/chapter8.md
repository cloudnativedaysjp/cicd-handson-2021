# Chapter 8 CD pipeline by Argo CD

## 8-1 Argo CDのInstall

GitOpsでCDを実現するArgo CDをインストールします。

最初にArgo CD専用の `argocd` というNamespaceを作成します。

```bash
$ kubectl create namespace argocd
namespace/argocd created
```

Argo CDをインストールします。

```bash
# 現時点での最新はv2.0.5になります
$ kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/v2.0.5/manifests/install.yaml
```

以下のPodがRunningになっていることを確認します。

```bash
$ kubectl get pods,services -n argocd
NAME                                      READY   STATUS    RESTARTS   AGE
pod/argocd-application-controller-0       1/1     Running   0          3m3s
pod/argocd-dex-server-68c7bf5fdd-b9l6v    1/1     Running   0          3m3s
pod/argocd-redis-7547547c4f-q2kb5         1/1     Running   0          3m3s
pod/argocd-repo-server-58f87478b8-r52fw   1/1     Running   0          3m3s
pod/argocd-server-6f4fcdc5dc-qfnjg        1/1     Running   0          3m3s

NAME                            TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)                      AGE
service/argocd-dex-server       ClusterIP   10.102.102.212   <none>        5556/TCP,5557/TCP,5558/TCP   3m4s
service/argocd-metrics          ClusterIP   10.98.135.108    <none>        8082/TCP                     3m4s
service/argocd-redis            ClusterIP   10.101.136.116   <none>        6379/TCP                     3m4s
service/argocd-repo-server      ClusterIP   10.101.22.255    <none>        8081/TCP,8084/TCP            3m4s
service/argocd-server           ClusterIP   10.110.97.130    <none>        80/TCP,443/TCP               3m4s
service/argocd-server-metrics   ClusterIP   10.103.136.238   <none>        8083/TCP                     3m3s
```

Argo CDのWebUIにアクセスするために、プロキシ接続の設定を行います。

```bash
# kubectl port-forward は Ctrl + C でキャンセルしない限りプロンプトが戻ってきません
$ kubectl port-forward service/argocd-server 8080:443 -n argocd
```

ブラウザを起動して、`https://localhost:8080/` にアクセスします。

初回は、「この接続ではプライバシーが保護されません」と表示されますが、［詳細設定］をクリックして［EXTERNAL-IP にアクセスする（安全ではありません）］をクリックしてアクセスしてください。

次にWebUIの初期パスワードを変更します。

プロキシ接続用にターミナルを利用しているので、新規ターミナルを起動します。

最初に初期パスワードを確認します。

※Windowsの場合、以下コマンドはコマンドプロンプトではbase64が無いというエラーが出るので、Git Bashで実行します。

```bash
$ kubectl -n argocd get secret argocd-initial-admin-secret -o jsonpath="{.data.password}" | base64 --decode; echo
xxxxxxxxxxxxxxxxxxx
```

Argo CD Serverにログインします。

※Windowsの場合は、以下コマンドはGit Bashでは証明書関連のエラーがでるので、コマンドプロンプトで実行します。

```bash
$ argocd login localhost:8080 --username admin --insecure
Password: # 先程のパスワードを入力します
'admin:login' logged in successfully
Context 'localhost:8080' updated
```

「logged in successfully」と表示されれば成功です。

続いて先程の初期パスワードを変更します。  
「Password updated」と表示されれば成功です。

```bash
$ argocd account update-password --account admin
*** Enter current password: # 先程のパスワードを入力します
*** Enter new password: # 任意のパスワードを入力します
*** Confirm new password: # もう一度任意のパスワードを入力します
Password updated
Context 'localhost:8080' updated
```

### TIPS

- 設定した admin パスワードを忘れてしまったら？

  新しく設定した admin のパスワードは argocd の namespace 内にある `argocd-secret` という Secret に格納されています。  
  Secret 内の `admin.password` 及び `admin.passwordMtime` を削除し、argocd-server の Pod を再作成することで初期パスワードが再作成されます。  
  その後再度 `argocd-initial-admin-secret` から初期パスワードを抽出し、「argocd login」をし直すことができます。

`https://localhost:8080/` WebUI画面で、「Username」は「admin」、「Password」は設定した任意のパスワードを入力して、「SIGN IN」をクリックしてログインします。

![Argo CD WebUI Login](images/chapter8/chapter08-001.png)

画面左上の「+ NEW APP」ボタンをクリックします。

![Argo CD Create App1](images/chapter8/chapter08-002.png)

GENERAL内の以下の項目を入力および設定します。

| 項目              | 値                  |
| ---------------- | ------------------- |
| Application Name | cicd-confernce-2021 |
| Project          | default             |
| SYNC POLICY      | Autmatic            |

![Argo CD Create App2](images/chapter8/chapter08-003.png)

続いて SOURCE 内の以下の項目を入力および設定します。

| 項目              | 値                         |
| ---------------- | -------------------------- |
| Repository URL   | ご自身のconfigリポジトリのURL |
| Path             | manifests                  |

![Argo CD Create App3](images/chapter8/chapter08-004.png)

続いてDESTINATION内の以下の項目を入力および設定します。

| 項目           | 値                             |
| ------------- | ------------------------------ |
| Cluster URL   | https://kubernetes.default.svc |
| Namespace     | default                        |

![Argo CD WebUI Login](images/chapter8/chapter08-005.png)

上部の「Create」ボタンをクリックします。

![Argo CD WebUI Login](images/chapter8/chapter08-006.png)

configリポジトリとの連携設定は終了です。

### TIPS

- WebUIを使わずにargocdコマンドから設定することも可能です
  ```bash
  $ argocd app create cicd-confernce-2021 --repo https://github.com/YOUR_GITHUB/cicd-handson-2021-config --path manifests --dest-namespace default --dest-server https://kubernetes.default.svc --sync-policy automatic
  application 'cicd-confernce-2021' created
  ```
- 設定した内容は「argocd app list」で確認することができます
  ```bash
  $ argocd app list
  NAME                 CLUSTER                         NAMESPACE  PROJECT  STATUS  HEALTH   SYNCPOLICY  CONDITIONS  REPO                                                     PATH       TARGET
  cicd-confernce-2021  https://kubernetes.default.svc  default    default  Synced  Healthy  Auto        <none>      https://github.com/YOUR_GITHUB/cicd-handson-2021-config  manifests 
  ```


## 8-2 codeリポジトリ内の「main.yml」にconfigリポジトリへのプルリクエスト処理の追加

GitHub Actionsの「main.yml」にコンテナイメージタグの更新を契機にプルリクエストをconfigリポジトリに出す処理を追加します。

```yaml
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

        # コンテナイメージをGitHub Packagesにプッシュ
      - name: GitHub Packages login
        uses: docker/login-action@v1
        with:
          registry: docker.pkg.github.com
          username: ${{ secrets.USERNAME }}
          password: ${{ secrets.PERSONAL_ACCESS_TOKEN }}

      - name: Push image to GitHub Packages
        run: docker image push docker.pkg.github.com/${{ github.repository }}/gitops-go-app:${{ github.run_number }}

  # manifest とかはまだ仮だと思うので一旦こんな感じにしておく

  # プルリクエストを作る job を新規で定義します
  # 「needs: build」を書いておくことで、build の job が終わった後に実行されるようにします
  create-pr-k8s-manifest:
    needs: build
    runs-on: ubuntu-latest
    steps:
      # config repo を checkout します
      - name: Checkout code
        uses: actions/checkout@v2
        with:
          token: ${{ secrets.PERSONAL_ACCESS_TOKEN }}
          repository: YOUR_GITHUB/cicd-handson-2021-config

      # image tagを書き換えます
      - name: Update image tag
        run: |
          sed -i -e "s|image: docker.pkg.github.com/${{ github.repository }}/gitops-go-app:.*|image: docker.pkg.github.com/${{ github.repository }}/gitops-go-app:${{ github.run_number }}|" manifests/deployment.yaml

      # プルリクエストを作成します
      - name: Create Pull Request
        uses: peter-evans/create-pull-request@v3
        with:
          token: ${{ secrets.PERSONAL_ACCESS_TOKEN }}
          commit-message: "Update tag ${{ github.run_number }}"
          title: "Update tag ${{ github.run_number }}"
          body: "Please Merge !!"
          branch: "feature/${{ github.run_number }}"

      #  # values.yamlの更新、新規ブランチ作成、プッシュ、プルリクエスト
      #- name: Update values.yaml & Pull Request to Config Repository
      #  run: |
      #    # GitHubログイン
      #    echo -e "machine github.com\nlogin ${{ secrets.USERNAME }}\npassword ${{ secrets.GH_PASSWORD }}" > ~/.netrc
      #    # 「config」リポジトリからクローン
      #    git clone https://github.com/${{ secrets.USERNAME }}/cicd-handson-2021-config.git
      #    # GitHub Email/Username セットアップ
      #    cd cicd-handson-2021-config/gitops-helm
      #    git config --global user.email "${{ secrets.EMAIL }}"
      #    git config --global user.name "${{ secrets.USERNAME }}"
      #    # 新規ブランチ作成
      #    git branch feature/${{ github.run_number }}
      #    git checkout feature/${{ github.run_number }}
      #    # values.yamlのタグ番号を更新
      #    sed -i 's/tag: [0-9]*/tag: ${{ github.run_number }}/g' values.yaml
      #    # プッシュ処理
      #    git add values.yaml
      #   git commit -m "Update tag ${{ github.run_number }}"
      #    git push origin feature/${{ github.run_number }}
      #    # プルリクエスト処理
      #    echo ${{ secrets.PERSONAL_ACCESS_TOKEN }} > token.txt
      #    gh auth login --with-token < token.txt
      #    gh pr create  --title "Update Tag ${{ github.run_number }}" --body "Please Merge !!"
```

修正した「main.yml」をcodeリポジトリに「git push」して、CIが通って、configリポジトリにプルリクエストがあることを確認します。

```git
$ git add .github/workflows/main.yml
$ git commit -m "Pull Request to config repogitry add main.yml"
$ git push origin main
```

configリポジトリのプルリクエストの内容を確認して、マージします。

Argo CD の WebUI を確認して、しばらくするとステータスが「Out Of Sync」になり、その後自動的にプルリクエストの内容が反映されることを確認します。  
Argo CD は定期的に対象リポジトリと Sync しますが、タイミングによっては少し時間がかかる場合があります。  
その場合は「REFRESH」を押してステータスを手動で確認させることが可能です。
