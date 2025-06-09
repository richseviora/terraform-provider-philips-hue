package device

import (
	"fmt"
	"github.com/richseviora/huego/pkg/resources/zigbee_connectivity"
	"strings"
)

func formatName(name string) string {
	lowerName := strings.ToLower(name)
	return strings.ReplaceAll(lowerName, " ", "_")
}

func GenerateImportOutput(entries []DeviceMappingEntry, missingEntries []zigbee_connectivity.Data) string {
	resourceResult := ""
	result := ""

	for _, entry := range entries {
		if !(entry.IsLight() || entry.IsMotion()) {
			continue
		}
		newResult, newResourceResult := generateEntryOutput(entry)

		result += newResult
		resourceResult += newResourceResult
	}

	for _, entry := range missingEntries {
		result += fmt.Sprintf(`
/* 
Could not resolve MAC address:
%+v
*/
`, entry)
	}
	return result + resourceResult
}

func generateEntryOutput(entry DeviceMappingEntry) (string, string) {
	if entry.IsLight() {
		return generateLightOutput(entry)
	} else if entry.IsMotion() {
		return generateMotionOutput(entry)
	}
	return "", ""
}

func generateLightOutput(entry DeviceMappingEntry) (string, string) {
	formattedName := formatName(entry.Name)
	result := fmt.Sprintf("\nimport {\n  # Name = %s\n  id = \"%s\"\n  to = philips_light.%s\n}\n", entry.Name, entry.MacAddress, formattedName)
	resourceResult := fmt.Sprintf(`
resource philips_light "%s" {
  name = "%s"
  type = "decorative"
}
`, formattedName, entry.Name)
	return result, resourceResult
}

func generateMotionOutput(entry DeviceMappingEntry) (string, string) {
	formattedName := formatName(entry.Name)
	result := fmt.Sprintf("\nimport {\n  # Name = %s\n  id = \"%s\"\n  to = philips_motion.%s\n}\n", entry.Name, entry.MacAddress, formattedName)
	resourceResult := fmt.Sprintf(`
resource philips_motion "%s" {
  enabled = true
}
`, formattedName)
	return result, resourceResult
}
