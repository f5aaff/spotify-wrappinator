package device

type Device struct {
	id                 string
	is_active          string
	is_private_session string
	name               string
	device_type        string
	volume_percent     string
	supports_volume    string
}

type DeviceOpt func(device *Device)

func withId(s string) DeviceOpt {
	return func(d *Device) {
		d.id = s
	}
}

func withIs_active(s string) DeviceOpt {
	return func(d *Device) {
		d.is_active = s
	}
}
func withIs_private_session(s string) DeviceOpt {
	return func(d *Device) {
		d.is_private_session = s
	}
}
func withName(s string) DeviceOpt {
	return func(d *Device) {
		d.name = s
	}
}
func withDevice_type(s string) DeviceOpt {
	return func(d *Device) {
		d.device_type = s
	}
}
func withVolume_Percent(s string) DeviceOpt {
	return func(d *Device) {
		d.volume_percent = s
	}
}
func withSupports_volume(s string) DeviceOpt {
	return func(d *Device) {
		d.supports_volume = s
	}
}

func New(deviceOpts ...DeviceOpt) *Device {
	d := &Device{}
	for _, opt := range deviceOpts {
		opt(d)
	}
	return d
}
