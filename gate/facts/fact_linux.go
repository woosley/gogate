// +build linux
package facts

import (
	"bufio"
	"fmt"
	"github.com/woosley/gogate/gate/types"
	"github.com/woosley/gogate/gate/utils"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

var osfiles []string = []string{
	"/etc/redhat-release",
	"/etc/os-release",
}

const (
	AllBlocksDir     = "/sys/block"
	VirtualBlocksDir = "/sys/devices/virtual/block"
)

const unknown string = "Unknown"
const zero int64 = 0

func getOsFile(osfs []string) string {
	for _, f := range osfs {
		if exists, err := utils.IsFile(f); err == nil && exists {
			return f
		}
	}
	return ""
}

func GetOs() string {

	switch f := getOsFile(osfiles); f {
	case "/etc/redhat-release":
		if bytes, err := ioutil.ReadFile("/etc/redhat-release"); err == nil {
			return strings.TrimSpace(string(bytes))
		}
	case "/etc/os-release":
		return getFromOsRelease(f)

	default:
		return unknown
	}
	return unknown
}

//getFromOsRelease return distribution name from /etc/os-release
func getFromOsRelease(f string) string {
	fd, err := os.Open(f)
	if err != nil {
		return unknown
	}

	fscanner := bufio.NewScanner(fd)

	for fscanner.Scan() {
		line := fscanner.Text()
		if strings.Contains(line, "=") {
			kv := strings.SplitN(line, "=", 2)
			k, v := strings.TrimSpace(kv[0]), strings.TrimSpace(kv[1])
			if k == "PRETTY_NAME" {
				return strings.Trim(v, "\"")
			}
		}
	}
	return unknown
}

func GetCpu() (types.CpuInfo, error) {
	var cpu types.CpuInfo
	if fd, err := os.Open("/proc/cpuinfo"); err != nil {
		return cpu, err
	} else {
		cpuinfo := make(map[string]string)
		fscanner := bufio.NewScanner(fd)
		for fscanner.Scan() {
			line := fscanner.Text()
			if strings.Contains(line, ":") {
				kv := strings.SplitN(fscanner.Text(), ":", 2)
				cpuinfo[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
			}
		}
		cpu.Count, _ = strconv.Atoi(cpuinfo["physical id"])
		cpu.Count += 1
		cpu.Cores, _ = strconv.Atoi(cpuinfo["cpu cores"])
		return cpu, nil
	}
}

//getMemoryLinux returns both MemTotal and SwapTotal
func GetMemory() (string, string, error) {
	if fd, err := os.Open("/proc/meminfo"); err != nil {
		return "", "", err
	} else {
		mem := make(map[string]string)
		fscanner := bufio.NewScanner(fd)
		for fscanner.Scan() {
			kv := strings.SplitN(fscanner.Text(), ":", 2)
			mem[kv[0]] = strings.Trim(kv[1], " ")
		}
		return mem["MemTotal"], mem["SwapTotal"], nil
	}
}

func GetUptime() int64 {
	fname := "/proc/uptime"
	if bytes, err := ioutil.ReadFile(fname); err == nil {
		up := strings.SplitN(strings.TrimSpace(string(bytes)), " ", 2)[0]
		upt, err := strconv.ParseFloat(up, 64)
		if err != nil {
			return zero
		} else {
			return int64(upt)
		}
	}
	return zero
}

func GetDiskInfo() []types.DiskInfo {
	// lsdir /sys/block
	// lsdir /sys/devices/virtual/block
	// /sys/block/sda/queue/hw_sector_size disk sector size
	// /sys/block/sda/size disk sectors
	pdisks := make([]types.DiskInfo, 0)
	disks, _ := utils.LsDir(AllBlocksDir)
	vdisks, _ := utils.LsDir(VirtualBlocksDir)

	// physical is what listed in AllBlocksDir but not in VirtualBlocksDir
	for _, v := range disks {
		if !utils.ListHasString(vdisks, v) {
			size := GetDiskSize(v)
			pdisks = append(pdisks, types.DiskInfo{Name: v, Size: size})
		}
	}
	return pdisks
}

func GetDiskSize(disk string) string {
	sectors, _ := utils.SlurpFile(fmt.Sprintf("%s/%s/size", AllBlocksDir, disk))
	sectorSize, _ := utils.SlurpFile(fmt.Sprintf("%s/%s/queue/hw_sector_size", AllBlocksDir, disk))
	s1, _ := strconv.Atoi(sectors[0])
	s2, _ := strconv.Atoi(sectorSize[0])
	return fmt.Sprintf("%vMB", s1*s2/1000/1000)

}
