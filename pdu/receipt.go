package pdu

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type (
	DeliveryStat string
	DeliveryErr  int
)

// DeliveryReceipt in format
//Receipt in short_message field
//Many pre-v3.4 APIs and Message Centers supporting v3.3 are likely to have a means of passing receipt information
//within the short_message field. This applies to MC Delivery Receipts and Intermediate Notifications.
//The format specifics of this information are SMS gateway and SMSC platform specific and beyond the scope of the specification.
//However, the following shows the approach typically taken:
//id:123A456B sub:1 dlvrd:1 submit date:1702281424 done date:1702281424 stat:DELIVRD err:0 text:
type DeliveryReceipt struct {
	Id         string       //The message ID allocated to the message by the SMSC when originally submitted.
	Sub        int          //Number of short messages originally submitted. The value may be padded with leading zeros.
	Dlvrd      int          //Number of short messages delivered. The value may be padded with leading zeros.
	SubmitDate time.Time    //The time and date at which the short message was submitted. In the case of a message which has been replaced, this is the date that the original message was replaced.
	DoneDate   time.Time    //The time and date at which the short message reached it’s final state. The format is the same as for the submit date.
	Stat       DeliveryStat //The final status of the message. See Message states below. State text may be abbreviated.
	Err        DeliveryErr  //A network or SMSC error code for the message. See Error codes below.
	Text       string       //Unused field, result will be blank.
}

const (
	DelStatEnRoute DeliveryStat = "ENROUTE" //The message is in enroute state.
	//This is a general state used to describe a message as being active within the MC.
	//The message may be in retry or dispatched to a mobile network for delivery to the mobile.
	DelStatDelivered DeliveryStat = "DELIVRD" //Message is delivered to destination.
	// The message has been delivered to the destination.
	//No further deliveries will occur.
	DelStatExpired DeliveryStat = "EXPIRED" // Message validity period has expired. The message has failed to be
	// delivered within its validity period and/or retry period.
	//No further delivery attempts will be made.
	DelStatDeleted DeliveryStat = "DELETED" //Message has been deleted. The message has been cancelled or deleted from the MC.
	// No further delivery attempts will take place.
	DelStatUndeliverable DeliveryStat = "UNDELIV"
	DelStatAccepted      DeliveryStat = "ACCEPTD"
	DelStatUnknown       DeliveryStat = "UNKNOWN"
	DelStatRejected      DeliveryStat = "REJECTD"
)

var DelStatMap = map[uint8]DeliveryStat{
	1: DelStatEnRoute,
	2: DelStatDelivered,
	3: DelStatExpired,
	4: DelStatDeleted,
	5: DelStatUndeliverable,
	6: DelStatAccepted,
	7: DelStatUnknown,
	8: DelStatRejected,
}

func (dr *DeliveryReceipt) String() string {

	return fmt.Sprintf(
		"id:%s sub:%d dlvrd:%d submit date:%s done date:%s stat:%s err:%d text:%s",
		dr.Id, dr.Sub, dr.Dlvrd, dr.SubmitDate.Format(recDateLayout), dr.DoneDate.Format(recDateLayout), dr.Stat, dr.Err, dr.Text,
	)
}

var (
	recDateLayout    = "0601021504"
	secRecDateLayout ="060102150405";
)

// ParseDeliveryReceipt parses delivery receipt format defined in smpp 3.4 specification
func ParseDeliveryReceipt(sm string) (*DeliveryReceipt, error) {
	var receipt DeliveryReceipt
	textI := strings.Index(sm, " text:")
	if textI == -1 {
		return &DeliveryReceipt{}, errors.New("smpp: invalid receipt txt ")
	}
	textMsg := sm[textI+6:]
	receipt.Text = textMsg
	formatSm := sm[:textI]
	formatSm = strings.Replace(formatSm, "done date:", "done_date:", 1)
	formatSm = strings.Replace(formatSm, "submit date:", "submit_date:", 1)
	receiptField := strings.Split(formatSm, " ")
	i := -1
	if len(receiptField) < 7 {
		return &DeliveryReceipt{}, errors.New("smpp: receipt miss key ")

	}
	for _, fieldWithValue := range receiptField {
		i = strings.Index(fieldWithValue, ":")
		if i == -1 {
			return &DeliveryReceipt{}, errors.New("smpp: invalid receipt format field " + fieldWithValue)

		}
		switch fieldWithValue[:i] {
		case "id":
			receipt.Id = fieldWithValue[i+1:]
		case "sub", "dlvrd":
			count, err := strconv.Atoi(fieldWithValue[i+1:])
			if err != nil {
				return &DeliveryReceipt{}, err
			}
			if fieldWithValue[:i] == "sub" {
				receipt.Sub = count
			} else {
				receipt.Dlvrd = count
			}

		case "submit_date", "done_date":

			date, err := time.Parse(recDateLayout, fieldWithValue[i+1:])
			if err != nil {
				date, err = time.Parse(secRecDateLayout, fieldWithValue[i+1:])
				if err != nil {
					return &DeliveryReceipt{}, err
				}

			}
			if fieldWithValue[:i] == "submit_date" {
				receipt.SubmitDate = date
			} else {
				receipt.DoneDate = date
			}
		case "stat":
			receipt.Stat = DeliveryStat(fieldWithValue[i+1:])
		case "err":
			count, err := strconv.Atoi(fieldWithValue[i+1:])
			if err != nil {
				return &DeliveryReceipt{}, err

			}
			receipt.Err = DeliveryErr(count)
		default:
			return &DeliveryReceipt{}, errors.New("smpp: invalid receipt format field " + fieldWithValue)
		}

	}
	return &receipt, nil
}
