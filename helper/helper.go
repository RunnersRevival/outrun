package helper

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/RunnersRevival/outrun/config"
	"github.com/RunnersRevival/outrun/cryption"
	"github.com/RunnersRevival/outrun/db"
	"github.com/RunnersRevival/outrun/emess"
	"github.com/RunnersRevival/outrun/netobj"
	"github.com/RunnersRevival/outrun/netobj/constnetobjs"
	"github.com/RunnersRevival/outrun/requests"
	"github.com/RunnersRevival/outrun/responses"
	"github.com/RunnersRevival/outrun/responses/responseobjs"
	"github.com/RunnersRevival/outrun/status"
)

const (
	PrefixErr            = "ERR"
	PrefixOut            = "OUT"
	PrefixWarn           = "WARN"
	PrefixUncatchableErr = "UNCATCHABLE ERR"
	PrefixDebugOut       = "DEBUG (OUT)"

	LogOutBase = "[%s] (%s) %s\n"
	LogErrBase = "[%s] (%s) %s: %s\n"

	InternalServerError = "Internal Server Error"
	BadRequest          = "Bad Request"

	DefaultIV                           = "FoundDeadInMiami"
	RandomizeIV                         = true // highly experimental option; may slow Outrun down slightly
	RandomizeIVAfterUpToHowManyRequests = 16
)

var (
	RequestsLeft = 0
	CurrentIV    = DefaultIV
)

type Helper struct {
	CallerName string
	RespW      http.ResponseWriter
	Request    *http.Request
}

func MakeHelper(callerName string, r http.ResponseWriter, request *http.Request) *Helper {
	return &Helper{
		callerName,
		r,
		request,
	}
}

