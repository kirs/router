package router

import (
	"fmt"
	"os/exec"
)

type Router struct {
	Alias     string
	Source    string
	LaunchCmd string
	RailsPort string
	Running   bool
}

func (r *Router) Launch() {
	if !r.Running {

		go func(command, port string) {
			// /usr/local/Cellar/rbenv/0.4.0/bin/
			// var launchString = "rbenv exec ruby " + command +  " server -p " + port
			// fmt.Println(launchString)
			fmt.Println("starting " + r.Alias + " process in gorutine")
			cmd := exec.Command("rbenv", "exec", "ruby", command, "server", "-p", port)

			cmd.Run()

		}(r.LaunchCmd, r.RailsPort)

		r.Running = true
	}
}

func PopulateRouters(path string) []Router {
	config := ReadConfig(path)

	routers := make([]Router, len(config.Routers))

	for i, router := range config.Routers {
		routers[i] = NewRouterFromConfig(&router)
	}

	return routers
}

func NewRouterFromConfig(cr *configRouter) Router {
	return Router{Alias: cr.Alias, Source: cr.Source, LaunchCmd: cr.LaunchCmd, Running: false, RailsPort: cr.RailsPort}
}
