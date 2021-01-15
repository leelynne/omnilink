package home

import (
	"fmt"
	"log"
	"time"

	"github.com/leelynne/omnilink/omni"
)

func New(logger *log.Logger, endpoint, key string) (*Home, error) {
	c, err := omni.NewClient(fmt.Sprintf("%s:4369", endpoint), key)
	if err != nil {
		return nil, err
	}

	h := Home{}

	si, err := c.GetSystemInformation()
	if err != nil {
		panic(err)
	}

	h.ModelNumber = int(si.ModelNumber)
	switch h.ModelNumber {
	case 16:
		h.ModelName = "HAI OmniPro II"
	case 30:
		h.ModelName = "HAI Omni IIe"
	case 36:
		h.ModelName = "HAI Lumina"
	case 37:
		h.ModelName = "HAI Lumina Pro"
	}

	h.Version = fmt.Sprintf("%d.%d", si.MajorVersion, si.MinorVersion)
	h.PhoneNumber = string(si.LocalPhoneNumber[:])

	sf, err := c.GetSystemFeatures()
	if err != nil {
		return nil, err
	}
	h.features = sf.Features

	ss, err := c.GetSystemStatus()
	if err != nil {
		return nil, err
	}

	st := Status{}
	if int(ss.DateValid) > 0 {
		st.DateSet = true
		st.Date = time.Date(2000+int(ss.Year), time.Month(ss.Month), int(ss.Day), int(ss.Hour), int(ss.Minute), int(ss.Second), 0, time.UTC)
	}
	st.Battery = ss.Battery
	h.latestStatus = st
	return &h, nil
}
