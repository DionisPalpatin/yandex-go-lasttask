package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/go-resty/resty/v2"
)

func main() {
	url := "http://srv.msk01.gigacorp.local/_stats"
	//url := "http://localhost:8080/_stats"
	client := resty.New()
	dataGetErr := 0

	i := 0
	for {
		// fmt.Printf("=== Итерация %d ===", i)
		i++

		resp, err := client.R().Get(url)
		if err != nil {
			dataGetErr++
			if dataGetErr >= 3 {
				fmt.Println("Unable to fetch server statistic")
			}
			continue
		}

		if resp.StatusCode() != 200 || strings.TrimSpace(resp.String()) == "" {
			dataGetErr++
			if dataGetErr >= 3 {
				fmt.Println("Unable to fetch server statistic")
			}
			continue
		}

		parts := strings.Split(resp.String(), ",")
		if len(parts) != 7 {
			dataGetErr++
			if dataGetErr >= 3 {
				fmt.Println("Unable to fetch server statistic")
			}
			continue
		}

		var vals [7]int64
		bad := false
		for i, el := range parts {
			vStr := strings.TrimSpace(el)
			v, perr := strconv.ParseInt(vStr, 10, 64)
			if perr != nil {
				bad = true
				break
			}
			vals[i] = v
		}
		if bad {
			dataGetErr++
			if dataGetErr >= 3 {
				fmt.Println("Unable to fetch server statistic")
			}
			continue
		}

		load := vals[0]
		totalMem := vals[1]
		usedMem := vals[2]
		totalDisk := vals[3]
		usedDisk := vals[4]
		totalNet := vals[5]
		usedNet := vals[6]
		dataGetErr = 0

		// new comment1
		// fmt.Printf("Распарсенные значения: %d, %d, %d, %d, %d, %d, %d\n\n\n",
		// 	load, totalMem, usedMem, totalDisk, usedDisk, totalNet, usedNet)

		if load > 30 {
			fmt.Printf("Load Average is too high: %d\n", load)
		}

		memPct := float64(usedMem) * 100.0 / float64(totalMem)
		if memPct > 80 {
			fmt.Printf("Memory usage too high: %d%%\n", int(memPct))
		}

		if float64(usedDisk) * 100.0 / float64(totalDisk) > 90.0 {
			freeBytes := max(totalDisk - usedDisk, 0)
			freeMb := freeBytes / (1024 * 1024)
			fmt.Printf("Free disk space is too low: %d Mb left\n", freeMb)
		}

		if float64(usedNet) * 100.0 / float64(totalNet) > 90.0 {
			availBytes := max(totalNet - usedNet, 0)
			availMbit := availBytes / 1_000_000
			fmt.Printf("Network bandwidth usage high: %d Mbit/s available\n", availMbit)
		}
	}
}