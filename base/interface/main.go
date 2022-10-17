package main

import (
	"fmt"
	"practise/base/interface/mock"
	"practise/base/interface/real"
	"time"
)

// Retriever 定义get方法
type Retriever interface {
	Get(url string) string
}

// Poster 定义Poster 方法
type Poster interface {
	Post(url string,
		form map[string]string) string
}

// k 接口组合
type k interface {
	Retriever
	Poster
}

func session(s k) string {
	s.Post("http://www.baidu.com", map[string]string{
		"contents": "another faked imooc.com",
	})
	return s.Get("http://www.baidu.com")
}

func main() {
	var r Retriever
	// mock
	mockRetriever := mock.Retriever{
		Contents: "傻逼",
	}
	r = &mockRetriever
	a := r.Get("http://www.baidu.com")
	fmt.Println("======第一======", a)
	// real
	realRetriever := real.Retriever{
		UserAgent: "Mozilla/5.0",
		TimeOut:   time.Minute,
	}
	r = &realRetriever
	b := r.Get("http://www.baidu.com")
	fmt.Println("======第二======", b)

	//aa := []int{1, 2, 3}
	//bb := []int{4, 5, 6}
	//fmt.Println(append(aa, bb...))
	fmt.Println(session(&mockRetriever))
}
