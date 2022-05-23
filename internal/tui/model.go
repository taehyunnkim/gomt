package tui

import (
	"github.com/go-routeros/routeros"
	"github.com/charmbracelet/bubbles/progress"
)

type MtModel struct {
	deviceInfo string
	client *routeros.Client
	sub chan dataMessage
	data *dataMessage
	cpu cpuData
}

type cpuData struct {
	count int
	bar map[int] progress.Model
	data map[int] float64
}

type dataMessage struct {
	resourceData data
	cpuData data 
}

type data struct {
	reply *routeros.Reply
	err error
}

func New(client *routeros.Client, deviceInfo string, cpuCoreCount int) MtModel {
	bars := make(map[int] progress.Model)

	for i := 0; i < cpuCoreCount; i++ {
		bars[i] = progress.New(progress.WithScaledGradient("#FF7CCB", "#FDFF8C"), progress.WithWidth(20))
	}

	return MtModel {
		deviceInfo: deviceInfo,
		client: client,
		sub: make(chan dataMessage),
		data: nil,
		cpu: cpuData{
			count: cpuCoreCount,
			bar: bars,
			data: make(map[int] float64),
		},
	}
}