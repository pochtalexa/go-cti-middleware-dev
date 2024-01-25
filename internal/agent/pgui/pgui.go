package pgui

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/rs/zerolog/log"

	"github.com/pochtalexa/go-cti-middleware/internal/agent/auth"
	"github.com/pochtalexa/go-cti-middleware/internal/agent/flags"
	"github.com/pochtalexa/go-cti-middleware/internal/agent/storage"
)

var (
	Action     *tview.Form
	Call       *tview.Form
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
	body.Rid = getRandRid()
	body.Login = flags.Login
	body.Cid = storage.AgentEvents.NewCall.Cid
	body.On = checked

	log.Info().Str("body", fmt.Sprintln(body)).Msg("mute")
	log.Info().Str("AgentEvents.NewCall", fmt.Sprintln(storage.AgentEvents.NewCall)).Msg("mute")

	enc := json.NewEncoder(&buf)
	enc.SetIndent("", "  ")
	if err := enc.Encode(body); err != nil {
		log.Error().Str("op", op).Err(err).Msg("Encode")
		storage.AppConfig.DisplayErrCh <- fmt.Sprintf("%s: json error. err: %s", op, err)
		return
	}

	res, err := sendControlCommand(buf, op)
	if err != nil {
		log.Error().Str("op", op).Err(err).Msg("sendControlCommand")
		return
	}

	if err = checkStatusCode(res.StatusCode, op); err != nil {
		log.Debug().Str("op", op).Err(err).Msg("checkStatusCode")
	} else {
		log.Debug().Str("op", op).Msg("ok")
	}
}

func answer() {
	const op = "pgui.answer"

	buf := bytes.Buffer{}

	body := storage.NewWsCommand()
	body.Rid = getRandRid()
	body.Name = "Answer"
	body.Login = flags.Login
	body.Cid = storage.AgentEvents.NewCall.Cid

	enc := json.NewEncoder(&buf)
	enc.SetIndent("", "  ")
	if err := enc.Encode(body); err != nil {
		log.Error().Str("op", op).Err(err).Msg("Encode")
		storage.AppConfig.DisplayErrCh <- fmt.Sprintf("%s: json error. err: %s", op, err)
		return
	}

	res, err := sendControlCommand(buf, op)
	if err != nil {
		log.Error().Str("op", op).Err(err).Msg("sendControlCommand")
		return
	}

	if err = checkStatusCode(res.StatusCode, op); err != nil {
		log.Debug().Str("op", op).Err(err).Msg("checkStatusCode")
	} else {
		log.Debug().Str("op", op).Msg("ok")
	}
}

func hangup() {
	const op = "pgui.hangup"

	buf := bytes.Buffer{}

	body := storage.NewWsCommand()
	body.Name = "Hangup"
	body.Rid = getRandRid()
	body.Login = flags.Login
	body.Cid = storage.AgentEvents.NewCall.Cid

	enc := json.NewEncoder(&buf)
	enc.SetIndent("", "  ")
	if err := enc.Encode(body); err != nil {
		log.Error().Str("op", op).Err(err).Msg("Encode")
		storage.AppConfig.DisplayErrCh <- fmt.Sprintf("%s: json error. err: %s", op, err)
		return
	}

	res, err := sendControlCommand(buf, op)
	if err != nil {
		log.Error().Str("op", op).Err(err).Msg("sendControlCommand")
		return
	}

	if err = checkStatusCode(res.StatusCode, op); err != nil {
		log.Debug().Str("op", op).Err(err).Msg("checkStatusCode")
	} else {
		log.Debug().Str("op", op).Msg("ok")
	}
}

func closeWrapUp() {
	const op = "pgui.closeWrapUp"

	buf := bytes.Buffer{}

	body := storage.NewWsCommand()
	body.Name = "Close"
	body.Rid = getRandRid()
	body.Login = flags.Login
	body.Cid = storage.AgentEvents.NewCall.Cid

	enc := json.NewEncoder(&buf)
	enc.SetIndent("", "  ")
	if err := enc.Encode(body); err != nil {
		log.Error().Str("op", op).Err(err).Msg("Encode")
		storage.AppConfig.DisplayErrCh <- fmt.Sprintf("%s: json error. err: %s", op, err)
		return
	}

	res, err := sendControlCommand(buf, op)
	if err != nil {
		log.Error().Str("op", op).Err(err).Msg("sendControlCommand")
		return
	}

	if err = checkStatusCode(res.StatusCode, op); err != nil {
		log.Debug().Str("op", op).Err(err).Msg("checkStatusCode")
	} else {
		log.Debug().Str("op", op).Msg("ok")
	}
}

func status(status string, index int) {
	const op = "pgui.statusWork"

	if status == "init" {
		return
	}

	buf := bytes.Buffer{}

	body := storage.NewWsCommand()
	body.Name = "ChangeUserState"
	body.Rid = getRandRid()
	body.Login = flags.Login
	body.State = status

	enc := json.NewEncoder(&buf)
	enc.SetIndent("", "  ")
	if err := enc.Encode(body); err != nil {
		log.Error().Str("op", op).Err(err).Msg("Encode")
		storage.AppConfig.DisplayErrCh <- fmt.Sprintf("%s: json error. err: %s", op, err)
		return
	}

	res, err := sendControlCommand(buf, op)
	if err != nil {
		log.Error().Str("op", op).Err(err).Msg("sendControlCommand")
		return
	}

	if err = checkStatusCode(res.StatusCode, op); err != nil {
		log.Debug().Str("op", op).Err(err).Msg("checkStatusCode")
	} else {
		log.Debug().Str("op", op).Str("status", status).Msg("status changed")
	}

}

