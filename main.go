package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

func main() {
	d := flag.String("d", "", "-d domain name")
	f := flag.String("f", "", "-f folder to be serve")
	r := flag.String("r", "", "-r remove virtual host")
	flag.Parse()
	dEmpty := len(*d) < 1
	fEmpty := len(*f) < 1

	if len(*r) > 0 {
		removeHost(*r)
		os.Exit(0)
	}
	if dEmpty && fEmpty {
		fmt.Println(usage)
		os.Exit(0)
	}

	if fEmpty {
		log.Fatal("please provide the name of the folder to be serve\n")
	}
	if dEmpty {
		log.Fatal("domain name are required")
	}

	if err := createServerDir(*d, *f); err != nil {
		log.Fatal(err)
	}
	if err := createVirtualConfig(*d); err != nil {
		log.Fatal(err)
	}
	setupHostConfig(*d)
	restartNginx()

	fmt.Println("Virtual Host deployed ", "http://"+*d)
}

func restartNginx() {
	if err := exec.Command("nginx", "-t").Run(); err != nil {
		log.Print(err.Error())
	}
	if err := exec.Command("systemctl", "restart", "nginx").Run(); err != nil {
		log.Print(err.Error())
	}

}

func createServerDir(domain, folder string) (err error) {
	domain = SERVER_FILES + domain
	_, exist := isExist(domain)
	if exist {
		fmt.Println("Updating server directory :" + domain)
		if err := os.RemoveAll(domain); err != nil {
			return err
		}
	}
	if err := os.Mkdir(domain, 0777); err != nil {
		return err
	}

	if dir, _ := isExist(folder); !dir {
		path, _ := os.Getwd()
		if before, _, ok := strings.Cut(path, folder); !ok {
			return errors.New("folder you provide doesnt exist")
		} else {
			folder = before + folder
		}
	}

	if err := exec.Command("cp", "-r", folder, domain+"/html").Run(); err != nil {
		return err
	}

	return nil
}

func createVirtualConfig(domain string) error {
	configName := domain + ".conf"
	available := SITES_AVAILABLE + configName
	enabled := SITES_ENABLED + configName
	config = strings.ReplaceAll(config, "example.com", domain)

	if _, exist := isExist(available); !exist {
		file, err := os.Create(available)
		if err != nil {
			return err
		}
		file.Close()

		if err := ioutil.WriteFile(available, []byte(config), 0664); err != nil {
			err = os.Remove(available)
			return err
		}
		if err := os.Symlink(available, enabled); err != nil {
			err = os.Remove(available)
			return err
		}
	}

	return nil
}
func setupHostConfig(domain string) {
	data, _ := ioutil.ReadFile(HOSTS_FILE)
	comment := "# virtual hosts added by clinx"
	host := LOCAL_IP + "\t" + domain
	commented := false
	commentLine := 0
	hostAdded := false
	stringN := strings.Split(string(data), "\n")

	for index, v := range stringN {
		if strings.EqualFold(v, host) {
			fmt.Println(domain, " already added to hosts file")
			hostAdded = true
			break
		}
		if !commented {
			if strings.EqualFold(comment, v) {
				commented = true
				commentLine = index
			}
		}
	}

	if !hostAdded {
		var texts []string
		texts = append(texts, comment, host)
		if !commented {
			stringN = append(texts, stringN...)
		} else {
			stringN = append(texts, stringN[commentLine+1:]...)
		}
		dataToWrite := strings.Join(stringN, "\n")
		ioutil.WriteFile(HOSTS_FILE, []byte(dataToWrite), 0644)
	}
}

func isExist(name string) (folder bool, exist bool) {
	info, err := os.Stat(name)
	if os.IsNotExist(err) {
		return false, false
	}
	return info.IsDir(), true
}

func removeHost(domain string) {
	fmt.Println("removing virtual host config for " + "http://" + domain)
	file := domain + ".conf"
	available := SITES_AVAILABLE + file
	enabled := SITES_ENABLED + file
	serverFolder := SERVER_FILES + domain
	hosts, _ := ioutil.ReadFile(HOSTS_FILE)
	hostsN := strings.Split(string(hosts), "\n")

	for index, v := range hostsN {
		if strings.EqualFold(v, LOCAL_IP+"\t"+domain) {
			hostsN = append(hostsN[:index], hostsN[index+1:]...)
			break
		}
	}
	ioutil.WriteFile(HOSTS_FILE, []byte(strings.Join(hostsN, "\n")), 0644)

	os.RemoveAll(available)
	os.RemoveAll(enabled)
	os.RemoveAll(serverFolder)
}
