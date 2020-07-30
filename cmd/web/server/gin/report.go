package gin

import (
	"encoding/json"
	"fmt"
	"github.com/L1ghtman2k/ScoreTrak/pkg/api/client"
	"github.com/L1ghtman2k/ScoreTrak/pkg/config"
	"github.com/L1ghtman2k/ScoreTrak/pkg/logger"
	"github.com/L1ghtman2k/ScoreTrak/pkg/report"
	"github.com/L1ghtman2k/ScoreTrak/pkg/round"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"math"
	"net/http"
	"sync"
	"time"
)

type reportController struct {
	log          logger.LogInfoFormat
	reportClient report.Serv
	configClient config.Serv
	roundClient  round.Serv
	mu           sync.RWMutex
	report       *report.Report
}

func NewReportController(log logger.LogInfoFormat, rs report.Serv, cs config.Serv, rrs round.Serv) *reportController {
	return &reportController{log: log, reportClient: rs, configClient: cs, roundClient: rrs}
}

func (u *reportController) Get(c *gin.Context) {
	u.mu.RLock()
	lr := report.Report{}
	err := copier.Copy(&lr, u.report)
	u.mu.RUnlock()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	var tID uint32
	if val, ok := c.Get("team_id"); ok && val != nil {
		tID, _ = val.(uint32)
	}
	simpleReport := &report.SimpleReport{}
	err = json.Unmarshal([]byte(lr.Cache), simpleReport)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if tID != 0 {
		for t := range simpleReport.Teams {
			if t != tID {
				for h := range simpleReport.Teams[t].Hosts {
					for s := range simpleReport.Teams[t].Hosts[h].Services {
						simpleReport.Teams[t].Hosts[h].Services[s].Err = ""
						simpleReport.Teams[t].Hosts[h].Services[s].Log = ""
						simpleReport.Teams[t].Hosts[h].Services[s].Properties = map[string]string{}
					}
				}
			}
		}
	}
	c.JSON(200, simpleReport)
}

func (u *reportController) LazyUpdate(client client.ScoretrakClient) {

	currentScoretrakTime := &time.Time{}
	err := client.GenericGet(currentScoretrakTime, "/time")
	if err != nil {
		panic("failed to retrieve time from scoretrak")
	}
	dsync := time.Since(*currentScoretrakTime)
	if float64(time.Second*2) < math.Abs(float64(dsync)) {
		panic(fmt.Sprintf("time difference between web host, and scoretrak host is too large. Please synchronize time\n(The difference should not exceed 2 seconds)\nTime on web:%s\nTime on master:%s", currentScoretrakTime.String(), time.Now()))
	}

	var sleep time.Duration
	for {
		time.Sleep(sleep)
		conf, err := u.configClient.Get()
		if err != nil {
			panic(err)
		}

		lastRound, err := u.roundClient.GetLastRound()
		if err != nil {
			sleep = config.MinRoundDuration
		} else {
			nextRoundStart := lastRound.Start.Add(time.Duration(conf.RoundDuration)*time.Second + time.Second*2)
			if time.Until(nextRoundStart) < config.MinRoundDuration {
				sleep = time.Until(nextRoundStart)
			} else {
				sleep = config.MinRoundDuration
			}
		}

		if conf.Enabled == nil || !*conf.Enabled {
			sleep = config.MinRoundDuration
			continue
		}
		r, err := u.reportClient.Get()
		if err != nil {
			panic(err)
		}
		u.mu.Lock()
		u.report = r
		u.mu.Unlock()
		nextRoundStart := lastRound.Start.Add(time.Duration(conf.RoundDuration)*time.Second + time.Second*2)
		if time.Until(nextRoundStart) > 0 && time.Until(nextRoundStart) < config.MinRoundDuration {
			sleep = time.Until(nextRoundStart)
		} else {
			sleep = config.MinRoundDuration
		}

	}
}

//TODO: Implement time Desync function between ScoretrakWeb, and Scoretrak
