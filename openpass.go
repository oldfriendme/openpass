package main

import (
	"net"
	"os"
	"fmt"
	"io"
	"runtime"
	"time"
	"encoding/hex"
	"strconv"
)
var ic string;var caesarNum int;
func caesarShift(c net.Conn,srv bool) {
	r, err := net.DialTimeout("tcp", ic, time.Duration(time.Second*30))
	if err != nil {
		fmt.Println("Connect remote :", err)
		c.Close()
		return
	}
	ends := make(chan bool, 2)
	go dc(c,r,ends);
	if !srv{
		go func(){
			caesar:=caesarNum&1023
			k:=32768-caesar
			bf:=make([]byte,2*k)
			eb:=make([]byte,1)
			caesarKey:=byte(caesarNum%23)
			for ;k>0;{
			n,err:=r.Read(bf[:2*k])
			if err!=nil{ends<-true;return;}
			if n&1==1{
				_,err:=r.Read(eb)
				if err!=nil{ends<-true;return;}
				bf[n]=eb[0]
				n++;
			}
			ks:=bf[:n]
			for i:=0;i<len(ks);i++{
				ks[i]-=caesarKey
				caesarKey++
				caesarKey%=23
			}
			bt, err := hex.DecodeString(string(ks))
			if err!=nil{ends<-true;return;}
			k-=n/2;
			c.Write(bt)
			}
			io.Copy(c,r)
			ends<-true
		}()
		caesar:=caesarNum&1023
		k:=32768-caesar
		bf:=make([]byte,k)
		caesarKey:=byte(caesarNum%23)
for ;k>0;{
		n,err:=c.Read(bf[:k])
		if err!=nil{ends<-true;return;}
		k-=n;
		he := hex.EncodeToString(bf[:n])
		ks:=[]byte(he)
		for i:=0;i<len(ks);i++{
			ks[i]+=caesarKey
			caesarKey++
			caesarKey%=23
		}
		r.Write(ks)
}	
io.Copy(r,c)
ends<-true
	}else{
		go func(){
			caesar:=caesarNum&1023
			k:=32768-caesar
			bf:=make([]byte,2*k)
			eb:=make([]byte,1)
			caesarKey:=byte(caesarNum%23)
			for ;k>0;{
			n,err:=c.Read(bf[:2*k])
			if err!=nil{ends<-true;return;}
			if n&1==1{
				_,err:=c.Read(eb)
				if err!=nil{ends<-true;return;}
				bf[n]=eb[0]
				n++;
			}
			ks:=bf[:n]
			for i:=0;i<len(ks);i++{
				ks[i]-=caesarKey
				caesarKey++
				caesarKey%=23
			}
			bt, err := hex.DecodeString(string(ks))
			if err!=nil{ends<-true;return;}
			k-=n/2;
			r.Write(bt)
			}
			io.Copy(r,c)
			ends<-true
		}()
		caesar:=caesarNum&1023
		k:=32768-caesar
		bf:=make([]byte,k)
		caesarKey:=byte(caesarNum%23)
for ;k>0;{
		n,err:=r.Read(bf[:k])
		if err!=nil{ends<-true;return;}
		k-=n;
		he := hex.EncodeToString(bf[:n])
		ks:=[]byte(he)
		for i:=0;i<len(ks);i++{
			ks[i]+=caesarKey
			caesarKey++
			caesarKey%=23
		}
		c.Write(ks)
}	
io.Copy(c,r)
ends<-true
	}
}

func dc(c,r net.Conn,e chan bool){
    <-e
	r.Close()
	c.Close()
}

func main() {
	argc:=len(os.Args)
	if argc < 5 {
		fmt.Println("openpass v0.1")
		fmt.Println("Usage: openpass [server/client] [Listen] [Connect] [caesarNum]")
		os.Exit(0)
	}
	runtime.GOMAXPROCS(runtime.NumCPU())
	srv:=false
	if os.Args[1][0]=='s'{srv=true;}
	ic=os.Args[3]
	l, err := net.Listen("tcp", os.Args[2])
	if err != nil {
	fmt.Println(err)
	os.Exit(-1)
	}
	caesarNum, err = strconv.Atoi(os.Args[4])
	if err!=nil{caesarNum=0;}
	if caesarNum <0{caesarNum=0;}
	fmt.Printf("INFO: Listening %s \n", os.Args[2])
	for {
		c, err := l.Accept()
		if err != nil {
		fmt.Println(err)
		} else {
	 	fmt.Println("INFO: new client:", c.RemoteAddr())
		go caesarShift(c,srv)
		}
	}

}
