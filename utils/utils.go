package utils

import (
	"net/http"
	"time"

	gocache "github.com/patrickmn/go-cache"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

// Utils can be used as a value in the pipeline to provide several usefull functions
type Utils struct {
	ScriptName       string
	cache            *gocache.Cache
	httpClient       *http.Client
	logrusEntry      *logrus.Entry
	promErrorCounter *prometheus.CounterVec
}

// InitUtils returns a new instance of *Utils
// httpClient, logrusEntry, and cache params will used default values if nil
// a
func InitUtils(httpClient *http.Client, logrusEntry *logrus.Entry, cache *gocache.Cache) *Utils {
	u := &Utils{
		httpClient:  httpClient,
		logrusEntry: logrusEntry,
		cache:       cache,
	}

	if u.httpClient == nil {
		u.httpClient = http.DefaultClient
	}

	if u.cache == nil {
		u.cache = gocache.New(gocache.DefaultExpiration, time.Minute*5)
	}

	return u
}
