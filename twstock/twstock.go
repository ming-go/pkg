package twstock

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/ming-go/pkg/mhttp"
	"github.com/ming-go/pkg/mstring"
	"github.com/ming-go/pkg/mtime"
)

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
	mhc mhttp.Client
}

var mhc *mhttp.Client = mhttp.NewClient()

const urlGetStockInfo string = "https://mis.twse.com.tw/stock/api/getStockInfo.jsp?ex_ch=%s&_=%d"

func joinStockCodes(stockCodes ...string) string {
	buf := make([]byte, 0, len(stockCodes)*(4+7)+(len(stockCodes)-1))

	flag := false

	for _, stockCode := range stockCodes {
		if !flag {
			flag = true
		} else {
			buf = append(buf, '|')
		}

		buf = append(buf, "tse_"...)
		buf = append(buf, stockCode...)
		buf = append(buf, ".tw"...)
	}

	return mstring.BytesToString(buf)
}

func GetStockInfo(ctx context.Context, stockTime time.Time, stockCodes ...string) (*TwStock, error) {
	url := fmt.Sprintf(urlGetStockInfo, joinStockCodes(stockCodes...), mtime.UnixMilli(stockTime)) // TODO, func
	log.Println(url)
	mr, err := mhttp.HttpResponseToMHttpResponse(
		mhc.GetWithContext(ctx, url, nil, nil),
	)
	if err != nil {
		return nil, err
	}

	twStock := TwStock{}
	err = json.Unmarshal(mr.RespBody, &twStock)
	if err != nil {
		return nil, err
	}

	return &twStock, nil
}
