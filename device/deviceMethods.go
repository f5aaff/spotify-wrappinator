package device

import (
	"encoding/json"
	agent "wrappinator.agent"
	requests "wrappinator.requests"
)

const (
	BaseURL string = "https://api.spotify.com/v1/"
)

type Func func(device *Device)

func changeVolume(increment int, a *agent.Agent) Func {
	return func(d *Device) {
		newVol := d.VolumePercent + increment
		req := requests.New(requests.WithRequestURL("me/player/"), requests.WithBaseURL(BaseURL))
		requests.ParamRequest(a, req, requests.Fields("volume", string(rune(newVol))), requests.Fields("device_id", d.ID))
	}
}

type AutoGenerated struct {
	Devices []Device `json:"devices"`
}

// GetCurrentDevice
//makes a request for a list of the available devices, returns a device struct of the active device
//marshals spotify response into struct containing list of devices, changes method caller to response from request./**
func (d *Device) GetCurrentDevice(a *agent.Agent) error {
	c := requests.New(requests.WithRequestURL("me/player/devices"), requests.WithBaseURL(BaseURL))
	requests.GetRequest(a, c)
	var list AutoGenerated
	err := json.Unmarshal(c.Response, &list)
	if err != nil {
		return err
	}
	for _, x := range list.Devices {
		if x.IsActive == true {
			*d = x
			return nil
		}
	}
	*d = list.Devices[0]
	return nil
}
