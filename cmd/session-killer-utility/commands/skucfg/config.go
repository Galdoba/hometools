package skucfg

import "encoding/json"

type Config struct {
	DailyAllowance int `json:"daily allowance"`
}

func (cfg *Config) ToBytes() ([]byte, error) {
	bt, err := json.Marshal(cfg)
	return bt, err
}

func (cfg *Config) FromBytes(bt []byte) error {
	err := json.Unmarshal(bt, cfg)
	return err
}
