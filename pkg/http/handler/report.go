package handler

import (
	"encoding/json"
	"github.com/ScoreTrak/ScoreTrak/pkg/logger"
	"github.com/ScoreTrak/ScoreTrak/pkg/report"
	"github.com/ScoreTrak/Web/pkg/role"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"net/http"
)

type reportController struct {
	log    logger.LogInfoFormat
	client *ClientStore
}

func NewReportController(log logger.LogInfoFormat, client *ClientStore) *reportController {
	return &reportController{log: log, client: client}
}

func (u *reportController) Get(c *gin.Context) {
	lr, err := u.client.ReportClient.Get()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		u.log.Error(err)
	}
	if lr.Cache == "" {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "The report was not generated. Ensure that rounds are being scored"})
		u.log.Error(err)
		return
	}

	var r string
	if val, ok := c.Get("role"); ok && val != nil {
		r, _ = val.(string)
	}
	var tID uuid.UUID
	if val, ok := c.Get("team_id"); ok && val != nil {
		tID, _ = val.(uuid.UUID)
	}
	simpleReport := &report.SimpleReport{}
	err = json.Unmarshal([]byte(lr.Cache), simpleReport)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		u.log.Error(err)
		return
	}

	if r == role.Blue || r == role.Anonymous {
		p := u.client.PolicyClient.GetPolicy()
		for t := range simpleReport.Teams {
			if !simpleReport.Teams[t].Enabled {
				delete(simpleReport.Teams, t)
				continue
			}
			for h := range simpleReport.Teams[t].Hosts {
				if !simpleReport.Teams[t].Hosts[h].Enabled {
					delete(simpleReport.Teams[t].Hosts, h)
					continue
				}
				for s := range simpleReport.Teams[t].Hosts[h].Services {
					if !simpleReport.Teams[t].Hosts[h].Services[s].Enabled {
						delete(simpleReport.Teams[t].Hosts[h].Services, s)
						continue
					}
					if t != tID {
						simpleReport.Teams[t].Hosts[h].Services[s].Err = ""
						simpleReport.Teams[t].Hosts[h].Services[s].Log = ""
						prop := map[string]string{}

						if port, ok := simpleReport.Teams[t].Hosts[h].Services[s].Properties["Port"]; ok && !*p.ShowAddresses {
							prop["Port"] = port
						}

						simpleReport.Teams[t].Hosts[h].Services[s].Properties = prop
						if !*p.ShowPoints {
							simpleReport.Teams[t].Hosts[h].Services[s].Points = 0
							simpleReport.Teams[t].Hosts[h].Services[s].PointsBoost = 0
						}
					}
				}
				if t != tID {
					if !*p.ShowAddresses {
						simpleReport.Teams[t].Hosts[h].Address = ""
					}
				}
			}
		}
	}
	c.JSON(200, simpleReport)

}

func (u *reportController) GetByTeamID(c *gin.Context) {
	lr, err := u.client.ReportClient.Get()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		u.log.Error(err)
	}
	if lr.Cache == "" {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "The report was not generated. Ensure that rounds are being scored"})
		u.log.Error(err)
		return
	}

	id, err := UuidResolver(c, "id")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		u.log.Error(err)
	}

	var tID uuid.UUID
	var r string

	if val, ok := c.Get("role"); ok && val != nil {
		r, _ = val.(string)
	}
	if val, ok := c.Get("team_id"); ok && val != nil {
		tID, _ = val.(uuid.UUID)
	}

	if r == role.Black || (r == role.Blue && tID == id) {
		simpleReport := &report.SimpleReport{}
		err = json.Unmarshal([]byte(lr.Cache), simpleReport)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			u.log.Error(err)
			return
		}
		simpleTeamReport := &report.SimpleReport{Round: simpleReport.Round, Teams: map[uuid.UUID]*report.SimpleTeam{id: simpleReport.Teams[id]}}
		c.JSON(200, simpleTeamReport)
	}
}

//func (u *reportController) LazyReportLoader(client client.ScoretrakClient) {
//	currentScoretrakTime := &time.Time{}
//	err := client.GenericGet(currentScoretrakTime, "/time")
//	if err != nil {
//		panic("failed to retrieve time from scoretrak")
//	}
//	dsync := time.Since(*currentScoretrakTime)
//	if float64(time.Second*2) < math.Abs(float64(dsync)) {
//		panic(fmt.Sprintf( "time difference between web host, and scoretrak host is too large. Please synchronize time\n(The difference should not exceed 2 seconds)\nTime on web:%s\nTime on master:%s", currentScoretrakTime.String(), time.Now()))
//	}
//
//	var sleep time.Duration
//	for {
//		time.Sleep(sleep)
//		conf, err := u.client.ConfigClient.Get()
//		if err != nil {
//			panic(err)
//		}
//
//		lastRound, err := u.client.RoundClient.GetLastRound()
//		if err != nil {
//			sleep = config.MinRoundDuration
//		} else {
//			nextRoundStart := lastRound.Start.Add(time.Duration(conf.RoundDuration)*time.Second + time.Second*2)
//			if time.Until(nextRoundStart) < config.MinRoundDuration {
//				sleep = time.Until(nextRoundStart)
//			} else {
//				sleep = config.MinRoundDuration
//			}
//		}
//
//		if conf.Enabled == nil || !*conf.Enabled {
//			sleep = config.MinRoundDuration
//			continue
//		}
//		r, err := u.client.ReportClient.Get()
//		if err != nil {
//			panic(err)
//		}
//		u.mu.Lock()
//		u.report = r
//		u.mu.Unlock()
//		nextRoundStart := lastRound.Start.Add(time.Duration(conf.RoundDuration)*time.Second + time.Second*2)
//		if time.Until(nextRoundStart) > 0 && time.Until(nextRoundStart) < config.MinRoundDuration {
//			sleep = time.Until(nextRoundStart)
//		} else {
//			sleep = config.MinRoundDuration
//		}
//
//	}
//}
