package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
	"gopkg.in/yaml.v2"
)

// Config structure to hold YAML configuration
type Config struct {
	Interval int `yaml:"interval"`
}

// Function to load configuration from a YAML file
func loadConfig(configFile string) (Config, error) {
	var config Config
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		return config, err
	}
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return config, err
	}
	return config, nil
}

func main() {
	// Define command-line flags
	updateInterval := flag.String("interval", "", "Update interval in seconds")
	configFile := flag.String("config", "config.yaml", "Path to configuration file")
	flag.Parse()

	// Load configuration from file
	config, err := loadConfig(*configFile)
	if err != nil {
		log.Fatalf("Error loading config file: %v", err)
	}

	// Determine interval (from command-line flag or config file)
	interval := config.Interval
	if *updateInterval != "" {
		interval, err = strconv.Atoi(*updateInterval)
		if err != nil || interval <= 0 {
			log.Fatalf("Invalid interval: %s. Please provide a positive integer.", *updateInterval)
		}
	}

	// Start monitoring in a loop with the user-defined interval
	for {
		monitorSystem()
		time.Sleep(time.Duration(interval) * time.Second)
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
