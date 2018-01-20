package main

import (
	"bytes"
	"crypto/tls"
	"encoding/gob"
	"log"
	"reflect"

	"github.com/MarcStorm/secureconnection/generateRemote"
	"github.com/MarcStorm/secureconnection/protocol"
)

func main() {

	// cert, err := tls.LoadX509KeyPair("../certs/client.pem", "../certs/client.key")
	//
	// if err != nil {
	// 	log.Fatalf("server: loadkeys: %s", err)
	// }
	//
	// // create a pool of trusted certs
	// certPool := x509.NewCertPool()
	// pemFile, errRead := ioutil.ReadFile("../certs/server.pem")
	// if err != nil {
	// 	log.Fatalf("server: load pem: %s", errRead)
	// }
	// certPool.AppendCertsFromPEM(pemFile)
	//
	// config := &tls.Config{
	// 	RootCAs:            certPool,
	// 	Certificates:       []tls.Certificate{cert},
	// 	InsecureSkipVerify: false,
	// }

	_, config := generateRemote.Gr()
	//_, config := generate.GenerateCerts()

	// Everything below this comment should be handled by goSpace.
	conn, err := tls.Dial("tcp", "127.0.0.1:8000", config)

	if err != nil {
		log.Fatalf("client: dial: %s", err)
	}

	defer conn.Close()

	log.Println("client: connected to: ", conn.RemoteAddr())

	//state := conn.ConnectionState()

	// for _, v := range state.PeerCertificates {
	// 	fmt.Println("Client: Server public key is:")
	// 	fmt.Println(x509.MarshalPKIXPublicKey(v.PublicKey))
	// }
	//
	// log.Println("client: handshake: ", state.HandshakeComplete)
	// log.Println("client: mutual: ", state.NegotiatedProtocolIsMutual)

	// With gob
	registerTypes()

	// Convert message to byte array with gob
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	message := createMessage()

	errEnc := enc.Encode(message)
	if err != nil {
		log.Fatalf("client: gob.encode: %s", errEnc)
	}

	// Send and receive message.
	//message := "Hello\n"

	println("Client conn type: ", reflect.TypeOf(conn).String())

	n, errWrite := conn.Write(buf.Bytes())

	//n, err := io.WriteString(conn, message)

	if errWrite != nil {
		log.Fatalf("client: write: %s", errWrite)
	}

	log.Printf("client: wrote %q (%d bytes)", message, n)
	//
	// reply := make([]byte, 256)
	//
	// n, err = conn.Read(reply)
	//
	// log.Printf("client: read %q (%d bytes)", string(reply[:n]), n)
	log.Print("client: exiting")
}

func createMessage() protocol.Message {
	return protocol.CreateMessage("Put", createTuple())
}

func createTuple() protocol.Tuple {
	testFields := make([]interface{}, 5)
	testFields[0] = "Field 1"
	testFields[1] = 2
	testFields[2] = 3.14
	testFields[3] = false
	testFields[4] = 991

	return protocol.CreateTuple(testFields...)
}

func registerTypes() {
	// Register default structures for communication.
	gob.Register(protocol.Message{})
	gob.Register(protocol.Tuple{})
	gob.Register([]interface{}{})
}
