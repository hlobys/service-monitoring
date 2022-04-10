# System Monitor Utility

This is a simple system monitoring utility written in Go that gathers real-time data about your system, including CPU usage, memory usage, disk usage, and network statistics.

## Features

- Monitor CPU usage
- Monitor memory usage (total, used, and free)
- Monitor disk usage (total, used, and free)
- Monitor network usage (bytes sent and received)

## Clone the repository:

```bash
git clone https://github.com/yourusername/system-monitor.git
cd system-monitor
```
## Install the required dependencies:

```bash
go get github.com/shirou/gopsutil/v3/...
```

## Build the project:

```bash
go build -o system-monitor
```

## To run the system monitor, simply execute the compiled binary:
```bash
./system-monitor
```
