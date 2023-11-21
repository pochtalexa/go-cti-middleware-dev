package storage

import (
	"fmt"
)

var (
	AgentEvents = NewAgentEvents()
)

// StAgentEvents возможные события относительно оператора
type StAgentEvents struct {
	UserState              StUserState
	NewCall                StNewCall
	CallStatus             StCallStatus
	CurrentCall            StCurrentCall
	CallParams             StCallParams
	Calls                  StCalls
	OnClose                StOnClose
	OnTransferCall         StOnTransferCall
	TransferSucceed        StTransferSucceed
	TransferFailed         StTransferFailed
	TransferCallReturned   StTransferCallReturned
	SetSessionModeResponse StSetSessionModeResponse
	CallParamsUpdated      StCallParamsUpdated
	LocalParamsUpdated     StLocalParamsUpdated
	Conferences            StConferences
	Ok                     StOk
	Error                  StError
	ParseError             StParseError
}

// StUserState Информация об изменении состояния программного телефона пользователя
type StUserState struct {
	Rid       string   `json:"rid,omitempty"`
	Name      string   `json:"name"`
	Login     string   `json:"login"`
	State     string   `json:"state"`
	Substates []string `json:"substates,omitempty"`
	Time      int      `json:"time,omitempty"`
	Reason    string   `json:"reason,omitempty"`
}

// StNewCall Оповещение о создании в программном телефоне телефонного вызова или текстового сообщения
type StNewCall struct {
	Rid          string   `json:"rid,omitempty"`
	Name         string   `json:"name"`
	Login        string   `json:"login"`
	Cid          int      `json:"cid,omitempty"`
	Type         string   `json:"type"`
	State        string   `json:"state"`
	Direction    string   `json:"direction"`
	Hold         bool     `json:"hold,omitempty"`
	Muted        bool     `json:"muted,omitempty"`
	DisplayName  string   `json:"displayName,omitempty"`
	SrcAddr      string   `json:"srcAddr,omitempty"`
	DstAddr      string   `json:"dstAddr,omitempty"`
	CreationTime int      `json:"creationTime,omitempty"`
	AnswerTime   int      `json:"answerTime,omitempty"`
	HangupTime   int      `json:"hangupTime,omitempty"`
	Params       StParams `json:"params,omitempty"`
}

// StParams Параметры вызова. Если параметры вызова еще не известны, то поле отсутствует.
// передается в составе других структур
type StParams struct {
	SessionID string `json:"sessionId"`
	URL       string `json:"url,omitempty"`
	Caller    string `json:"caller,omitempty"`
	Called    string `json:"called,omitempty"`
	SrcAddr   string `json:"srcAddr,omitempty"`
}

// StCallStatus Оповещение о изменении состояния вызова или текстового сообщения
type StCallStatus struct {
	Rid          string `json:"rid,omitempty"`
	Name         string `json:"name"`
	Login        string `json:"login"`
	Cid          int    `json:"cid,omitempty"`
	Type         string `json:"type"`
	State        string `json:"state"`
	EndedBySide  string `json:"endedBySide,omitempty"`
	Direction    string `json:"direction"`
	Hold         bool   `json:"hold,omitempty"`
	HoldTarget   string `json:"holdTarget,omitempty"`
	Muted        bool   `json:"muted,omitempty"`
	DisplayName  string `json:"displayName,omitempty"`
	DstAddr      string `json:"dstAddr,omitempty"`
	SrcAddr      string `json:"srcAddr,omitempty"`
	CreationTime int    `json:"creationTime,omitempty"`
	AnswerTime   int    `json:"answerTime,omitempty"`
	HangupTime   int    `json:"hangupTime,omitempty"`
	LocalParams  struct {
		Name string `json:"name,omitempty"`
	} `json:"localParams,omitempty"`
	Params       StParams `json:"params,omitempty"`
	TransferInfo struct {
		Called    string `json:"called,omitempty"`
		StartTime int    `json:"startTime,omitempty"`
		State     string `json:"state,omitempty"`
		Type      string `json:"type,omitempty"`
	} `json:"transferInfo,omitempty"`
}

// StCurrentCall Оповещение о изменении активного вызова в программном телефоне
type StCurrentCall struct {
	Name  string `json:"name"`
	Login string `json:"login"`
	Cid   int    `json:"cid,omitempty"`
}

