// 管理Session的类
// 		import (
//			"yuensoft.com/net/session"
//		)
//
//		var sessionManager *session.Manager
//		func init() {
//			sessionManager = session.NewManager("GOSESSIONID", 10)
//			sessionManager.SessionTransmitType(session.SessionTransmitTypeURL)
//
//			//GC
//			go sessionManager.SessionGC(sessionManager.MaxLifeTime)
//		}
//
// Created by Yuen Cheng on 13-11-4.
// Copyright (c) 2013年 YuenSoft.com. All rights reserved.
package session

import (
	"fmt"
	"github.com/nu7hatch/uuid"
	"net/http"
	"net/url"
	"sync"
	"time"
)

// 用来指定SessionID是通过何种方式传递的（URL GET/Cookie)
const (
	SessionTransmitTypeURL = iota
	SessionTransmitTypeCookie
)

type Manager struct {
	CookieName          string
	MaxLifeTime         int64
	lock                sync.Mutex
	sessionProvider     map[string]*Session
	sessionTransmitType int //指定SessionID是通过何种方式传递的（URL GET/Cookie)
}

// 创建新的Session管理器
func NewManager(cookieName string, maxlifetime int64) *Manager {
	return &Manager{CookieName: cookieName, sessionProvider: make(map[string]*Session), MaxLifeTime: maxlifetime}
}

// Manager's Interface
type IManager interface {
	SessionTransmitType(t int)
	SessionStart(w http.ResponseWriter, r *http.Request) ISession
	SessionInit(sid string) (ISession, error)
	SessionRead(sid string) (ISession, error)
	SessionDestroy(w http.ResponseWriter, r *http.Request)
	SessionGC(maxLifeTime int64)
}

// (私有方法)获取UUID字符串
func (manager *Manager) sessionID() string {
	u5, _ := uuid.NewV4()
	return u5.String()
}

func (manager *Manager) SessionTransmitType(t int) {
	manager.sessionTransmitType = t
}

// 为每个来访的用户分配或获取与它相关的Session。
// 检测是否已经有某个Session与当前来访用户发生了关联，如果没有，创建之。
func (manager *Manager) SessionStart(w http.ResponseWriter, r *http.Request) (*Session, error) {
	manager.lock.Lock()
	defer manager.lock.Unlock()

	switch {
	case manager.sessionTransmitType == SessionTransmitTypeURL:
		return manager.sessionStartWithURL(w, r), nil
	case manager.sessionTransmitType == SessionTransmitTypeCookie:
		return manager.sessionStartWithCookie(w, r), nil
	}

	return nil, fmt.Errorf("Session Start Error")
}

// GET方式传递SessionID
func (manager *Manager) sessionStartWithURL(w http.ResponseWriter, r *http.Request) *Session {
	var session *Session
	if sid := r.FormValue(manager.CookieName); sid != "" {
		sid, _ = url.QueryUnescape(sid)
		session, _ = manager.SessionRead(sid)
		//无效SID
		if session.SessionID() == "" {
			sid := manager.sessionID()
			session = manager.SessionInit(sid)
		}
	} else {
		sid := manager.sessionID()
		session = manager.SessionInit(sid)
	}
	return session
}

// Cookie方式传递SessionID
func (manager *Manager) sessionStartWithCookie(w http.ResponseWriter, r *http.Request) *Session {
	var session *Session
	sessionCookie, err := r.Cookie(manager.CookieName)

	if err != nil || sessionCookie.Value == "" {
		sid := manager.sessionID()
		session = manager.SessionInit(sid)
		cookie := http.Cookie{Name: manager.CookieName, Value: url.QueryEscape(sid), Path: "/", HttpOnly: true, MaxAge: int(manager.MaxLifeTime)}
		http.SetCookie(w, &cookie)
	} else {
		sid, _ := url.QueryUnescape(sessionCookie.Value)
		session, _ = manager.SessionRead(sid)
		//无效SID
		if session.SessionID() == "" {
			sid = manager.sessionID()
			session = manager.SessionInit(sid)
		}
	}
	return session
}

// 初始化Session
//
// 没有Session记录的情况下的初始化。
// 其会添加至Manager的Map
func (manager *Manager) SessionInit(sid string) *Session {
	session := &Session{sID: sid, timeAccessed: time.Now()}
	manager.sessionProvider[sid] = session
	return session
}

// 根据SessionID，读取Session对象
//
// 根据SessionID，读取Session对象
func (manager *Manager) SessionRead(sid string) (*Session, error) {
	session, ok := manager.sessionProvider[sid]
	if ok {
		return session, nil
	}

	return session, fmt.Errorf("Get Session:%s Error", sid)
}

// 销毁Session对象
func (manager *Manager) SessionDestroy(w http.ResponseWriter, r *http.Request) {
	switch {
	case manager.sessionTransmitType == SessionTransmitTypeURL:
		manager.sessionDestroyWithURL(w, r)
	case manager.sessionTransmitType == SessionTransmitTypeCookie:
		manager.sessionDestroyWithCookie(w, r)
	}
}

// GET方式传递SessionID的销毁
func (manager *Manager) sessionDestroyWithURL(w http.ResponseWriter, r *http.Request) error {
	sid := r.FormValue(manager.CookieName)
	if sid == "" {
		return nil
	}

	manager.lock.Lock()
	defer manager.lock.Unlock()
	session, err := manager.SessionRead(sid)
	if err != nil {
		return fmt.Errorf("SessionRead Wrong in sessionDestoryWithURL func")
	}
	delete(manager.sessionProvider, session.SessionID())
	return nil
}

// Cookie方式传递SessionID的销毁
func (manager *Manager) sessionDestroyWithCookie(w http.ResponseWriter, r *http.Request) error {
	var session *Session
	sessionCookie, err := r.Cookie(manager.CookieName)
	if err != nil || sessionCookie.Value == "" {
		return nil
	} else {
		sid, _ := url.QueryUnescape(sessionCookie.Value)
		session, _ = manager.SessionRead(sid)

		if err != nil {
			return fmt.Errorf("SessionRead Wrong in sessionDestroyWithCookie func")
		}
		delete(manager.sessionProvider, session.SessionID())
	}

	return nil
}

// 根据Session的有效期的时间，来定期回收Session
func (manager *Manager) SessionGC(maxLifeTime int64) {
	fmt.Println("Start GC...")
	manager.lock.Lock()
	defer manager.lock.Unlock()

	for sid, session := range manager.sessionProvider {
		if int64(time.Since(session.timeAccessed)/time.Second) > manager.MaxLifeTime {
			fmt.Printf("SID:%s, will GC.\n", session.SessionID())
			delete(manager.sessionProvider, sid)
		}
	}

	//到时间后就执行，再其自身的Routine中
	time.AfterFunc(time.Duration(manager.MaxLifeTime)*time.Second, func() {
		manager.SessionGC(manager.MaxLifeTime)
	})
}
