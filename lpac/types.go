package lpac

type Information struct {
	EID         string `json:"eid,omitempty"`
	DefaultSMDS string `json:"default_smds,omitempty"`
	DefaultSMDP string `json:"default_smdp,omitempty"`
}

type Profile struct {
	ICCID               string       `json:"iccid,omitempty"`
	ISDPAID             string       `json:"isdpAid,omitempty"`
	ProfileName         string       `json:"profileName,omitempty"`
	ProfileState        ProfileState `json:"profileState,omitempty"`
	ProfileClass        ProfileClass `json:"profileClass,omitempty"`
	ServiceProviderName string       `json:"serviceProviderName,omitempty"`
}

type Notification struct {
	Index      int    `json:"seqNumber,omitempty"`
	Management int    `json:"profileManagementOperation,omitempty"`
	Address    string `json:"notificationAddress,omitempty"`
	ICCID      string `json:"iccid,omitempty"`
}
