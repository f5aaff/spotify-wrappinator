package device

import (
	"encoding/json"
	"errors"
	"strconv"
	agent "wrappinator.agent"
	requests "wrappinator.requests"
)

const (
	BaseURL string = "https://api.spotify.com/v1/"
)

type DevList struct {
	Devices []Device `json:"devices"`
}

// GetCurrentDevice
//makes a request for a list of the available devices, returns a device struct of the active device
//marshals spotify response into struct containing list of devices, changes method caller to response from request./**
func (d *Device) GetCurrentDevice(a *agent.Agent) error {
	c := requests.New(requests.WithRequestURL("me/player/devices"), requests.WithBaseURL(BaseURL))
	requests.GetRequest(a, c)
	var list DevList
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

func (d *Device) ChangeVolume(a *agent.Agent, c *requests.ClientRequest, increment int) error {

	c = requests.New(requests.WithRequestURL("me/player/volume"), requests.WithBaseURL(BaseURL))
	newVolume := strconv.Itoa(d.VolumePercent + increment)
	requests.PutRequest(a, c, requests.Fields("volume_percent", newVolume))
	err := errors.New("")
	d.VolumePercent, err = strconv.Atoi(newVolume)
	if err != nil {
		err = errors.New(err.Error() + c.RequestURL)
		return err
	}
	return nil
}
