package prometheus

import (
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	v2 "github.com/wx-up/coding/web/v2"
)

type MiddlewareBuilder struct {
	Name        string
	Subsystem   string
	ConstLabels map[string]string
	Help        string
}

func (mb *MiddlewareBuilder) Build() v2.Middleware {
	// Summary 类型：https://blog.csdn.net/hugo_lei/article/details/113743597
	summaryVec := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Name:        mb.Name,
		Subsystem:   mb.Subsystem,
		ConstLabels: mb.ConstLabels,
		Help:        mb.Help,
	}, []string{"pattern", "method", "status"})
	prometheus.MustRegister(summaryVec) // 需要注册
	return func(next v2.HandleFunc) v2.HandleFunc {
		return func(ctx *v2.Context) {
			startTime := time.Now()
			next(ctx)
			endTime := time.Now()
			status := ctx.RespStatusCode
			route := "unknown"
			if ctx.MatchPath != "" {
				route = ctx.MatchPath
			}
			summaryVec.WithLabelValues(route, ctx.Req.Method, strconv.Itoa(status)).Observe(float64(endTime.Sub(startTime).Milliseconds()))
			// go report(endTime.Sub(startTime), ctx, summaryVec)
		}
	}
}

func report(dur time.Duration, ctx *v2.Context, vec prometheus.ObserverVec) {
	status := ctx.RespStatusCode
	route := "unknown"
	if ctx.MatchPath != "" {
		route = ctx.MatchPath
	}
	ms := dur / time.Millisecond
	vec.WithLabelValues(route, ctx.Req.Method, strconv.Itoa(status)).Observe(float64(ms))
}
