package main

import (
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/gorilla/mux"
)

var (
	hostPort = 3999
)

// Assumes that present is running at $GOPATH/src
func main() {
	var port = flag.String("port", "3998", "The port to accept web requests.")
	// var origin = flag.String("origin", "127.0.0.1", "The origin host.")
	flag.Parse()

	r := mux.NewRouter()
	r.HandleFunc("/static/{any:.*}", proxy)
	r.HandleFunc("/play.js", proxy)
	r.HandleFunc("/github.com/{user:.+}/{name:[^/]+}{path:.*}", func(w http.ResponseWriter, r *http.Request) {
		repo := "github.com/" + mux.Vars(r)["user"] + "/" + mux.Vars(r)["name"]
		if err := maybeCloneGitRepo(repo); err != nil {
			handleErr(w, "Failed to clone git repo: "+repo+": "+err.Error(), "Failed clone repo.", http.StatusInternalServerError)
			return
		}
		proxy(w, r)
	})
	log.Println("Listening on :" + *port)
	http.ListenAndServe(":"+*port, r)
}

func handleErr(w http.ResponseWriter, logMsg string, clientMsg string, status int) {
	log.Println(logMsg)
	http.Error(w, clientMsg, status)
}

func proxy(w http.ResponseWriter, r *http.Request) {
	log.Println("Proxying request:", r.URL.Path)
	resp, err := http.Get("http://127.0.0.1:3999" + r.URL.Path)
	if err != nil {
		handleErr(w, "Failed to get path from present: "+err.Error(), "Failed to load presentation.", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	for key, values := range resp.Header {
		w.Header()[key] = values
	}
	io.Copy(w, resp.Body)
}

func maybeCloneGitRepo(repo string) error {
	dir := os.Getenv("GOPATH") + "/src/" + repo
	if _, err := os.Stat(dir); err != nil {
		if os.IsNotExist(err) {
			log.Println("Cloning:", "git", "clone", "https://"+repo+".git", dir)
			cmd := exec.Command("git", "clone", "https://"+repo+".git", dir)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			return cmd.Run()
		}
		return err
	}
	return nil
}
