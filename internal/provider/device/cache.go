package device

import (
	"context"
	"errors"
	"github.com/richseviora/huego/pkg/resources/behavior_instance"
	"github.com/richseviora/huego/pkg/resources/behavior_script"
	"github.com/richseviora/huego/pkg/resources/client"
	"github.com/richseviora/huego/pkg/resources/device"
	"github.com/richseviora/huego/pkg/resources/light"
	"github.com/richseviora/huego/pkg/resources/motion"
	"github.com/richseviora/huego/pkg/resources/room"
	"github.com/richseviora/huego/pkg/resources/scene"
	"github.com/richseviora/huego/pkg/resources/zigbee_connectivity"
	"github.com/richseviora/huego/pkg/resources/zone"
	"slices"
	"sync"
)

type ResultCache struct {
	cache map[string]interface{}
	mutex *sync.Mutex
}

func NewResultCache() *ResultCache {
	return &ResultCache{
		cache: map[string]interface{}{},
		mutex: &sync.Mutex{},
	}
}

type ClientWithCache struct {
	client              client.HueServiceClient
	deviceCache         map[string]DeviceMappingEntry
	behaviorScriptCache *ResultCache
	zigbeeErrors        []zigbee_connectivity.Data
	cacheBuilt          bool
	mutex               *sync.Mutex
}

type ClientWithLightIDCache interface {
	client.HueServiceClient
	GetLightIDForMacAddress(macAddress string) (string, error)
	GetMotionIDForMacAddress(macAddress string) (string, error)
	GetBehaviorScriptIDForMetadataName(name string) (string, error)
}

func NewClientWithCache(client client.HueServiceClient) *ClientWithCache {
	return &ClientWithCache{
		client:              client,
		behaviorScriptCache: NewResultCache(),
		deviceCache:         make(map[string]DeviceMappingEntry),
		mutex:               &sync.Mutex{},
	}
}

func (c *ClientWithCache) buildDeviceMap(ctx context.Context) (map[string]DeviceMappingEntry, []zigbee_connectivity.Data, error) {
	devices, err := c.client.DeviceService().GetAllDevices(ctx)
	if err != nil {
		return nil, nil, err
	}
	zigbees, err := c.client.ZigbeeConnectivityService().GetAllZigbeeConnectivity(ctx)
	if err != nil {
		return nil, nil, err
	}

	deviceMap := make(map[string]DeviceMappingEntry)

	for _, d := range devices.Data {
		entry := DeviceMappingEntry{
			DeviceID: d.ID,
			Name:     d.Metadata.Name,
		}

		for _, service := range d.Services {
			switch service.Rtype {
			case "light":
				entry.LightID = service.Rid
			case "zigbee_connectivity":
				entry.ZigbeeConnectivityID = service.Rid
			case "motion":
				entry.MotionID = service.Rid
			}
		}
		deviceMap[d.ID] = entry
	}
	zigbeeEntries := make([]zigbee_connectivity.Data, 0)
	for _, zigbee := range zigbees.Data {
		deviceEntry, ok := deviceMap[zigbee.Owner.RID]
		if ok {
			deviceEntry.MacAddress = zigbee.MacAddress
			deviceMap[zigbee.Owner.RID] = deviceEntry
		} else {
			zigbeeEntries = append(zigbeeEntries, zigbee)
		}
	}
	return deviceMap, zigbeeEntries, nil
}

func (c *ClientWithCache) GetBehaviorScriptIDForMetadataName(name string) (string, error) {
	c.behaviorScriptCache.mutex.Lock()
	defer c.behaviorScriptCache.mutex.Unlock()

	if cached, ok := c.behaviorScriptCache.cache[name]; ok {
		if script, ok := cached.(behavior_script.Data); ok {
			return script.ID, nil
		}
	}

	scripts, err := c.client.BehaviorScriptService().GetAllBehaviorScripts(context.Background())
	if err != nil {
		return "", err
	}

	var result behavior_script.Data
	for _, script := range scripts.Data {
		c.behaviorScriptCache.cache[script.Metadata.Name] = script
		if script.Metadata.Name == name {
			result = script
		}
	}
	if result.ID != "" {
		return result.ID, nil
	}

	return "", errors.New("could not find behavior script with name: " + name)
}

func (c *ClientWithCache) GetLightIDForMacAddress(macAddress string) (string, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	_, _, _, err := c.buildCache()
	if err != nil {
		return "", err
	}
	for _, d := range c.deviceCache {
		if d.MacAddress == macAddress {
			return d.LightID, nil
		}
	}
	return "", errors.New("could not find Mac Address in cache: " + macAddress + "")
}

func (c *ClientWithCache) GetMotionIDForMacAddress(macAddress string) (string, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	_, _, _, err := c.buildCache()
	if err != nil {
		return "", err
	}
	for _, d := range c.deviceCache {
		if d.MacAddress == macAddress {
			return d.MotionID, nil
		}
	}
	return "", errors.New("could not find Mac Address in cache: " + macAddress + "")
}

func (c *ClientWithCache) GetAllDevices() ([]DeviceMappingEntry, []zigbee_connectivity.Data, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	devices, entries, data, err := c.buildCache()
	if err != nil {
		return entries, data, err
	}
	for _, d := range c.deviceCache {
		devices = append(devices, d)
	}
	slices.SortFunc(devices, func(i, j DeviceMappingEntry) int {
		if i.Name < j.Name {
			return -1
		} else if i.Name > j.Name {
			return 1
		}
		return 0
	})
	return devices, c.zigbeeErrors, nil
}

func (c *ClientWithCache) buildCache() ([]DeviceMappingEntry, []DeviceMappingEntry, []zigbee_connectivity.Data, error) {
	if !c.cacheBuilt {
		deviceMap, zigbeeErrors, err := c.buildDeviceMap(context.Background())
		if err != nil {
			return nil, nil, nil, err
		}
		c.zigbeeErrors = zigbeeErrors
		c.deviceCache = deviceMap
		c.cacheBuilt = true
	}
	devices := make([]DeviceMappingEntry, 0)
	return devices, nil, nil, nil
}

var (
	_ client.HueServiceClient = &ClientWithCache{}
	_ ClientWithLightIDCache  = &ClientWithCache{}
)

// region Services
func (c *ClientWithCache) ZoneService() zone.ZoneService {
	return c.client.ZoneService()
}

func (c *ClientWithCache) RoomService() room.RoomService {
	return c.client.RoomService()
}

func (c *ClientWithCache) SceneService() scene.SceneService {
	return c.client.SceneService()
}

func (c *ClientWithCache) LightService() light.LightService {
	return c.client.LightService()
}

func (c *ClientWithCache) DeviceService() device.Service {
	return c.client.DeviceService()
}

func (c *ClientWithCache) ZigbeeConnectivityService() zigbee_connectivity.Service {
	return c.client.ZigbeeConnectivityService()
}

func (c *ClientWithCache) BehaviorInstanceService() behavior_instance.Service {
	return c.client.BehaviorInstanceService()
}

func (c *ClientWithCache) MotionService() motion.Service {
	return c.client.MotionService()
}

func (c *ClientWithCache) BehaviorScriptService() behavior_script.Service {
	return c.client.BehaviorScriptService()
}

//endregion
