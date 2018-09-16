package cmd

import (
	"math/rand"
	"time"
)

// PowerSuppliesListOutput 代表供电设备
type PowerSuppliesListOutput struct {
	ID                int    // 标识符
	PowerRating       int    // 额定功率
	PowerSupplyMethod int    // 供电方式
	RunningState      int    // 运行功率
	SystemGrant       int    // 机型容量
	SystemType        string // 系统型号
	WorkMode          int    // 工作模式
	In                Input
	Out               Output
	Battery           Battery
	Environment       Environment
}

// PowerSuppliesList list power supplies
func PowerSuppliesList() ([]PowerSuppliesListOutput, error) {
	return randomPowerSuppliesList(), nil
}

func randomPowerSuppliesList() []PowerSuppliesListOutput {
	ups := []PowerSuppliesListOutput{}
	for i := 0; i < randRange(3, 5); i++ {
		ups = append(ups, randomPowerSupplyList(i))
	}
	return ups
}

func randomPowerSupplyList(ID int) PowerSuppliesListOutput {
	rand.Seed(time.Now().UnixNano())
	return PowerSuppliesListOutput{
		ID:                ID,
		PowerRating:       randRange(1, 9),
		PowerSupplyMethod: randRange(1, 7),
		RunningState:      randRange(1, 5),
		SystemGrant:       randRange(1, 100000),
		SystemType:        randString("UPS5000", "UPS8000"),
		WorkMode:          rand.Intn(4),
		In: Input{
			PowerFactor: randRange(-100, 100),
			Frequency:   randRange(0, 100),
		},
		Out: Output{
			Voltage:    randRange(0, 10000),
			Current:    randRange(0, 10000),
			Crequerycy: randRange(0, 100),
		},
		Battery: Battery{
			State:        randString("PowerNormal", "PowerLow"),
			Voltage:      randRange(0, 10000),
			Current:      randRange(0, 10000),
			Temperature:  randRange(-20, 80),
			BackupTime:   randRange(0, 172800),
			CapacityLeft: randRange(0, 100),
		},
		Environment: Environment{
			Temperature: randRange(-20, 80),
			Humidty:     randRange(0, 1000),
		},
	}
}
