package cmd

import (
	"fmt"
	"math/rand"

	"github.com/google/uuid"
)

// REFsListOutput 制冷设备
type REFsListOutput struct {
	ID                string `json:"ID"`                // 标识符
	PowerRating       int    `json:"PowerRating"`       // 额定功率
	PowerSupplyMethod int    `json:"PowerSupplyMethod"` // 供电方式
	RunningState      int    `json:"RunningState"`      // 运行功率
	SystemGrant       int    `json:"SystemGrant"`       // 机型容量
	SystemType        string `json:"SystemType"`        // 系统型号
	WorkMode          int    `json:"WorkMode"`          // 工作模式
	In                Input
	Out               Output
	Battery           Battery
	Environment       Environment
}

// Input 输入参数
type Input struct {
	PowerFactor int `json:"PowerFactor"` // 输入功率因数
	Frequency   int `json:"Frequency"`   //  输入频率
}

// Output 输出参数
type Output struct {
	Voltage    int `json:"Voltage"`    // 电压
	Current    int `json:"Current"`    // 电流
	Crequerycy int `json:"Crequerycy"` // 输出频率
}

// Environment 环境参数
type Environment struct {
	Temperature int `json:"Temperature"` // 环境温度
	Humidty     int `json:"Humidty"`     // 环境湿度
}

// Battery 电池参数
type Battery struct {
	State        string `json:"State"`        // 状态
	Voltage      int    `json:"Voltage"`      // 电压
	Current      int    `json:"Current"`      // 电流
	Temperature  int    `json:"Temperature"`  // 温度
	BackupTime   int    `json:"BackupTime"`   // 后备时间
	CapacityLeft int    `json:"CapacityLeft"` // 剩余容量
}

// REFGet get ref by id
func REFGet(id string) (REFsListOutput, error) {
	for _, item := range refs {
		if item.ID == id {
			return item, nil
		}
	}
	return REFsListOutput{}, fmt.Errorf("REF %s Not Found", id)
}

// REFsList list refs
func REFsList() ([]REFsListOutput, error) {
	return refs, nil
}

// REFsCreate create refs
func REFsCreate(item REFsListOutput) ([]REFsListOutput, error) {
	if item.ID == "" {
		item.ID = uuid.New().String()
	}
	refs = append(refs, item)
	return refs, nil
}

// REFUpdate update refs
func REFUpdate(newItem REFsListOutput) error {
	for index, item := range refs {
		if item.ID == newItem.ID {
			refs = append(refs[:index], append([]REFsListOutput{newItem}, refs[index+1:]...)...)
			return nil
		}
	}
	return fmt.Errorf("REF %s Not Found", newItem.ID)
}

// REFDelete delete ups by id
func REFDelete(id string) error {
	for index, item := range refs {
		if item.ID == id {
			refs = append(refs[:index], refs[index+1:]...)
			return nil
		}
	}
	return fmt.Errorf("REF %s Not Found", id)
}

func randomREFsList(ID string) REFsListOutput {
	return REFsListOutput{
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