func (r *Helper) GetGameRequest() []byte {
	recv := cryption.GetReceivedMessage(r.Request)
	return recv
}
func (r *Helper) SendResponse(i interface{}) error {
	out, err := json.Marshal(i)
	if err != nil {
		return err
	}
	r.Respond(out)
	return nil
}
func (r *Helper) SendInsecureResponse(i interface{}) error {
	r.RespondInsecure(i)
	return nil
}
func (r *Helper) RespondRaw(out []byte, secureFlag, iv string) {
	if config.CFile.LogAllResponses {
		nano := time.Now().UnixNano()
		nanoStr := strconv.Itoa(int(nano))
		filename := r.Request.RequestURI + "--" + nanoStr
		filename = strings.ReplaceAll(filename, ".", "-")
		filename = strings.ReplaceAll(filename, "/", "-") + ".txt"
		filepath := "logging/all_responses/" + filename
		r.Out("DEBUG: Saving request to " + filepath)
		err := ioutil.WriteFile(filepath, out, 0644)
		if err != nil {
			r.Out("DEBUG ERROR: Unable to write file '" + filepath + "'")
		}
	}
	response := map[string]string{}
	if secureFlag != "0" && secureFlag != "1" {
		r.Warn("Improper secureFlag in call to RespondRaw!")
	}
	response["secure"] = secureFlag
	response["key"] = iv
	if secureFlag == "1" {
		encrypted := cryption.Encrypt(out, cryption.EncryptionKey, []byte(iv))
		encryptedBase64 := cryption.B64Encode(encrypted)
		response["param"] = encryptedBase64
	} else {
		response["param"] = string(out)
	}
	toClient, err := json.Marshal(response)
	if err != nil {
		r.InternalErr("Error marshalling in RespondRaw", err)
		return
	}
	r.RespW.Write(toClient)
}
func (r *Helper) Respond(out []byte) {
	iv := DefaultIV
	if RandomizeIV {
		if RequestsLeft <= 0 {
			RequestsLeft = rand.Intn(RandomizeIVAfterUpToHowManyRequests)
			r.Out("DEBUG: Rerolling random IV...")
			randChar := func(charset string, length int64) string {
				runes := []rune(charset)
				final := make([]rune, length)
				for i := range final {
					final[i] = runes[rand.Intn(len(runes))]
				}
				return string(final)
			}
			CurrentIV = randChar("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890", 16)
		} else {
			RequestsLeft--
		}
		iv = CurrentIV
	}
	r.RespondRaw(out, "1", iv)
}
func (r *Helper) RespondInsecure(out interface{}) {
	if config.CFile.LogAllResponses {
		nano := time.Now().UnixNano()
		nanoStr := strconv.Itoa(int(nano))
		filename := r.Request.RequestURI + "--" + nanoStr
		filename = strings.ReplaceAll(filename, ".", "-")
		filename = strings.ReplaceAll(filename, "/", "-") + ".txt"
		filepath := "logging/all_responses/" + filename
		r.Out("DEBUG: Saving request to " + filepath)
		outS, err := json.Marshal(out)
		if err != nil {
			r.Out("DEBUG ERROR: Unable to marshal param")
		} else {
			err := ioutil.WriteFile(filepath, outS, 0644)
			if err != nil {
				r.Out("DEBUG ERROR: Unable to write file '" + filepath + "'")
			}
		}
	}
	response := map[string]interface{}{"secure": "0", "param": out}
	toClient, err := json.Marshal(response)
	if err != nil {
		r.InternalErr("Error marshalling in RespondInsecure", err)
		return
	}
	r.RespW.Write(toClient)
}
func (r *Helper) Out(s string, a ...interface{}) {
	msg := fmt.Sprintf(s, a...)
	log.Printf(LogOutBase, PrefixOut, r.CallerName, msg)
}
func (r *Helper) DebugOut(s string, a ...interface{}) {
	if config.CFile.DebugPrints {
		msg := fmt.Sprintf(s, a...)
		log.Printf(LogOutBase, PrefixDebugOut, r.CallerName, msg)
	}
}
func (r *Helper) Warn(s string, a ...interface{}) {
	msg := fmt.Sprintf(s, a...)
	log.Printf(LogOutBase, PrefixWarn, r.CallerName, msg)
}
func (r *Helper) WarnErr(msg string, err error) {
	log.Printf(LogErrBase, PrefixWarn, r.CallerName, msg, err.Error())
}
func (r *Helper) Uncatchable(msg string) {
	log.Printf(LogOutBase, PrefixOut, r.CallerName, msg)
}
func (r *Helper) InternalErr(msg string, err error) {
	log.Printf(LogErrBase, PrefixErr, r.CallerName, msg, err.Error())
	r.RespW.WriteHeader(http.StatusInternalServerError)
	//r.RespW.Write([]byte(BadRequest))
	r.SendResponse(responses.NewBaseResponse(r.BaseInfo(emess.OK, status.InternalServerError)))
}
func (r *Helper) Err(msg string, err error) {
	log.Printf(LogErrBase, PrefixErr, r.CallerName, msg, err.Error())
	r.RespW.WriteHeader(http.StatusBadRequest)
	//r.RespW.Write([]byte(BadRequest))
	r.SendResponse(responses.NewBaseResponse(r.BaseInfo(emess.OK, status.ClientError)))
}
func (r *Helper) ErrRespond(msg string, err error, response string) {
	// TODO: remove if never used in stable builds
	log.Printf(LogErrBase, PrefixErr, r.CallerName, msg, err.Error())
	r.RespW.WriteHeader(http.StatusInternalServerError) // ideally include an option for this, but for now it's inconsequential
	r.RespW.Write([]byte(response))
}
func (r *Helper) InternalFatal(msg string, err error) {
	log.Fatalf(LogErrBase, PrefixErr, r.CallerName, msg, err.Error())
	r.RespW.WriteHeader(http.StatusInternalServerError)
	//	r.RespW.Write([]byte(BadRequest))
	r.SendResponse(responses.NewBaseResponse(r.BaseInfo(emess.OK, status.InternalServerError)))
}
func (r *Helper) Fatal(msg string, err error) {
	log.Fatalf(LogErrBase, PrefixErr, r.CallerName, msg, err.Error())
	r.RespW.WriteHeader(http.StatusBadRequest)
	//	r.RespW.Write([]byte(BadRequest))
	r.SendResponse(responses.NewBaseResponse(r.BaseInfo(emess.OK, status.ClientError)))
}
func (r *Helper) BaseInfo(em string, statusCode int64) responseobjs.BaseInfo {
	return responseobjs.NewBaseInfo(em, statusCode)
}
func (r *Helper) InvalidRequest() {
	//	r.RespW.WriteHeader(http.StatusBadRequest)
	//	r.RespW.Write([]byte(BadRequest))
	r.SendResponse(responses.NewBaseResponse(r.BaseInfo(emess.OK, status.ClientError)))
}
func (r *Helper) CheckSession(sendResponseOnFalseResult bool) bool {
	recv := r.GetGameRequest()
	request := requests.Base{
		RevivalVerID: 0,
	}
	err := json.Unmarshal(recv, &request)
	if err != nil {
		// likely malformed request
		if sendResponseOnFalseResult {
			r.RespW.WriteHeader(http.StatusBadRequest)
			r.SendResponse(responses.NewBaseResponse(r.BaseInfo(emess.OK, status.ClientError)))
		}
		return false
	}
	sid := []byte(request.SessionID)
	validsession, err := db.BoltIsValidSessionID(sid)
	if err != nil {
		if sendResponseOnFalseResult {
			r.RespW.WriteHeader(http.StatusInternalServerError)
			r.SendResponse(responses.NewBaseResponse(r.BaseInfo(emess.OK, status.InternalServerError)))
		}
		return false
	}
	if !validsession {
		r.DebugOut("Invalid session ID!")
		if sendResponseOnFalseResult {
			r.SendResponse(responses.NewBaseResponse(r.BaseInfo(emess.OK, status.ExpiredSession)))
		}
		return false
	}
	return true
}
func (r *Helper) GetCallingPlayer() (netobj.Player, error) {
	// Powerful function to get the player directly from the response
	recv := r.GetGameRequest()
	var request requests.Base
	err := json.Unmarshal(recv, &request)
	if err != nil {
		return constnetobjs.BlankPlayer, err
	}
	sid := request.SessionID
	player, err := db.BoltGetPlayerBySessionID(sid)
	if err != nil {
		return constnetobjs.BlankPlayer, err
	}
	if config.CFile.PrintPlayerNames {
		r.Out("Player '" + player.Username + "' (" + player.ID + ")")
	}
	return player, nil
}
func (r *Helper) GetCallingPlayerID() (string, error) {
	recv := r.GetGameRequest()
	var request requests.Base
	err := json.Unmarshal(recv, &request)
	if err != nil {
		return "0", err
	}
	sid := request.SessionID
	uid, err := db.BoltGetPlayerIDBySessionID(sid)
	if err != nil {
		return "0", err
	}
	return uid, nil
}
func (r *Helper) GetSessionID() (string, bool) {
	recv := r.GetGameRequest()
	var request requests.Base
	err := json.Unmarshal(recv, &request)
	if err != nil {
		// likely malformed request
		return "", false
	}
	sid := []byte(request.SessionID)
	validsession, err := db.BoltIsValidSessionID(sid)
	if err != nil {
		return string(sid), false
	}
	return string(sid), validsession
}
