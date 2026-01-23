package report

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

type StatusPayload struct {
	GameID         string `json:"game_id"`
	State          string `json:"state"`
	PlayerCount    int    `json:"playerCount"`
	MaxPlayerCount int    `json:"maxPlayerCount"`
}

type Reporter struct {
	backendURL string
	gameID     string
	client     *http.Client

	updateChan chan StatusPayload

	stop    chan struct{}
	stopped sync.WaitGroup
}

func NewReporter(backendURL string, gameID string) *Reporter {
	reporter := &Reporter{
		backendURL: backendURL,
		gameID:     gameID,
		client:     &http.Client{Timeout: 3 * time.Second},
		updateChan: make(chan StatusPayload, 16),
		stop:       make(chan struct{}),
	}
	reporter.stopped.Add(1)
	go reporter.run()
	return reporter
}

func (r *Reporter) Stop() {
	close(r.stop)
	r.stopped.Wait()
}

func (r *Reporter) Update(payload StatusPayload) {
	select {
	case r.updateChan <- payload:
	default:
		fmt.Printf("status update dropped: %v", payload)
	}
}

func (r *Reporter) run() {
	defer r.stopped.Done()

	// debounce window: coalesce rapid updates within this duration
	const debounce = 200 * time.Millisecond
	var pending *StatusPayload
	var timer *time.Timer

	send := func(payload StatusPayload) {
		// build URL and POST
		url := r.backendURL + "/games/" + payload.GameID + "/status"
		b, _ := json.Marshal(payload)

		// retry policy: try once, then one retry after short backoff
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(b))
		if err != nil {
			log.Println("reporter: build request error:", err)
			return
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := r.client.Do(req)
		if err != nil {
			log.Println("reporter: send error, retrying once:", err)
			// simple retry
			time.Sleep(250 * time.Millisecond)
			ctx2, cancel2 := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel2()
			req2, err2 := http.NewRequestWithContext(ctx2, "POST", url, bytes.NewReader(b))
			if err2 == nil {
				req2.Header.Set("Content-Type", "application/json")
				resp2, err2 := r.client.Do(req2)
				if err2 != nil {
					log.Println("reporter: retry failed:", err2)
				} else {
					_ = resp2.Body.Close()
				}
			}
			return
		}
		_ = resp.Body.Close()
	}

	for {
		select {
		case <-r.stop:
			// send any pending immediately before exit
			if pending != nil {
				send(*pending)
			}
			return
		case upd := <-r.updateChan:
			// coalesce into pending
			p := upd
			p.GameID = r.gameID // ensure gameId is set
			pending = &p

			// reset timer
			if timer == nil {
				timer = time.NewTimer(debounce)
			} else {
				if !timer.Stop() {
					<-timer.C
				}
				timer.Reset(debounce)
			}
		case <-func() <-chan time.Time {
			if timer == nil {
				// never fires
				ch := make(chan time.Time)
				return ch
			}
			return timer.C
		}():
			if pending != nil {
				send(*pending)
				pending = nil
			}
		}
	}
}
