// Package smpp implements SMPP protocol v3.4.
//
// It allows easier creation of SMPP clients and servers by providing utilities for PDU and session handling.
// In order to do any kind of interaction you first need to create an SMPP [Session](https://godoc.org/github.com/pentolbakso/smpp-go#Session). Session is the main carrier of the protocol and enforcer of the specification rules.
//
// Naked session can be created with:
//
//	// You must provide already established connection and configuration struct.
//	sess := smpp.NewSession(conn, conf)
//
// But it's much more convenient to use helpers that would do the binding with the remote SMSC and return you session prepared for sending:
//
//	// Bind with remote server by providing config structs.
//	sess, err := smpp.BindTRx(sessConf, bindConf)
//
// And once you have the session it can be used for sending PDUs to the binded peer.
//
//	sm := smpp.SubmitSm{
//	    SourceAddr:      "11111111",
//	    DestinationAddr: "22222222",
//	    ShortMessage:    []byte("Hello from SMPP!"),
//	}
//	// Session can then be used for sending PDUs.
//	resp, err := sess.Send(p)
//
// Session that is no longer used must be closed:
//
//	sess.Close()
//
// If you want to handle incoming requests to the session specify SMPPHandler in session configuration when creating new session similarly to HTTPHandler from _net/http_ package:
//
//	conf := smpp.SessionConf{
//	    Handler: smpp.HandlerFunc(func(ctx *smpp.Context) {
//	        switch ctx.CommandID() {
//	        case pdu.UnbindID:
//	            ubd, err := ctx.Unbind()
//	            if err != nil {
//	                t.Errorf(err.Error())
//	            }
//	            resp := ubd.Response()
//	            if err := ctx.Respond(resp, pdu.StatusOK); err != nil {
//	                t.Errorf(err.Error())
//	            }
//	        }
//	    }),
//	}
//
// Detailed examples for SMPP client and server can be found in the examples dir.
package smpp

import (
	"context"
	"net"
	"time"

	"github.com/pentolbakso/smpp-go/pdu"
)

const (
	// Version of the supported SMPP Protocol. Only supporting 3.4 for now.
	Version = 0x34
	// SequenceStart is the starting reference for sequence number.
	SequenceStart = 0x00000001
	// SequenceEnd s sequence number upper boundary.
	SequenceEnd = 0x7FFFFFFF
)

// BindConf is the configuration for binding to smpp servers.
type BindConf struct {
	// Bind will be attempted to this addr.
	Addr string
	// Mandatory fields for binding PDU.
	SystemID   string
	Password   string
	SystemType string
	AddrTon    int
	AddrNpi    int
	AddrRange  string
}

