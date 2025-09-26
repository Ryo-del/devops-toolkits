package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	monitor "github.com/ryo-del/devops-toolkit/ui/layout/monitoring"

	// parser "github.com/ryo-del/devops-toolkit/ui/layout/parsering"
	docker "github.com/ryo-del/devops-toolkit/ui/layout/docker"
	ping "github.com/ryo-del/devops-toolkit/ui/layout/ping"
	port "github.com/ryo-del/devops-toolkit/ui/layout/portscanner"
	settings "github.com/ryo-del/devops-toolkit/ui/layout/settings"
)

func main() {
	a := app.New()
	w := a.NewWindow("DevOps Toolkit")

	tabs := container.NewAppTabs(
		container.NewTabItem("📊 Мониторинг", monitor.NewMonitorTab()),
		//container.NewTabItem("📁 Парсер логов", parser.NewParserTab()),
		container.NewTabItem("📡Порт сканер", port.NewScannerTab()),
		container.NewTabItem("🛜Сеть", ping.NewPingTab()),
		container.NewTabItem("🐳Docker", docker.NewDockerTab()),
		container.NewTabItem("⚙️Настройки", settings.NewSettingsTab()),
	)

	w.SetContent(tabs)
	w.Resize(fyne.NewSize(800, 600))
	w.ShowAndRun()
}
