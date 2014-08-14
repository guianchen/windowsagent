package main

import (
	"bytes"
	"encoding/json"
	"strconv"
	//"fmt"
	"github.com/codegangsta/martini"
	"log"
	"net/http"
	"os/exec"
	//"strings"
	//"github.com/codegangsta/martini-contrib/auth"
)

var m *martini.Martini
var sessions = make(map[int]*Session)

func init() {
	m = martini.New()
	// Setup middleware
	m.Use(martini.Recovery())
	m.Use(martini.Logger())
	m.Use(func(c martini.Context, w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
	})
	//m.Use(auth.Basic(AuthToken, ""))
	//m.Use(MapEncoder)

	// Setup routes
	r := martini.NewRouter()
	r.Get("/execute/:sid", execute)
	r.Get("/session", sessionList)
	r.Post("/session", sessionNew)
	r.Delete("/session", sessionDelete)

	// Add the router action
	m.Action(r.Handle)
}

func sessionDelete(w http.ResponseWriter, r *http.Request) []byte {
	sessionId := r.FormValue("id")
	if sessionId == "" {
		return jsonify("ID field can not be empty")
	}
	id, err := strconv.Atoi(sessionId)
	if err != nil {
		log.Panic(err)
		return jsonify("ID field can not be parsed")
	}
	delete(sessions, id)
	return jsonify(true)
}

func sessionNew(w http.ResponseWriter, r *http.Request) []byte {
	env := jsonList(r.FormValue("env"))
	session := NewSession(r.FormValue("name"), r.FormValue("type"), r.FormValue("workingpath"), env)
	sessions[session.id] = session
	return jsonify(session.id)
}

func sessionList(w http.ResponseWriter, r *http.Request) []byte {
	result := []map[string]interface{}{}
	for _, session := range sessions {
		sessionData := map[string]interface{}{"id": session.id, "name": session.name, "workingpath": session.workingPath, "type": session.cmdType, "env": session.env}
		result = append(result, sessionData)
	}
	return jsonify(result)
}

func jsonify(data interface{}) []byte {
	result, err := json.Marshal(data)
	if err != nil {
		log.Panic(err)
	}
	return result
}

func jsonList(data string) []string {
	result := []string{}
	if data != "" {
		if err := json.Unmarshal([]byte(data), &result); err != nil {
			log.Panic(err)
		}
	}
	return result
}

func execute(w http.ResponseWriter, r *http.Request, parms martini.Params) (string, int) {
	workingPath, cmd, args, env, cmdType := r.FormValue("workingpath"), r.FormValue("cmd"), r.FormValue("args"), r.FormValue("env"), r.FormValue("type")
	sid, err := strconv.Atoi(parms["sid"])
	if err == nil {
		session := sessions[sid]
		workingPath, env, cmdType = session.workingPath, string(jsonify(session.env)), session.cmdType
	}

	var pscmd string
	if cmdType == "ps" {
		pscmd = cmd
		cmd = "powershell"
	}

	command := exec.Command(cmd)

	cmdArgs := []string{cmd}
	if pscmd != "" {
		cmdArgs = append(cmdArgs, pscmd)
	}
	cmdArgs2 := jsonList(args)
	for _, v := range cmdArgs2 {
		cmdArgs = append(cmdArgs, v)
		command.Args = cmdArgs
	}

	envVars := jsonList(env)
	command.Env = envVars

	var out bytes.Buffer
	command.Stdout = &out
	command.Dir = workingPath
	if err := command.Run(); err != nil {
		log.Panic(err.Error() + ": " + out.String())
		return err.Error() + ": " + out.String(), 500
	}

	return out.String(), 200
}

func main() {
	if err := http.ListenAndServe(":3000", m); err != nil {
		log.Fatal(err)
	}
}
