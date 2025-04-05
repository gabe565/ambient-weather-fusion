package discovery

type Payload struct {
	AvailabilityTopic string              `json:"avty_t"`
	Device            Device              `json:"dev"`
	Origin            Origin              `json:"o"`
	StateTopic        string              `json:"stat_t"`
	Components        map[Topic]Component `json:"cmps"`
}

type Device struct {
	Identifiers string `json:"ids"`
	Name        string `json:"name"`
	SWVersion   string `json:"sw"`
}

type Origin struct {
	Name       string `json:"name"`
	SWVersion  string `json:"sw"`
	SupportURL string `json:"url"`
}

type Component struct {
	Name                      string      `json:"name,omitempty"`
	Platform                  Platform    `json:"p,omitempty"`
	ObjectID                  string      `json:"obj_id,omitempty"`
	UniqueID                  string      `json:"uniq_id,omitempty"`
	ValueTemplate             string      `json:"val_tpl,omitempty"`
	UnitOfMeasurement         Unit        `json:"unit_of_meas,omitempty"`
	DeviceClass               DeviceClass `json:"dev_cla,omitempty"`
	StateClass                StateClass  `json:"stat_cla,omitempty"`
	SuggestedDisplayPrecision int         `json:"sug_dsp_prc,omitempty"`
	EnabledByDefault          *bool       `json:"en,omitempty"`
	Icon                      string      `json:"ic,omitempty"`
}
