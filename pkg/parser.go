package pkg

import "sync"

type Parser interface {
	Next() string
	Send(value string) bool
	SetPlain(bool2 bool)
	GetPlain() bool
	Len() int
	Close()
	Running() bool
	Start()
}

type MyParser struct {
	CharChan  chan string
	IsPlain   chan bool
	Closed    bool
	State     bool
	mu        sync.Mutex
	sendMu    sync.Mutex
	plainMu   sync.Mutex
	done      chan struct{}
	initOnce  sync.Once
	closeOnce sync.Once
}

func (m *MyParser) Done() <-chan struct{} {
	m.initOnce.Do(func() { m.done = make(chan struct{}) })
	return m.done
}

func (m *MyParser) Next() string { return <-m.CharChan }

func (m *MyParser) Len() int { return len(m.IsPlain) }

func (m *MyParser) Send(value string) bool {
	m.sendMu.Lock()
	defer m.sendMu.Unlock()
	done := m.Done()
	select {
	case <-done:
		return false
	default:
	}
	select {
	case <-done:
		return false
	case m.CharChan <- value:
		return true
	}
}

func (m *MyParser) GetPlain() bool { return <-m.IsPlain }

func (m *MyParser) SetPlain(plain bool) {
	m.plainMu.Lock()
	defer m.plainMu.Unlock()
	done := m.Done()
	select {
	case <-done:
		return
	default:
	}
	select {
	case <-done:
	case m.IsPlain <- plain:
	}
}

func (m *MyParser) Close() {
	m.Done()
	m.closeOnce.Do(func() {
		m.mu.Lock()
		m.Closed = true
		m.mu.Unlock()
		close(m.done)
		m.sendMu.Lock()
		defer m.sendMu.Unlock()
		m.plainMu.Lock()
		defer m.plainMu.Unlock()
		if m.CharChan != nil {
			close(m.CharChan)
		}
		if m.IsPlain != nil {
			close(m.IsPlain)
		}
	})
}

func (m *MyParser) Running() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.State
}

func (m *MyParser) Start() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.State = true
}
