package data_splitting

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/elastic/beats/v7/libbeat/beat"
	conf "github.com/elastic/elastic-agent-libs/config"
	"github.com/elastic/elastic-agent-libs/mapstr"
)

func TestProcessor(t *testing.T) {
	tests := []struct {
		name   string
		c      map[string]interface{}
		fields mapstr.M
		values map[string]string
	}{
		//{
		//	name:   "default field/default target",
		//	c:      map[string]interface{}{"max_byte_size": 113, "retain_original_field": true},
		//	fields: mapstr.M{"log": mapstr.M{"offset": int64(100000)}, "message": "没有共产党就没有新中国没有共产党就没有新中国共产党，辛劳为民族共产党他一心救中国他指给了人民解放的道路他领导中国走向光明他坚持了抗战八年多他改善了人民的生活他建设了敌后根据地他实行了民主好处多没有共产党就没有新中国没有共产党就没有新中国（间奏）没有共产党就没有新中国没有共产党就没有新中国共产党，辛劳为民族共产党他一心救中国他指给了人民解放的道路他领导中国走向光明他坚持了抗战八年多他改善了人民的生活他建设了敌后根据地他实行了民主好处多没有共产党就没有新中国没有共产党就没有新中国"},
		//	values: map[string]string{"messageChunks": ""},
		//},
		{
			name:   "default field/default target",
			c:      map[string]interface{}{"max_byte_size": 1, "retain_original_field": true},
			fields: mapstr.M{"log": mapstr.M{"offset": int64(100000)}, "message": "新中国"},
			values: map[string]string{"messageChunks": ""},
		},
	}

	//charSize := len("你好")
	//fmt.Println("字节大小:", charSize)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c, err := conf.NewConfigFrom(test.c)
			if !assert.NoError(t, err) {
				return
			}

			processor, err := NewProcessor(c)
			if !assert.NoError(t, err) {
				return
			}

			e := beat.Event{Fields: test.fields}
			newEvent, err := processor.Run(&e)
			if !assert.NoError(t, err) {
				return
			}

			fmt.Println("输出结果:", newEvent)

		})
	}
}
