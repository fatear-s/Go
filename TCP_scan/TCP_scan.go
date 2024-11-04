package main

import (
	"flag"
	"fmt"
	"net"
	"sync"
	"time"
)

func main() {
	host := flag.String("host_name:", "baidu.com", "string")
	startPort := flag.Int("start_port:", 50, "int")
	endPort := flag.Int("end_port", 200, "int")
	timeOut := flag.Duration("timeout", time.Millisecond*1000, "Duration")
	//指针
	flag.Parse()

	wg := &sync.WaitGroup{}
	mu := &sync.Mutex{}
	ports := []int{}
	for port := *startPort; port <= *endPort; port++ {
		wg.Add(1)
		go func(p int) {
			if isOpen(*host, p, *timeOut) {
				mu.Lock()
				ports = append(ports, p)
				mu.Unlock()
			}
			wg.Done()
		}(port)
	}
	wg.Wait()
	fmt.Printf("connection successful:%s:%d\n", *host, ports[0])

}
func isOpen(host string, port int, timeout time.Duration) bool {
	time.Sleep(time.Millisecond * 1)
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", host, port), timeout)
	if err == nil {
		conn.Close()
		fmt.Println("Connect successful")
		return true
	} else {
		//fmt.Println(err)
		return false
	}
}
