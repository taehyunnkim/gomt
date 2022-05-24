package tui

import (
	"github.com/go-routeros/routeros"
	"github.com/charmbracelet/bubbles/progress"
)

type state int

const (
	fetching state = iota
	ready
)

type MtModel struct {
	deviceInfo DeviceInfo
	client *routeros.Client
	state state
	sub chan dataMessage
	resource *resourceData
	cpu *cpuData
	height int
	width int
	minWidth int
}

type DeviceInfo struct {
	Platform string
	BoardName string
	OsVersion string
	CpuCoreCount int
}

type resourceData struct {
	uptime string
	freeMem uint64
	totalMem uint64
	memoryBar *progress.Model
	freeHdd uint64
	totalHdd uint64
	err error
}

type cpuData struct {
	count int
	bar map[int] *progress.Model
	data map[int] float64
	err error
}

type dataMessage struct {
	resourceData data
	cpuData data 
}

type data struct {
	reply *routeros.Reply
	err error
}

func New(client *routeros.Client, deviceInfo DeviceInfo, minWidth int) MtModel {
	cpuBars := make(map[int] *progress.Model)

	for i := 0; i < deviceInfo.CpuCoreCount; i++ {
		bar := progress.New(progress.WithDefaultGradient())
		cpuBars[i] = &bar
	}

	memoryBar := progress.New(progress.WithDefaultGradient())

	return MtModel {
		deviceInfo: deviceInfo,
		client: client,
		state: fetching,
		sub: make(chan dataMessage),
		resource: &resourceData{
			memoryBar: &memoryBar,
		},
		cpu: &cpuData{
			count: deviceInfo.CpuCoreCount,
			bar: cpuBars,
			data: make(map[int] float64),
		},
		minWidth: minWidth,
	}
}