// Session类
// 		import (
//			"yuensoft.com/net/session"
//		)
//
// Created by Yuen Cheng on 13-11-4.
// Copyright (c) 2013年 YuenSoft.com. All rights reserved.
package session

import (
	"time"
)

type ISession interface {
	Set(key, value interface{}) error
	Get(key interface{}) interface{}
	Delete(key interface{}) error
	SessionID() string
}

type Session struct {
	sID          string
	timeAccessed time.Time
	values       map[interface{}]interface{}
}

func (session *Session) Set(key, value interface{}) error {
	if session.values == nil {
		session.values = map[interface{}]interface{}{}
	}
	session.values[key] = value

	//Update LifeTime
	session.updateSessionLifeTime()
	return nil
}
func (session *Session) Get(key interface{}) interface{} {
	//Update LifeTime
	session.updateSessionLifeTime()
	v, ok := session.values[key]
	if !ok {
		return nil
	}
	return v
}
func (session *Session) Delete(key interface{}) error {
	//Update LifeTime
	session.updateSessionLifeTime()
	delete(session.values, key)

	return nil
}
func (session *Session) SessionID() string {
	return session.sID
}
func (session *Session) updateSessionLifeTime() error {
	session.timeAccessed = time.Now()
	return nil
}
