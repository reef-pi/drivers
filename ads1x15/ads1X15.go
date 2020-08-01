package ads1x15

const (
	configOsSingle  = 0x8000
	configOsBusy    = 0x0000
	configOsNotBusy = 0x8000

	configMuxSingle0 = 0x4000
	configMuxSingle1 = 0x5000
	configMuxSingle2 = 0x6000
	configMuxSingle3 = 0x7000
	configMuxDiff01  = 0x0000
	configMuxDiff03  = 0x1000
	configMuxDiff13  = 0x2000
	configMuxDiff23  = 0x3000

	configGainTwoThirds = 0x0000
	configGainOne       = 0x0200
	configGainTwo       = 0x0400
	configGainFour      = 0x0600
	configGainEight     = 0x0800
	configGainSixteen   = 0x0A00

	configModeSingle     = 0x0100
	configModeContinuous = 0x0000

	configDataRate128  = 0x0000
	configDataRate250  = 0x0020
	configDataRate490  = 0x0040
	configDataRate920  = 0x0060
	configDataRate1600 = 0x0080
	configDataRate2400 = 0x00A0
	configDataRate3300 = 0x00C0

	configComparatorModeTraditional = 0x0000
	configComparatorModeWindow      = 0x0010

	configComparitorPolarityActiveLow  = 0x0000
	configComparitorPolarityActiveHigh = 0x0008

	configComparitorNonLatching = 0x0000
	configComparitorLatching    = 0x0004

	configComparitorQueue1    = 0x0000
	configComparitorQueue2    = 0x0001
	configComparitorQueue4    = 0x0002
	configComparitorQueueNone = 0x0003

	configDefault = configOsSingle |
		configMuxDiff01 |
		configGainTwo |
		configModeSingle |
		configDataRate1600 |
		configComparatorModeTraditional |
		configComparitorPolarityActiveLow |
		configComparitorNonLatching |
		configComparitorQueueNone
)
