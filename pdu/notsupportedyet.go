package pdu

import (
	"fmt"
)

// ReplaceSm Not supported yet.
type ReplaceSm struct{}

func (p ReplaceSm) CommandID() CommandID {
	return ReplaceSmID
}

func (p ReplaceSm) MarshalBinary() ([]byte, error) {
	return nil, fmt.Errorf("Command %s is not supported yet", p.CommandID())
}

func (p *ReplaceSm) UnmarshalBinary(body []byte) error {
	return fmt.Errorf("Command %s is not supported yet", p.CommandID())
}

// ReplaceSmResp Not supported yet.
type ReplaceSmResp struct{}

func (p ReplaceSmResp) CommandID() CommandID {
	return ReplaceSmRespID
}

func (p ReplaceSmResp) MarshalBinary() ([]byte, error) {
	return nil, fmt.Errorf("Command %s is not supported yet", p.CommandID())
}

func (p *ReplaceSmResp) UnmarshalBinary(body []byte) error {
	return fmt.Errorf("Command %s is not supported yet", p.CommandID())
}

// CancelSm Not supported yet.
type CancelSm struct{}

func (p CancelSm) CommandID() CommandID {
	return CancelSmID
}

func (p CancelSm) MarshalBinary() ([]byte, error) {
	return nil, fmt.Errorf("Command %s is not supported yet", p.CommandID())
}

func (p *CancelSm) UnmarshalBinary(body []byte) error {
	return fmt.Errorf("Command %s is not supported yet", p.CommandID())
}

// CancelSmResp Not supported yet.
type CancelSmResp struct{}

func (p CancelSmResp) CommandID() CommandID {
	return CancelSmRespID
}

func (p CancelSmResp) MarshalBinary() ([]byte, error) {
	return nil, fmt.Errorf("Command %s is not supported yet", p.CommandID())
}

func (p *CancelSmResp) UnmarshalBinary(body []byte) error {
	return fmt.Errorf("Command %s is not supported yet", p.CommandID())
}

// Outbind Not supported yet.
type Outbind struct{}

func (p Outbind) CommandID() CommandID {
	return OutbindID
}

func (p Outbind) MarshalBinary() ([]byte, error) {
	return nil, fmt.Errorf("Command %s is not supported yet", p.CommandID())
}

func (p *Outbind) UnmarshalBinary(body []byte) error {
	return fmt.Errorf("Command %s is not supported yet", p.CommandID())
}

// SubmitMulti Not supported yet.
type SubmitMulti struct{}

func (p SubmitMulti) CommandID() CommandID {
	return SubmitMultiID
}

func (p SubmitMulti) MarshalBinary() ([]byte, error) {
	return nil, fmt.Errorf("Command %s is not supported yet", p.CommandID())
}

func (p *SubmitMulti) UnmarshalBinary(body []byte) error {
	return fmt.Errorf("Command %s is not supported yet", p.CommandID())
}

// SubmitMultiResp Not supported yet.
type SubmitMultiResp struct{}

func (p SubmitMultiResp) CommandID() CommandID {
	return SubmitMultiRespID
}

func (p SubmitMultiResp) MarshalBinary() ([]byte, error) {
	return nil, fmt.Errorf("Command %s is not supported yet", p.CommandID())
}

func (p *SubmitMultiResp) UnmarshalBinary(body []byte) error {
	return fmt.Errorf("Command %s is not supported yet", p.CommandID())
}

// AlertNotification Not supported yet.
type AlertNotification struct{}

func (p AlertNotification) CommandID() CommandID {
	return AlertNotificationID
}

func (p AlertNotification) MarshalBinary() ([]byte, error) {
	return nil, fmt.Errorf("Command %s is not supported yet", p.CommandID())
}

func (p *AlertNotification) UnmarshalBinary(body []byte) error {
	return fmt.Errorf("Command %s is not supported yet", p.CommandID())
}

// DataSm Not supported yet.
type DataSm struct {
	ServiceType        string
	SourceAddrTon      int
	SourceAddrNpi      int
	SourceAddr         string
	DestAddrTon        int
	DestAddrNpi        int
	DestinationAddr    string
	EsmClass           EsmClass
	RegisteredDelivery RegisteredDelivery
	DataCoding         int
	Options            *Options
}

func (p DataSm) CommandID() CommandID {
	return DataSmID
}

func (p DataSm) MarshalBinary() ([]byte, error) {
	out := append(
		[]byte(p.ServiceType),
		0,
		byte(p.SourceAddrTon),
		byte(p.SourceAddrNpi),
	)
	out = append(out, append([]byte(p.SourceAddr), 0)...)
	out = append(out, byte(p.DestAddrTon), byte(p.DestAddrNpi))
	out = append(out, append([]byte(p.DestinationAddr), 0)...)
	out = append(out, p.EsmClass.Byte(), p.RegisteredDelivery.Byte(), byte(p.DataCoding))
	if p.Options == nil {
		return out, nil
	}
	opts, err := p.Options.MarshalBinary()
	if err != nil {
		return nil, err
	}
	return append(out, opts...), nil
}

func (p *DataSm) UnmarshalBinary(body []byte) error {
	fmt.Println(body)
	return fmt.Errorf("Command %s is not supported yet", p.CommandID())
}

// DataSmResp Not supported yet.
type DataSmResp struct {
	MessageID string
	Options   *Options
}

func (p DataSmResp) CommandID() CommandID {
	return DataSmRespID
}

func (p DataSmResp) MarshalBinary() ([]byte, error) {
	return nil, fmt.Errorf("Command %s is not supported yet", p.CommandID())
}

func (p *DataSmResp) UnmarshalBinary(body []byte) error {
	var err error
	p.MessageID, p.Options, err = cStringOptsRespUnmarshal(body)
	return err
}
