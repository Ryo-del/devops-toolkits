package monitorui

import (
	"encoding/json"
	"os"
	"regexp"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/ryo-del/devops-toolkit/internal/monitor"
)

type Config struct {
	Monitor MonitorConfig `json:"monitor"`
}
type MonitorConfig struct {
	Interval int `json:"interval"` // секунды
}

var config MonitorConfig

var numberRe = regexp.MustCompile(`\d+(\.\d+)?`) // ищем число вида 12 или 12.34

func NewMonitorTab() fyne.CanvasObject {
	// читаем конфиг
	file, err := os.Open("internal/settings/main.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var cfg Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		panic(err)
	}
	config = cfg.Monitor

	// widgets
	cpuBar := widget.NewProgressBar()
	memBar := widget.NewProgressBar()
	diskBar := widget.NewProgressBar()

	cpuLabel := widget.NewLabel("CPU: loading...")
	memLabel := widget.NewLabel("Memory: loading...")
	diskLabel := widget.NewLabel("Disk: loading...")
	netLabel := widget.NewLabel("Network: loading...")
	hostLabel := widget.NewLabel("Host: loading...")

	// общий апдейтер: получает строку от getter, обновляет label + опционально progress bar
	update := func(label *widget.Label, bar *widget.ProgressBar, getter func() (string, error)) {
		go func() {
			// делаем первый запрос сразу, а потом по таймеру
			for {
				text, err := getter()

				// все изменения UI — внутри fyne.Do
				fyne.Do(func() {
					if err != nil {
						label.SetText("Error")
						if bar != nil {
							bar.SetValue(0)
						}
						return
					}

					label.SetText(text)

					if bar != nil {
						m := numberRe.FindString(text)
						if m == "" {
							bar.SetValue(0)
						} else {
							v, err := strconv.ParseFloat(m, 64)
							if err != nil {
								bar.SetValue(0)
							} else {
								// ProgressBar принимает значение 0.0..1.0
								bar.SetValue(v / 100.0)
							}
						}
					}
				})

				time.Sleep(time.Duration(config.Interval) * time.Second)
			}
		}()
	}

	// запускаем обновления
	update(cpuLabel, cpuBar, monitor.GetCPUUsage)
	update(memLabel, memBar, monitor.GetMemoryUsage)
	update(diskLabel, diskBar, monitor.GetDiskUsage)
	update(netLabel, nil, monitor.GetNetworkIO)
	update(hostLabel, nil, monitor.GetHostInfo)

	// простой аккуратный layout
	cpuBox := container.NewVBox(cpuLabel, cpuBar)
	memBox := container.NewVBox(memLabel, memBar)
	diskBox := container.NewVBox(diskLabel, diskBar)
	rightCol := container.NewVBox(netLabel, hostLabel)

	header := widget.NewLabelWithStyle("📊 Системный мониторинг", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	grid := container.NewGridWithColumns(2, cpuBox, memBox, diskBox, rightCol)

	return container.NewVBox(header, grid)
}
