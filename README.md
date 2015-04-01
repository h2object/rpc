httpRPC
---

RPC for normal HTTP request with a special response Analyser interface.
鉴于各种基于HTTP的Restful服务, 抽象出常用请求接口, 同时提供HTTP响应分析接口, 适配各种不同的Restful应用的响应。

### 请求接口

当前实现 http 常规请求接口,如下:

*	PostBinary()
*	PutBinary()

*	PostJson()
*	PutJson()
*	PatchJson()

*	PostForm()
*	PutForm()
*	PatchForm()

*	Get()
*	Delete()


### 响应分析

响应分析接口定义如下:

````go

type Analyzer interface {
	//! ret, 响应分析后填充的对象
	//! resp, http response 指针对象
	//! 返回: 失败信息
	Analyse(ret interface{}, resp *http.Response) error
}

````

#### h2o 服务响应分析

系统默认实现了 [h2o](https://github.com/h2object/h2o) 服务的响应分析

#### api.dnspod.com API服务响应分析

请参考[GoDNSPOD](https://github.com/h2object/GoDNSPOD)
