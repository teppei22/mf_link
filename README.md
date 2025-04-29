# マネーフォワードにCSVデータを自動で取り込むGo言語スクリプト

## 概要

このスクリプトは、**Selenium** を使用してCSVデータをマネーフォワードに自動で取り込むGo言語スクリプトです。Python版からGo言語に移植されています。

## 機能

- Chromeブラウザを起動
- マネーフォワードにログインして登録する口座のページを開く
- CSVファイルを読み込み、各取引データを入力
- 入力内容を保存（テストモードでは保存せずに閉じる）

## 必要な環境

- **Go** 1.16以上
- **Selenium WebDriver for Go** (`github.com/tebeka/selenium`)
- **ChromeDriver** (Seleniumを使ってChromeを操作するためのドライバ)
- **Google Chrome** 最新版
- **適切なCSVファイル**（Web版マネーフォワードのCSVエクスポート機能で取得してください）

## インストール

1. Go言語をインストール：[Go公式サイト](https://golang.org/dl/)からダウンロードしてインストール

2. Selenium WebDriverパッケージのインストール：

```sh
go get github.com/tebeka/selenium
```

3. ChromeDriverのインストール:
   - macOSの場合: `brew install chromedriver`
   - Linuxの場合: 各ディストリビューションのパッケージマネージャを使用
   - Windowsの場合: [ChromeDriverダウンロードページ](https://sites.google.com/a/chromium.org/chromedriver/downloads)からダウンロードし、PATH環境変数に追加

## 実行方法

1. **スクリプトを実行する前に、マネーフォワードのユーザー名、パスワード、インポート先の口座URLを編集してください。**
   - コード内の `url`、`user`、`password` 変数を編集します：
   ```go
   url := "https://moneyforward.com/accounts/show_manual/xxxxxxxxxxxxxxx" // インポート先の口座URL
   user := "<自分のアカウント>"
   password := "<自分のパスワード>"
   ```

2. CSVファイルをスクリプトと同じフォルダに配置する。
    - Web版マネーフォワードのCSVエクスポート機能で取得できるCSVと同じフォーマットです。

```
[0] "計算対象", 
[1] "日付", 
[2] "内容", 
[3] "金額（円）", 
[4] "保有金融機関", 
[5] "大項目", 
[6] "中項目", 
[7] "メモ", 
[8] "振替", 
[9] "ID"
```

3. スクリプトを実行

```sh
go run mf_import_csv.go data.csv
```

または、ビルドしてから実行：

```sh
go build mf_import_csv.go
./mf_import_csv data.csv
```

- `data.csv` はインポートするCSVファイルのパスです。

## 注意点

- `time.Sleep` の値を調整することで、環境に合わせて動作を安定させられます。
- マネーフォワードのUI変更により、要素のIDやクラスが変わる可能性があります。
- 本スクリプトの使用は自己責任でお願いします。本スクリプトを使用したことによるいかなる損害についても、作者は責任を負いません。
- 本スクリプトは個人による開発であり、マネーフォワードとは一切関係ありません。
- 自動化ツールの利用は、マネーフォワードの利用規約に違反しないように注意してください。

## Go言語版の特徴と違い

- エラーハンドリングが強化されています
- ログ出力が詳細になっています
- Goのコンカレント機能を活用することで将来的に高速化が可能です

## ライセンス

このプロジェクトはMITライセンスのもとで公開されています。

## 原作者

元のPythonスクリプトはnanosnsによって開発・管理されています。
