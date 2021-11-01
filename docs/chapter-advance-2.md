# Chapter Advance 2 Conftest & Open Policy Agent (OPA)

ここでは、Conftestを利用したポリシーチェックについて学びます。  
Conftestではポリシーの定義をOpen Policy Agent(OPA)でも使われているRegoという言語を用いて行います。

# Open Policy Agent（OPA）

OPAは、オープンソースの汎用的なポリシーエンジンです。ポリシーに違反した情報を発見し、事前に定義されたアクションを実行する仕組みで、Regoという言語を使用してポリシーを定義します。
主な特徴は以下となります。

* 軽量で汎用性のあるOSSのポリシーエンジン
* Kubernetes専用というわけではなく、YAML、JSONなど構造化データのポリシーエンジン
* KubernetesではCI時に、Conftestとの組み合わせで導入するケースが多い
* CNCFのGraduatedプロジェクト

APIに送信するQueryと、Policyを参照して、評価したDecision（結果）を返す仕組みです。

![Policy Decoupling](images/chapter-advance/chapter-advance-003.png)

[OPA公式ドキュメント](https://www.openpolicyagent.org/docs/latest/)

# Conftest & Rego

Conftestは、バイナリファイルを所定のディレクトリに格納して、パスを通すことで簡単に利用できます。
Rego言語で定義したポリシーファイルと実際にチェックするファイルを所定のディレクトリに格納して、「$ conftest test <チェックするファイルまたはディレクトリのパス>」という形式でコマンドを実行して、定義したポリシーに違反していなければOK、NG場合はFAILを返します。

以下の例では、マニフェストファイルに「runAsNonRoot」が設定されていればOK、設定されていなければNGを返すポリシーをRegoで定義、実行しているものになります。

![Conftest & Rego](images/chapter-advance/chapter-advance-004.png)

 [Conftest GitHub](https://github.com/open-policy-agent/conftest/)

本ハンズオンでは、マニフェスト内のイメージタグに「latest」がある場合をNGとする定義をRegoで作成します。
以下は、マニフェスト内のイメージタグに「latest」が定義されていたら、NGとして「Cannot use latest tag !!」を返すというポリシーです。

```go
package main

#イメージタグにlatestがある場合
deny[msg] {
  input.kind == "Deployment"
  image := input.spec.template.spec.containers[_].image
  endswith(image, "latest")
  msg := sprintf("latest tag: %s", [ image ])
}

#イメージタグに指定がない場合
deny[msg] {
  input.kind == "Deployment"
  image := input.spec.template.spec.containers[_].image
  not contains(image, ":")
  msg := sprintf("latest tag: %s", [ image ])
}
```

このRegoで定義したファイルを「policy」というディレクトリを作成して、格納します。

```bash
$ mkdir ./policy
$ vi ./policy/latest-tag-check.rego
```

Regoも使いこなすにはそれなりの学習コストが必要となります。本ハンズオンでは初歩的な定義ですが、詳細は以下公式ドキュメントをご参照ください。

 [Rego Official Documents](https://www.openpolicyagent.org/docs/latest/policy-language/)

 これまでの変更をConfigリポジトリにプッシュします。

```git
$ git add .
$ git commit -m "Conftest and Rego"
$ git push origin main
```

プッシュ後に、CodeリポジトリからConfigリポジトリにプルリクエストが発行されたことをトリガーに、ポリシーチェックのCIが実行されて、OKとなります。マージするとこれまでと同じ処理の流れとなります。

# Conftest NG Attempt

実際に、ConftestでNGを確認するには、以下の流れとなります。

1. `cicd-handson-2021-config`ディレクトリに移動
2. `goapp.yaml`ファイルのイメージタグを`latest`に変更して保存
3. Configリポジトリにプルリクエスト
4. GitHub Actions CI処理の自動稼働
5. NGの確認

`cicd-handson-2021-config`ディレクトリに移動して、`goapp.yaml`ファイルのイメージタグを`latest`に変更して保存します。

```bash
$ cd cicd-handson-2021-config
$ vi ./manifests/goapp.yaml
```

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: goapp-deployment
spec:
  replicas: 3
  selector:
    matchLabels:
      app: goapp
  template:
    metadata:
      labels:
        app: goapp
    spec:
      containers:
      - name: goapp
        image: docker.pkg.github.com/<GITHUB_USER>/cicd-handson-2021-code/go-image:latest #変更箇所
        ports:
        - containerPort: 9090
      imagePullSecrets:
      - name: dockerconfigjson-github-com
```

Configリポジトリにプッシュします。

```git
$ git branch feature/latest
$ git checkout feature/latest
$ git add manifests
$ git commit -m "Update tag latest"
$ git push origin feature/latest
```

GitHubのConfigリポジトリの画面からプルリクエストを作成します。
`Compare & pull request`ボタンをクリックします。

![Compare & pull request](images/chapter-advance/chapter-advance-005.png)

マージ先を `<GitHub Account/cicd-handson-2021-config>` リポジトリを指定します。

![base:main](images/chapter-advance/chapter-advance-006.png)

`latest tag`と入力して、`Create pull request`ボタンをクリックします。

![Create pull request](images/chapter-advance/chapter-advance-007.png)

プルリクエストが発行されたことをトリガーに、ポリシーチェックのCIが実行されて、NGとなります。  
上部メニュー「Actions」を選択して、CIでエラーになっていることを確認します。  
内容を確認する場合は、「Feature/latest」の部分をクリックします。

![Actions NG1](images/chapter-advance/chapter-advance-008.png)

「validation」をクリックします。

![Actions NG2](images/chapter-advance/chapter-advance-009.png)

「validate manifest-dir with custom policies」のプルダウンを開きます。
「return code: 1」になっていることを確認します。

![Actions NG3](images/chapter-advance/chapter-advance-010.png)

上部メニュー「Pull requests」をクリックして、対象のプルリクエストを選択します。

![Actions NG3](images/chapter-advance/chapter-advance-011.png)

`[conftest] Validation Failed` の「detail」をクリックして、failureの内容を確認します。

![Actions NG3](images/chapter-advance/chapter-advance-012.png)

Codeリポジトリ側のmain.ymlにおいては、イメージタグは自動で番号が付与されるようになっているので、latestが付与されることはありません。  
本ハンズオンでは、確認のためにマニフェストファイルにlatestタグを付けて、直接Configリポジトリにプッシュ、プルリクエストというやり方を行いました。

以上で完了となります。  
使用したリソースは忘れずに削除をしておいてください。  

---
[Chapter 10 Clean up](chapter10.md)へ
