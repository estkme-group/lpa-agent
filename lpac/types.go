package lpac

<<<<<<< HEAD
import (
	"encoding/json"
)
=======
import "encoding/json"
>>>>>>> 0d9746d (fix: don't ignore lpac)

type lpaResponse struct {
	Type    string      `json:"type"`
	Payload *lpaPayload `json:"payload"`
}

type lpaPayload struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

type Information struct {
	EID         string `json:"eid"`
	DefaultSMDS string `json:"default_smds"`
	DefaultSMDP string `json:"default_smdp,omitempty"`
}

type Profile struct {
	ICCID               string       `json:"iccid"`
	ISDPAID             string       `json:"isdpAid"`
	DisplayName         string       `json:"profileNickname,omitempty"`
	BuiltinName         string       `json:"profileName"`
	State               ProfileState `json:"profileState"`
	Class               ProfileClass `json:"profileClass"`
	ServiceProviderName string       `json:"serviceProviderName"`
}

type Notification struct {
	Index      int    `json:"seqNumber"`
	Management int    `json:"profileManagementOperation"`
	Address    string `json:"notificationAddress"`
	ICCID      string `json:"iccid"`
}
