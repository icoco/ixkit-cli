package core

import (
	"fmt"
	"github.com/satori/go.uuid"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
)

func GetUUID() string {
	u1 := uuid.Must(uuid.NewV1())
	return u1.String()
}

//https://stackoverflow.com/questions/23558425/how-do-i-get-the-local-ip-address-in-go
// Get preferred outbound ip of this machine
func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

func GetMacAddress() string {
	// 获取本机的MAC地址
	interfaces, err := net.Interfaces()
	if err != nil {
		panic("Poor soul, here is what you got: " + err.Error())
	}
	var result string
	for _, inter := range interfaces {
		mac := inter.HardwareAddr //获取本机MAC地址
		fmt.Println("%s,MAC = %s", inter.Name, mac)
		result = mac.String()

	}
	return result
}

func Response2Json(r *http.Response) string {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
	}
	return string(body)
}

func GetCurrentPath() string {
	dir, _ := os.Getwd()
	return dir
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		//log.Fatal(err)
	}
	return dir
}

func GetTempFolder() string {
	dir, err := ioutil.TempDir("", "com.ixkit.cli.")
	if err != nil {
		fmt.Println(err)
	}
	return dir
}
