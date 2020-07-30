###  接口：

1.登录 (POST)： /login 

传入格式（form）  

```
{
“name”:string,
"password":string
}
```

2.注册(POST)：/register

传入格式与登录一样

3.创建红包抽奖(POST)：/red/creat

传入格式（json）

```json
{
    "id":int,//红包的id
    "total":float64,//总金额
    "number":int//个数
}
```

4.抽取红包(POST)：/red/consume/:id

需携带cooike获取名字

5.创建弹幕抽奖(POST)：/creatlottery

传入格式（json）

```json
{
    "num":int,//个数
    "minute":int,//持续分钟，不能少于三分钟
    "content":string//匹配的内容
}
```

6.添加黑名单(POST)：/addblack

传入格式（form):

```
{
"name":string
}
```

7.弹幕（websocket）：/ws

携带cookie获取名字

发送信息和接收的信息都为json反序列化后的

```
{"name":"as","color:{"R":0,"G":0,"B":0,"A":1},
"time":1596010334,"content":"no way"}
color -> rgbint8(255,255,255)
time  -> 时间戳
content -> 内容
```

8.返回格式（对应的都是code第一个数字）

```
0 -> 错误  
1 -> 正确  
2 -> 关于红包  cod为200 为到  201 为未抢到  
3 -> 登录或注册错误 
5 -> 数据库错误
```

[**postman地址**](https://documenter.getpostman.com/view/9500172/T1Dtdarf?version=latest#be7f638a-0b0e-46f9-844c-a8bf84384289)

dockerhub地址

部署地址：101.201.140.26 