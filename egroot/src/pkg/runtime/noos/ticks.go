package noos

import "log"

// Uptime returns how long system is running (in nanosecond). It includes or not
// the time when system is in deep sleep state - this is implementation specific.
func Uptime() uint64 {
	if MaxTasks() == 0 {
		log.Panic("noos.Uptime not supported (MaxTasks==0)")
	}
	return uptime()
}

var sysClk uint

// SetSysClk informs runtime about current system clock frequency. It should be
// called at every system clock change.
func SetSysClk(hz uint) {
	sysClk = hz
	setTickPeriod()
}

var tickEvent = AssignEvent()

// TickEvent returns event that is sended at every tasker interrupt.
func TickEvent() Event {
	return tickEvent
}