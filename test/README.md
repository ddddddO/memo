### JMeter
- テスト計画
    - スレッドグループ
        - サンプラー(の`HTTP Request`)
        - リスナー(の`Summary Report` でテスト結果出力)

### GUI
- .jmxファイルは一旦GUIから作成する(`C:\Program Files\apache-jmeter-5.2\bin\jmeter.bat`を管理者で実行)
    - 認証画面へのリクエスト
    - 認証画面から名前・パスワードをPOSTリクエスト
### コマンド
`/mnt/c/Program\ Files/apache-jmeter-5.2/bin/jmeter.sh -t tag-mng.jmx` で、
```
================================================================================
Don't use GUI mode for load testing !, only for Test creation and Test debugging.
For load testing, use CLI Mode (was NON GUI):
   jmeter -n -t [jmx file] -l [results file] -e -o [Path to web report folder]
& increase Java Heap to meet your test requirements:
   Modify current env variable HEAP="-Xms1g -Xmx1g -XX:MaxMetaspaceSize=256m" in the jmeter batch file
Check : https://jmeter.apache.org/usermanual/best-practices.html
================================================================================
An error occurred:
No X11 DISPLAY variable was set, but this program performed an operation which requires it.
```

`-n` で、`NON GUI` ぽい。以下で実行すること。
`/mnt/c/Program\ Files/apache-jmeter-5.2/bin/jmeter.sh -n -t tag-mng.jmx`
