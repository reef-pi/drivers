package ads1x15

import (
	"sync"
	"time"

	"github.com/reef-pi/hal"
)

type ads1015Factory struct {
	ads1X15Factory
}

var factory1015 *ads1015Factory
var once sync.Once

// Ads1015Factory returns a singleton ADS1015 factory
func Ads1015Factory() hal.DriverFactory {

	once.Do(func() {
		factory1015 = &ads1015Factory{
			ads1X15Factory{
				meta: hal.Metadata{
					Name:         "ADS1015",
					Description:  "Supports ADS1015 ADC",
					Capabilities: []hal.Capability{hal.AnalogInput},
				},
			},
		}

		factory1015.appendParameters()
	})

	return factory1015
}

func (f *ads1015Factory) NewDriver(parameters map[string]interface{}, hardwareResources interface{}) (hal.Driver, error) {
	return f.newDriver(parameters, hardwareResources, 4, 1*time.Millisecond)
}
