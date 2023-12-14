package cti

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/pochtalexa/go-cti-middleware/internal/server/config"
	"github.com/pochtalexa/go-cti-middleware/internal/server/storage"
	"github.com/pochtalexa/go-cti-middleware/internal/server/ws"
	"github.com/rs/zerolog/log"
	"math/rand"
	"net/url"
	"slices"
	"strconv"
	"time"
)

var Conn *websocket.Conn

func getRandRid() string {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	// generate a random 8-digit ID
	id := ""
	for i := 0; i < 8; i++ {
		id += strconv.Itoa(rand.Intn(10))
	}

	return id
}

func Init() error {
	var err error

	uCTI := url.URL{
		Scheme: config.ServerConfig.CtiAPI.Scheme,
		Host:   config.ServerConfig.CtiAPI.Host + ":" + config.ServerConfig.CtiAPI.Port,
		Path:   config.ServerConfig.CtiAPI.Path,
	}
	log.Info().Str("ws connecting to", uCTI.String()).Msg("")

	config.ServerConfig.WsConn, _, err = websocket.DefaultDialer.Dial(uCTI.String(), nil)
	if err != nil {
		return fmt.Errorf("websocket dial: %w", err)
	}
	log.Info().Str("ws connected", uCTI.String()).Msg("")

	return nil
}

func InitCTISess() error {

	messInitConn := storage.NewWsCommand()
	messInitConn.Rid = getRandRid()
	messInitConn.Name = "SetProtocolVersion"
	messInitConn.ProtocolVersion = "13"

	body, err := json.Marshal(messInitConn)
	if err != nil {
		return fmt.Errorf("messInitConn - marshal: %w", err)
	}

	if err := ws.SendCommand(config.ServerConfig.WsConn, body); err != nil {
		return fmt.Errorf("initCTISess: %w", err)
	}

	log.Info().Msg("InitCTISess - ok")

	return nil
}

func AttachUser(login string) error {

	messAttachUser := storage.NewWsCommand()
	messAttachUser.Rid = getRandRid()
	messAttachUser.Name = "AttachToUser"
	messAttachUser.Login = login

	body, err := json.Marshal(messAttachUser)
	if err != nil {
		return fmt.Errorf("messAttachUser - marshal: %w", err)
	}

	if err := ws.SendCommand(config.ServerConfig.WsConn, body); err != nil {
		return fmt.Errorf("AttachUser: %w", err)
	}

	log.Info().Str("login", login).Msg("AttachUser - ok")
	return nil
}

func ChageStatus(login string, status string) error {
	if !slices.Contains(storage.AgentsInfo.ValidStatuses, status) {
		return fmt.Errorf("ChageStatus: bad status val: %s", status)
	}

	messChageStatus := storage.NewWsCommand()
	messChageStatus.Name = "ChangeUserState"

	messChageStatus.Rid = getRandRid()
	messChageStatus.Login = login
	messChageStatus.State = status

	body, err := json.Marshal(messChageStatus)
	if err != nil {
		return fmt.Errorf("ChageStatus - marshal: %w", err)
	}

	if err := ws.SendCommand(config.ServerConfig.WsConn, body); err != nil {
		return fmt.Errorf("ChageStatus: %w", err)
	}

	log.Info().Str("login", login).Msg("ChageStatus - ok")

	return nil
}

func Answer(login string, cid int) error {

	messAnswer := storage.NewWsCommand()
	messAnswer.Rid = getRandRid()
	messAnswer.Name = "Answer"
	messAnswer.Login = login
	messAnswer.Cid = cid

	body, err := json.Marshal(messAnswer)
	if err != nil {
		return fmt.Errorf("answer - marshal: %w", err)
	}

	if err := ws.SendCommand(config.ServerConfig.WsConn, body); err != nil {
		return fmt.Errorf("Answer: %w", err)
	}

	log.Info().Str("login", login).Msg("Answer - ok")

	return nil
}

func Hangup(login string, cid int) error {

	messHangup := storage.NewWsCommand()
	messHangup.Rid = getRandRid()
	messHangup.Name = "Hangup"
	messHangup.Login = login
	messHangup.Cid = cid

	body, err := json.Marshal(messHangup)
	if err != nil {
		return fmt.Errorf("hangup - marshal: %w", err)
	}

	if err := ws.SendCommand(config.ServerConfig.WsConn, body); err != nil {
		return fmt.Errorf("Hangup: %w", err)
	}

	log.Info().Str("login", login).Msg("Hangup - ok")

	return nil
}

func Mute(login string, cid int, on bool) error {

	messMute := storage.NewWsCommandMute()
	messMute.Rid = getRandRid()
	messMute.Name = "Mute"
	messMute.Login = login
	messMute.Cid = cid
	messMute.On = on

	body, err := json.Marshal(messMute)
	if err != nil {
		return fmt.Errorf("mute - marshal: %w", err)
	}

	log.Info().Str("messMute", messMute.String()).Msg("Mute")

	if err := ws.SendCommand(config.ServerConfig.WsConn, body); err != nil {
		return fmt.Errorf("Mute: %w", err)
	}

	log.Info().Str("login", login).Msg("Mute - ok")

	return nil
}
