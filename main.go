/* This program is free software. It comes without any warranty, to
     * the extent permitted by applicable law. You can redistribute it
     * and/or modify it under the terms of the Do What The Fuck You Want
     * To Public License, Version 2, as published by Sam Hocevar. See
     * http://www.wtfpl.net/ for more details. */

package main
import (
	"log"
	"os"
	"bytes"
	"net/http"
	"net/url"
	"io/ioutil"
	"flag"
	"time"
	"strconv"
)

//return true if network is connected
func detectNetwork() bool {

	resp, err := http.Get("http://www.msftconnecttest.com/connecttest.txt")
	if err != nil {
		log.Println("Error: "+err.Error());
		return false
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err!=nil {
		log.Println("Error: "+err.Error());
		return false
	}

	return bytes.Compare(body,[]byte ("Microsoft Connect Test"))==0;

}

func LoginPost(scheme,server,rawUsername,rawPassword string) {
	postUrl := scheme+"://"+server+"/srun_portal_pc.php?ac_id=1&"
	
	values := url.Values{}
	values.Set("action","login")
	values.Set("ac_id","1")
	values.Set("user_ip","")
	values.Set("nas_ip","")
	values.Set("user_mac","")
	values.Set("save_me","1")
	values.Set("username",rawUsername)
	values.Set("password",rawPassword)

	resp, err := http.PostForm(postUrl,values)
	if (err != nil){
		log.Println("Error: "+err.Error());
		return
	}
	defer resp.Body.Close()

}

func main(){
	scheme := flag.String("scheme","http","scheme (optional): http or https")
	server := flag.String("s","","server (required): server's hostname and port (port is optional)")
	username := flag.String("u","","username (required)")
	password := flag.String("p","","password (required)")
	logFilename := flag.String("log","","log (optional): the output log file. Otherwise it will be output to stdout")
	interval := flag.Int("interval",60,"interval (optional): time interval (seconds)")
	flag.Parse()

	//required arguments
	if (*server == ""||*username==""||*password==""){
		flag.Usage()
		return
	}

	//write log to stdout
	log.SetOutput(os.Stdout)

	//openfile
	if (*logFilename != ""){
		logFile,err := os.OpenFile(*logFilename,os.O_RDWR|os.O_CREATE|os.O_APPEND,0644)
		defer logFile.Close()
		if err != nil {
			log.Panicln(err)
			return
		}
		log.SetOutput(logFile)
	}
	
	//say hello
	log.Println("SDU-Qingdao SRUN3000 Portal Client")

	//loop
	for {
		if (! detectNetwork()){
			log.Println("Network is down. Log in via web portal...")
			LoginPost(*scheme,*server,*username,*password)
			time.Sleep(time.Duration(5)*time.Second)
			if (! detectNetwork()){
				log.Println("Network is still down. Retry after "+strconv.Itoa(*interval)+" seconds." )
			}else{
				log.Println("Network is up." )
			}
		}else{
			log.Println("Network is up. Nothing to do." )
		}
		time.Sleep(time.Duration(*interval)*time.Second)
	}
	
}
