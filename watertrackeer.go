package main

import (
	"log"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func dateKey(t time.Time) string {
	return t.Format("2006-01-02") // ISO 8601 date, YYYY-MM-DD
}

func dateForMonday() time.Time {
	day := time.Now().Weekday()
	if day == time.Sunday {
		return time.Now().Add(-1 * time.Hour * 24 * 6) // Monday of previous week
	}

	daysSinceMonday := time.Duration(day - 1)
	dayLength := time.Hour * 24
	return time.Now().Add(-1 * dayLength * daysSinceMonday) // Monday is day 1
}

func historyLabel(date time.Time, p fyne.Preferences) fyne.CanvasObject {
	data := binding.BindPreferenceInt(dateKey(date), p)

	str := binding.IntToStringWithFormat(data, "%dml")
	num := widget.NewLabelWithData(str)
	num.Alignment = fyne.TextAlignTrailing
	return num
}

func makeuii(p fyne.Preferences) fyne.CanvasObject {
	total := binding.BindPreferenceInt(dateKey(time.Now()), p)

	label := canvas.NewText("0ml", theme.PrimaryColor())
	label.TextSize = 42
	label.Alignment = fyne.TextAlignCenter
	totalStr := binding.IntToStringWithFormat(total, "%dml")
	totalStr.AddListener(binding.NewDataListener(
		func() {
			label.Text, _ = totalStr.Get()
			label.Refresh()
		}))

	date := widget.NewLabel(time.Now().Format("Mon Jan 2 2006"))
	date.Alignment = fyne.TextAlignCenter

	amount := widget.NewEntry()
	amount.SetText("250")
	input := container.NewBorder(nil, nil, nil, widget.NewLabel("ml"), amount)
	add := widget.NewButton("Add", func() {
		inc, err := strconv.Atoi(amount.Text)
		if err != nil {
			log.Println("Failed to parse integer: " + amount.Text)
			return
		}

		current, _ := total.Get()
		total.Set(current + inc)
	})

	weekStart := dateForMonday()
	dayLength := time.Hour * 24
	history := container.NewGridWithColumns(2,
		widget.NewLabel("Monday"), historyLabel(weekStart, p),
		widget.NewLabel("Tuesday"), historyLabel(weekStart.Add(dayLength), p),
		widget.NewLabel("Wednesday"), historyLabel(weekStart.Add(dayLength*2), p),
		widget.NewLabel("Thursday"), historyLabel(weekStart.Add(dayLength*3), p),
		widget.NewLabel("Friday"), historyLabel(weekStart.Add(dayLength*4), p),
		widget.NewLabel("Saturday"), historyLabel(weekStart.Add(dayLength*5), p),
		widget.NewLabel("Sunday"), historyLabel(weekStart.Add(dayLength*6), p),
	)

	return container.NewVBox(label, date,
		container.NewGridWithColumns(2, input, add),
		widget.NewCard("History", "Total this week", history))
}

func main(){
	a := app.NewWithID("com.example.watertracker")
	w := a.NewWindow("Water Tracker")

	pref := a.Preferences()
	w.SetContent(makeuii(pref))
	w.ShowAndRun()
}