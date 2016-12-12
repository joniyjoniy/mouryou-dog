# mouryou-dog

サーバの負荷量をwebsocket経由で通知するプログラムです.
1秒ごとに測定値を送信します．

# measured value

現在試作段階のため，測定値は未定です.

# Install

Go言語をインストールしておく必要があります．

```
git clone git://github.com/joniyjoniy/mouryou-dog.git
cd mouryou-dog
make gom link
make build
```

# Run

サーバプログラムを先に実行してください.

## Server

```
bin/sample-server
```

## plugin

```
bin/mouryou-dog
```
