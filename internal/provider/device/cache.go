package device

import (
	"context"
	"errors"
	"github.com/richseviora/huego/pkg/resources/client"
	"github.com/richseviora/huego/pkg/resources/device"
	"github.com/richseviora/huego/pkg/resources/light"
	"github.com/richseviora/huego/pkg/resources/room"
	"github.com/richseviora/huego/pkg/resources/scene"
	"github.com/richseviora/huego/pkg/resources/zigbee_connectivity"
	"github.com/richseviora/huego/pkg/resources/zone"
)

type DeviceMappingEntry struct {
	DeviceID             string
	LightID              string
	ZigbeeConnectivityID string
	MacAddress           string
}

type ClientWithCache struct {
	client     client.HueServiceClient
	cache      map[string]DeviceMappingEntry
	cacheBuilt bool
}

func NewClientWithCache(client client.HueServiceClient) *ClientWithCache {
	return &ClientWithCache{
		client: client,
		cache:  make(map[string]DeviceMappingEntry),
	}
}

func (c *ClientWithCache) BuildDeviceMap(ctx context.Context) (map[string]DeviceMappingEntry, error) {
	devices, err := c.client.DeviceService().GetAllDevices(ctx)
	if err != nil {
		return nil, err
	}
	zigbees, err := c.client.ZigbeeConnectivityService().GetAllZigbeeConnectivity(ctx)
	if err != nil {
		return nil, err
	}

	deviceMap := make(map[string]DeviceMappingEntry)
	for _, d := range devices.Data {
		entry := DeviceMappingEntry{
			DeviceID: d.ID,
		}

		for _, service := range d.Services {
			switch service.Rtype {
			case "light":
				entry.LightID = service.Rid
			case "zigbee_connectivity":
				entry.ZigbeeConnectivityID = service.Rid
			}
		}
		deviceMap[d.ID] = entry
	}
	for _, zigbee := range zigbees.Data {
		deviceEntry, ok := deviceMap[zigbee.Owner.RID]
		if ok {
			deviceEntry.MacAddress = zigbee.MacAddress
		}
	}
	return deviceMap, nil
}

func (c *ClientWithCache) GetLightIDForMacAddress(macAddress string) (string, error) {
	if !c.cacheBuilt {
		deviceMap, err := c.BuildDeviceMap(context.Background())
		if err != nil {
			return "", err
		}
		c.cache = deviceMap
		c.cacheBuilt = true
	}
	for _, d := range c.cache {
		if d.MacAddress == macAddress {
			return d.LightID, nil
		}
	}
	return "", errors.New("could not find Mac Address in cache: " + macAddress + "")
}

var (
	_ client.HueServiceClient = &ClientWithCache{}
)

func (c ClientWithCache) ZoneService() zone.ZoneService {
	return c.client.ZoneService()
}

func (c ClientWithCache) RoomService() room.RoomService {
	return c.client.RoomService()
}

func (c ClientWithCache) SceneService() scene.SceneService {
	return c.client.SceneService()
}

func (c ClientWithCache) LightService() light.LightService {
	return c.client.LightService()
}

func (c ClientWithCache) DeviceService() device.Service {
	return c.client.DeviceService()
}

func (c ClientWithCache) ZigbeeConnectivityService() zigbee_connectivity.Service {
	return c.client.ZigbeeConnectivityService()
}

type ClientWithLightIDCache interface {
	client.HueServiceClient
	GetLightIDForMacAddress(macAddress string) (string, error)
}