func call() {
	const op = "pgui.call"

	inputField := Call.GetFormItem(0).(*tview.InputField)
	text := inputField.GetText()

	log.Debug().Str("op", op).Str("text", text).Msg("PhoneNumber")

	if text == "" {
		return
	}

	buf := bytes.Buffer{}

	body := storage.NewWsCommand()
	body.Name = "MakeCall"
	body.Rid = getRandRid()
	body.Login = flags.Login
	body.PhoneNumber = text

	enc := json.NewEncoder(&buf)
	enc.SetIndent("", "  ")
	if err := enc.Encode(body); err != nil {
		log.Error().Str("op", op).Err(err).Msg("Encode")
		storage.AppConfig.DisplayErrCh <- fmt.Sprintf("%s: json error. err: %s", op, err)
		return
	}

	res, err := sendControlCommand(buf, op)
	if err != nil {
		log.Error().Str("op", op).Err(err).Msg("sendControlCommand")
		return
	}

	if err = checkStatusCode(res.StatusCode, op); err != nil {
		log.Debug().Str("op", op).Err(err).Msg("MakeCall")
	} else {
		log.Debug().Str("op", op).Msg("MakeCall - ok")
	}

	return
}

func reconnect() {
	const op = "pgui.reconnect"

	if err := auth.GetAuthorization(); err != nil {
		log.Debug().Str("op", op).Err(err).Msg("")
	}
}

func getRandRid() string {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	// generate a random 8-digit ID
	id := ""
	for i := 0; i < 8; i++ {
		id += strconv.Itoa(rand.Intn(10))
	}

	return id
}

func FooterSetText(text string, color tcell.Color) {
	app.QueueUpdateDraw(func() {
		Footer.SetText(text).SetTextColor(color)
	},
	)
}

func sendControlCommand(buf bytes.Buffer, opCalling string) (*http.Response, error) {
	const op = "sendControlCommand"

	req, _ := http.NewRequest(http.MethodPost, storage.AppConfig.ApiRoutes.Control, &buf)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", storage.AppConfig.GetToken()))
	res, err := storage.AppConfig.HTTPClient.Do(req)
	if err != nil {
		storage.AppConfig.DisplayErrCh <- fmt.Sprintf("%s-%s: connection error. err: %s", op, opCalling, err)
		return nil, fmt.Errorf("%s-%s: %w", op, opCalling, err)
	}
	defer res.Body.Close()

	return res, nil
}

func checkStatusCode(statusCode int, op string) error {
	if statusCode != http.StatusOK {
		if statusCode == http.StatusUnauthorized {
			if err := auth.Refresh(); err != nil {
				log.Error().Str("op", op).Err(err).Msg("Refresh")
				storage.AppConfig.DisplayErrCh <- fmt.Sprintf("%s: connection error. statusCode: %d", op, statusCode)
				return fmt.Errorf("auth refresh. %s: %w", op, err)
			} else {
				log.Debug().Str("op", op).Msg("Refresh done")
				return nil
			}
		}
		log.Error().Str("op", op).Int("StatusCode", statusCode).Msg("checkStatusCode")
		storage.AppConfig.DisplayErrCh <- fmt.Sprintf("%s: connection error. statusCode: %d", op, statusCode)
		return fmt.Errorf("%s: statusCode: %d", op, statusCode)
	}

	return nil
}

func Init() {
	// TODO добавить управление курсорами
	app = tview.NewApplication()

	Header = tview.NewTextView().SetText("CTI Demo Control Board")
	Footer = tview.NewTextView().SetText("")
	UserState = newTextView("UserState")
	NewCall = newTextView("NewCall")
	CallStatus = newTextView("CallStatus")
	Action = newActionForm("Action", "agent")
	Call = newCallForm("Call")

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
	grid.AddItem(Action, 1, 0, 2, 1, 1, 1, true).
		AddItem(Call, 3, 0, 1, 1, 1, 1, false).
		AddItem(UserState, 1, 1, 1, 2, 1, 1, false).
		AddItem(NewCall, 2, 1, 1, 2, 1, 1, false).
		AddItem(CallStatus, 3, 1, 1, 2, 1, 1, false)

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

func newActionForm(title string, login string) *tview.Form {
	form := tview.NewForm()

	form.SetTitle(title).SetBorder(true)
	form.AddTextView("login", login, 10, 1, true, false)
	form.AddDropDown("Status", []string{"init", "normal", "away", "dnd"}, 0, status)
	form.AddCheckbox("Mute", false, mute)
	form.AddButton("Answ", answer)
	form.AddButton("Hang", hangup)
	form.AddButton("Wrap", closeWrapUp)
	form.MouseHandler()

	return form
}

// TODO: добавить конпку Reconnect
func newCallForm(title string) *tview.Form {
	form := tview.NewForm()

	form.SetTitle(title).SetBorder(true)
	form.AddInputField("Num:", "", 20, newCallFormValidateInput, nil)
	form.AddButton("Call", call)
	form.AddButton("Hang", hangup)
	form.AddButton("ReConn", reconnect)
	form.MouseHandler()

	return form
}

func newCallFormValidateInput(textToCheck string, lastChar rune) bool {
	if m, _ := regexp.MatchString("^\\d+$", textToCheck); m {
		return true
	}
	return false
}
