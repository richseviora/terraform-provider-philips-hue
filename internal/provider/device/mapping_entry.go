package device

type DeviceMappingEntry struct {
	Name                 string
	DeviceID             string
	LightID              string
	ZigbeeConnectivityID string
	MotionID             string
	MacAddress           string
}

func (d DeviceMappingEntry) IsLight() bool {
	return d.LightID != ""
}

func (d DeviceMappingEntry) IsMotion() bool {
	return d.MotionID != ""
}
