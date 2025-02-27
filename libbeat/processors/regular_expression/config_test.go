package regular_expression

import (
	"testing"

	"github.com/stretchr/testify/assert"

	conf "github.com/elastic/elastic-agent-libs/config"
)

func TestConfig(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		c, err := conf.NewConfigFrom(map[string]interface{}{
			"tokenizer": "%{value1}",
			"field":     "message",
		})
		if !assert.NoError(t, err) {
			return
		}

		cfg := config{}
		err = c.Unpack(&cfg)
		if !assert.NoError(t, err) {
			return
		}
	})

	t.Run("invalid", func(t *testing.T) {
		c, err := conf.NewConfigFrom(map[string]interface{}{
			"tokenizer": "%value1}",
			"field":     "message",
		})
		if !assert.NoError(t, err) {
			return
		}

		cfg := config{}
		err = c.Unpack(&cfg)
		if !assert.Error(t, err) {
			return
		}
	})

	t.Run("with tokenizer missing", func(t *testing.T) {
		c, err := conf.NewConfigFrom(map[string]interface{}{})
		if !assert.NoError(t, err) {
			return
		}

		cfg := config{}
		err = c.Unpack(&cfg)
		if !assert.Error(t, err) {
			return
		}
	})

	t.Run("with empty tokenizer", func(t *testing.T) {
		c, err := conf.NewConfigFrom(map[string]interface{}{
			"tokenizer": "",
		})
		if !assert.NoError(t, err) {
			return
		}

		cfg := config{}
		err = c.Unpack(&cfg)
		if !assert.Error(t, err) {
			return
		}
	})

	t.Run("tokenizer with no field defined", func(t *testing.T) {
		c, err := conf.NewConfigFrom(map[string]interface{}{
			"tokenizer": "hello world",
		})
		if !assert.NoError(t, err) {
			return
		}

		cfg := config{}
		err = c.Unpack(&cfg)
		if !assert.Error(t, err) {
			return
		}
	})

	t.Run("with wrong trim_mode", func(t *testing.T) {
		c, err := conf.NewConfigFrom(map[string]interface{}{
			"tokenizer":   "hello %{what}",
			"field":       "message",
			"trim_values": "bananas",
		})
		if !assert.NoError(t, err) {
			return
		}

		cfg := config{}
		err = c.Unpack(&cfg)
		if !assert.Error(t, err) {
			return
		}
	})

	t.Run("with valid trim_mode", func(t *testing.T) {
		c, err := conf.NewConfigFrom(map[string]interface{}{
			"tokenizer":   "hello %{what}",
			"field":       "message",
			"trim_values": "all",
		})
		if !assert.NoError(t, err) {
			return
		}

		cfg := config{}
		err = c.Unpack(&cfg)
		if !assert.NoError(t, err) {
			return
		}
	})
}

func TestConfigForDataType(t *testing.T) {
	t.Run("valid data type", func(t *testing.T) {
		c, err := conf.NewConfigFrom(map[string]interface{}{
			"tokenizer": "%{value1|integer} %{value2|float} %{value3|boolean} %{value4|long} %{value5|double}",
			"field":     "message",
		})
		if !assert.NoError(t, err) {
			return
		}

		cfg := config{}
		err = c.Unpack(&cfg)
		if !assert.NoError(t, err) {
			return
		}
	})
	t.Run("invalid data type", func(t *testing.T) {
		c, err := conf.NewConfigFrom(map[string]interface{}{
			"tokenizer": "%{value1|int} %{value2|short} %{value3|char} %{value4|void} %{value5|unsigned} id=%{id|xyz} status=%{status|abc} msg=\"%{message}\"",
			"field":     "message",
		})
		if !assert.NoError(t, err) {
			return
		}

		cfg := config{}
		err = c.Unpack(&cfg)
		if !assert.Error(t, err) {
			return
		}
	})
	t.Run("missing data type", func(t *testing.T) {
		c, err := conf.NewConfigFrom(map[string]interface{}{
			"tokenizer": "%{value1|} %{value2|}",
			"field":     "message",
		})
		if !assert.NoError(t, err) {
			return
		}

		cfg := config{}
		err = c.Unpack(&cfg)
		if !assert.Error(t, err) {
			return
		}
	})
}
