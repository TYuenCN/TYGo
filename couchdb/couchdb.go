// CouchDB DAO Level, Wrapper
//
// Created by Yuen Cheng on 13-11-4.
// Copyright (c) 2013年 YuenSoft.com. All rights reserved.
//
// //定义对应于CouchDB文档对象的Struct注意。如：
//
// 		type DocUser struct {
// 			ID   string `json:"_id,omitempty"`  //,omitempty为空时，不写入JSON，对于CouchDB很重要
// 			Rev  string `json:"_rev,omitempty"` //,omitempty为空时，不写入JSON，对于CouchDB很重要
// 			Name string `json:"name,omitempty"` //,omitempty为空时，不写入JSON，对于CouchDB很重要
// 			Age  int8   `json:"age,omitempty"`  //,omitempty为空时，不写入JSON，对于CouchDB很重要
// 		}
//
// 		func (docUser *DocUser) GetDBName() string {
// 			return "hhcehua_users"
// 		}
//
// 		func (docUser *DocUser) GetID() string {
// 			return docUser.ID
// 		}
//
// 		func (docUser *DocUser) GetRev() string {
// 			return docUser.Rev
// 		}
//
// 		func (docUser *DocUser) GetJSONBytes() []byte {
// 			bytes, err := json.Marshal(docUser)
// 			if err != nil {
// 				handleError(err, "GetJSONBytes")
// 				return nil
// 			} else {
// 				return bytes
// 			}
// 			return nil
// 		}
package couchdb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	//"strconv"
)

type CouchDBDAO struct {
	COUCH_DB_HOST string
}

// 所有对应CouchDB中“文档”的结构要实现的接口
type IDoc interface {
	GetDBName() string
	GetID() string
	GetRev() string
	GetJSONBytes() []byte
}

// Couch 返回的 UUID
type UUIDs struct {
	UUIDs []string `json:"uuids"`
}

// Update后的单条，影响，结果状态行
//		结果可用做的操作，包含如下：Insert, Update, Delete
type EffectRowResult struct {
	OK  bool   `json:"ok"`
	ID  string `json:"id"`
	Rev string `json:"rev"`
}

// 查询文档的返回结果，的，基础struct
//		结果可用做的操作，包含如下：QueryView
//
//		视图的查询结果模版为：{"total_rows":13,"offset":3,"rows":[{"id":"","key":"","value":{}]}
//		rows := &couchdb.ResultRows{}
//		rows.Rows = []map[string]interface{}{} //默认初始化，为nil，不为零值
//		json.Unmarshal(bytes, rows)
type ResultRows struct {
	TotalRows int32                    `json:"total_rows"`
	Offset    int32                    `json:"offset"`
	Rows      []map[string]interface{} `json:"rows"`
}

// 创建CouchDB DAO对象
//
// 		//host: `http://user:password@192.168.0.101:5984/`
//
// 		dao := couchdb.NewCouchDBDAO(`http://user:password@192.168.0.101:5984/`)
func NewCouchDBDAO(host string) *CouchDBDAO {
	dao := &CouchDBDAO{COUCH_DB_HOST: host}
	return dao
}

// 处理Error
//
// err: Error对象
//
// panicMessage: 抛出异常时的消息内容
func handleError(err error, panicMessage string) {
	if err != nil {
		fmt.Println(err)
		if panicMessage != "" {
			panic(panicMessage)
		}
	}
}

// 让CouchDB返回UUID
// count: 请求多少条UUID
func (dao *CouchDBDAO) GetUUIDFromCouchDB(count uint8) []string {
	queryUUIDsURLString := dao.COUCH_DB_HOST + "_uuids" + "?count=" + fmt.Sprintf("%d", count)
	queryUUIDsURL, _ := url.Parse(queryUUIDsURLString)
	r := http.Request{
		Method: `GET`,
		URL:    queryUUIDsURL,
		Header: map[string][]string{},
		Close:  true,
	}
	if password, ok := r.URL.User.Password(); ok {
		r.SetBasicAuth(r.URL.User.Username(), password)
	}
	c := http.Client{}
	resp, err := c.Do(&r)
	if err != nil {
		handleError(err, "从CouchDB请求创建UUID,获取响应部分错误。")
	}
	bytes, _ := ioutil.ReadAll(resp.Body)
	uuids := &UUIDs{}
	err = json.Unmarshal(bytes, uuids)
	if err != nil {
		handleError(err, "从CouchDB请求创建UUID,解码JSON部分错误。")
	} else {
		return uuids.UUIDs
	}
	return nil
}

