package plugin

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"
	"os"

	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
)

// DHCPNetworkOptions contains options for the DHCP network driver
type DHCPNetworkOptions struct {
	Bridge          string
	IPv6            bool
	LeaseTimeout    time.Duration `mapstructure:"lease_timeout"`
	IgnoreConflicts bool          `mapstructure:"ignore_conflicts"`
	SkipRoutes      bool          `mapstructure:"skip_routes"`
}

func decodeOpts(input interface{}) (DHCPNetworkOptions, error) {
	var opts DHCPNetworkOptions
	optsDecoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:           &opts,
		ErrorUnused:      true,
		WeaklyTypedInput: true,
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			mapstructure.StringToTimeDurationHookFunc(),
		),
	})
	if err != nil {
		return opts, fmt.Errorf("failed to create options decoder: %w", err)
	}

	if err := optsDecoder.Decode(input); err != nil {
		return opts, err
	}

	return opts, nil
}

func DeleteNetworkOptions(networkID string) {
	configPath := "/home/" + networkID + ".json"
	os.Remove(configPath)
}

func LoadNetworkOptions(networkID string) (DHCPNetworkOptions, error) {
	result := DHCPNetworkOptions{}

	configPath := "/home/" + networkID + ".json"
	jsonBlob, err := ioutil.ReadFile(configPath)
	if err != nil {
		return result, fmt.Errorf("failed to read network options from file: %w", err)
	}

	err = json.Unmarshal(jsonBlob, &result)
    if err != nil {
		return result, fmt.Errorf("network options invalid: %w", err)
    }

	return result, nil
}

func PersistNetworkOptions(networkID string, opts DHCPNetworkOptions) {
	configPath := "/home/" + networkID + ".json"
	jsonData, _ := json.Marshal(opts)
    err := ioutil.WriteFile(configPath, jsonData, 0640)
	
	if err != nil {
		log.WithField("err", err).Error("Failed to write JSON file with network options")
	} else {
		log.Info("Network options persisted successfully to: " + configPath)
	}
}