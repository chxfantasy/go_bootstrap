package middleware

import (
	"bytes"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/chxfantasy/go_bootstrap/conf"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

/*
{
    trace: {
        from:  // trace来源，如某个业务的某个模块
        step:  // trace深度，请求到第几层了
        id:    // trace id，每次请求都有的一个唯一的id，比如可以采用每次生成不重复的UUID
    },
    addr: {
        rmt:  // 对端的ip和端口，比如10.2.3.4:8000
        loc:  // 本地的ip和端口，比如10.2.3.5:8020
    },
    time: // 毫秒 201809221830100 「global」
    params:{} // 请求参数
    elapse: // 耗时，毫秒
    result: { // 返回结果
        ret: 1,
        data: {}
    },
    level: //INFO、ERROR、WARNING、DEBUG 「global」
    path: “filename.go:300:/favor/add", // 打点位置  「global」
    ext: {} // 自定义字段。不超过5个
}
*/

var localIP string

type addr struct {
	RemoteIP string
	LocalIP  string
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}
func (w bodyLogWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

type traceInfo struct {
	ServiceName string `json:"serviceName"`
	StartTime   int64  `json:"startTime"`
	TraceID     string `json:"traceID"`
}

func init() {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Errorf("Oops, InterfaceAddrs err:\n%v", err)
	}
	for _, a := range addrs {
		if ipNet, ok := a.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				localIP = ipNet.IP.String()
				return
			}
		}
	}
}

//Boss 打印日志中间件
func Boss(serviceName string, logger *conf.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		URLPath := c.Request.URL.Path
		if URLPath == "/health" {
			c.Next()
			return
		}

		start := time.Now()
		traceID := c.GetHeader("traceID")
		if traceID == "" {
			traceID = uuid.NewV4().String()
		}
		c.Set("traceID", traceID)
		raw := c.Request.URL.RawQuery
		traceInfo := &traceInfo{}
		traceInfo.ServiceName = serviceName
		traceInfo.TraceID = traceID
		traceInfo.StartTime = start.UnixNano() / int64(time.Microsecond)

		if raw != "" {
			URLPath = URLPath + "?" + raw
		}
		clientIP := c.ClientIP()
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw
		c.Next()

		//过滤 404 请求，不写trace日志
		if c.Writer.Status() == http.StatusNotFound {
			return
		}

		end := time.Now()
		latency := end.Sub(start)
		_addr := &addr{
			RemoteIP: clientIP,
			LocalIP:  localIP,
		}
		elapse := latency.Milliseconds()
		m := map[string]interface{}{
			"addr":      _addr,
			"elapse":    elapse,
			"url":       URLPath,
			"status":    c.Writer.Status(),
			"response":  blw.body.String(),
			"traceInfo": traceInfo,
			"traceID":   traceID,
		}
		msg := "jaeger-trace"
		logger.Infow(msg, m)
		if elapse > 100 {
			logger.Errorw(msg, m)
		}
	}
}
