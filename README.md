# アプリの実行方法
1. Go1.17をインストールする
1. 必要に応じて，`.env`ファイルのDBのホスト名，ユーザー名，パスワード，DB名，アプリの起動ポートを書き換える
1. `make init-db`を実行する(このコマンドでDBのテーブル，ビューを作成し，初期データを挿入する)
1. `make run`を実行する
1. `http://localhost:`[`.env`で指定したアプリの起動ポート]にアクセスするとアプリが表示される

# ディレクトリ構成
```
├── config # 環境変数を取得する
├── db     # DBとのコネクションを構成し提供する
├── domain # 各エンティティの構造体とそのメソッドを定義する
├── lib    # 各パッケージで共通して利用される関数のうち，domainに含まれないものを定義する
├── pkg    # main.goで振り分けられた各HTTPリクエストを処理して，DBを操作し，レスポンスを返す
├── sql    # テーブルやビュー，初期データを定義する
└── view   # HTMLテンプレートファイルを定義する
```

# 初期ユーザー一覧
| user_id | name    | password | role   |
|:-------:|:-------:|:--------:|:------:|
| 1       | owner1  | owner    | owner  |
| 2       | admin1  | admin    | admin  |
| 3       | member1 | member   | member |
| 4       | member2 | member   | member |
| 5       | member3 | member   | member |

