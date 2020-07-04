# 第17章 Go言語とコンテナ

* コンテナを実現するためのOSカーネルの機能
    * コントロールグループ(cgroups)
        * CPU、メモリ、ブロックデバイス、ネットワーク、デバイスファイル(`/dev/`)の使用量とアクセスを制限する
    * 名前空間(Namespaces)
        * 名前空間の分離
            * プロセスID
            * ネットワーク(インタフェース、ルーティングテーブル、ソケットなど)
            * マウント
            * UTS(ホスト名)
            * IPC (プロセス間通信)
            * ユーザ(UID, GID)
* libcontainerを使えばとりあえずできるっぽい
    * [README](https://github.com/opencontainers/runc/tree/master/libcontainer)にやり方が書いてあるらしいが、本の写経で進める

以下、Mac上での作業メモ

| environment | version |
|--- | --- |
| macOS Catalina | 10.15.5（19F101） |
| Docker for mac| 2.3.0.3(45519) |


ファイルシステムの取り出し

```
% docker pull alpine  # IMAGE ID: a24bb4013296
% docker run --name alpine alpine
% docker export alpine > alpine.tar
% docker rm alpine
% mkdir rootfs
% tar -C rootfs -xvf alpine.tar 
```

必要なパッケージのインストール

```
% go get github.com/opencontainers/runc/libcontainer
% go get golang.org/x/sys/unix
```

ソースコードを写経すると、
```
libcontainer.New()
```
のところで、 `Unresolved reference 'New'` となる。
[参照元を見る](https://github.com/opencontainers/runc/blob/2c632d1a2de0192c3f18a2542ccb6f30a8719b1f/libcontainer/init_linux.go)と、 `go get` する時にLinuxでしかインストールされないっぽい。

Dockerコンテナがない