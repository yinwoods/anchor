package cmd

import (
	"fmt"
	"math/rand"

	"github.com/google/uuid"
)

// UPSListOutput 代表供电设备
type UPSListOutput struct {
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

// UPSGet get ups by id
func UPSGet(id string) (UPSListOutput, error) {
	for _, item := range ups {
		if item.ID == id {
			return item, nil
		}
	}
	return UPSListOutput{}, fmt.Errorf("UPS %s Not Found", id)
}

// UPSList list power supplies
func UPSList() ([]UPSListOutput, error) {
	return ups, nil
}

// UPSCreate create power supplies
func UPSCreate(item UPSListOutput) ([]UPSListOutput, error) {
	if item.ID == "" {
		item.ID = uuid.New().String()
	}
	ups = append(ups, item)
	return ups, nil
}

// UPSUpdate update power supplies
func UPSUpdate(newItem UPSListOutput) error {
	for index, item := range ups {
		if item.ID == newItem.ID {
			ups = append(ups[:index], append([]UPSListOutput{newItem}, ups[index+1:]...)...)
			return nil
		}
	}
	return fmt.Errorf("UPS %s Not Found", newItem.ID)
}

// UPSDelete delete ups by id
func UPSDelete(id string) error {
	for index, item := range ups {
		if item.ID == id {
			ups = append(ups[:index], ups[index+1:]...)
			return nil
		}
	}
	return fmt.Errorf("UPS %s Not Found", id)
}

func randomUPSList(ID string) UPSListOutput {
	return UPSListOutput{
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
