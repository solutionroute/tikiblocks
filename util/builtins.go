package util

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
)

// blockId is automatically allocated
// send channel is used to update blocks data
// Gothreads are waked up by messages on rec channels
// action is a map of whatever was in action json object for corressponding action
var FunctionMap = map[string]func(blockId int, send chan Change, rec chan bool, action map[string]interface{}){
	"#Date":       Date,
	"#Memory":     Memory,
	"#MemoryUsed": MemoryUsed,
	"#Cpu":        Cpu,
	"#Load":       Load,
	"#Uptime":     Uptime,
}

// Update time according to "format" property
func Date(blockId int, send chan Change, rec chan bool, action map[string]interface{}) {
	run := true
	for run {
		send <- Change{blockId, time.Now().Format(action["format"].(string)), true}
		// Block until other thread will ping you
		run = <-rec
	}
}

// Load returns the current 1, 5, 15 minute load averages
func Load(blockId int, send chan Change, rec chan bool, action map[string]interface{}) {
	run := true
	for run {
		v, _ := load.Avg()
		send <- Change{blockId, fmt.Sprintf(action["format"].(string), v.Load1, v.Load5, v.Load15), true}
		// Block until other thread will ping you
		run = <-rec
	}
}

// Uptime returns the system uptime as a formatted string
func Uptime(blockId int, send chan Change, rec chan bool, action map[string]interface{}) {
	run := true
	for run {
		v, _ := host.Uptime()
		send <- Change{blockId, fmt.Sprintf(action["format"].(string), HumanizeDuration(time.Duration(v)*time.Second)), true}
		// Block until other thread will ping you
		run = <-rec
	}
}

// Get current % of used memory
func Memory(blockId int, send chan Change, rec chan bool, action map[string]interface{}) {
	run := true
	for run {
		v, _ := mem.VirtualMemory()
		send <- Change{blockId, fmt.Sprintf(action["format"].(string), v.UsedPercent), true}
		// Block until other thread will ping you
		run = <-rec
	}
}

// Exposes Active and Total Memory Used, in GB
func MemoryUsed(blockId int, send chan Change, rec chan bool, action map[string]interface{}) {
	run := true
	for run {
		v, _ := mem.VirtualMemory()
		send <- Change{blockId, fmt.Sprintf(action["format"].(string), float64(v.Used)/1000000000, float64(v.Total)/1000000000), true}
		// Block until other thread will ping you
		run = <-rec
	}
}

// Get current % of used CPU
func Cpu(blockId int, send chan Change, rec chan bool, action map[string]interface{}) {
	run := true
	for run {
		val, _ := cpu.Percent(time.Second, false)
		send <- Change{blockId, fmt.Sprintf(action["format"].(string), val[0]), true}
		// Block until other thread will ping you
		run = <-rec
	}
}
