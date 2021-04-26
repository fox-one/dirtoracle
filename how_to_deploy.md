# DirtOracle 部署

DirtOracle 是基于 BLST 多签的 Price Oracle 解决方案。

- 在 DirtOracle 中，每个节点代表一个或多个价格来源。价格来源是有序的，查询价格时，将使用第一个成功价格的数据
- 价格的协商将使用群聊功能完成
- 当业务方有价格需求时，节点将发送 Price Proposal 至群聊中
- 其他节点收到 Proposal 后将与其节点价格进行比较，若价格相差小于 1% 则将进行签名并回复给发起节点
- 发起节点收到足够签名后，将生成 Price Data，并通过转账形式发送给业务方

## Docker Pull

配合 Github Action, 每次有版本更新时将打包发送至 Github GCR. 要使用 GCR 需在 Github 中生成对应的 token 进行登陆。

- 在 [GITHUB](https://github.com/settings/tokens) 创建一个有 packages/read 权限的 token 作为密码，登录 docker.pkg.github.com
- echo $TOKEN | docker login docker.pkg.github.com -u USERNAME --password-stdin
- docker pull docker.pkg.github.com/fox-one/dirtoracle/dirtoracle:0.0.22

## 配置文件

- 创建数据库
- Mixin Developer 后台创建 dapp
- 生成 blst 私钥对 sign key 与其 verify key, docker run --rm -it docker.pkg.github.com/fox-one/dirtoracle/dirtoracle:0.0.22 keys. 其中私钥便是 Sign Key, 公钥便是 verify key
- 将 dapp.client_id, verify key 发送给需要的业务方
- 编辑配置文件

**config.yaml template:**

```yaml
db:
  dialect: "mysql"
  host: "127.0.0.1"
  port: 3306
  user: "xxx"
  password: "xxxxxx"
  database: "dirtoracle"

dapp:
  pin: ""
  client_id: ""
  session_id: ""
  pin_token: ""
  private_key: ""

group:
  conversation_id: ""
  sign_key: "xxx__blst_private_key__xx"

bwatch:
  api_base: "https://f1-bwatch-api.firesbox.com"

gas:
  amount: "1"
  asset: "965e5c6e-434c-3fa9-b780-c50f43cd955c"
```

## RUN

- 每个节点需指定其价格来源，目前支持来源包括: coinbase, binance, bitfinix, bitstamp, bittrex, huobi, okex, exinswap 以及 4swap. 建议节点使用普通交易所 + 4swap/exinswap, 以正常同步 XIN, MOB 等代币价格
- 编辑 docker-compose.yaml
- 运行, docker-compose up -d

**docker-compose.yaml template:**

```yaml
version: "3.9"

services:
  node:
    image: docker.pkg.github.com/fox-one/dirtoracle/dirtoracle:0.0.22
    restart: always
    command: worker --port 7121 --config /app/config.yaml --exchanges coinbase --exchanges 4swap
    volumes:
      - ./config.yaml:/app/config.yaml
    ports:
      - "7121:7121"

volumes:
  data:
```
