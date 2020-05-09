// +build nintendoswitch

package runtime

type timeUnit int64

const tickMicros = 1000 // (1000 * 625) / 12
const asyncScheduler = false

func postinit() {}

// Entry point for Go. Initialize all packages and call main.main().
//export main
func main() int {
	preinit()
	run()

	return exit(0) // Call libc_exit to cleanup libnx
}

// sleepTicks argument are actually in microseconds
func sleepTicks(d timeUnit) {
	usleep(uint(d))
}

// getArmSystemTimeNs returns ARM cpu ticks converted to nanoseconds
func getArmSystemTimeNs() uint64 {
	t := getArmSystemTick()
	return armTicksToNs(t)
}

// armTicksToNs converts cpu ticks to nanoseconds
// Nintendo Switch CPU ticks has a fixed rate at 19200000
// It is basically 52 ns per tick
func armTicksToNs(tick uint64) uint64 {
	return tick * 52
}

func armNsToTicks(ns int64) int64 {
	return ns / 52
}

func ticks() timeUnit {
	return timeUnit(getArmSystemTimeNs())
}

var stdoutBuffer = make([]byte, 0, 120)

func putchar(c byte) {
	if c == '\n' || len(stdoutBuffer)+1 >= 120 {
		svcOutputDebugString(&stdoutBuffer[0], uint64(len(stdoutBuffer)))
		stdoutBuffer = stdoutBuffer[:0]
		return
	}

	stdoutBuffer = append(stdoutBuffer, c)
}

//export usleep
func usleep(usec uint) int

//export abort
func abort() {
	exit(1)
}

//go:export exit
func exit(code int) int

//go:export armGetSystemTick
func getArmSystemTick() uint64

// armGetSystemTickFreq returns the system tick frequency
// means how many ticks per second
//go:export armGetSystemTickFreq
func armGetSystemTickFreq() uint64
