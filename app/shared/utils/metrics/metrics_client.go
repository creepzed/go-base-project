package metrics

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
)

type metrics struct {
	mailSent   *prometheus.CounterVec
	mailStates *prometheus.CounterVec
}

const APP = "hermes"

var metricsCli *metrics

func init() {
	var mailSent = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: fmt.Sprintf("%s_mail_sent", APP),
		Help: "Total of mail sent",
	}, []string{"channel", "event"})

	var mailStates = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: fmt.Sprintf("%s_mail_state", APP),
		Help: "Total of mail per states",
	}, []string{"state"})
	prometheus.MustRegister(mailSent, mailStates)

	metricsCli = &metrics{
		mailSent:   mailSent,
		mailStates: mailStates,
	}
}

func IncrementMailSent(channel string, event string) { metricsCli.IncrementMailSent(channel, event) }
func (metrics *metrics) IncrementMailSent(channel string, event string) {
	metrics.mailSent.With(prometheus.Labels{"channel": channel, "event": event}).Inc()
}

func IncrementMailState(state string) { metricsCli.IncrementMailState(state) }
func (metrics *metrics) IncrementMailState(state string) {
	metrics.mailStates.With(prometheus.Labels{"state": state}).Inc()
}
