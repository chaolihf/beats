package regular_expression

import (
	"errors"
	"fmt"
	"github.com/elastic/beats/v7/libbeat/beat"
	"github.com/elastic/beats/v7/libbeat/processors"
	jsprocessor "github.com/elastic/beats/v7/libbeat/processors/script/javascript/module/processor"
	cfg "github.com/elastic/elastic-agent-libs/config"
	"github.com/elastic/elastic-agent-libs/mapstr"
)

const flagParsingError = "regular_expression_parsing_error"

type processor struct {
	config config
}

func init() {
	processors.RegisterPlugin("regular", NewProcessor)
	jsprocessor.RegisterPlugin("Regular", NewProcessor)
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
		m   []string
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

	//正则表达式解析字符串 返回 string数组
	m, err = p.config.Regular.Regular(s)

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

	return event, nil
}

func (p *processor) mapper(event *beat.Event, m []string) (*beat.Event, error) {
	prefix := ""
	if p.config.TargetPrefix != "" {
		prefix = p.config.TargetPrefix + "."
	}

	for k, v := range m {
		if k == 0 {
			continue
		}
		key := fmt.Sprintf("${%d}", k)
		prefixKey := prefix + key
		if _, err := event.GetValue(prefixKey); errors.Is(err, mapstr.ErrKeyNotFound) {
			_, _ = event.PutValue(prefixKey, v)
		} else {
			// When the target key exists but is a string instead of a map.
			if err != nil {
				return event, fmt.Errorf("cannot override existing key with `%s`: %w", prefixKey, err)
			}
			return event, fmt.Errorf("cannot override existing key with `%s`", prefixKey)
		}
	}

	return event, nil
}

func (p *processor) String() string {
	return "regexp=" + p.config.Regular.Raw() +
		",field=" + p.config.Field +
		",target_prefix=" + p.config.TargetPrefix
}
