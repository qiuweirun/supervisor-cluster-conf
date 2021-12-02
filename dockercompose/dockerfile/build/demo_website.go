package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
)

func GetIP() (string) {
	netInterfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("net.Interfaces failed, err:", err.Error())
	}

	var ip string
	Loop:
	for i := 0; i < len(netInterfaces); i++ {
		if (netInterfaces[i].Flags & net.FlagUp) != 0 {
			addrs, _ := netInterfaces[i].Addrs()

			for _, address := range addrs {
				if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					if ipnet.IP.To4() != nil {
						ip = ipnet.IP.String()
						break Loop
					}
				}
			}
		}
	}
	return ip
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	name, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	fmt.Fprintln(w, GetIP()+" demo.website1.com")
	fmt.Fprintln(w, "\r\n")
	fmt.Fprintln(w, "主机名："+name)
	fmt.Fprintln(w, "\r\n")
	fmt.Fprintln(w, "进程ID:"+strconv.Itoa(os.Getpid()));
}

func main() {
	http.HandleFunc("/", IndexHandler)
	http.ListenAndServe(":8088", nil)
}
