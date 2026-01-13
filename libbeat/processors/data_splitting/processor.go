package data_splitting

import (
	"fmt"
	"github.com/elastic/beats/v7/libbeat/beat"
	"github.com/elastic/beats/v7/libbeat/processors"
	jsprocessor "github.com/elastic/beats/v7/libbeat/processors/script/javascript/module/processor"
	cfg "github.com/elastic/elastic-agent-libs/config"
	"github.com/elastic/elastic-agent-libs/mapstr"
)

const flagParsingError = "data_splitting_error"

type processor struct {
	config config
}

func init() {
	processors.RegisterPlugin("splitting", NewProcessor)
	jsprocessor.RegisterPlugin("Splitting", NewProcessor)
}

func NewProcessor(c *cfg.C) (beat.Processor, error) {
	config := defaultConfig
	err := c.Unpack(&config)
	if err != nil {
		return nil, err
	}
	p := &processor{config: config}
	return p, nil
}

func (p *processor) Run(event *beat.Event) (*beat.Event, error) {
	var (
		m   []TargetMessage
		v   interface{}
		err error
	)

	//默认是message 获取一行数据
	v, err = event.GetValue(p.config.Field)
	if err != nil {
		return event, err
	}

	s, ok := v.(string)
	if !ok {
		return event, fmt.Errorf("field is not a string, value: `%v`, field: `%s`", v, p.config.Field)
	}

	// 检查总体大小是否不超过限制，如果是则直接返回
	totalSize := len([]byte(s))
	if int64(totalSize) <= p.config.MaxByteSize {
		return event, nil
	}

	//超过限制，做分割
	//处理offset
	offsets, err := event.GetValue("log.offset")
	if err != nil {
		return event, err
	}
	offset, ok := offsets.(int64)
	if !ok {
		return event, fmt.Errorf("field is not a int64, value: `%v`, field: `%s`", offsets, "offset")
	}

	//分割字符
	m, err = split(s, p.config.MaxByteSize, offset)

	if err != nil {
		if err := mapstr.AddTagsWithKey(
			event.Fields,
			beat.FlagField,
			[]string{flagParsingError},
		); err != nil {
			return event, fmt.Errorf("cannot add new flag the event: %w", err)
		}
		if p.config.IgnoreFailure {
			return event, nil
		}
		return event, err
	}

	backup := event.Clone()
	event, err = p.mapper(event, m)
	if err != nil {
		return backup, err
	}

	//如果不保留原始字段，删除
	if !p.config.RetainOriginalField {
		err := event.Delete(p.config.Field)
		if err != nil {
			return backup, err
		}
	}
	event.PutValue("data_splitting", true)
	event.PutValue("maxByteSize", p.config.MaxByteSize)
	return event, nil
}

func (p *processor) mapper(event *beat.Event, m []TargetMessage) (*beat.Event, error) {
	prefix := ""
	if p.config.Target != "" {
		prefix = p.config.Target
	}
	_, _ = event.PutValue(prefix, m)
	return event, nil
}

func (p *processor) String() string {
	return "field=" + p.config.Field +
		",target=" + p.config.Target +
		",max_byte_size=" + fmt.Sprint(p.config.MaxByteSize)
}
