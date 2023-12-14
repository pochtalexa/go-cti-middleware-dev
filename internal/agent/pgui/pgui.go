package pgui

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/pochtalexa/go-cti-middleware/internal/agent/auth"
	"github.com/pochtalexa/go-cti-middleware/internal/agent/flags"
	"github.com/pochtalexa/go-cti-middleware/internal/agent/storage"
	"github.com/rivo/tview"
	"github.com/rs/zerolog/log"
	"net/http"
)

var (
	Action     *tview.Form
	Header     *tview.TextView
	Footer     *tview.TextView
	UserState  *tview.TextView
	NewCall    *tview.TextView
	CallStatus *tview.TextView
	app        *tview.Application
)

func mute(checked bool) {
	const op = "pgui.mute"

	buf := bytes.Buffer{}

	body := storage.NewWsCommand()
	body.Name = "Mute"
	body.Login = flags.Login
	body.Cid = storage.AgentEvents.NewCall.Cid
	body.On = checked

	log.Info().Str("body", fmt.Sprintln(body)).Msg("mute")
	log.Info().Str("AgentEvents.NewCall", fmt.Sprintln(storage.AgentEvents.NewCall)).Msg("mute")

	enc := json.NewEncoder(&buf)
	enc.SetIndent("", "  ")
	if err := enc.Encode(body); err != nil {
		log.Error().Str("op", op).Err(err).Msg("Encode")
		return
	}

	// TODO добавить уведомление об ошибке отправки команды в CTI API
	req, _ := http.NewRequest(http.MethodPost, storage.AppConfig.ApiRoutes.Control, &buf)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", storage.AppConfig.TokenString))
	res, err := storage.AppConfig.HTTPClient.Do(req)
	if err != nil {
		log.Error().Str("op", op).Err(err).Msg("Do")
		FooterSetText("Connection error", tcell.ColorRed)
		return
	}
	defer res.Body.Close()

	checkStatusCode(res.StatusCode, op)
}

func answer() {
	const op = "pgui.answer"

	buf := bytes.Buffer{}

	body := storage.NewWsCommand()
	body.Name = "Answer"
	body.Login = flags.Login
	body.Cid = storage.AgentEvents.NewCall.Cid

	enc := json.NewEncoder(&buf)
	enc.SetIndent("", "  ")
	if err := enc.Encode(body); err != nil {
		log.Error().Str("op", op).Err(err).Msg("Encode")
		return
	}

	// TODO добавить уведомление об ошибке отправки команды в CTI API
	req, _ := http.NewRequest(http.MethodPost, storage.AppConfig.ApiRoutes.Control, &buf)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", storage.AppConfig.TokenString))
	res, err := storage.AppConfig.HTTPClient.Do(req)
	if err != nil {
		log.Error().Str("op", op).Err(err).Msg("Do")
		FooterSetText("Connection error", tcell.ColorRed)
		return
	}
	defer res.Body.Close()

	checkStatusCode(res.StatusCode, op)
}

func hangup() {
	const op = "pgui.hangup"

	buf := bytes.Buffer{}

	body := storage.NewWsCommand()
	body.Name = "Hangup"
	body.Login = flags.Login
	body.Cid = storage.AgentEvents.NewCall.Cid

	enc := json.NewEncoder(&buf)
	enc.SetIndent("", "  ")
	if err := enc.Encode(body); err != nil {
		log.Error().Str("op", op).Err(err).Msg("Encode")
		return
	}

	// TODO добавить уведомление об ошибке отправки команды в CTI API
	req, _ := http.NewRequest(http.MethodPost, storage.AppConfig.ApiRoutes.Control, &buf)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", storage.AppConfig.TokenString))
	res, err := storage.AppConfig.HTTPClient.Do(req)
	if err != nil {
		log.Error().Str("op", op).Err(err).Msg("Do")
		FooterSetText("Connection error", tcell.ColorRed)
		return
	}
	defer res.Body.Close()

	checkStatusCode(res.StatusCode, op)
}

