# mouryou-dog

サーバの負荷量をwebsocket経由で通知するプログラムです.
1秒ごとに測定値を送信します．

# measured value

現在試作段階のため，測定値は未定です.

# sample-server

送信先のサーバには[gorilla/websocketのechoサンプル](https://github.com/gorilla/websocket/blob/master/examples/echo/server.go)を利用しています．

# Run

サーバプログラムを先に実行してください.

## Server

```
$ go run sample-server.go
```

## plugin

```
$ go run main.go
```