// 插入单条Document，返回CouchDB结果JSON的[]byte
//
// 文档的struct，不用设定ID字段，本方法会填充一个新的
func (dao *CouchDBDAO) InsertDoc(doc IDoc) *EffectRowResult {
	docBytes := doc.GetJSONBytes()
	uuid := dao.GetUUIDFromCouchDB(1)[0]
	dbWithNewIdURLString := dao.COUCH_DB_HOST + doc.GetDBName() + "/" + uuid
	url, _ := url.Parse(dbWithNewIdURLString)
	r := http.Request{
		Method:        `PUT`,
		URL:           url,
		Header:        map[string][]string{},
		ContentLength: int64(len(docBytes)),
		Close:         true,
	}
	r.Body = ioutil.NopCloser(bytes.NewBuffer(docBytes))
	if password, ok := r.URL.User.Password(); ok {
		r.SetBasicAuth(r.URL.User.Username(), password)
	}
	c := http.Client{}
	resp, err := c.Do(&r)
	if err != nil {
		handleError(err, "从CouchDB请求创建新Doc,获取响应部分错误。")
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		handleError(err, "从CouchDB请求创建新Doc,解码JSON部分错误。")
	} else {
		effect := &EffectRowResult{}
		err := json.Unmarshal(bytes, effect)
		if err != nil {
			handleError(err, "从CouchDB请求插入Doc,编码JSON bytes为EffectRowResult部分错误。")
			return nil
		} else {
			return effect
		}
	}

	return nil
}

// 根据文档struct的ID字段，查询，并返回文档的数据
func (dao *CouchDBDAO) QueryDoc(doc IDoc) []byte {
	dbURLString := dao.COUCH_DB_HOST + doc.GetDBName() + "/"
	docURLString := dbURLString
	if id := doc.GetID(); id != "" {
		docURLString += id
	}
	//
	docURL, _ := url.Parse(docURLString)
	r := http.Request{
		Method: `GET`,
		URL:    docURL,
		Header: map[string][]string{},
		Close:  true,
	}
	if password, ok := r.URL.User.Password(); ok {
		r.SetBasicAuth(r.URL.User.Username(), password)
	}
	c := http.Client{}
	resp, err := c.Do(&r)
	if err != nil {
		handleError(err, "从CouchDB请求查询Doc,获取响应部分错误。")
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		handleError(err, "从CouchDB请求查询Doc,解码JSON部分错误。")
	} else {
		return bytes
	}

	return nil
}

// 查询某数据库的DesignDocument(DD)
//
// 		dbName：数据库名
//		ddName： Is the Design Document ID's value,that without "_design/"
//		viewName：Is the view name that without "_view/"
//		queryString：查询参数；(key、startkey、endkey .et)
//
//		视图的查询结果模版为：{"total_rows":13,"offset":3,"rows":[{"id":"","key":"","value":{}]}
//		rows := &couchdb.ResultRows{}
//		rows.Rows = []map[string]interface{}{} //默认初始化，为nil，不为零值
//		json.Unmarshal(bytes, rows)
func (dao *CouchDBDAO) QueryView(dbName string, ddName string, viewName string, queryString string) []byte {
	dbURLString := dao.COUCH_DB_HOST + dbName + "/"
	ddURLString := dbURLString + "_design/" + ddName + "/"
	vwURLString := ddURLString + "_view/" + viewName
	vwQueryURLString := vwURLString + "?"
	if queryString != "" {
		vwQueryURLString += queryString
	}

	//
	vwQueryURL, _ := url.Parse(vwQueryURLString)
	r := http.Request{
		Method: `GET`,
		URL:    vwQueryURL,
		Header: map[string][]string{},
		Close:  true,
	}
	if password, ok := r.URL.User.Password(); ok {
		r.SetBasicAuth(r.URL.User.Username(), password)
	}
	c := http.Client{}
	resp, err := c.Do(&r)
	if err != nil {
		handleError(err, "从CouchDB请求DB的View,获取响应部分错误。")
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		handleError(err, "从CouchDB请求DB的View,解码JSON部分错误。")
	} else {
		fmt.Println(string(bytes))
		return bytes
	}

	return nil
}

