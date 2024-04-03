package spotify_wrappinator.device

type Device struct {
	ID               string `json:"id"`
	IsActive         bool   `json:"is_active"`
	IsPrivateSession bool   `json:"is_private_session"`
	IsRestricted     bool   `json:"is_restricted"`
	Name             string `json:"name"`
	SupportsVolume   bool   `json:"supports_volume"`
	Type             string `json:"type"`
	VolumePercent    int    `json:"volume_percent"`
}

type Opt func(device *Device)

func withId(s string) Opt {
	return func(d *Device) {
		d.ID = s
	}
}

func withIsActive(s bool) Opt {
	return func(d *Device) {
		d.IsActive = s
	}
}
func withIsPrivateSession(s bool) Opt {
	return func(d *Device) {
		d.IsPrivateSession = s
	}
}
func withName(s string) Opt {
	return func(d *Device) {
		d.Name = s
	}
}
func withDeviceType(s string) Opt {
	return func(d *Device) {
		d.Type = s
	}
}
func withVolumePercent(s int) Opt {
	return func(d *Device) {
		d.VolumePercent = s
	}
}
func withSupportsVolume(s bool) Opt {
	return func(d *Device) {
		d.SupportsVolume = s
	}
}

func New(deviceOpts ...Opt) *Device {
	d := &Device{}
	for _, opt := range deviceOpts {
		opt(d)
	}
	return d
}
