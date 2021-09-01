# Chapter 9 CI/CD pipeline

本チャプターまでたどり着いた皆様、おめでとうございます。  
ここまでの手順が達成できていれば、CI/CDは完成していると言えるでしょう。  
しかしここで満足してはいけません。本ドキュメントは基本的なCI/CDを紹介したに過ぎません。  
世の中には様々なシステムがあり、そのシステムに合ったCI/CDもそれぞれ存在します。  
是非さまざまなあなた自身のCI/CDを組んでみて下さい。  
そして上手くいった時は是非カンファレンスやコミュニティイベントで発表してください。

最後のステップとして、本チャプターではアプリケーションを更新し、CIからCDまできちんと処理されることを再度確認します。

## 9-1 アプリケーションの更新

アプリケーションを以下の手順で更新します。  
cssを編集して「Simple CICD Landscape!」の文字の色を変更します。

```bash
# git clone した code リポジトリで作業します
$ cd ./cicd-handson-2021-code

# vi 等で color の部分を #ff33cc に変更します
$ vi apps/web/static/style.css
```

```diff
diff --git a/apps/web/static/style.css b/apps/web/static/style.css
index 273df93..faa11d3 100644
--- a/apps/web/static/style.css
+++ b/apps/web/static/style.css
@@ -1,7 +1,7 @@
 h1{
     text-align: center;
     font-size: 40px;
-    color: #111111;
+    color: #ff33cc;
 }
```

```git
# リポジトリへプッシュします
$ git add apps/web/static/style.css
$ git commit -m "change title color"
$ git push origin main
```

GitHub ActionsのCIの実行、GitHub Packagesへのコンテナイメージ保存、configリポジトリへのプルリクエストまでできることを確認します。

その後、configリポジトリのプルリクエストをマージします。

Argo CDによって同期されるのを待ちます。(もしくは手動でSyncさせます)  
Syncが終わったらアプリケーションを表示させます。  
「kubectl port-forward」した後に `http://localhost:9090/` にアクセスします。

```bash
$ kubectl port-forward deployment.apps/goapp-deployment 9090:9090
```

タイトルの色が変わっていれば成功です。
（更新されない場合は、ブラウザのハードリロードも試してみてください。）

![New Application Title](images/chapter9/chapter09-001.png)

## 9-2 アプリケーション変更のプルリクエスト時にテストを流す

9-1 の手順ではアプリケーションの更新後、プルリクエストを作らずにそのままmainへプッシュしました。  
本来であればこのような手順は好ましくありません。テストが通らなくなってしまう変更は避けるべきです。  
そこでcodeリポジトリのGitHub Actionsを変更し、プルリクエスト時にテストが流れるように変更します。

```bash
# 新規ファイルで pr-check.yml を用意します
$ vi .github/workflows/pr-check.yml
```

下記の内容を pr-check.yml に書き込みます。

```yaml
name: Pull Request Check

# mainブランチへのプルリクエスト時にトリガー
on:
  pull_request:
    branches: [ main ]

jobs:
  check:
    name: Application test
    runs-on: ubuntu-latest
    steps:
      - name: Checkout code
        uses: actions/checkout@v2

        # アプリケーションテスト
      - name: Application test
        run: |
          cd apps
          make run-test
```

pr-check.yml をリポジトリにプッシュします。

```bash
$ git add .github/workflows/pr-check.yml
$ git commit -m "add pr-check action"
$ git push origin main
```

## TIPS

- pr-check.yml を追加する時のプッシュによってGitHub Actionsが動くことになると思います。  
  本来であればアプリケーションの変更に関係ないファイルの更新時にはアプリケーションのビルドはしなくていい場合が多いと思います。  
  その場合はGitHub Actionsの設定で、特定のパスのファイルが更新されたらアプリケーションをビルドするGitHub Actionsを動かす、もしくはその逆で特定のパスのファイルの更新を無視する、といった柔軟な設定をすることができます。  
  上手く使いこなして下さい。([公式Document](https://docs.github.com/en/actions/reference/workflow-syntax-for-github-actions#onpushpull_requestpaths))

アプリケーションを更新するプルリクエストを作成します。  
今回はテストに失敗するような変更にするため、chap 6-2で直したコードを元に戻すプルリクエストを作成します。

```bash
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
100                                                                 Crunchbase:  proj.Twitter, #ここを Twitterに戻す
101                                                                 StarCount:   getStarCount(proj.RepoUrl),
102                                                         }
103                                                 }(i, proj)

# プルリクエスト用のbranchを作成し、プッシュします
$ git branch pr-check
$ git checkout pr-check
$ git add apps/server/landscape.go
$ git commit -m "it will be failed..."
$ git push origin pr-check
```

プルリクエストを作成します。  
codeリポジトリをブラウザで開き、「branches」をクリックします。

![PR Create1](images/chapter9/chapter09-002.png)

`pr-check` ブランチの「New pull request」をクリックします。

![PR Create2](images/chapter9/chapter09-003.png)

「base」の部分を**ご自身のリポジトリの main**に変更します。  
その後「Create pull request」をクリックします。

![PR Create3](images/chapter9/chapter09-004.png)

プルリクエストを作成後、GitHub Actionsが実行されます。  
プルリクエストの画面を見ると失敗していることが確認できると思います。

![PR Create4](images/chapter9/chapter09-005.png)

このようにプルリクエスト時にチェックをすることは、アプリケーションの品質向上に役立つことになります。  
ぜひやり方を覚えて、ご自身のシステム開発に役立てて下さい。

### TIPS

- コマンドラインからプルリクエストを作成する

  gh コマンドが使えるのであればコマンドラインからプルリクエストを作成することができます。  

   ```bash
   $ gh pr create --title "pr-check test" --body "pr-check test" --base main
   ```

ハンズオンは以上になります。お疲れさまでした！  
使用したリソースは忘れずに削除をしておいてください。  

---
[Chapter 10 Clean up](chapter10.md)へ
