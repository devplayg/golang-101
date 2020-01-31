package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"

	"github.com/pkg/sftp"
)

var (
	USER = flag.String("user", os.Getenv("USER"), "ssh username")
	HOST = flag.String("host", "localhost", "ssh server hostname")
	PORT = flag.Int("port", 22, "ssh server port")
	PASS = flag.String("pass", os.Getenv("SOCKSIE_SSH_PASSWORD"), "ssh password")
	SIZE = flag.Int("s", 1<<15, "set max packet size")
)

func init() {
	flag.Parse()
}

func main() {
	var auths []ssh.AuthMethod
	if aconn, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK")); err == nil {
		auths = append(auths, ssh.PublicKeysCallback(agent.NewClient(aconn).Signers))

	}
	if *PASS != "" {
		auths = append(auths, ssh.Password(*PASS))
	}

	config := ssh.ClientConfig{
		User:            *USER,
		Auth:            auths,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	addr := fmt.Sprintf("%s:%d", *HOST, *PORT)
	conn, err := ssh.Dial("tcp", addr, &config)
	if err != nil {
		log.Fatalf("unable to connect to [%s]: %v", addr, err)
	}
	defer conn.Close()

	c, err := sftp.NewClient(conn, sftp.MaxPacket(*SIZE))
	if err != nil {
		log.Fatalf("unable to start sftp subsytem: %v", err)
	}
	defer c.Close()

	src := "D:/utils/ubuntu-18.04.3-live-server-amd64.iso"

	w, err := c.OpenFile(filepath.Base(src), os.O_CREATE|os.O_RDWR|os.O_TRUNC)

	//w, err := c.OpenFile("/dev/null", syscall.O_WRONLY)
	if err != nil {
		log.Fatal(err)
	}
	defer w.Close()

	f, err := os.Open(src)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		log.Fatal(err)
	}

	size := fi.Size()

	log.Printf("writing %v bytes", size)
	t1 := time.Now()
	n, err := io.Copy(w, io.LimitReader(f, size))
	if err != nil {
		log.Fatal(err)
	}
	if n != size {
		log.Fatalf("copy: expected %v bytes, got %d", size, n)
	}
	log.Printf("wrote %v bytes in %s", size, time.Since(t1))
}

//package main
//
//import (
//    "flag"
//    "fmt"
//    "log"
//    "net"
//    "os"
//    "path/filepath"
//    "time"
//
//    "golang.org/x/crypto/ssh"
//    "golang.org/x/crypto/ssh/agent"
//
//    "github.com/pkg/sftp"
//)
//
//var (
//    USER = flag.String("user", os.Getenv("USER"), "ssh username")
//    HOST = flag.String("host", "localhost", "ssh server hostname")
//    PORT = flag.Int("port", 22, "ssh server port")
//    PASS = flag.String("pass", os.Getenv("SOCKSIE_SSH_PASSWORD"), "ssh password")
//    SIZE = flag.Int("s", 1<<15, "set max packet size")
//)
//
//func init() {
//    flag.Parse()
//}
//
//func main() {
//    var auths []ssh.AuthMethod
//    if aconn, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK")); err == nil {
//        auths = append(auths, ssh.PublicKeysCallback(agent.NewClient(aconn).Signers))
//
//    }
//    if *PASS != "" {
//        auths = append(auths, ssh.Password(*PASS))
//    }
//
//    config := ssh.ClientConfig{
//        User: *USER,
//        Auth: auths,
//        HostKeyCallback: ssh.InsecureIgnoreHostKey(),
//    }
//    addr := fmt.Sprintf("%s:%d", *HOST, *PORT)
//    conn, err := ssh.Dial("tcp", addr, &config)
//    if err != nil {
//        log.Fatalf("unable to connect to [%s]: %v", addr, err)
//    }
//    defer conn.Close()
//
//    c, err := sftp.NewClient(conn, sftp.MaxPacket(*SIZE))
//    if err != nil {
//        log.Fatalf("unable to start sftp subsytem: %v", err)
//    }
//    defer c.Close()
//
//    src := "D:/utils/ubuntu-18.04.3-live-server-amd64.iso"
//
//    w, err := c.OpenFile(filepath.Base(src), os.O_CREATE|os.O_RDWR|os.O_TRUNC)
//    if err != nil {
//        log.Fatal(err)
//    }
//    defer w.Close()
//
//    f, err := os.Open(src)
//    if err != nil {
//        log.Fatal(err)
//    }
//    defer f.Close()
//
//    //const size = 1e9
//    //
//    //log.Printf("writing %v bytes", size)
//    t1 := time.Now()
//    n ,  err := w.ReadFrom(f)
//    if err != nil {
//        log.Fatal(err)
//    }
//    //n, err := w.Write(make([]byte, size))
//    //if err != nil {
//    //    log.Fatal(err)
//    //}
//    //if n != size {
//    //   log.Fatalf("copy: expected %v bytes, got %d", size, n)
//    //}
//    log.Printf("wrote %v bytes in %s", n, time.Since(t1))
//}
