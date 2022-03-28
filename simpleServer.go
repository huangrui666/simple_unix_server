package main

import (
	"log"
	"net"
	"os"
)

func main() {
	var sockAddr string
	if os.Getenv("K8S_ENV") == "true" {
		sockAddr = "/var/run/cri-resmgr/cri-resmgr-fps.sock"
	} else {
		sockAddr = "/home/vicky/ww05/image/workdir/ipc/resAllocSocket"
	}
	log.Println("sockAddr:", sockAddr)
	os.Remove(sockAddr)
	unixAddr, err := net.ResolveUnixAddr("unix", sockAddr)
	if err != nil {
		log.Println(err)
		return
	}
	listener, err := net.ListenUnix("unix", unixAddr)
	if err != nil {
		log.Println("listening error:", err)
	}
	os.Chmod(sockAddr, 0777)
	log.Println("listening... :", unixAddr)
	for {
		c, err := listener.Accept()
		if err != nil {
			log.Println("Accept: " + err.Error())
		} else {
			// go HandleServerConn(c)
			go HandleServerConn(c)
		}
	}
}

// msg
// "cpu" : "1-6",   [process cpuset]
// "fps" : "28.687267",  [fps]
// "game" : "stackball", [game name]
// "isFpsDrop" : "false", [whether fps satisfys QoS]
// "keySched" : "5.90,21.79,4", [key_threads: running, runnable, threads_count]
// "microTimeStamp" : "1648448084957443", [timestamp/us]
// "pod" : "",  [pod_name]
// "totalSched" : "6.70,49.33,238"  [total_threads: running, runnable, threads_count]

func HandleServerConn(c net.Conn) {
	log.Println("### HandleServerConn start", c)
	defer c.Close()
	for {
		buf := make([]byte, 10480)
		nr, err := c.Read(buf)
		if err != nil {
			log.Println("Read: " + err.Error())
			break
		} else {
			log.Println("msg:" + string(buf[0:nr]))
		}
	}
	log.Println("### HandleServerConn end")
}
