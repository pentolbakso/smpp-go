package main

import (
	"context"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/pentolbakso/smpp-go"
	"github.com/pentolbakso/smpp-go/pdu"
)

var (
	serverAddr string
	dstAddr    string
	srcAddr    string
	msg        string
)

func main() {
	flag.StringVar(&serverAddr, "addr", "localhost:2775", "server will listen on this address.")
	flag.StringVar(&dstAddr, "dst_addr", "111111", "destination to which you are sending the message.")
	flag.StringVar(&srcAddr, "src_addr", "222222", "source from which the message is comming from.")
	flag.StringVar(&msg, "msg", "hello world", "contents of the message.")
	flag.Parse()

	bc := smpp.BindConf{
		Addr:     serverAddr,
		SystemID: "ExampleClient",
		Password: "password",
	}
	sc := smpp.SessionConf{}
	sess, err := smpp.BindTRx(sc, bc)
	if err != nil {
		fail("Can't bind: %v", err)
	}

	// test simple message -------

	sm := &pdu.SubmitSm{
		SourceAddr:      srcAddr,
		DestinationAddr: dstAddr,
		ShortMessage:    []byte(msg),
	}
	_, resp, err := sess.Send(context.Background(), sm)
	if err != nil {
		log.Printf("Can't send message: %+v", err)
	}
	log.Printf("Message sent\n")
	log.Printf("Received response %s %+v\n", resp.CommandID(), resp)

	// test with UDH ------------

	msgWithUdh, _ := hex.DecodeString("0B0504158200000003AA030174657374") // UDH + "test"

	smWithUdh := &pdu.SubmitSm{
		SourceAddr:      srcAddr,
		DestinationAddr: dstAddr,
		EsmClass: pdu.EsmClass{
			Feature: pdu.UDHIEsmFeat,
		},
		ShortMessage: []byte(msgWithUdh),
	}
	_, resp2, err2 := sess.Send(context.Background(), smWithUdh)
	if err2 != nil {
		log.Printf("Can't send message: %+v", err2)
	}
	log.Printf("Message sent\n")
	log.Printf("Received response %s %+v\n", resp2.CommandID(), resp2)

	if err := smpp.Unbind(context.Background(), sess); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
}

func fail(msg string, params ...interface{}) {
	log.Printf(msg+"\n", params...)
	os.Exit(1)
}
