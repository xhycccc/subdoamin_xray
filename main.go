package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
)

var (
	file      = flag.String("f", "conf/domain.txt", "域名文件")
	thread   = flag.Int("t", 3, "进程数 例如:-n=3")
	h         = flag.Bool("h", false, "帮助信息")

)

func main() {
	flag.Parse()
	if *h == true {
		flag.PrintDefaults()
		return
	}
	domain_file, err := os.Open(*file)
	if err != nil{
		fmt.Println(err)
	}
	defer domain_file.Close()
	sc := bufio.NewScanner(domain_file)
	var domains []string
	for sc.Scan() {
		domain := sc.Text()
		domains = append(domains, domain)
	}
	c_task := make(chan struct{}, len(domains))
	c_thread := make(chan struct{}, *thread)
	defer close(c_task)
	defer close(c_thread)
	for _, domain := range domains{
		c_thread <- struct{}{}
		c_task <- struct{}{}
		go func(domain string) {
			cmd := exec.Command("./xray_linux_amd64", "subdomain", "--target", domain, "--json-output", domain+".json")
			output, err := cmd.CombinedOutput()
			if err != nil {
				fmt.Println(fmt.Sprint(err) + ":" + string(output))
				return
			}
			<- c_thread

		}(domain)
	}
	for i:=len(domains); i>0;i-- {
		<- c_task
	}
	fmt.Println("主协程")

}
