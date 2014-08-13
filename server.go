package main

import (
	"bytes"
	"encoding/json"
	//"fmt"
	"github.com/codegangsta/martini"
	"log"
	"net/http"
	"os/exec"
	//"strings"
	//"github.com/codegangsta/martini-contrib/auth"
)

var m *martini.Martini

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
	r.Get("/execute", execute)

	// Add the router action
	m.Action(r.Handle)
}

func execute(w http.ResponseWriter, r *http.Request) (string, int) {
	workingPath, cmd, args, env, cmdType := r.FormValue("workingpath"), r.FormValue("cmd"), r.FormValue("args"), r.FormValue("env"), r.FormValue("type")

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
	if args != "" {
		cmdArgs2 := []string{}
		if err := json.Unmarshal([]byte(args), &cmdArgs2); err != nil {
			log.Panic(err)
			return err.Error(), 500
		}
		for _, v := range cmdArgs2 {
			cmdArgs = append(cmdArgs, v)
		}
	}
	command.Args = cmdArgs

	if env != "" {
		envVars := []string{}
		if err := json.Unmarshal([]byte(env), &envVars); err != nil {
			log.Panic(err)
			return err.Error(), 500
		}
		command.Env = envVars
	}
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