func bind(req pdu.PDU, sc SessionConf, bc BindConf) (*Session, error) {
	conn, err := net.Dial("tcp", bc.Addr)
	if err != nil {
		return nil, err
	}
	sess := NewSession(conn, sc)
	timeout := sc.WindowTimeout
	if timeout == 0 {
		timeout = time.Second * 5
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	_, _, err = sess.Send(ctx, req)
	if err != nil {
		return sess, err
	}
	return sess, nil
}

// BindTx binds transmitter session.
func BindTx(sc SessionConf, bc BindConf) (*Session, error) {
	return bind(&pdu.BindTx{
		SystemID:         bc.SystemID,
		Password:         bc.Password,
		SystemType:       bc.SystemType,
		InterfaceVersion: Version,
		AddrTon:          bc.AddrTon,
		AddrNpi:          bc.AddrNpi,
		AddressRange:     bc.AddrRange,
	}, sc, bc)
}

// BindRx binds receiver session.
func BindRx(sc SessionConf, bc BindConf) (*Session, error) {
	return bind(&pdu.BindRx{
		SystemID:         bc.SystemID,
		Password:         bc.Password,
		SystemType:       bc.SystemType,
		InterfaceVersion: Version,
		AddrTon:          bc.AddrTon,
		AddrNpi:          bc.AddrNpi,
		AddressRange:     bc.AddrRange,
	}, sc, bc)
}

// BindTRx binds transreceiver session.
func BindTRx(sc SessionConf, bc BindConf) (*Session, error) {
	return bind(&pdu.BindTRx{
		SystemID:         bc.SystemID,
		Password:         bc.Password,
		SystemType:       bc.SystemType,
		InterfaceVersion: Version,
		AddrTon:          bc.AddrTon,
		AddrNpi:          bc.AddrNpi,
		AddressRange:     bc.AddrRange,
	}, sc, bc)
}

// Unbind session will initiate session unbinding and close the session.
// First it will try to notify peer with unbind request.
// If there was any error during unbinding an error will be returned.
// Session will be closed even if there was an error during unbind.
func Unbind(ctx context.Context, sess *Session) error {
	defer func() {
		sess.Close()
	}()
	_, _, err := sess.Send(ctx, pdu.Unbind{})
	if err != nil {
		return err
	}
	return nil
}

// SendGenericNack is a helper function for sending GenericNack PDU.
func SendGenericNack(ctx context.Context, sess *Session, p *pdu.GenericNack) error {
	_, _, err := sess.Send(ctx, p)
	if err != nil {
		return err
	}
	return nil
}

// SendBindRx is a helper function for sending BindRx PDU.
func SendBindRx(ctx context.Context, sess *Session, p *pdu.BindRx) (*pdu.BindRxResp, error) {
	var tresp *pdu.BindRxResp
	_, resp, err := sess.Send(ctx, p)
	if resp != nil {
		tresp = resp.(*pdu.BindRxResp)
	}
	if err != nil {
		return tresp, err
	}
	return tresp, nil
}

// SendBindRxResp is a helper function for sending BindRxResp PDU.
func SendBindRxResp(ctx context.Context, sess *Session, p *pdu.BindRxResp) error {
	_, _, err := sess.Send(ctx, p)
	if err != nil {
		return err
	}
	return nil
}

// SendBindTx is a helper function for sending BindTx PDU.
func SendBindTx(ctx context.Context, sess *Session, p *pdu.BindTx) (*pdu.BindTxResp, error) {
	var tresp *pdu.BindTxResp
	_, resp, err := sess.Send(ctx, p)
	if resp != nil {
		tresp = resp.(*pdu.BindTxResp)
	}
	if err != nil {
		return tresp, err
	}
	return tresp, nil
}

// SendBindTxResp is a helper function for sending BindTxResp PDU.
func SendBindTxResp(ctx context.Context, sess *Session, p *pdu.BindTxResp) error {
	_, _, err := sess.Send(ctx, p)
	if err != nil {
		return err
	}
	return nil
}

// SendQuerySm is a helper function for sending QuerySm PDU.
func SendQuerySm(ctx context.Context, sess *Session, p *pdu.QuerySm) (*pdu.QuerySmResp, error) {
	var tresp *pdu.QuerySmResp
	_, resp, err := sess.Send(ctx, p)
	if resp != nil {
		tresp = resp.(*pdu.QuerySmResp)
	}
	if err != nil {
		return tresp, err
	}
	return tresp, nil
}

// SendQuerySmResp is a helper function for sending QuerySmResp PDU.
func SendQuerySmResp(ctx context.Context, sess *Session, p *pdu.QuerySmResp) error {
	_, _, err := sess.Send(ctx, p)
	if err != nil {
		return err
	}
	return nil
}

// SendSubmitSm is a helper function for sending SubmitSm PDU.
func SendSubmitSm(ctx context.Context, sess *Session, p *pdu.SubmitSm) (*pdu.SubmitSmResp, error) {
	var tresp *pdu.SubmitSmResp
	_, resp, err := sess.Send(ctx, p)
	if resp != nil {
		tresp = resp.(*pdu.SubmitSmResp)
	}
	if err != nil {
		return tresp, err
	}
	return tresp, nil
}

// SendSubmitSmResp is a helper function for sending SubmitSmResp PDU.
func SendSubmitSmResp(ctx context.Context, sess *Session, p *pdu.SubmitSmResp) error {
	_, _, err := sess.Send(ctx, p)
	if err != nil {
		return err
	}
	return nil
}

// SendDeliverSm is a helper function for sending DeliverSm PDU.
func SendDeliverSm(ctx context.Context, sess *Session, p *pdu.DeliverSm) (*pdu.DeliverSmResp, error) {
	var tresp *pdu.DeliverSmResp
	_, resp, err := sess.Send(ctx, p)
	if resp != nil {
		tresp = resp.(*pdu.DeliverSmResp)
	}
	if err != nil {
		return tresp, err
	}
	return tresp, nil
}

// SendDeliverSmResp is a helper function for sending DeliverSmResp PDU.
func SendDeliverSmResp(ctx context.Context, sess *Session, p *pdu.DeliverSmResp) error {
	_, _, err := sess.Send(ctx, p)
	if err != nil {
		return err
	}
	return nil
}

// SendUnbind is a helper function for sending Unbind PDU.
func SendUnbind(ctx context.Context, sess *Session, p *pdu.Unbind) (*pdu.UnbindResp, error) {
	var tresp *pdu.UnbindResp
	_, resp, err := sess.Send(ctx, p)
	if resp != nil {
		tresp = resp.(*pdu.UnbindResp)
	}
	if err != nil {
		return tresp, err
	}
	return tresp, nil
}

// SendUnbindResp is a helper function for sending UnbindResp PDU.
func SendUnbindResp(ctx context.Context, sess *Session, p *pdu.UnbindResp) error {
	_, _, err := sess.Send(ctx, p)
	if err != nil {
		return err
	}
	return nil
}

// SendReplaceSm is a helper function for sending ReplaceSm PDU.
func SendReplaceSm(ctx context.Context, sess *Session, p *pdu.ReplaceSm) (*pdu.ReplaceSmResp, error) {
	var tresp *pdu.ReplaceSmResp
	_, resp, err := sess.Send(ctx, p)
	if resp != nil {
		tresp = resp.(*pdu.ReplaceSmResp)
	}
	if err != nil {
		return tresp, err
	}
	return tresp, nil
}

// SendReplaceSmResp is a helper function for sending ReplaceSmResp PDU.
func SendReplaceSmResp(ctx context.Context, sess *Session, p *pdu.ReplaceSmResp) error {
	_, _, err := sess.Send(ctx, p)
	if err != nil {
		return err
	}
	return nil
}

// SendCancelSm is a helper function for sending CancelSm PDU.
func SendCancelSm(ctx context.Context, sess *Session, p *pdu.CancelSm) (*pdu.CancelSmResp, error) {
	var tresp *pdu.CancelSmResp
	_, resp, err := sess.Send(ctx, p)
	if resp != nil {
		tresp = resp.(*pdu.CancelSmResp)
	}
	if err != nil {
		return tresp, err
	}
	return tresp, nil
}

// SendCancelSmResp is a helper function for sending CancelSmResp PDU.
func SendCancelSmResp(ctx context.Context, sess *Session, p *pdu.CancelSmResp) error {
	_, _, err := sess.Send(ctx, p)
	if err != nil {
		return err
	}
	return nil
}

// SendBindTRx is a helper function for sending BindTRx PDU.
func SendBindTRx(ctx context.Context, sess *Session, p *pdu.BindTRx) (*pdu.BindTRxResp, error) {
	var tresp *pdu.BindTRxResp
	_, resp, err := sess.Send(ctx, p)
	if resp != nil {
		tresp = resp.(*pdu.BindTRxResp)
	}
	if err != nil {
		return tresp, err
	}
	return tresp, nil
}

// SendBindTRxResp is a helper function for sending BindTRxResp PDU.
func SendBindTRxResp(ctx context.Context, sess *Session, p *pdu.BindTRxResp) error {
	_, _, err := sess.Send(ctx, p)
	if err != nil {
		return err
	}
	return nil
}

// SendOutbind is a helper function for sending Outbind PDU.
func SendOutbind(ctx context.Context, sess *Session, p *pdu.Outbind) error {
	_, _, err := sess.Send(ctx, p)
	if err != nil {
		return err
	}
	return nil
}

// SendEnquireLink is a helper function for sending EnquireLink PDU.
func SendEnquireLink(ctx context.Context, sess *Session, p *pdu.EnquireLink) (*pdu.EnquireLinkResp, error) {
	var tresp *pdu.EnquireLinkResp
	_, resp, err := sess.Send(ctx, p)
	if resp != nil {
		tresp = resp.(*pdu.EnquireLinkResp)
	}
	if err != nil {
		return tresp, err
	}
	return tresp, nil
}

// SendEnquireLinkResp is a helper function for sending EnquireLinkResp PDU.
func SendEnquireLinkResp(ctx context.Context, sess *Session, p *pdu.EnquireLinkResp) error {
	_, _, err := sess.Send(ctx, p)
	if err != nil {
		return err
	}
	return nil
}

// SendSubmitMulti is a helper function for sending SubmitMulti PDU.
func SendSubmitMulti(ctx context.Context, sess *Session, p *pdu.SubmitMulti) (*pdu.SubmitMultiResp, error) {
	var tresp *pdu.SubmitMultiResp
	_, resp, err := sess.Send(ctx, p)
	if resp != nil {
		tresp = resp.(*pdu.SubmitMultiResp)
	}
	if err != nil {
		return tresp, err
	}
	return tresp, nil
}

// SendSubmitMultiResp is a helper function for sending SubmitMultiResp PDU.
func SendSubmitMultiResp(ctx context.Context, sess *Session, p *pdu.SubmitMultiResp) error {
	_, _, err := sess.Send(ctx, p)
	if err != nil {
		return err
	}
	return nil
}

// SendAlertNotification is a helper function for sending AlertNotification PDU.
func SendAlertNotification(ctx context.Context, sess *Session, p *pdu.AlertNotification) error {
	_, _, err := sess.Send(ctx, p)
	if err != nil {
		return err
	}
	return nil
}

// SendDataSm is a helper function for sending DataSm PDU.
func SendDataSm(ctx context.Context, sess *Session, p *pdu.DataSm) (*pdu.DataSmResp, error) {
	var tresp *pdu.DataSmResp
	_, resp, err := sess.Send(ctx, p)
	if resp != nil {
		tresp = resp.(*pdu.DataSmResp)
	}
	if err != nil {
		return tresp, err
	}
	return tresp, nil
}

// SendDataSmResp is a helper function for sending DataSmResp PDU.
func SendDataSmResp(ctx context.Context, sess *Session, p *pdu.DataSmResp) error {
	_, _, err := sess.Send(ctx, p)
	if err != nil {
		return err
	}
	return nil
}
