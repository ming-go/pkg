package twstock

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestJoinStockCodes(t *testing.T) {
	assert.Equal(t, "tse_2330.tw|tse_3008.tw", joinStockCodes("2330", "3008"))
	assert.Equal(t, "tse_2330.tw", joinStockCodes("2330"))
	assert.Equal(t, "tse_2330.tw|tse_3008.tw|tse_2498.tw", joinStockCodes("2330", "3008", "2498"))
}

func TestGetStockInfo(t *testing.T) {
	twStock, err := getStockInfo(context.Background(), time.Now(), "2330", "3008", "2498")
	if err != nil {
		log.Println(err)
	}

	log.Println(twStock)
}
