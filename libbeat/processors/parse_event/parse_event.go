package parse_event

import (
	"strings"
	"time"

	"github.com/elastic/beats/v7/libbeat/logp"
	"github.com/pkg/errors"

	"github.com/elastic/beats/v7/libbeat/beat"
	"github.com/elastic/beats/v7/libbeat/common"
	"github.com/elastic/beats/v7/libbeat/processors"
	jsprocessor "github.com/elastic/beats/v7/libbeat/processors/script/javascript/module/processor"
)

func init() {
	processors.RegisterPlugin("parse_event", New)
	jsprocessor.RegisterPlugin("PARSE_EVENT", New)
}

type parse_event struct {
	config Config
	fields [] string
	log    *logp.Logger
}

const (
	processorName = "parse_event"
	logName       = "processor.parse_event"
)

// New constructs a new parse_event processor.
func New(cfg *common.Config) (processors.Processor, error) {
	// 初始化配置文件
	config := defaultConfig()
	//logp.Info("config New is %v", config)
	if err := cfg.Unpack(&config); err != nil {
		return nil, errors.Wrapf(err, "fail to unpack the %v configuration", processorName)
	}

	p := &parse_event{
		config: config,
		// 待分割的每段日志对应的 key
		log:    logp.NewLogger(logName),
	}

	return p, nil
}
func (p *parse_event) checkPathValue(pathValue string) (result bool) {
	if len([]rune(pathValue)) == 0 {
		return false
	}
	if len(strings.Split(pathValue, p.config.separator)) <= p.config.Index {
		return false
	}
	return true
}

// 真正的日志处理逻辑
// 为了保证后面的 processor 正常处理，这里面没有 return 任何 error，只是简单的打印
func (p *parse_event) Run(event *beat.Event) (*beat.Event, error) {
	// 尝试获取 message，理论上这一步不应该出现问题
	//msg, err := event.GetValue("message")
	logp.Info("begin to exec processors parse_event %v", event)
	msg, err := event.GetValue(p.config.fieldPath)
	if err != nil {
		p.log.Error(err)
		return event, nil
	}

	message, ok := msg.(string)
	if !ok {
		p.log.Error("failed to parse message")
		return event, nil
	}

	switch p.config.mode {
	case "manual":
		if p.checkPathValue(message){
			_, _ = event.PutValue(p.config.keyName, strings.Split(message, p.config.separator))
		}else {
			if !p.config.ignoreError {
				return event, nil
			}
		}
	case "auto":
		fieldsValue := strings.Split(message, p.config.separator)
		p.log.Debugf("message fields: %v", fieldsValue)
		if len(fieldsValue) < 8 {
			p.log.Errorf("incorrect field length: %d, expected source_path: %v", len(fieldsValue), fieldsValue)
			return event, nil
		}
		_, _ = event.PutValue("app_code", fieldsValue[5])
		if p.config.enableEnvType {
			_, _ = event.PutValue("pod_name", fieldsValue[7])
			_, _ = event.PutValue("source_host", fieldsValue[7])
			_, _ = event.PutValue("env_type", fieldsValue[6])
		}
		fullLogName := fieldsValue[len(fieldsValue) - 1]
		fullLogNameSplit := strings.Split(fullLogName, ".")
		if fullLogNameSplit[0] == "access"{
			_, _ = event.PutValue("log_name", "access")
		}else {
			logName := fullLogNameSplit[0:len(fullLogNameSplit)-1]
			_, _ = event.PutValue("log_name", strings.Join([]string(logName), ","))
		}

		if p.config.deleteUnuseFields {
			_ = event.Delete("@metadata")
			_ = event.Delete("host");
			_ = event.Delete("beat");
			_ = event.Delete("log");
			_ = event.Delete("offset");
			_ = event.Delete("prospector");
			_ = event.Delete("input");
			_ = event.Delete("ecs");
			_ = event.Delete("agent");
		}
		if p.config.enableTime {
			timestamp, _ := event.GetValue("@timestamp")
			_, _ = event.PutValue("timestamp", timestamp)
			_, _ = event.PutValue("send_time", time.Now().Unix())
		}
		logp.Info("push event parse_event is %v", *event)
		return event, nil
	}
	return event, errors.New("event split err")
}

func (p *parse_event) String() string {
	return processorName
}

