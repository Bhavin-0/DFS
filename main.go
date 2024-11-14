package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/Bhavin-0/DFS/p2p"  // Importing the p2p package
	"github.com/Bhavin-0/DFS/DFS"  // Import the DFS package for FileServer and related structs
)

// Adjusted makeServer function to fully qualify types and functions
func makeServer(listenAddr string, nodes ...string) *DFS.FileServer {
	// Initializing the TCP transport with the p2p package options
	tcptransportOpts := p2p.TCPTransportOpts{
		ListenAddr:    listenAddr,
		HandshakeFunc: p2p.NOPHandshakeFunc,
		Decoder:       p2p.DefaultDecoder{},
	}
	tcpTransport := p2p.NewTCPTransport(tcptransportOpts)

	// Using the DFS package to qualify FileServerOpts and related functions
	fileServerOpts := DFS.FileServerOpts{
		EncKey:            DFS.NewEncryptionKey(),     // Assuming newEncryptionKey is in DFS as NewEncryptionKey
		StorageRoot:       listenAddr + "_network",
		PathTransformFunc: DFS.CASPathTransformFunc,   // Qualify CASPathTransformFunc if defined in DFS
		Transport:         tcpTransport,
		BootstrapNodes:    nodes,
	}

	s := DFS.NewFileServer(fileServerOpts)  // Assuming NewFileServer is in DFS

	tcpTransport.OnPeer = s.OnPeer

	return s
}

func main() {
	// Initialize servers on different ports
	s1 := makeServer(":3000")
	s2 := makeServer(":7000")
	s3 := makeServer(":5000", ":3000", ":7000")

	// Start servers in separate goroutines
	go func() { log.Fatal(s1.Start()) }()
	time.Sleep(500 * time.Millisecond)
	go func() { log.Fatal(s2.Start()) }()

	time.Sleep(2 * time.Second)

	go s3.Start()
	time.Sleep(2 * time.Second)

	// Perform file storage operations
	for i := 0; i < 20; i++ {
		key := fmt.Sprintf("picture_%d.png", i)
		data := bytes.NewReader([]byte("my big data file here!"))
		s3.Store(key, data)

		if err := s3.Store.Delete(s3.ID, key); err != nil {  // Using fully qualified methods if necessary
			log.Fatal(err)
		}

		r, err := s3.Get(key)
		if err != nil {
			log.Fatal(err)
		}

		b, err := ioutil.ReadAll(r)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(string(b))
	}
}