func status(status string, index int) {
	const op = "pgui.statusWork"

	if status == "init" {
		return
	}

	buf := bytes.Buffer{}

	body := storage.NewWsCommand()
	body.Name = "ChangeUserState"
	body.Login = flags.Login
	body.State = status

	enc := json.NewEncoder(&buf)
	enc.SetIndent("", "  ")
	if err := enc.Encode(body); err != nil {
		log.Error().Str("op", op).Err(err).Msg("Encode")
		return
	}

	// TODO добавить уведомление об ошибке отправки команды в CTI API
	req, _ := http.NewRequest(http.MethodPost, storage.AppConfig.ApiRoutes.Control, &buf)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", storage.AppConfig.TokenString))
	res, err := storage.AppConfig.HTTPClient.Do(req)
	if err != nil {
		log.Error().Str("op", op).Err(err).Msg("Do")
		FooterSetText("Connection error", tcell.ColorRed)
		return
	}
	defer res.Body.Close()

	checkStatusCode(res.StatusCode, op)

	log.Debug().Str("op", op).Str("status", status).Msg("status changed")
}

func FooterSetText(text string, color tcell.Color) {
	app.QueueUpdateDraw(func() {
		Footer.SetText(text).SetTextColor(color)
	},
	)
}

func checkStatusCode(statusCode int, op string) {
	if statusCode != http.StatusOK {
		if statusCode == http.StatusUnauthorized {
			if err := auth.Refresh(); err != nil {
				log.Error().Str("op", op).Err(err).Msg("Refresh")
				FooterSetText("Connection error", tcell.ColorRed)
				return
			} else {
				log.Debug().Str("op", op).Msg("Refresh done")
				return
			}
		}
		log.Error().Str("op", op).Int("StatusCode", statusCode).Msg("")
		FooterSetText("Connection error", tcell.ColorRed)
		return
	}
}

func Init() {
	// TODO добавить управление курсорами
	app = tview.NewApplication()

	Header = tview.NewTextView().SetText("CTI Demo Control Board")
	Footer = tview.NewTextView().SetText("")
	UserState = newTextView("UserState")
	NewCall = newTextView("NewCall")
	CallStatus = newTextView("CallStatus")
	Action = newForm("Action", "agent")

	grid := tview.NewGrid().
		SetRows(1, 0, 0, 0, 1).
		SetColumns(30, 0, 30).
		SetBorders(true).
		AddItem(Header, 0, 0, 1, 3, 0, 0, false).
		AddItem(Footer, 4, 0, 1, 3, 0, 0, false)

	//Layout for screens narrower than 100 cells (menu and side bar are hidden).
	//grid.AddItem(actions, 0, 0, 0, 0, 0, 0, true).
	//	AddItem(main, 1, 0, 1, 3, 0, 0, false)

	// Layout for screens wider than 100 cells.
	grid.AddItem(Action, 1, 0, 3, 1, 1, 1, true).
		AddItem(UserState, 1, 1, 1, 2, 1, 1, false).
		AddItem(NewCall, 2, 1, 1, 2, 1, 1, false).
		AddItem(CallStatus, 3, 1, 1, 2, 1, 1, false)

	//go refresh()

	if err := app.SetRoot(grid, true).SetFocus(grid).EnableMouse(true).Run(); err != nil {
		log.Fatal().Err(err).Msg("run PguiApp")
	}
}

func newTextView(title string) *tview.TextView {
	textView := tview.NewTextView()

	textView.SetTextAlign(tview.AlignLeft)
	textView.SetScrollable(true)
	textView.SetTitle(title).SetBorder(true)
	textView.SetChangedFunc(func() {
		app.Draw()
	})
	return textView
}

func newForm(title string, login string) *tview.Form {
	form := tview.NewForm()

	form.SetTitle(title).SetBorder(true)
	form.AddTextView("login", login, 10, 1, true, false)
	form.AddDropDown("Status", []string{"init", "normal", "away", "dnd"}, 0, status)
	form.AddCheckbox("Mute", false, mute)
	form.AddButton("Answer", answer)
	form.AddButton("Hangup", hangup)
	form.MouseHandler()

	return form
}