// StCallParams Предоставление списка параметров вызова
type StCallParams struct {
	Rid    string   `json:"rid,omitempty"`
	Name   string   `json:"name"`
	Login  string   `json:"login"`
	Cid    int      `json:"cid,omitempty"`
	Params StParams `json:"params,omitempty"`
}

// StCalls Предоставление списка текущих вызовов, обслуживаемых программным телефоном оператора
type StCalls struct {
	Rid   string         `json:"rid,omitempty"`
	Name  string         `json:"name"`
	Login string         `json:"login"`
	Calls []StCallStatus `json:"calls"`
}

// StOnClose Оповещение о завершении обработки вызова
type StOnClose struct {
	Name  string `json:"name"`
	Login string `json:"login"`
	Cid   int    `json:"cid,omitempty"`
}

// StOnTransferCall Оповещение о начале перенаправления вызова
type StOnTransferCall struct {
	Name         string `json:"name"`
	Login        string `json:"login"`
	Cid          int    `json:"cid,omitempty"`
	WasConnected bool   `json:"wasConnected,omitempty"`
	Type         string `json:"type,omitempty"`
	Called       string `json:"called,omitempty"`
	StartTime    int    `json:"startTime,omitempty"`
}

// StTransferSucceed Оповещение об успешном перенаправлении вызова
type StTransferSucceed struct {
	Name  string `json:"name"`
	Login string `json:"login"`
	Cid   int    `json:"cid,omitempty"`
}

// StTransferFailed Оповещение о неуспешном перенаправлении вызова
type StTransferFailed struct {
	Name   string `json:"name"`
	Login  string `json:"login"`
	Cid    int    `json:"cid,omitempty"`
	Status string `json:"status,omitempty"`
}

// StTransferCallReturned Оповещение о возврате вызова при перенаправлении с возвратом
type StTransferCallReturned struct {
	Name  string `json:"name"`
	Login string `json:"login"`
	Cid   int    `json:"cid,omitempty"`
}

// StSetSessionModeResponse Оповещение о изменениях режима безопасности
type StSetSessionModeResponse struct {
	Name  string `json:"name"`
	Login string `json:"login"`
	Cid   int    `json:"cid,omitempty"`
	Error string `json:"error,omitempty"`
}

// StCallParamsUpdated Оповещение об успешном изменении или добавлении параметра вызова
type StCallParamsUpdated struct {
	Name   string            `json:"name"`
	Login  string            `json:"login"`
	Cid    int               `json:"cid,omitempty"`
	Params map[string]string `json:"params,omitempty"`
}

// StLocalParamsUpdated Оповещение об успешном обновлении дополнительных параметров вызова
type StLocalParamsUpdated struct {
	Name   string            `json:"name"`
	Login  string            `json:"login"`
	Cid    int               `json:"cid,omitempty"`
	Params map[string]string `json:"params,omitempty"`
}

// StConferences Предоставление списка конференций
type StConferences struct {
	Rid         string             `json:"rid,omitempty"`
	Name        string             `json:"name"`
	Login       string             `json:"login"`
	Conferences []StConferenceInfo `json:"conferences,omitempty"`
}

// StConferenceInfo Предоставление списка вызовов, участвующих в конференции
type StConferenceInfo struct {
	Rid                  string   `json:"rid,omitempty"`
	Name                 string   `json:"name"`
	Login                string   `json:"login"`
	ConfID               string   `json:"confId,omitempty"`
	DisplayName          string   `json:"displayName,omitempty"`
	Calls                []string `json:"calls,omitempty"`
	OperatorInConference bool     `json:"operatorInConference,omitempty"`
}

// StOk Успешная обработка команды
type StOk struct {
	Rid  string `json:"rid,omitempty"`
	Name string `json:"name"`
}

