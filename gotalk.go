package gotalk

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"os/exec"
	"strings"
	"sync"

	"golang.org/x/net/context"
)

var (
	runningContainers = map[string]*ContainerStatus{} // Keyed on repo
	mu                sync.Mutex
)

func GetContainerStatus(repo string) *ContainerStatus {
	return runningContainers[repo]
}

type ContainerStatus struct {
	Repo     string
	HostPort string
	Name     string
}

func StartContainer(_ context.Context, repo, originHost, originPort string) (*ContainerStatus, error) {
	hostPort, err := freePort()
	if err != nil {
		return nil, err
	}

	name := strings.Replace(repo, "/", "_", -1)
	cmd := exec.Command("docker", "run", "-i", "-p", hostPort+":"+originPort, "--name="+name, "kevincantwell/gotalk:latest", repo, originHost, originPort)
	pipe, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}

	go func(cmd *exec.Cmd) {
		if err := cmd.Run(); err != nil {
			log.Println("Container exited:", runningContainers[repo], "error:", err)
			mu.Lock()
			delete(runningContainers, repo)
			mu.Unlock()
		}
	}(cmd)

	scanner := bufio.NewScanner(pipe)
	for scanner.Scan() {
		text := scanner.Text()
		log.Println(repo+":"+hostPort+">", text)
		if strings.Contains(text, "Open your web browser and visit") {
			status := ContainerStatus{Repo: repo, HostPort: hostPort, Name: name}
			mu.Lock()
			runningContainers[repo] = &status
			mu.Unlock()
			return &status, nil
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return nil, errors.New("container failed to start")
}

func StopContainer(status *ContainerStatus) error {
	if status == nil {
		return nil
	}
	if err := exec.Command("docker", "stop", "-t", "0", status.Name).Run(); err != nil {
		return err
	}
	mu.Lock()
	delete(runningContainers, status.Repo)
	mu.Unlock()
	// Must also remove the container if we wish to re-use the name in the future
	return exec.Command("docker", "rm", status.Name).Run()
}

func freePort() (string, error) {
	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:0")
	if err != nil {
		return "", err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return "", err
	}
	defer l.Close()
	return fmt.Sprint(l.Addr().(*net.TCPAddr).Port), nil
}
