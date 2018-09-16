package cmd

import (
	"math/rand"
	"time"
)

// RefrigerationsListOutput 制冷设备
type RefrigerationsListOutput struct {
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

// Input 输入参数
type Input struct {
	PowerFactor int // 输入功率因数
	Frequency   int //  输入频率
}

// Output 输出参数
type Output struct {
	Voltage    int // 电压
	Current    int // 电流
	Crequerycy int // 输出频率
}

// Environment 环境参数
type Environment struct {
	Temperature int // 环境温度
	Humidty     int // 环境湿度
}

// Battery 电池参数
type Battery struct {
	State        string // 状态
	Voltage      int    // 电压
	Current      int    // 电流
	Temperature  int    // 温度
	BackupTime   int    // 后备时间
	CapacityLeft int    // 剩余容量
}

// RefrigerationsList list refrigerations
func RefrigerationsList() ([]RefrigerationsListOutput, error) {
	return randomRefgerations(), nil
}

func randomRefgerations() []RefrigerationsListOutput {
	ref := []RefrigerationsListOutput{}
	for i := 0; i < randRange(3, 5); i++ {
		ref = append(ref, randomRefgeration(i))
	}
	return ref
}

func randomRefgeration(ID int) RefrigerationsListOutput {
	rand.Seed(time.Now().UnixNano())
	return RefrigerationsListOutput{
		ID:                ID,
		PowerRating:       randRange(1, 9),
		PowerSupplyMethod: randRange(1, 7),
		RunningState:      randRange(1, 5),
		SystemGrant:       randRange(1, 100000),
		SystemType:        randString("REF5300", "REF6600"),
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
