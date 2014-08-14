package main

var counter int = 0

type Session struct {
	id          int
	name        string
	cmdType     string
	env         []string
	workingPath string
}

func NewSession(name string, cmdType string, workingPath string, env []string) *Session {
	counter++
	return &Session{id: counter, name: name, cmdType: cmdType, workingPath: workingPath, env: env}
}
