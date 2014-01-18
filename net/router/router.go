// 支持RESTful的路由器
//
// Created by Yuen Cheng on 13-11-4.
// Copyright (c) 2013年 YuenSoft.com. All rights reserved.

package router

import (
	//"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

//type SingleRouteInterface interface {
//}

type SingleRoute struct {
	pattern string
	handler http.Handler
}

type Router struct {
	routers     []SingleRoute
	staticPaths map[string]string
}

func New() *Router {
	router := &Router{}
	router.routers = []SingleRoute{}
	router.staticPaths = make(map[string]string)
	return router
}

func (this *Router) Add(pattern string, handler http.Handler) {
	singleRoute := SingleRoute{}
	singleRoute.pattern = pattern
	singleRoute.handler = handler
	this.routers = append(this.routers, singleRoute)
}

func (this *Router) AddStatic(pattern string, realPath string) {
	this.staticPaths[pattern] = realPath
}

func (this *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 请求静态文件？
	// 将匹配到的模式（URL文件名外的前半截）用模式实际对应的真实路径替换，即为，文件真实路径。
	for urlPattern, realPath := range this.staticPaths {
		if isStaticRequest := strings.HasPrefix(r.URL.Path, urlPattern); isStaticRequest {
			part4RealPath := strings.Replace(r.URL.Path, urlPattern, realPath, 1)
			staticFileRealPath := filepath.Dir(os.Args[0]) + part4RealPath
			http.ServeFile(w, r, staticFileRealPath)
			return
		}
	}

	// 请求路径
	for _, router := range this.routers {
		reg, err := regexp.Compile(`^` + router.pattern + `$`)
		if err != nil {
			panic("通过地址注册的Handler，在匹配时，创建Regexp出错。")
		}
		if matched := reg.MatchString(filepath.Dir(r.URL.Path)); matched {
			router.handler.ServeHTTP(w, r)
			return
		}
	}

	http.NotFound(w, r)
}
