package regular_expression

type config struct {
	Regular       *regular `config:"regexp" validate:"required"`
	Field         string   `config:"field"`
	TargetPrefix  string   `config:"target_prefix"`
	IgnoreFailure bool     `config:"ignore_failure"`
}

var defaultConfig = config{
	Field:        "message",
	TargetPrefix: "regular",
}

type regular = Regular

// Unpack a tokenizer into a dissector this will trigger the normal validation of the dissector.
func (t *regular) Unpack(v string) error {
	d, err := New(v)
	if err != nil {
		return err
	}
	*t = *d
	return nil
}
