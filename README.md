# 说明

## 目录

1. [服务端请求方式](./docs/cn/1.服务端请求方式.md)
2. [客户端请求方式](./docs/cn/2.客户端请求方式.md)
3. [登录相关接口](./docs/cn/3.登录相关接口.md)
4. [支付相关接口](./docs/cn/4.支付相关接口.md)
5. [附录](./docs/cn/5.附录.md)
6. [错误码](./docs/cn/6.错误码.md)

## 对接流程

因为我们的对接会先在测试环境进行，故以下的流程说明的地址均是测试环境的地址，上线前需在正式环境同样走一遍

1. 联系DragonEx工作人员，添加访问测试环境的权限

2. [注册](https://test.dragonex.co/zh-hans/account/register)两个账号：一个用于设置APP的一些信息，如支付回调地址等；一个用户利润划转（若无利润划转需求，此账号可不需要）。（建议使用公司邮箱或公司专用手机号注册，避免因公司人员变动产生不必要的麻烦。）
   
    ![](./docs/cn/images/注册.png)

3. 将注册好的账号提供给DragonEx工作人员（区分好两个账号各自的用户，App管理员账号会被冻结交易、提现等与资金相关的权限），添加APP信息

4. 登录管理员账号，进入[“个人中心”-->“开放平台”](https://test.dragonex.co/zh-hans/asset/open/app)，创建AccessKey与SecretKey、配置支付回调地址等

    ![](./docs/cn/images/配置APP信息.png)

    1. **注意**：更新**响应数据校验秘钥**后，前端展示会立即修改，但是DragonEx实际进行签名的会有一小段延时（1H内全部切换为新的秘钥），在此期间可能有部分使用新的秘钥，部分使用旧的秘钥，需接入方在更改秘钥时注意下。

5. 阅读本文档，进行开发（建议读完本README后，按上述[目录](#目录)顺序，了解我们提供的接口是否满足您的需求）

6. 测试、对账

7. 上线

## 登录授权时序图

![登录授权时序图](./docs/cn/images/DragonEx开放平台-登录授权流程图.png)

 1. 用户同意授权后，DragonEx会提供一个AccessCode，并携带此参数跳转到接入方指定的地址（获取AccessCode的流程也可以使用DragonEx统一提供的页面处理，但是这样的话用户要进行登录授权就必须从DragonEx的开放平台列表进入了）

 2. 接入方拿到AccessCode后，需要由接入方服务端带着此AccessCode，以及AccessKey及其签名校验请求DragonEx，获取AccessToken，签名方式见后续说明

 3. AccessKey与SecretKey是重要信息，请勿跟随客户端分发
   
## 开放账户支付时序图

![开放账户支付](./docs/cn/images/DragonEx开放平台-开放账户支付流程图.png)

## 币币账户支付时序图

![支付时序图](./docs/cn/images/DragonEx开放平台-通过H5支付时序图.png)

 1. 支付成功后，除了跟随HTTP请求返回支付状态外，还会有异步回调告诉接入方Server端
 2. DragonEx无法保证HTTP请求与回调到达的先后顺序，需接入方自行处理可能出现的情况

## 接口列表

1. 登录相关接口列表
   
    ![登录相关接口列表](./docs/cn/images/login_apis.svg)

2. 支付相关接口列表
   
    ![支付相关接口列表](./docs/cn/images/payment_apis.svg)

