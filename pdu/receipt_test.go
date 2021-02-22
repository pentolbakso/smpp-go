package pdu

import (
	"errors"
	"regexp"
	"strings"
	"testing"
	"time"
)

func TestParsingGoodDeliveryReceipt(t *testing.T) {
	good := "id:123123123 sub:0 dlvrd:0 submit date:1507011202 done date:1507011101 stat:DELIVRD err:0 text:Test information"
	dr, err := ParseDeliveryReceipt(good)
	if err != nil {
		t.Errorf("Error parsing good receipt %s", err)
	}
	if dr.Id != "123123123" {
		t.Errorf("Receipt id is wrong %s expected 123123123", dr.Id)
	}
	extime, _ := time.Parse(time.RFC3339, "2015-07-01T12:02:00Z")
	if dr.SubmitDate != extime {
		t.Errorf("Receipt submit date is wrong %s expected %s",
			dr.SubmitDate.Format(time.RFC3339),
			extime.Format(time.RFC3339),
		)
	}
	if dr.String() != good {
		t.Errorf("Receipt string representation is wrong %s", dr)
	}
}

func TestParsingBadDeliveryReceipt(t *testing.T) {
	keys := "id:123123123 dfdfsub:0 dlvrd:0 submit date:1507011202 done date:1507011101 stat:DELIVRD err:0 text:Test information"
	_, err := ParseDeliveryReceipt(keys)
	if err == nil {
		t.Errorf("Parsing bad receipt with wrong key name returned no error")
	}
	missingkeys := "id:123123123 sub:0 dlvrd:0 submit date:1507011202 stat:DELIVRD err:0 text:Test information"
	_, err = ParseDeliveryReceipt(missingkeys)
	if err == nil {
		t.Errorf("Parsing bad receipt with missing keys returned no error")
	}
	date := "id:123123123 sub:0 dlvrd:0 submit date:150701adsfas1202 done date:1507011101 stat:DELIVRD err:0 text:Test information"
	_, err = ParseDeliveryReceipt(date)
	if err == nil {
		t.Errorf("Parsing bad receipt with wrong date format returned no error")
	}
}

func TestParsingUUIDDeliveryReceipt(t *testing.T) {
	dlr := "id:a03ea27b-9bb4-4d5e-b87f-3f578ab46153 sub:001 dlvrd:001 submit date:161003211236 done date:161003211236 stat:DELIVRD err:000 text:-"
	r, err := ParseDeliveryReceipt(dlr)
	if err != nil {
		t.Fatalf("Error parsing UUID delivery receipt %v", err)
	}
	if r.Id != "a03ea27b-9bb4-4d5e-b87f-3f578ab46153" {
		t.Errorf("ParseDeliveryReceipt() => %s expected %s", r.Id, "a03ea27b-9bb4-4d5e-b87f-3f578ab46153")
	}
	if r.Stat != "DELIVRD" {
		t.Errorf("ParseDeliveryReceipt() => %s expected %s", r.Stat, "DELIVRD")
	}
}

func BenchmarkParseDeliveryReceipt(b *testing.B) {
	good := "id:123123123 sub:0 dlvrd:0 submit date:1507011202 done date:1507011101 stat:DELIVRD err:0 text:Test information"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ParseDeliveryReceipt(good)
	}
	b.StopTimer()
}
func BenchmarkParseDeliveryReceipt01(b *testing.B)  {
	good := "id:123123123 sub:0 dlvrd:0 submit date:1507011202 done date:1507011101 stat:DELIVRD err:0 text:Test information"
	var rule = regexp.MustCompile(`(\w+ ?\w+)+:([\w\-]+)`)
	type DeliveryReceipt struct {
		Id         string
		Sub        string
		Dlvrd      string
		SubmitDate time.Time
		DoneDate   time.Time
		Stat       DeliveryStat
		Err        string
		Text       string

	}
	f:= func(sm string)(*DeliveryReceipt, error){
		e := errors.New("smpp: invalid receipt format")
		i := strings.Index(sm, "text:")
		if i == -1 {
			i = strings.Index(sm, "Text:")
			if i == -1 {
				return nil, e
			}
		}
		delRec := DeliveryReceipt{}
		match := rule.FindAllStringSubmatch(sm[:i], -1)
		for idx, m := range match {
			if len(m) != 3 {
				return nil, e
			}
			// TODO improve error with more details
			switch idx {
			case 0:
				if m[1] != "id" {
					return nil, e
				}
				delRec.Id = m[2]
			case 1:
				if m[1] != "sub" {
					return nil, e
				}
				delRec.Sub = m[2]
			case 2:
				if m[1] != "dlvrd" {
					return nil, e
				}
				delRec.Dlvrd = m[2]
			case 3:
				if m[1] != "submit date" {
					return nil, e
				}
				t, err := time.Parse(recDateLayout, m[2])
				if err != nil {
					t, err = time.Parse(secRecDateLayout, m[2])
					if err != nil {
						return nil, e
					}
				}
				delRec.SubmitDate = t
			case 4:
				if m[1] != "done date" {
					return nil, e
				}
				t, err := time.Parse(recDateLayout, m[2])
				if err != nil {
					t, err = time.Parse(secRecDateLayout, m[2])
					if err != nil {
						return nil, e
					}
				}
				delRec.DoneDate = t
			case 5:
				if m[1] != "stat" {
					return nil, e
				}
				// TODO validate status value
				delRec.Stat = DeliveryStat(m[2])
			case 6:
				if m[1] != "err" {
					return nil, e
				}
				delRec.Err = m[2]
			default:
				return nil, e
			}
		}
		delRec.Text = sm[i+5:]
		return &delRec, nil
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f(good)
	}
	b.StopTimer()
}
