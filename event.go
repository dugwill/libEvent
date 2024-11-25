package libevent

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/Comcast/gots/v2"
	"github.com/Comcast/gots/v2/scte35"
)

// Event is use to create a xml representation of the scte35 event
type Event struct {
	Time          time.Time                `xml:"manifestTime" json:"manifestTime,omitempty"`
	StreamName    string                   `xml:"streamName,attr" json:"streamName,omitempty"`
	Origin        string                   `xml:"origin,attr" json:"origin,omitempty"`
	PeriodID      string                   `xml:"periodId,attr" json:"upid,omitempty"`
	EventTime     time.Time                `xml:"eventTime,attr" json:"eventTime,omitempty"`
	EventID       uint32                   `xml:"eventID,attr" json:"eventID,omitempty"`
	PTS           gots.PTS                 `xml:"pts,attr" json:"pts,omitempty"`
	Command       scte35.SpliceCommandType `xml:"command,attr" json:"command,omitempty"`
	TypeID        scte35.SegDescType       `xml:"typeID,attr" json:"typeID,omitempty"`
	Signal        string                   `xml:"signal,attr" json:"signal,omitempty"`
	BreakDuration gots.PTS                 `xml:"duration,attr" json:"duration,omitempty"`
	TimeToSplice  uint64                   `xml:"timeToSplice,attr" json:"timeToSplice,omitempty"`
	EventJSON     []byte                   `xml:"eventJSON,attr" json:"-"`
}

// JMarshalEvent returns a []byte containing the marshalled event
func (e *Event) JMarshalEvent() (s []byte, err error) {

	s, err = json.MarshalIndent(e, "  ", "  ")
	if err != nil {
		return nil, err
	}
	return s, nil
}

// JUnMarshalEvent UnMarshals a []byte into an event struct
func (e *Event) JUnMarshalEvent(s []byte) (err error) {

	err = json.Unmarshal(s, e)
	if err != nil {
		return err
	}
	return nil
}

// NewEvent create an event instance and populates it with the decoded SCTE-35 data
func NewEvent(d []byte) (e *Event, err error) {

	scte, err := scte35.NewSCTE35(d)
	if err != nil {
		return nil, err
	}

	e = &Event{
		PTS:     scte.PTS(),
		Command: scte.Command(),
	}

	for _, desc := range scte.Descriptors() {
		if desc.TypeID() >= 0x30 && desc.TypeID() <= 0x37 {
			e.TypeID = desc.TypeID()
			e.EventID = desc.EventID()
			if desc.HasDuration() {
				e.BreakDuration = desc.Duration()
			}

			// Set the event signal value, acconting for MID structure in descriptor
			if desc.UPIDType() == scte35.SegUPIDADI {
				e.Signal = fmt.Sprintf("%s", desc.UPID())

			} else if desc.UPIDType() == scte35.SegUPIDMID {
				for _, m := range desc.MID() {
					if m.UPIDType() == scte35.SegUPIDADI {
						e.Signal = fmt.Sprintf("%s", m.UPID())
					}
				}
			} else {
				err = fmt.Errorf("could not find signal")
			}

		}
	}

	return e, nil
}

// Store Event writes event to given file
func (e *Event) StoreEvent(f *os.File) {

	jsonString, err := json.Marshal(e)
	if err != nil {
		e := fmt.Sprintf("error marshalling event %v ", err)
		f.WriteString(e)
		return
	}
	f.WriteString(string(jsonString) + "\n")
}
