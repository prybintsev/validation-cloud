package samples

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

const gethUrl = "https://javelin.validationcloud.io/v1/GDqpM9s4ED8zhy4awDsrPjCu_9SKeI4EsGC5q9AKg34"

type HeightSamplesCollector struct {
	frequency   time.Duration
	samplesRepo Repository
}

type Repository interface {
	InsertSample(ctx context.Context, height uint64, createdAt time.Time) error
	DeleteOldSamples(ctx context.Context, deleteAfter time.Duration) error
}

func NewHeightSamplesCollector(frequency time.Duration, samplesRepo Repository) HeightSamplesCollector {
	return HeightSamplesCollector{frequency: frequency, samplesRepo: samplesRepo}
}

func (h HeightSamplesCollector) Run(ctx context.Context) {
	log.Info("starting height samples collection")
	h.collectAndStoreSample(ctx)
	for {
		select {
		case <-time.After(h.frequency):
			h.collectAndStoreSample(ctx)
		case <-ctx.Done():
			log.Info("gracefully stopping samples collection")
			return
		}
	}
}

type SampleResp struct {
	Result string `json:"result"`
}

func (h HeightSamplesCollector) collectAndStoreSample(ctx context.Context) {
	log.Info("collecting blockchain height sample")
	height, err := getBlockNumberFromGeth()
	if err != nil {
		log.WithError(err).Error("failed to retrieve blockchain height from geth")
		return
	}
	now := time.Now().UTC()
	err = h.samplesRepo.InsertSample(ctx, height, now)
	if err != nil {
		log.WithError(err).Error("failed to insert sample to DB")
		return
	}

	// Delete samples older than one hour
	err = h.samplesRepo.DeleteOldSamples(ctx, time.Hour)
	if err != nil {
		log.WithError(err).Error("failed to delete old samples")
		return
	}
}

func getBlockNumberFromGeth() (uint64, error) {
	req := `{"jsonrpc": "2.0", "method": "eth_blockNumber", "params": [], "id": 1}`
	reqReader := strings.NewReader(req)
	resp, err := http.Post(gethUrl, "application/json", reqReader)
	if err != nil {
		return 0, errors.New("failed to send request to geth")
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, errors.New("failed to read geth response")
	}
	var respStruct SampleResp
	err = json.Unmarshal(respBody, &respStruct)
	if err != nil {
		return 0, errors.New("failed to unmarshal geth response")
	}

	blockNumBArr, err := hex.DecodeString(strings.TrimLeft(respStruct.Result, "0x"))
	if err != nil {
		return 0, errors.New("could not convert block number hex")
	}
	exp := 0
	res := uint64(0)
	for idx := len(blockNumBArr) - 1; idx >= 0; idx-- {
		res += uint64(blockNumBArr[idx]) << exp
		exp += 8
	}
	return res, nil
}