// StError Ошибка
type StError struct {
	Rid         string `json:"rid,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

// StParseError Ошибка разбора сообщения
type StParseError struct {
	Rid         string `json:"rid,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

// StWsCommand команда в сторону CTI API
type StWsCommand struct {
	Name            string `json:"name"` // название команды или события
	Login           string `json:"login,omitempty"`
	Rid             string `json:"rid,omitempty"` // Значение данного поля ответного сообщения соответствует значению этого же поля команды. Если событие не является ответом на команду, то поле отсутствует.
	Cid             int    `json:"cid,omitempty"` // Идентификатор обращения (call identifier) в контексте оператора. Указывается как идентификатор обращения во всех сообщениях, которые касаются обработки обращения.
	ProtocolVersion string `json:"protocolVersion,omitempty"`
	PhoneNumber     string `json:"phoneNumber,omitempty"`
	ParamName       string `json:"paramName,omitempty"`
	ParamValue      string `json:"paramValue,omitempty"`
	State           string `json:"state,omitempty"` // Состояние программного телефона, которое необходимо установить
	On              bool   `json:"on,omitempty"`
	Enable          bool   `json:"enable,omitempty"`
	Target          string `json:"target,omitempty"`
	DTMFString      string `json:"DTMFString,omitempty"`
	Url             string `json:"url,omitempty"`
}

func NewAgentEvents() *StAgentEvents {
	return &StAgentEvents{}
}

func NewWsCommand() *StWsCommand {
	return &StWsCommand{}
}

// ToString Формируем строку для отображения в GUI
func (a *StAgentEvents) ToString(name string) (string, error) {
	var result string

	switch name {
	case "UserState":
		if a.UserState.Name != "" {
			event := a.UserState
			result = fmt.Sprintf("state: %v, substates: %v, reason: %v",
				event.State,
				event.Substates,
				event.Reason,
			)
		} else {
			result = "-"
		}
	case "NewCall":
		if a.NewCall.Name != "" {
			event := a.NewCall
			result = fmt.Sprintf("state: %v, direction: %v, displayName: %v",
				event.State,
				event.Direction,
				event.DisplayName,
			)
		} else {
			result = "-"
		}
	case "CallStatus":
		if a.CallStatus.Name != "" {
			event := a.CallStatus
			result = fmt.Sprintf(
				"state: %v, "+
					"params: %v, "+
					"creationTime: %v, "+
					"answerTime: %v, "+
					"hangupTime: %v, ",
				event.State,
				event.Params,
				event.CreationTime,
				event.AnswerTime,
				event.HangupTime,
			)
		} else {
			result = "-"
		}
	default:
		return "", fmt.Errorf("ToString - can not find event name: %s", name)
	}

	return result, nil
}

func (a *StAgentEvents) UpdateEvents(newEvents *StAgentEvents) {

	if newEvents.UserState.Name != "" {
		a.UserState = newEvents.UserState
	}

	if newEvents.NewCall.Name != "" {
		a.NewCall = newEvents.NewCall
	}

	if newEvents.CallStatus.Name != "" {
		a.CallStatus = newEvents.CallStatus
	}

	if newEvents.CurrentCall.Name != "" {
		a.CurrentCall = newEvents.CurrentCall
	}

	if newEvents.CallParams.Name != "" {
		a.CallParams = newEvents.CallParams
	}

	if newEvents.Calls.Name != "" {
		a.Calls = newEvents.Calls
	}

	if newEvents.OnClose.Name != "" {
		a.OnClose = newEvents.OnClose
	}

	if newEvents.OnTransferCall.Name != "" {
		a.OnTransferCall = newEvents.OnTransferCall
	}

	if newEvents.TransferSucceed.Name != "" {
		a.TransferSucceed = newEvents.TransferSucceed
	}

	if newEvents.TransferFailed.Name != "" {
		a.TransferFailed = newEvents.TransferFailed
	}

	if newEvents.TransferCallReturned.Name != "" {
		a.TransferCallReturned = newEvents.TransferCallReturned
	}

	if newEvents.SetSessionModeResponse.Name != "" {
		a.SetSessionModeResponse = newEvents.SetSessionModeResponse
	}

	if newEvents.CallParamsUpdated.Name != "" {
		a.CallParamsUpdated = newEvents.CallParamsUpdated
	}

	if newEvents.LocalParamsUpdated.Name != "" {
		a.LocalParamsUpdated = newEvents.LocalParamsUpdated
	}

	if newEvents.Conferences.Name != "" {
		a.Conferences = newEvents.Conferences
	}

	if newEvents.Ok.Name != "" {
		a.Ok = newEvents.Ok
	}

	if newEvents.Error.Name != "" {
		a.Error = newEvents.Error
	}

	if newEvents.ParseError.Name != "" {
		a.ParseError = newEvents.ParseError
	}
}
