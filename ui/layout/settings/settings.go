package settings

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type Config struct {
	Scaner  ScanerConfig  `json:"scaner"`
	Parser  ParserConfig  `json:"parser"`
	Monitor MonitorConfig `json:"monitor"`
	Ping    PingConfig    `json:"ping"`
}
type PingConfig struct {
	Packages int    `json:"packages"`
	Interval int    `json:"interval"`
	Source   string `json:"source"`
}
type ScanerConfig struct {
	Protocol string `json:"protocol"`
	Ports    []int  `json:"ports"`
	IP       string `json:"ip"`
	Work     int    `json:"Work"`
	Time     int    `json:"time"`
}

type ParserConfig struct {
	LogPath string `json:"log_path"`
	Format  string `json:"format"`
	SaveTo  string `json:"save_to"`
}

type MonitorConfig struct {
	CPU      bool `json:"cpu"`
	Interval int  `json:"interval"` // в секундах, например
}

func saveConfig(cfg Config) {
	file, err := os.OpenFile("internal/settings/main.json", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Println("Ошибка открытия для записи:", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(cfg)
	if err != nil {
		fmt.Println("Ошибка при записи JSON:", err)
	}
}

func NewSettingsTab() fyne.CanvasObject {
	// читаем json
	file, err := os.Open("internal/settings/main.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var cfg Config
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&cfg)
	if err != nil {
		panic(err)
	}

	// --- Scaner
	protocolRadio := widget.NewRadioGroup([]string{"tcp", "udp"}, func(value string) {
		cfg.Scaner.Protocol = value
		saveConfig(cfg)
	})
	protocolRadio.SetSelected(cfg.Scaner.Protocol)

	ipEntry := widget.NewEntry()
	ipEntry.SetText(cfg.Scaner.IP)
	ipEntry.OnChanged = func(s string) {
		cfg.Scaner.IP = s
		saveConfig(cfg)
	}

	portsEntry := widget.NewEntry()
	portsStr := make([]string, len(cfg.Scaner.Ports))
	for i, p := range cfg.Scaner.Ports {
		portsStr[i] = strconv.Itoa(p)
	}
	portsEntry.SetText(strings.Join(portsStr, ","))
	portsEntry.OnChanged = func(s string) {
		var ports []int
		for _, part := range strings.Split(s, ",") {
			if p, err := strconv.Atoi(strings.TrimSpace(part)); err == nil {
				ports = append(ports, p)
			}
		}
		cfg.Scaner.Ports = ports
		saveConfig(cfg)
	}

	workEntry := widget.NewEntry()
	workEntry.SetText(strconv.Itoa(cfg.Scaner.Work))
	workEntry.OnChanged = func(s string) {
		if v, err := strconv.Atoi(s); err == nil {
			cfg.Scaner.Work = v
			saveConfig(cfg)
		}
	}

	timeEntry := widget.NewEntry()
	timeEntry.SetText(strconv.Itoa(cfg.Scaner.Time))
	timeEntry.OnChanged = func(s string) {
		if v, err := strconv.Atoi(s); err == nil {
			cfg.Scaner.Time = v
			saveConfig(cfg)
		}
	}

	// --- Parser
	logPath := widget.NewEntry()
	logPath.SetText(cfg.Parser.LogPath)
	logPath.OnChanged = func(s string) {
		cfg.Parser.LogPath = s
		saveConfig(cfg)
	}

	format := widget.NewEntry()
	format.SetText(cfg.Parser.Format)
	format.OnChanged = func(s string) {
		cfg.Parser.Format = s
		saveConfig(cfg)
	}

	saveTo := widget.NewEntry()
	saveTo.SetText(cfg.Parser.SaveTo)
	saveTo.OnChanged = func(s string) {
		cfg.Parser.SaveTo = s
		saveConfig(cfg)
	}

	// --- Monitor
	cpuCheck := widget.NewCheck("CPU monitor", func(b bool) {
		cfg.Monitor.CPU = b
		saveConfig(cfg)
	})
	cpuCheck.SetChecked(cfg.Monitor.CPU)

	monitorInterval := widget.NewEntry()
	monitorInterval.SetText(strconv.Itoa(cfg.Monitor.Interval))
	monitorInterval.OnChanged = func(s string) {
		if v, err := strconv.Atoi(s); err == nil {
			cfg.Monitor.Interval = v
			saveConfig(cfg)
		}
	}

	// --- Ping
	packagesEntry := widget.NewEntry()
	packagesEntry.SetText(strconv.Itoa(cfg.Ping.Packages))
	packagesEntry.OnChanged = func(s string) {
		if v, err := strconv.Atoi(s); err == nil {
			cfg.Ping.Packages = v
			saveConfig(cfg)
		}
	}

	pingInterval := widget.NewEntry()
	pingInterval.SetText(strconv.Itoa(cfg.Ping.Interval))
	pingInterval.OnChanged = func(s string) {
		if v, err := strconv.Atoi(s); err == nil {
			cfg.Ping.Interval = v
			saveConfig(cfg)
		}
	}

	sourceEntry := widget.NewEntry()
	sourceEntry.SetText(cfg.Ping.Source)
	sourceEntry.OnChanged = func(s string) {
		cfg.Ping.Source = s
		saveConfig(cfg)
	}

	restartBtn := widget.NewButton("🔄 Restart", func() {
		go func() {
			exe, err := os.Executable()
			if err != nil {
				fyne.LogError("Не могу найти путь к бинарю", err)
				return
			}

			// копируем аргументы текущего процесса
			args := os.Args

			// запускаем новый процесс
			proc, err := os.StartProcess(exe, args, &os.ProcAttr{
				Files: []*os.File{
					os.Stdin,
					os.Stdout,
					os.Stderr,
				},
			})
			if err != nil {
				fyne.LogError("Не удалось перезапустить", err)
				return
			}

			// сразу завершаем текущий
			if proc != nil {
				os.Exit(0)
			}
		}()
	})

	return container.NewVBox(
		widget.NewLabel("Настройки Scaner"),
		protocolRadio,
		widget.NewForm(
			widget.NewFormItem("IP", ipEntry),
			widget.NewFormItem("Ports (comma separated)", portsEntry),
			widget.NewFormItem("Workers", workEntry),
			widget.NewFormItem("Time", timeEntry),
		),

		widget.NewLabel("Настройки Parser"),
		widget.NewForm(
			widget.NewFormItem("Log path", logPath),
			widget.NewFormItem("Format", format),
			widget.NewFormItem("Save to", saveTo),
		),

		widget.NewLabel("Настройки Monitor"),
		cpuCheck,
		widget.NewForm(
			widget.NewFormItem("Interval (s)", monitorInterval),
		),

		widget.NewLabel("Настройки Ping"),
		widget.NewForm(
			widget.NewFormItem("Packages", packagesEntry),
			widget.NewFormItem("Interval", pingInterval),
			widget.NewFormItem("Source", sourceEntry),
		),
		restartBtn,
	)
}