// 更新文档，_rev必须存在;
//
//		bytes = dao.QueryView("dbname", "users", "keyIsName", `key="tsengyuen"`)
//		rows := &couchdb.ResultRows{}
//		rows.Rows = []map[string]interface{}{}
//		json.Unmarshal(bytes, rows)
func (dao *CouchDBDAO) UpdateDoc(doc IDoc) *EffectRowResult {
	docBytes := doc.GetJSONBytes()
	dbURLString := dao.COUCH_DB_HOST + doc.GetDBName() + "/"
	docURLString := dbURLString
	if id := doc.GetID(); id != "" {
		docURLString += id
	}
	//
	jsonBytes, err := json.Marshal(docBytes)
	if err != nil {
		fmt.Println(err)
		panic("更新文档时，编码结构体为JSON时，出错。")
	}
	//
	docURL, _ := url.Parse(docURLString)
	r := http.Request{
		Method:        `PUT`,
		URL:           docURL,
		Header:        map[string][]string{},
		ContentLength: int64(len(jsonBytes)),
		Close:         true,
	}
	r.Body = ioutil.NopCloser(bytes.NewBuffer(jsonBytes))
	if password, ok := r.URL.User.Password(); ok {
		r.SetBasicAuth(r.URL.User.Username(), password)
	}

	c := http.Client{}
	resp, err := c.Do(&r)
	if err != nil {
		handleError(err, "从CouchDB请求更新Doc,获取响应部分错误。")
		return nil
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		handleError(err, "从CouchDB请求更新Doc,解码JSON部分错误。")
		return nil
	} else {
		effect := &EffectRowResult{}
		err := json.Unmarshal(bytes, effect)
		if err != nil {
			handleError(err, "从CouchDB请求更新Doc,编码JSON bytes为EffectRowResult部分错误。")
			return nil
		} else {
			return effect
		}

		return nil
	}
}

// 删除文档，_rev必须存在;
//
//		docUser := &DocUser{
//			ID:  "123123123",
//			Rev: "1-fe27a81c070c197c1382e82cbf9f44d1",
//		}
//		effectRowResult := dao.DeleteDoc(docUser)
func (dao *CouchDBDAO) DeleteDoc(doc IDoc) *EffectRowResult {
	dbURLString := dao.COUCH_DB_HOST + doc.GetDBName() + "/"
	docURLString := dbURLString
	if id := doc.GetID(); id != "" {
		docURLString += id
	}
	docWithQueryStringRevURLString := docURLString
	if rev := doc.GetRev(); rev != "" {
		docWithQueryStringRevURLString += "?rev=" + rev
	}
	//
	docWithQueryStringRevURL, _ := url.Parse(docWithQueryStringRevURLString)
	r := http.Request{
		Method: `DELETE`,
		URL:    docWithQueryStringRevURL,
		Header: map[string][]string{},
		Close:  true,
	}
	if password, ok := r.URL.User.Password(); ok {
		r.SetBasicAuth(r.URL.User.Username(), password)
	}

	c := http.Client{}
	resp, err := c.Do(&r)
	if err != nil {
		handleError(err, "从CouchDB请求更新Doc,获取响应部分错误。")
		return nil
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		handleError(err, "从CouchDB请求更新Doc,解码JSON部分错误。")
		return nil
	} else {
		effect := &EffectRowResult{}
		err := json.Unmarshal(bytes, effect)
		if err != nil {
			handleError(err, "从CouchDB请求删除Doc,编码JSON bytes为EffectRowResult部分错误。")
			return nil
		} else {
			return effect
		}

		return nil
	}
}
