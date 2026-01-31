package filemaker

type ResponseData struct {
	Response Response  `json:"response"`
	Messages []Message `json:"messages"`
}

type Message struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type Response struct {
	RecordID       string                     `json:"recordId,omitempty"`
	ModID          string                     `json:"modId,omitempty"`
	Token          string                     `json:"token,omitempty"`
	DataInfo       DataInfo                   `json:"dataInfo,omitempty"`
	Data           []Datum                    `json:"data,omitempty"`
	ProductInfo    *ProductInfo               `json:"productInfo,omitempty"`
	Databases      []Database                 `json:"databases,omitempty"`
	Layouts        []Layout                   `json:"layouts,omitempty"`
	Scripts        []Script                   `json:"scripts,omitempty"`
	ScriptResult   string                     `json:"scriptResult,omitempty"`
	ScriptError    string                     `json:"scriptError,omitempty"`
	FieldMetaData  []FieldMetaData            `json:"fieldMetaData,omitempty"`
	PortalMetaData map[string][]FieldMetaData `json:"portalMetaData,omitempty"`
	ValueLists     []ValueList                `json:"valueLists,omitempty"`
}

type ProductInfo struct {
	Name            string `json:"name"`
	BuildDate       string `json:"buildDate"`
	Version         string `json:"version"`
	DateFormat      string `json:"dateFormat"`
	TimeFormat      string `json:"timeFormat"`
	TimeStampFormat string `json:"timeStampFormat"`
}

type Database struct {
	Name string `json:"name"`
}

type Layout struct {
	Name string `json:"name"`
}

type Script struct {
	Name     string `json:"name"`
	IsFolder bool   `json:"isFolder,omitempty"`
}

type FieldMetaData struct {
	Name            string `json:"name"`
	Type            string `json:"type"`
	DisplayType     string `json:"displayType,omitempty"`
	Result          string `json:"result,omitempty"`
	Global          bool   `json:"global,omitempty"`
	AutoEnter       bool   `json:"autoEnter,omitempty"`
	FourDigitYear   bool   `json:"fourDigitYear,omitempty"`
	MaxRepeat       int    `json:"maxRepeat,omitempty"`
	MaxCharacters   int    `json:"maxCharacters,omitempty"`
	NotEmpty        bool   `json:"notEmpty,omitempty"`
	Numeric         bool   `json:"numeric,omitempty"`
	TimeOfDay       bool   `json:"timeOfDay,omitempty"`
	RepetitionStart int    `json:"repetitionStart,omitempty"`
	RepetitionEnd   int    `json:"repetitionEnd,omitempty"`
}

type ValueList struct {
	Name   string `json:"name"`
	Type   string `json:"type"`
	Values []struct {
		Value   string `json:"value"`
		Display string `json:"display"`
	} `json:"values"`
}

type Datum struct {
	FieldData  map[string]interface{}   `json:"fieldData,omitempty"`
	PortalData map[string][]interface{} `json:"portalData,omitempty"`
	RecordID   string                   `json:"recordId,omitempty"`
	ModID      string                   `json:"modId,omitempty"`
}

type DataInfo struct {
	Database         string `json:"database,omitempty"`
	Layout           string `json:"layout,omitempty"`
	Table            string `json:"table,omitempty"`
	TotalRecordCount int64  `json:"totalRecordCount,omitempty"`
	FoundCount       int64  `json:"foundCount,omitempty"`
	ReturnedCount    int64  `json:"returnedCount,omitempty"`
}
