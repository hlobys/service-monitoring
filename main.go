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

// loadConfig loads configuration from a YAML file
func loadConfig(configFile string) (Config, error) {
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		return Config{}, fmt.Errorf("error reading config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return Config{}, fmt.Errorf("error parsing config file: %w", err)
	}

	return config, nil
}

func main() {
	interval, err := parseFlags()
	if err != nil {
		log.Fatalf("Error parsing flags: %v", err)
	}

	for {
		if err := monitorSystem(); err != nil {
			log.Printf("Error monitoring system: %v", err)
		}
		time.Sleep(time.Duration(interval) * time.Second)
	}
}

func parseFlags() (int, error) {
	updateInterval := flag.String("interval", "", "Update interval in seconds")
	configFile := flag.String("config", "config.yaml", "Path to configuration file")
	flag.Parse()

	config, err := loadConfig(*configFile)
	if err != nil {
		return 0, fmt.Errorf("error loading config file: %w", err)
	}

	interval := config.Interval
	if *updateInterval != "" {
		interval, err = strconv.Atoi(*updateInterval)
		if err != nil || interval <= 0 {
			return 0, fmt.Errorf("invalid interval: %s. Please provide a positive integer", *updateInterval)
		}
	}

	return interval, nil
}

func monitorSystem() error {
	cpuPercent, err := getCPUUsage()
	if err != nil {
		return fmt.Errorf("error getting CPU info: %w", err)
	}

	virtualMem, err := getMemoryUsage()
	if err != nil {
		return fmt.Errorf("error getting memory info: %w", err)
	}

	diskUsage, err := getDiskUsage()
	if err != nil {
		return fmt.Errorf("error getting disk info: %w", err)
	}

	netIO, err := getNetworkUsage()
	if err != nil {
		return fmt.Errorf("error getting network info: %w", err)
	}

	printSystemInfo(cpuPercent, virtualMem, diskUsage, netIO)
	return nil
}

func getCPUUsage() (float64, error) {
	cpuPercent, err := cpu.Percent(0, false)
	if err != nil {
		return 0, err
	}
	return cpuPercent[0], nil
}

func getMemoryUsage() (*mem.VirtualMemoryStat, error) {
	return mem.VirtualMemory()
}

func getDiskUsage() (*disk.UsageStat, error) {
	return disk.Usage("/")
}

func getNetworkUsage() ([]net.IOCountersStat, error) {
	return net.IOCounters(false)
}

func printSystemInfo(cpuPercent float64, virtualMem *mem.VirtualMemoryStat, diskUsage *disk.UsageStat, netIO []net.IOCountersStat) {
	fmt.Printf("\n--- System Information ---\n")
	fmt.Printf("CPU Usage: %.2f%%\n", cpuPercent)
	fmt.Printf("Memory Usage: %.2f%% (Total: %v, Free: %v)\n", virtualMem.UsedPercent, virtualMem.Total, virtualMem.Free)
	fmt.Printf("Disk Usage: %.2f%% (Total: %v, Free: %v)\n", diskUsage.UsedPercent, diskUsage.Total, diskUsage.Free)
	fmt.Printf("Network Sent: %v, Received: %v\n", netIO[0].BytesSent, netIO[0].BytesRecv)
}
