package main

import (
	"fmt"
	"log"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
)

func main() {
	// Start monitoring in a loop
	for {
		monitorSystem()
		time.Sleep(5 * time.Second) // Refresh every 5 seconds
	}
}

func monitorSystem() {
	// Get CPU usage information
	cpuPercent, err := cpu.Percent(0, false)
	if err != nil {
		log.Fatalf("Error getting CPU info: %v", err)
	}

	// Get memory usage information
	virtualMem, err := mem.VirtualMemory()
	if err != nil {
		log.Fatalf("Error getting memory info: %v", err)
	}

	// Get disk usage information
	diskUsage, err := disk.Usage("/")
	if err != nil {
		log.Fatalf("Error getting disk info: %v", err)
	}

	// Get network interface information
	netIO, err := net.IOCounters(false)
	if err != nil {
		log.Fatalf("Error getting network info: %v", err)
	}

	// Print the results
	fmt.Printf("\n--- System Information ---\n")
	fmt.Printf("CPU Usage: %.2f%%\n", cpuPercent[0])
	fmt.Printf("Memory Usage: %.2f%% (Total: %v, Free: %v)\n", virtualMem.UsedPercent, virtualMem.Total, virtualMem.Free)
	fmt.Printf("Disk Usage: %.2f%% (Total: %v, Free: %v)\n", diskUsage.UsedPercent, diskUsage.Total, diskUsage.Free)
	fmt.Printf("Network Sent: %v, Received: %v\n", netIO[0].BytesSent, netIO[0].BytesRecv)
}
