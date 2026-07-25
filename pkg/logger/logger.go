package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type lokiPushPayload struct {
	Streams []lokiStream `json:"streams"`
}

type lokiStream struct {
	Stream map[string]string `json:"stream"`
	Values [][2]string       `json:"values"`
}

type lokiWriter struct {
	lokiURL     string
	serviceName string
	client      *http.Client
	logCh       chan [2]string
	stopCh      chan struct{}
	wg          sync.WaitGroup
}

func newLokiWriter(lokiURL, serviceName string) *lokiWriter {
	lw := &lokiWriter{
		lokiURL:     lokiURL,
		serviceName: serviceName,
		client:      &http.Client{Timeout: 3 * time.Second},
		logCh:       make(chan [2]string, 1000),
		stopCh:      make(chan struct{}),
	}

	lw.wg.Add(1)
	go lw.worker()
	return lw
}

func (w *lokiWriter) Write(p []byte) (n int, err error) {
	ts := strconv.FormatInt(time.Now().UnixNano(), 10)
	val := string(bytes.TrimSpace(p))

	select {
	case w.logCh <- [2]string{ts, val}:
	default:
		// Drop log entry if buffer full to prevent blocking main service execution
	}
	return len(p), nil
}

func (w *lokiWriter) worker() {
	defer w.wg.Done()
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	var batch [][2]string

	for {
		select {
		case item := <-w.logCh:
			batch = append(batch, item)
			if len(batch) >= 100 {
				w.flush(batch)
				batch = nil
			}
		case <-ticker.C:
			if len(batch) > 0 {
				w.flush(batch)
				batch = nil
			}
		case <-w.stopCh:
			for item := range w.logCh {
				batch = append(batch, item)
			}
			if len(batch) > 0 {
				w.flush(batch)
			}
			return
		}
	}
}

func (w *lokiWriter) flush(values [][2]string) {
	payload := lokiPushPayload{
		Streams: []lokiStream{
			{
				Stream: map[string]string{
					"service_name": w.serviceName,
					"environment":  "development",
				},
				Values: values,
			},
		},
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", w.lokiURL, bytes.NewBuffer(data))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := w.client.Do(req)
	if err == nil && resp != nil {
		_ = resp.Body.Close()
	}
}

func (w *lokiWriter) Close() {
	close(w.stopCh)
	w.wg.Wait()
}

func NewLogger(serviceName string) (*zap.Logger, func(), error) {
	lokiURL := os.Getenv("LOKI_URL")
	if lokiURL == "" {
		lokiURL = "http://localhost:3100/loki/api/v1/push"
	}

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	jsonEncoder := zapcore.NewJSONEncoder(encoderConfig)

	stdoutSyncer := zapcore.AddSync(os.Stdout)
	lokiSyncer := newLokiWriter(lokiURL, serviceName)

	multiSyncer := zapcore.NewMultiWriteSyncer(stdoutSyncer, zapcore.AddSync(lokiSyncer))

	core := zapcore.NewCore(
		jsonEncoder,
		multiSyncer,
		zap.InfoLevel,
	)

	logger := zap.New(core).With(
		zap.String("service_name", serviceName),
		zap.String("environment", "development"),
	)

	cleanup := func() {
		_ = logger.Sync()
		lokiSyncer.Close()
	}

	return logger, cleanup, nil
}
