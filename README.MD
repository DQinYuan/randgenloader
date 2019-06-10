# randgenloader

## Quick Start


提供Randgen sql生成器的http接口。

```shell
docker run -d --name randgen ddfddf/randgen-server
```

查看启动容器的ip地址：

```shell
docker inspect --format '{{ .NetworkSettings.IPAddress }}' randgen
```

然后即可通过该ip地址进行访问。

## 接口描述

### 返回值约定

返回值为一个json，格式如下：

```json
{
    "Resp": "...",   //返回的内容，可以是任意数据类型
    "ErrorInfo": "Err",   //错误信息，字符串类型，若没有错误，则该字段为空字符串
}
```

客户端拿到该返回值后应该先检查"ErrorInfo"字段是否为空字符串，若为空字符串，
再去"Resp"字段读取具体的值。


### 接口调用顺序

提供四个接口：

 - `/init`
 - `/loaddata`
 - `/query`
 - `/compare`
 
请务必按上述顺序调用，不然的话接口会返回相应的错误信息

### 接口参数与返回值

所有接口都可以任意通过GET或者POST的方式调用

`/init`: 初始化loader

入参：

|    参数     |  类型   | 是否必须 | 默认值  | 含义|
| ---------- | -------- | -------- | ----- | ---- |
| testname | string | 是 |  | 本次测试的名称 |

出参：

正常情况下返回字符串OK

```json
{
    "Resp": "OK",   
    "ErrorInfo": "",  
}
```


`/loaddata`: 获得用于生成数据的sql

入参：

|    参数     |  类型   | 是否必须 | 默认值  | 含义 |
| ---------- | --------| -------- | ----- | ---- |
| zz     | string | 是 |  | zz文件的内容 |
| yy     | string | 是 | | yy文件的内容 |
| db     | string | 否 | test | 数据库名称 |
| queries| int    | 否|  1000    | 要生成的查询数量 |


出参：

一个数组，数组中每个元素为一条sql

```json
{
    "Resp": ["CREATE ...", "INSERT ...", ...],   
    "ErrorInfo": "",  
}
```


`/query`: 获得测试的sql

入参：无

出参：同`/loaddata`


`/compare`：比较两个结果之间的区别

入参：

|    参数     |  类型   | 是否必须 | 默认值  | 含义 |
| ---------- | --------| -------- | ----- | ----|
| mysql     | string | 是 |  | mysql的输出结果 |
| tidb     | string | 是 |   | tidb的输出结果 |


出参：

返回一个含有两个字段的对象，"Comment"字段表示比较器对这此次对比的评论，
"Consistent"表示此次比较的定性结果，false表示比较结果不一致，true表示比较结果一致

```json
{
    "Resp":
        {
            "Comment":"--- Mysql\n+++ Tidb\n@@ -1 +1 @@\n-aaaaa\n+bbbbmj\n",
            "Consistent":false,
        },
    "ErrorInfo":"",
}
```