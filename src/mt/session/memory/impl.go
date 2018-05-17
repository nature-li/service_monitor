package memory

import (
	"sync"
	"container/list"
	"mt/session"
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"net/url"
	"time"
	"errors"
)

var manager *Manager
var gLock sync.Mutex

type SessionMemory struct {
	sessionId  string
	dict       map[string]string
	accessTime time.Time

	manager *Manager
}

func newSession(sessionId string, manager *Manager, accessTime time.Time) *SessionMemory {
	return &SessionMemory{sessionId: sessionId, dict: make(map[string]string), accessTime: accessTime, manager: manager}
}

func (o *SessionMemory) SessionId() string {
	gLock.Lock()
	defer gLock.Unlock()

	return o.sessionId
}

func (o *SessionMemory) Set(key string, value string) error {
	gLock.Lock()
	defer gLock.Unlock()

	o.updateAccessTime()
	o.dict[key] = value
	return nil
}

func (o *SessionMemory) Get(key string) string {
	gLock.Lock()
	defer gLock.Unlock()

	o.updateAccessTime()
	if value, ok := o.dict[key]; ok {
		return value
	}

	return ""
}

func (o *SessionMemory) Del(key string) error {
	gLock.Lock()
	defer gLock.Unlock()

	o.updateAccessTime()
	delete(o.dict, key)

	return nil
}

func (o *SessionMemory) updateAccessTime() error {
	o.accessTime = time.Now()

	if e, ok := o.manager.dict[o.sessionId]; ok {
		o.manager.list.MoveToBack(e)
	}

	return nil
}

type Manager struct {
	cookieSessionIdName string
	maxLifeTime         int64
	dict                map[string]*list.Element
	list                list.List
}

func NewManager(cookieSessionIdName string, maxLifeTime int64) (session.Manager, error) {
	gLock.Lock()
	defer gLock.Unlock()

	if manager != nil {
		return manager, nil
	}

	if cookieSessionIdName == "" {
		return nil, errors.New("cookieSessionIdName can not be empty")
	}
	manager = &Manager{cookieSessionIdName: cookieSessionIdName, maxLifeTime: maxLifeTime, dict: make(map[string]*list.Element)}
	go manager.gc()
	return manager, nil
}

func (o *Manager) randomId() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return ""
	}

	return base64.URLEncoding.EncodeToString(b)
}

func (o *Manager) SessionStart(w http.ResponseWriter, r *http.Request) session.Session {
	gLock.Lock()
	gLock.Unlock()

	c, err := r.Cookie(o.cookieSessionIdName)
	if err == nil && c.Value != "" {
		sid, _ := url.QueryUnescape(c.Value)
		if e, ok := o.dict[sid]; ok {
			return e.Value.(session.Session)
		}
	}

	sid := o.randomId()
	s := newSession(sid, o, time.Now())
	e := o.list.PushBack(s)
	o.dict[sid] = e
	cookie := http.Cookie{Name: o.cookieSessionIdName, Value: sid, Path: "/", HttpOnly: true, MaxAge: int(o.maxLifeTime)}
	http.SetCookie(w, &cookie)
	return s
}

func (o *Manager) SessionDestroy(w http.ResponseWriter, r *http.Request) error {
	gLock.Lock()
	gLock.Unlock()

	c, err := r.Cookie(o.cookieSessionIdName)
	if err != nil || c.Value == "" {
		return nil
	}

	sid, _ := url.QueryUnescape(c.Value)
	if e, ok := o.dict[sid]; ok {
		delete(o.dict, sid)
		o.list.Remove(e)
		expiration := time.Now()
		cookie := http.Cookie{Name: o.cookieSessionIdName, Path: "/", HttpOnly: true, Expires: expiration, MaxAge: -1}
		http.SetCookie(w, &cookie)
	}

	return nil
}

func (o *Manager) gc() {
	gLock.Lock()
	gLock.Unlock()

	now := time.Now()
	for {
		e := o.list.Front()
		if e == nil {
			break
		}

		s := e.Value.(*SessionMemory)
		expect := s.accessTime.Add(time.Second * time.Duration(o.maxLifeTime))
		if expect.After(now) {
			break
		}

		delete(o.dict, s.sessionId)
		o.list.Remove(e)
	}

	time.AfterFunc(time.Duration(time.Minute), o.gc)
}
