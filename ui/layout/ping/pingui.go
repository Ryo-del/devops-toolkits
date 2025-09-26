package ping

import (
	"fmt"
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/ryo-del/devops-toolkit/internal/ping"
)

func NewPingTab() fyne.CanvasObject {
	interval := ping.GetInterval()

	hostLabel := widget.NewLabel(fmt.Sprintf("Target: %s", ping.GetSource()))
	statusLabel := widget.NewLabel("⏳ Checking...")
	lastCheckLabel := widget.NewLabel("Last check: ...")

	// Кружочек-индикатор
	statusCircle := canvas.NewCircle(color.NRGBA{R: 200, G: 200, B: 200, A: 255})
	statusCircle.Resize(fyne.NewSize(20, 20)) // задание размера

	// Кнопка ручного обновления
	refreshBtn := widget.NewButton("🔄 Refresh Now", func() {
		updateStatus(statusLabel, lastCheckLabel, statusCircle)
	})

	// Горутинка для автообновления
	go func() {
		for {
			updateStatus(statusLabel, lastCheckLabel, statusCircle)
			time.Sleep(time.Duration(interval) * time.Second)
		}
	}()

	// Оформление: сверху хост, в центре индикатор и статус, снизу кнопка
	top := container.NewVBox(
		canvas.NewText("Ping Monitor", color.NRGBA{R: 0, G: 120, B: 220, A: 255}),
		hostLabel,
	)

	center := container.NewHBox(
		statusCircle,
		container.NewVBox(statusLabel, lastCheckLabel),
	)

	return container.NewBorder(top, refreshBtn, nil, nil, center)
}

func updateStatus(statusLabel, lastCheckLabel *widget.Label, statusCircle *canvas.Circle) {
	ping.RunPing()
	status := ping.GetPingStatus()
	statusLabel.SetText(status)
	lastCheckLabel.SetText(fmt.Sprintf("Last check: %s", time.Now().Format("15:04:05")))

	switch status {
	case "✅ Alive":
		statusCircle.FillColor = color.NRGBA{R: 0, G: 200, B: 0, A: 255}
	case "❌ Unreachable":
		statusCircle.FillColor = color.NRGBA{R: 200, G: 0, B: 0, A: 255}
	default:
		statusCircle.FillColor = color.NRGBA{R: 200, G: 200, B: 0, A: 255}
	}
	statusCircle.Refresh()
}
