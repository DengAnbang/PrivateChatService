# PrivateChatService

#### 介绍

这是一款使用go编写的聊天软件的后台,使用socket实现推送的服务

#### 安装教程

1. 修改config/config中的ConfigBean,填入数据库地址账号等信息
2. 使用go build .\main.go 打包成对应平台的可执行文件
3. 使用 ./main i 安装启动
4. 不使用的时候 使用 ./main u 卸载
5. android app https://github.com/DengAnbang/PrivateChat
#### 使用说明

api 地址:https://www.showdoc.com.cn/2469196901102085/10963116568454583
可替换成api接口地址 api/showdoc_api.sh 中的api_key和api_token