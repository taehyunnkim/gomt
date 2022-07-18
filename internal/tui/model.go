package tui

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/go-routeros/routeros"
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
	keys keyMap
	help help.Model
	debug bool
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
	minBarWidth int
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

func New(client *routeros.Client, deviceInfo DeviceInfo, debug bool, minWidth int) MtModel {
	cpuBars := make(map[int] *progress.Model)

	for i := 0; i < deviceInfo.CpuCoreCount; i++ {
		bar := progress.New(progress.WithDefaultGradient(), progress.WithoutPercentage())
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
			minBarWidth: 26,
		},
		keys: keys,
		help: help.New(),
		minWidth: minWidth,
		debug: debug,
	}
}