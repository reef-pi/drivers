package ads1x15

import (
	"sync"
	"time"

	"github.com/reef-pi/hal"
)

type ads1115Factory struct {
	ads1X15Factory
}

var factory1115 *ads1115Factory
var once1115 sync.Once

// Ads1115Factory returns a singleton ADS1015 factory
func Ads1115Factory() hal.DriverFactory {

	once1115.Do(func() {
		factory1115 = &ads1115Factory{
			ads1X15Factory{
				meta: hal.Metadata{
					Name:         "ADS1115",
					Description:  "Supports ADS1115 ADC",
					Capabilities: []hal.Capability{hal.AnalogInput},
				},
			},
		}

		factory1115.appendParameters()
	})

	return factory1115
}

func (f *ads1115Factory) NewDriver(parameters map[string]interface{}, hardwareResources interface{}) (hal.Driver, error) {
	return f.newDriver(parameters, hardwareResources, 0, 9*time.Millisecond)
}
