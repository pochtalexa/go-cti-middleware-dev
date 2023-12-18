package storage

import (
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"sync"
)

var (
	AgentsInfo = NewAgentsInfo()
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

// StAgentsInfo мапа с ключом - логин оператора
type StAgentsInfo struct {
	Events        map[string]StAgentEvents // евенты - по операторам
	Updated       map[string]bool          // были ли обновления - по операторам
	Mutex         map[string]*sync.RWMutex // мьютекс чтения/записи событий по агенту
	ValidStatuses []string
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

// StWsCommandMute команда в сторону CTI API
// нужна отдельная структура, т.к. не анмаршаллит поле bool если оно не обязательное и значение false
type StWsCommandMute struct {
	Rid   string `json:"rid,omitempty"`
	Name  string `json:"name"`
	Login string `json:"login"`
	Cid   int    `json:"cid"`
	On    bool   `json:"on"`
}

// StWsEvent событие или ответ от CTI API
type StWsEvent struct {
	Name       string                 `json:"name"`
	Login      string                 `json:"login,omitempty"`
	Body       map[string]interface{} `json:"-"`
	ErrorNames []string               `json:"-"`
}

func NewWsCommand() *StWsCommand {
	return &StWsCommand{}
}

func NewWsCommandMute() *StWsCommandMute {
	return &StWsCommandMute{}
}

func NewWsEvent() *StWsEvent {
	return &StWsEvent{
		ErrorNames: []string{"Error", "ParseError"},
	}
}

func NewAgentsInfo() *StAgentsInfo {
	return &StAgentsInfo{
		Events:        make(map[string]StAgentEvents),
		Updated:       make(map[string]bool),
		Mutex:         make(map[string]*sync.RWMutex),
		ValidStatuses: []string{"normal", "dnd", "away"},
	}
}

func (a *StAgentsInfo) SetEvent(event *StWsEvent, message []byte) error {
	var eventLogin string

	if event.Login == "" {
		eventLogin = "noName"
	} else {
		eventLogin = event.Login
	}

	mutex, ok := a.Mutex[eventLogin]
	if !ok {
		var m sync.RWMutex
		a.Mutex[eventLogin] = &m
		mutex = a.Mutex[eventLogin]
	}
	mutex.Lock()
	defer mutex.Unlock()

	// сохраняем текущие события по оператору и обновляем
	curEvents := a.Events[eventLogin]

	switch event.Name {
	case "UserState":
		if err := json.Unmarshal(message, &curEvents.UserState); err != nil {
			return fmt.Errorf("can not Unmarshal body for name: %v", event.Name)
		}
	case "NewCall":
		if err := json.Unmarshal(message, &curEvents.NewCall); err != nil {
			return fmt.Errorf("can not Unmarshal body for name: %v", event.Name)
		}
	case "LocalParamsUpdated":
		if err := json.Unmarshal(message, &curEvents.LocalParamsUpdated); err != nil {
			return fmt.Errorf("can not Unmarshal body for name: %v", event.Name)
		}
	case "CallParamsUpdated":
		if err := json.Unmarshal(message, &curEvents.CallParamsUpdated); err != nil {
			return fmt.Errorf("can not Unmarshal body for name: %v", event.Name)
		}
	case "SetSessionModeResponse":
		if err := json.Unmarshal(message, &curEvents.SetSessionModeResponse); err != nil {
			return fmt.Errorf("can not Unmarshal body for name: %v", event.Name)
		}
	case "TransferCallReturned":
		if err := json.Unmarshal(message, &curEvents.TransferCallReturned); err != nil {
			return fmt.Errorf("can not Unmarshal body for name: %v", event.Name)
		}
	case "TransferFailed":
		if err := json.Unmarshal(message, &curEvents.TransferFailed); err != nil {
			return fmt.Errorf("can not Unmarshal body for name: %v", event.Name)
		}
	case "TransferSucceed":
		if err := json.Unmarshal(message, &curEvents.TransferSucceed); err != nil {
			return fmt.Errorf("can not Unmarshal body for name: %v", event.Name)
		}
	case "OnTransferCall":
		if err := json.Unmarshal(message, &curEvents.OnTransferCall); err != nil {
			return fmt.Errorf("can not Unmarshal body for name: %v", event.Name)
		}
	case "OnClose":
		if err := json.Unmarshal(message, &curEvents.OnClose); err != nil {
			return fmt.Errorf("can not Unmarshal body for name: %v", event.Name)
		}
	case "Calls":
		if err := json.Unmarshal(message, &curEvents.Calls); err != nil {
			return fmt.Errorf("can not Unmarshal body for name: %v", event.Name)
		}
	case "CallParams":
		if err := json.Unmarshal(message, &curEvents.CallParams); err != nil {
			return fmt.Errorf("can not Unmarshal body for name: %v", event.Name)
		}
	case "CurrentCall":
		if err := json.Unmarshal(message, &curEvents.CurrentCall); err != nil {
			return fmt.Errorf("can not Unmarshal body for name: %v", event.Name)
		}
	case "CallStatus":
		if err := json.Unmarshal(message, &curEvents.CallStatus); err != nil {
			return fmt.Errorf("can not Unmarshal body for name: %v", event.Name)
		}
	case "Conferences":
		if err := json.Unmarshal(message, &curEvents.Conferences); err != nil {
			return fmt.Errorf("can not Unmarshal body for name: %v", event.Name)
		}
	case "Ok":
		if err := json.Unmarshal(message, &curEvents.Ok); err != nil {
			return fmt.Errorf("can not Unmarshal body for name: %v", event.Name)
		}
	case "Error":
		if err := json.Unmarshal(message, &curEvents.Error); err != nil {
			return fmt.Errorf("can not Unmarshal body for name: %v", event.Name)
		}
	case "ParseError":
		if err := json.Unmarshal(message, &curEvents.ParseError); err != nil {
			return fmt.Errorf("can not Unmarshal body for name: %v", event.Name)
		}
	default:
		return fmt.Errorf("can not find case for key %v", event.Name)
	}

	a.Events[eventLogin] = curEvents
	a.Updated[eventLogin] = true

	return nil
}

func (a *StAgentsInfo) DropAgentEvents(login string) {
	a.Events[login] = StAgentEvents{}
}

func (a *StAgentsInfo) IsUpdated(login string) (value bool, isKey bool) {
	value, isKey = a.Updated[login]
	return
}

func (a *StAgentsInfo) SetUpdated(login string, val bool) {
	a.Updated[login] = val
}

func (a *StAgentsInfo) GetMutex(login string) (mutex *sync.RWMutex, isKey bool) {
	mutex, isKey = a.Mutex[login]
	return
}

func (a *StAgentsInfo) GetEvents(login string) (events StAgentEvents, isKey bool) {
	events, isKey = a.Events[login]
	return
}

func (w *StWsCommand) String() string {
	var tempMap map[string]interface{}

	jsonData, err := json.Marshal(w)
	if err != nil {
		log.Error().Err(err).Msg("Marshal StWsCommand String")
		return ""
	}

	err = json.Unmarshal(jsonData, &tempMap)
	if err != nil {
		log.Error().Err(err).Msg("Unmarshal StWsCommand String")
		return ""
	}

	return fmt.Sprintln(tempMap)
}

func (w *StWsCommandMute) String() string {
	var tempMap map[string]interface{}

	jsonData, err := json.Marshal(w)
	if err != nil {
		log.Error().Err(err).Msg("Marshal StWsCommandMute String")
		return ""
	}

	err = json.Unmarshal(jsonData, &tempMap)
	if err != nil {
		log.Error().Err(err).Msg("Unmarshal StWsCommandMute String")
		return ""
	}

	return fmt.Sprintln(tempMap)
}
