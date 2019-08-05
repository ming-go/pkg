package twstock

type TwStock struct {
	MsgArray []struct {
		Ts    string `json:"ts"`
		Fv    string `json:"fv"`
		Tk0   string `json:"tk0"`
		Tk1   string `json:"tk1"`
		Oa    string `json:"oa"`
		Ob    string `json:"ob"`
		Tlong string `json:"tlong"`
		Ot    string `json:"ot"`
		F     string `json:"f"`
		Ex    string `json:"ex"`
		G     string `json:"g"`
		Ov    string `json:"ov"`
		D     string `json:"d"`
		It    string `json:"it"`
		B     string `json:"b"`
		C     string `json:"c"`
		Mt    string `json:"mt"`
		A     string `json:"a"`
		N     string `json:"n"`
		O     string `json:"o"`
		L     string `json:"l"`
		Oz    string `json:"oz"`
		Io    string `json:"io"`
		H     string `json:"h"`
		IP    string `json:"ip"`
		I     string `json:"i"`
		W     string `json:"w"`
		V     string `json:"v"`
		U     string `json:"u"`
		T     string `json:"t"`
		S     string `json:"s"`
		Pz    string `json:"pz"`
		Tv    string `json:"tv"`
		P     string `json:"p"`
		Nf    string `json:"nf"`
		Ch    string `json:"ch"`
		Z     string `json:"z"`
		Y     string `json:"y"`
		Ps    string `json:"ps"`
	} `json:"msgArray"`
	UserDelay int    `json:"userDelay"`
	Rtmessage string `json:"rtmessage"`
	Referer   string `json:"referer"`
	QueryTime struct {
		SysTime           string `json:"sysTime"`
		SessionLatestTime int    `json:"sessionLatestTime"`
		SysDate           string `json:"sysDate"`
		SessionFromTime   int    `json:"sessionFromTime"`
		StockInfoItem     int    `json:"stockInfoItem"`
		ShowChart         bool   `json:"showChart"`
		SessionStr        string `json:"sessionStr"`
		StockInfo         int    `json:"stockInfo"`
	} `json:"queryTime"`
	Rtcode string `json:"rtcode"`
}

type Client struct {
	mhttp.client
}

func getStockInfo(stockCode int) {
}
