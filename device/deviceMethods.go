package device

import (
	"strconv"
	agent "wrappinator.agent"
	requests "wrappinator.requests"
)

type funcParam []any
type Func func(device *Device)

func changeVolume(increment int, a *agent.Agent) Func {
	return func(d *Device) {
		currentVol, _ := strconv.Atoi(d.volume_percent)
		newVol := currentVol + increment
		req := requests.New(requests.WithRequestURL(""), requests.WithBaseURL(""))
		requests.ParamRequest(a, req, requests.Fields("volume", string(rune(newVol))), requests.Fields("device_id", d.id))
	}
}

func getPlaybackState(a *agent.Agent, c *requests.ClientRequest) Func {
	return func(d *Device) {
		c = requests.New(requests.WithRequestURL(""), requests.WithBaseURL(""))
		requests.GetRequest(a, c)
	}
}
