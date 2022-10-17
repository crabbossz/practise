package main

import (
	"crypto/tls"
	"fmt"
	"reflect"
	"time"
)

//

// Server 配置问题
type Server struct {
	Addr     string        // 必填
	Port     int           // 必填
	Protocol string        // 非必填
	Timeout  time.Duration // 非必填
	MaxConns int           // 非必填
	TLS      *tls.Config   // 非必填
}

// 解决方式一 分离可选项
type Config struct {
	Protocol string
	Timeout  time.Duration
	MaxConns int
	TLS      *tls.Config
}

type Server1 struct {
	Addr string
	Port int
	Conf *Config
}

func NewServer(addr string, port int, conf *Config) (*Server1, error) {
	//...
	a := Server1{
		addr,
		port,
		conf,
	}
	return &a, nil
}

// 解决方法二 Builder模式
// 使用一个builder类来做包装

type ServerBuilder struct {
	Server
}

func (sb *ServerBuilder) Create(addr string, port int) *ServerBuilder {
	sb.Server.Addr = addr
	sb.Server.Port = port
	//其它代码设置其它成员的默认值
	return sb
}
func (sb *ServerBuilder) WithProtocol(protocol string) *ServerBuilder {
	sb.Server.Protocol = protocol
	return sb
}
func (sb *ServerBuilder) WithMaxConn(maxconn int) *ServerBuilder {
	sb.Server.MaxConns = maxconn
	return sb
}
func (sb *ServerBuilder) WithTimeOut(timeout time.Duration) *ServerBuilder {
	sb.Server.Timeout = timeout
	return sb
}
func (sb *ServerBuilder) WithTLS(tls *tls.Config) *ServerBuilder {
	sb.Server.TLS = tls
	return sb
}
func (sb *ServerBuilder) Build() Server {
	return sb.Server
}

// 方法三 使用闭包

type Option func(*Server)

func WithProtocol(proto string) Option {
	return func(s *Server) {
		s.Protocol = proto
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.Timeout = timeout
	}
}

func WithMaxConns(maxConn int) Option {
	return func(s *Server) {
		s.MaxConns = maxConn
	}
}

func NewServer3(addr string, port int, opts ...Option) *Server {
	//创建Server对象，并填写可选项的默认值
	s := &Server{
		Addr:     "127.0.0.1",
		Port:     8080,
		Protocol: "udp",
		Timeout:  time.Second * 10,
		MaxConns: 10,
	}

	//都选项列表中每项都应用
	for _, option := range opts {
		option(s)
	}

	return s
}

// 方法四 基于接口实现

//Options 接口，定义需要实现的apply方法
type Options interface {
	apply(server *Server)
}

//ProtoOption 实现Option接口
type ProtoOption string

func (p ProtoOption) apply(s *Server) {
	s.Protocol = string(p)
}

// 提供WithProtocol函数将string类型的proto转换为ProtoOption类型，方便以Option接口的形式传递给函数使用

func WithProtocol4(proto string) Options {
	return ProtoOption(proto)
}

//TimeoutOption 实现Option接口
type TimeoutOption time.Duration

func (t TimeoutOption) apply(s *Server) {
	s.Timeout = time.Duration(t)
}

func WithTimeout4(timeout time.Duration) Options {
	//fmt.Println("======1234",reflect.TypeOf(TimeoutOption(timeout)))
	return TimeoutOption(timeout)
}

//MaxConnOption 实现Option接口
type MaxConnOption int

func (m MaxConnOption) apply(s *Server) {
	s.MaxConns = int(m)
}

func WithMaxConn4(maxConn int) Options {
	return MaxConnOption(maxConn)
}

func NewServer4(addr string, port int, opts ...Options) *Server {
	//创建Server，并填写可选项的默认值
	s := &Server{
		Addr:     "127.0.0.1",
		Port:     8080,
		Protocol: "udp",
		Timeout:  time.Second * 10,
		MaxConns: 10,
	}

	for _, opt := range opts {
		opt.apply(s)
	}

	return s
}

func main() {
	// 方法一 分离参数
	//srv1, _ := NewServer("localhost", 9000, nil)               //传入nil，表示不传可选项采纳数
	//conf := Config{Protocol: "tcp", Timeout: 60 * time.Second} //可选项参数，Protocol和Timeout
	//srv2, _ := NewServer("locahost", 9000, &conf)
	//fmt.Println(srv1, srv2)

	// 方法二 build模式
	// 使用一个builder类来包装Server，然后分两步操作：
	//
	//1、Create方法根据必填项addr和port参数来创建ServerBuilder实例
	//
	//2、WithProtocol、WithMaxConn、WithTimeOut、WithTLS方法将函数参数传递给builder类中的Server对象
	//
	//3、提供Build()方法将builder类中的Server对象暴露出来
	//
	//于是就可以以如下的方式来使用了
	//sb := ServerBuilder{}
	//server1 := sb.Create("127.0.0.1", 8080).
	//	WithProtocol("udp").
	//	WithMaxConn(1024).
	//	WithTimeOut(30 * time.Second).
	//	Build()
	//fmt.Println(server1)

	// 方法三
	//s := NewServer3("xxxx", 1234, WithProtocol("tcp"), WithTimeout(time.Minute), WithMaxConns(100))
	//fmt.Printf("server:%v\n", s)
	//
	//// 方法四
	//s4 := NewServer4("xxxx", 1234, WithProtocol4("tcp1"), WithTimeout4(time.Minute), WithMaxConn4(100))
	//fmt.Println(s4)

	hh := WithProtocol4("tcp1")
	fmt.Println(hh, reflect.TypeOf(hh))
}
