package main

import (
	"flag"
	"log"
	"net/http"
	"strings"

	"github.com/kevin-cantwell/gotalk"

	"github.com/gorilla/mux"
	"golang.org/x/net/context"
)

var (
	hostPort = 3999
)

func main() {
	var port = flag.String("port", "5555", "The port to accept web requests.")
	flag.Parse()

	r := mux.NewRouter()
	r.HandleFunc("/{repo:[^!]+}{sep:!?}{path:.*}", func(w http.ResponseWriter, r *http.Request) {
		repo := mux.Vars(r)["repo"]
		if status := gotalk.GetContainerStatus(repo); status == nil {
			if err := StartContainer(w, r); err != nil {
				handleErr(w, "Failed to start container: "+err.Error(), "Failed to start presentation.", http.StatusInternalServerError)
				return
			}
		}
		if err := ServePresentation(w, r); err != nil {
			handleErr(w, "Failed to serve presentation: "+err.Error(), "The presentation is not available.", http.StatusInternalServerError)
			return
		}
	})
	log.Println("Listening on :" + *port)
	http.ListenAndServe(":"+*port, r)
}

func handleErr(w http.ResponseWriter, logMsg string, clientMsg string, status int) {
	log.Println(logMsg)
	http.Error(w, clientMsg, status)
}

func StartContainer(w http.ResponseWriter, r *http.Request) error {
	repo := mux.Vars(r)["repo"]
	log.Println("Starting containter for repo:", repo)
	hostInfo := strings.Split(r.Host, ":")
	originHost, originPort := hostInfo[0], hostInfo[1]
	if originPort == "" {
		originPort = "80"
	}
	_, err := gotalk.StartContainer(context.TODO(), repo, originHost, originPort)
	if err != nil {
		return err
	}
	return nil
}

func ServePresentation(w http.ResponseWriter, r *http.Request) error {
	// repo := mux.Vars("repo")
	// status := gotalk.GetContainerStatus(repo)
	path := mux.Vars(r)["path"]
	log.Println("Serving presentation path:", "/"+path)
	w.Write([]byte(path))

	// resp, err := http.Get("http://")
	return nil
}
