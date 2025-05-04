package main

import (
	"net"
	"sync"
)

type SessionManager struct {
	sess *sync.Map // KEY: IP, VALUE: Session
}

type Session struct {
	Socket net.Conn
}

func (m *SessionManager) FindOne(key string) (Session, bool) {
	if val, ok := m.sess.Load(key); ok {
		return val.(Session), ok
	} else {
		return Session{}, false
	}
}

func (m *SessionManager) DeleteOne(key string) {
	m.sess.Delete(key)
}

func (m *SessionManager) CreateOne(key string, value Session) {
	m.sess.Store(key, value)
}
