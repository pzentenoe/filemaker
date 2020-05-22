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
	Token    string   `json:"token,omitempty"`
	DataInfo DataInfo `json:"dataInfo,omitempty"`
	Data     []Datum  `json:"data,omitempty"`
}

type Datum struct {
	FieldData  interface{} `json:"fieldData,omitempty"`
	PortalData interface{} `json:"portalData,omitempty"`
	RecordID   string      `json:"recordId,omitempty"`
	ModID      string      `json:"modId,omitempty"`
}

type DataInfo struct {
	Database         string `json:"database,omitempty"`
	Layout           string `json:"layout,omitempty"`
	Table            string `json:"table,omitempty"`
	TotalRecordCount int64  `json:"totalRecordCount,omitempty"`
	FoundCount       int64  `json:"foundCount,omitempty"`
	ReturnedCount    int64  `json:"returnedCount,omitempty"`
}
