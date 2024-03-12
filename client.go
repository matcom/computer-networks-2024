package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"time"
)

func parse_get_connection_ftp(input string) (string, string){
	// Parse the input string to get the IP and Port
	// Example: "227 Entering Passive Mode (127,0,0,1,195,149)"
	// The port is calculated as (195*256+149)	
	// Output: "127.0.0.1", "195*256+149".
	split1 := strings.Split(input, "(")
	split2 := strings.Split(split1[1], ",")
	split3 := strings.Split(split2[5], ")")
	ip := split2[0] + "." + split2[1] + "." + split2[2] + "." + split2[3]
	first_part_port, _ := strconv.ParseInt(split2[4], 10, 32)
	second_part_port, _ := strconv.ParseInt(split3[0], 10, 32)

	port := strconv.FormatInt((first_part_port*256 + second_part_port), 10)
	return ip, port
} 

func wr(connConfig *net.Conn ,data string) (string, error) {
	reader := bufio.NewReader(*connConfig)
	(*connConfig).Write([]byte(data + "\r\n"))
	time.Sleep( 1* time.Second)
    var response strings.Builder
	for {
    	line, err := reader.ReadString()
    	response.WriteString(line)
    	if err != nil {
        	if err == io.EOF {
            	break
        }
        return "", err
		}
	}
	
	fmt.Println("Response:", response)
	return response.String(), nil
}

func read(connData *net.Conn) (string, error) {
	reader := bufio.NewReader(*connData)
	time.Sleep(1 * time.Second)
    var response strings.Builder
	for {
    	line, err := reader.ReadString('\n')
    	response.WriteString(line)
    	if err != nil {
        	if err == io.EOF {
            	break
        }
        return "", err
		}
	}
	fmt.Println("Data Received:", response.String())
	return response.String(), nil
}

func open_conection(connDataConfig string) (net.Conn, error){
	connDataIP , connDataPort := parse_get_connection_ftp(connDataConfig)
	connData, err := net.Dial("tcp", connDataIP + ":" + connDataPort)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return connData, nil
}

func main(){
	connConfig, err := net.Dial("tcp", "localhost:21")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer connConfig.Close()
	wr(&connConfig, "Habla")
	wr(&connConfig, "USER brito")
	wr(&connConfig, "PASS password")
	wr(&connConfig, "CWD /upload")

}