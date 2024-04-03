package wrappinator.device

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
//marshals spotify response into struct containing list owf devices, changes method caller to response from request./**
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

func (d *Device) PlayPause(a *agent.Agent, action string) error {
	actions := []string{"pause", "play", "next", "previous"}

	for _, v := range actions {
		if action == v {
			c := requests.New(requests.WithRequestURL("me/player/"+action), requests.WithBaseURL(BaseURL))
			requests.PutRequest(a, c)
			return nil
		}
	}
	return errors.New("invalid action")

}

func (d *Device) GetCurrentlyPlaying(a *agent.Agent) (string, error) {
	c := requests.New(requests.WithBaseURL(BaseURL), requests.WithRequestURL("me/player/currently-playing"))
	requests.GetRequest(a, c)
	if c.Response != nil {
		return string(c.Response), nil
	}
	return "", errors.New("error retrieving song")
}

func (d *Device) GetQueue(a *agent.Agent) (string, error) {
	c := requests.New(requests.WithBaseURL(BaseURL), requests.WithRequestURL("me/player/queue"))
	requests.GetRequest(a, c)
	if c.Response != nil {
		return string(c.Response), nil
	}
	return "", errors.New("could not retrieve queue")
}

func CustomField(field string, val string) requests.RequestOption {
	return func(reqOpt *requests.RequestOptions) {
		reqOpt.UrlParams.Set(field, val)
	}
}

func (d *Device) PlayCustom(a *agent.Agent, contextUri *string, position *int, position_ms *int) error {

	m := map[string]interface{}{"context_uri": contextUri, "position": position, "position_ms": position_ms}

	var reqopts []requests.RequestOption
	i := 0
	for k, v := range m {
		switch v.(type) {
		case int:
			reqopts[i] = CustomField(k, strconv.Itoa(v.(int)))
		case string:
			reqopts[i] = CustomField(k, v.(string))
		}
		i++
	}

	c := requests.New(requests.WithBaseURL(BaseURL), requests.WithRequestURL("me/player/play"))
	requests.PutRequest(a, c, reqopts...)
	if c.Response == nil || len(c.Response) == 0 {
		return errors.New("error playing song")
	}

	return nil
}
