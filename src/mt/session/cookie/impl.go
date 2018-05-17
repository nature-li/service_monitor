package cookie

import (
	"net/http"
	"time"
	"net/url"
	"encoding/base64"
	"mt/session"
	"crypto/rand"
	"errors"
	"strconv"
)

type SessionCookie struct {
	sessionId      string
	dict           map[string]string
	lastAccessTime time.Time
	w              http.ResponseWriter
	manager        *Manager
}

var manager *Manager

type Manager struct {
	secretKey            string
	cookieSessionIdName  string
	cookieAccessTimeName string
	maxLifeTime          int64
}

func (o *SessionCookie) timeToString(t time.Time) string {
	return strconv.FormatInt(t.Unix(), 10)
}

func (o *SessionCookie) stringToTime(s string) (time.Time, error) {
	value, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return time.Now(), err
	}

	when := time.Unix(value, 0)
	return when, nil
}

func (o *SessionCookie) encodeString(raw string) (string, error) {
	result, err := encrypt([]byte(o.manager.secretKey), raw)
	if err != nil {
		return "", err
	}

	url.QueryEscape(result)
	return result, nil
}

func (o *SessionCookie) decodeString(secret string) (string, error) {
	result, err := url.QueryUnescape(secret)
	if err != nil {
		return "", err
	}

	raw, err := decrypt([]byte(o.manager.secretKey), result)
	if err != nil {
		return "", err
	}

	return string(raw), nil
}
func newSessionCookie(manager *Manager, w http.ResponseWriter) *SessionCookie {
	return &SessionCookie{dict: make(map[string]string), manager: manager, w: w}
}

func (o *SessionCookie) SessionId() string {
	return o.sessionId
}

func (o *SessionCookie) Set(key string, value string) error {
	o.updateAccessTime()

	o.dict[key] = value
	result, err := o.encodeString(value)
	if err != nil {
		return err
	}
	cookie := http.Cookie{Name: key, Value: result, Path: "/", HttpOnly: true, MaxAge: 0}
	http.SetCookie(o.w, &cookie)
	return nil
}

func (o *SessionCookie) Get(key string) string {
	o.updateAccessTime()

	if v, ok := o.dict[key]; ok {
		return v
	}
	return ""
}

func (o *SessionCookie) Del(key string) error {
	o.updateAccessTime()

	delete(o.dict, key)
	expiration := time.Now()
	cookie := http.Cookie{Name: key, Path: "/", HttpOnly: true, Expires: expiration, MaxAge: -1}
	http.SetCookie(o.w, &cookie)
	return nil
}

func (o *SessionCookie) updateAccessTime() error {
	o.lastAccessTime = time.Now()

	timeString := o.timeToString(o.lastAccessTime)
	encodedTime, err := o.encodeString(timeString)
	if err != nil {
		return err
	}
	cookie := http.Cookie{Name: o.manager.cookieAccessTimeName, Value: encodedTime, Path: "/", HttpOnly: true, MaxAge: 0}
	http.SetCookie(o.w, &cookie)

	sessionId, err := o.encodeString(o.sessionId)
	if err != nil {
		return err
	}
	cookie = http.Cookie{Name: o.manager.cookieSessionIdName, Value: sessionId, Path: "/", HttpOnly: true, MaxAge: int(o.manager.maxLifeTime)}
	http.SetCookie(o.w, &cookie)
	return nil
}

func NewManager(secretKey string, cookieSessionIdName string, cookieAccessTimeName string, maxLifeTime int64) (session.Manager, error) {
	if manager != nil {
		return manager, nil
	}

	if len(secretKey) != 32 {
		return nil, errors.New("length of secretKey must be 32")
	}

	if cookieSessionIdName == "" {
		return nil, errors.New("cookieSessionIdName can not be empty")
	}

	if cookieSessionIdName == "" {
		return nil, errors.New("cookieAccessTimeName can not be empty")
	}

	manager = &Manager{secretKey: secretKey, cookieSessionIdName: cookieSessionIdName, cookieAccessTimeName: cookieAccessTimeName, maxLifeTime: maxLifeTime}
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
	s := newSessionCookie(o, w)

	// check session by cookieSessionIdName
	c, err := r.Cookie(o.cookieSessionIdName)
	if err != nil || c.Value == "" {
		s.Set(o.cookieSessionIdName, o.randomId())
		return s
	}

	// check session by cookieAccessTimeName
	c, err = r.Cookie(o.cookieAccessTimeName)
	if err != nil || c.Value == "" {
		s.Set(o.cookieSessionIdName, o.randomId())
		return s
	}

	// get last access time
	decodedString, err := s.decodeString(c.Value)
	if err != nil || decodedString == "" {
		s.Set(o.cookieSessionIdName, o.randomId())
		return s
	}
	when, err := s.stringToTime(decodedString)
	if err != nil {
		s.Set(o.cookieSessionIdName, o.randomId())
		return s
	}

	// check if it is expired
	now := time.Now()
	expect := when.Add(time.Second * time.Duration(o.maxLifeTime))
	if expect.Before(now) {
		s.Set(o.cookieSessionIdName, o.randomId())
		return s
	}

	// parse all cookie value and store in session
	cs := r.Cookies()
	for _, item := range cs {
		if value, ok := s.decodeString(item.Value); ok == nil {
			s.dict[item.Name] = value
		}
	}

	s.updateAccessTime()
	return s
}

func (o *Manager) SessionDestroy(w http.ResponseWriter, r *http.Request) error {
	// check if session exist
	c, err := r.Cookie(o.cookieSessionIdName)
	if err != nil || c.Value == "" {
		return nil
	}

	// delete all cookies
	expiration := time.Now()
	cs := r.Cookies()
	for _, item := range cs {
		cookie := http.Cookie{Name: item.Name, Path: "/", HttpOnly: true, Expires: expiration, MaxAge: -1}
		http.SetCookie(w, &cookie)
	}

	return nil
}
