package kuberneteslauncher

import (
	"path"
	"os"
	"io/ioutil"
	"os/exec"
	"time"
	"strconv"
	"syscall"
	"fmt"
	"strings"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"log"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/api/core/v1"
	"k8s.io/client-go/rest"
)

func ListOnConfigmap(dest string, namespace string, configmap string, command []string) {
	fmt.Printf("Monitoring changes in %s/%s", namespace, configmap)
	var config *rest.Config
	var err error
	if os.Getenv("KUBERNETES_CONFIG") != "" {
		config, err = clientcmd.BuildConfigFromFlags("", os.Getenv("KUBERNETES_CONFIG"))
		if err != nil {
			panic(err.Error())
		}
	} else {
		config, err = rest.InClusterConfig()
		if err != nil {
			panic(err.Error())
		}
	}

	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	supervisor := make(chan bool)

	started := false
	for {
		watch, err := clientset.CoreV1().ConfigMaps(namespace).Watch(meta_v1.ListOptions{Watch: true, FieldSelector: "metadata.name=" + configmap})
		if err != nil {
			panic(err)
		}
		events := watch.ResultChan()
		for {
			event, ok := <-events
			log.Printf("Service Event %v: %+v", event.Type, event.Object.GetObjectKind())
			if event.Type == "MODIFIED" || event.Type == "ADDED" {
				println("Configmap is added/modified")
				cm := event.Object.(*v1.ConfigMap)
				for key, value := range cm.Data {
					saveFile(dest, key, []byte(value))
				}
				if !started {
					go startProcess(command, supervisor)
					started = true
				} else {
					supervisor <- true
				}
			} else if event.Type == "DELETED" {
			} else if event.Type == "" {
				log.Printf("Service watch timed out")
			} else {
				log.Printf("Service watch unhandled event: %v", event.Type)
			}
			if !ok {
				break
			}
		}

	}
}

func kerub(supervisor chan bool, process *os.Process) {
	signal := <-supervisor
	if (signal) {
		fmt.Println("Killing process " + strconv.Itoa(process.Pid) + " ")
		process.Kill()
	}

}
func startProcess(command[] string, supervisor chan bool) {
	retry := true
	for retry {
		var cmd *exec.Cmd
		if len(command) > 1 {
			cmd = exec.Command(command[0], command[1:]...)
		} else {
			cmd = exec.Command(command[0])
		}
		os.Stdout.Sync()
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		fmt.Println("Starting process: " + strings.Join(command, " "))
		err := cmd.Start()
		if err != nil {
			panic(err)
		}
		go kerub(supervisor, cmd.Process)
		err = cmd.Wait()
		if err != nil {
			if exitError, ok := err.(*exec.ExitError); ok {
				waitStatus := exitError.Sys().(syscall.WaitStatus)
				fmt.Println(fmt.Sprintf("Exit code: %d", waitStatus.ExitStatus()))
			} else {
				fmt.Println("Other error: " + err.Error())
			}
		} else {
			retry = !cmd.ProcessState.Success()
			fmt.Println("Process has been stopped with exit code: " + strconv.Itoa(int(cmd.ProcessState.Sys().(syscall.WaitStatus))))
		}
		time.Sleep(5 * time.Second)
	}
	os.Exit(0)
}

func saveFile(directory string, relative_path string, bytes []byte) {
	dest_file := path.Join(directory, relative_path)
	dest_dir := path.Dir(dest_file)
	err := os.MkdirAll(dest_dir, 0777)
	if (err != nil) {
		panic(err)
	}
	err = ioutil.WriteFile(dest_file, bytes, 0644)
	if err != nil {
		panic(err)
	}
	fmt.Println(dest_file + " file is written")
}
