package main

type Retriver interface {
	Get(url string) string
}

type Poster interface {
	Post(url string, form map[string]string) string
}

// 组合接口
type RetrieverPoster interface {
	Retriver
	Poster
	// 当然这里还可以添加别的方法
}

// 假设有个session函数，它需要一个参数，既是一个Retriver，又是一个Poster，此时就可以用到接口的组合，可以同时调用组合接口中的所有方法
func session(s RetrieverPoster) string {
	s.Post("http://www.baidu.com", map[string]string{"key": "value"})
	return s.Get("http://www.baidu.com")
}
