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

以下、作業

まずは、ファイルシステムの取り出し

```
% docker pull alpine  # IMAGE ID: a24bb4013296
% docker run --name alpine alpine
% docker export alpine > alpine.tar
% docker rm alpine
% mkdir rootfs
% tar -C rootfs -xvf alpine.tar 
```

ビルドと実行はラズパイ上で行った。
ソースコードは書籍のだとビルドできなかったので、[libcontainerのREADME](https://github.com/opencontainers/runc/blob/b2bec98/libcontainer/README.md)を参考にした。 -> [main.go]

実行するとシェルが立ち上がり、そこで `/bin/hostname` を叩いて `testing` となればOKらしいが本体のホスト名が表示されてしまった。
