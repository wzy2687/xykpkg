package xutils

import (
	"encoding/json"
	"log/slog"
	"os"
	"strings"
)

func MustLoadConfPara1Cfg[T any]() *T {
	cfgPath := "./conf/" + os.Args[1]
	if strings.HasPrefix(os.Args[1], "./conf") {
		cfgPath = os.Args[1]
	}
	return MustLoadCfg[T](cfgPath)
}

func MustLoadCfg[T any](cfgPath string) *T {
	rt := new(T)
	b, err := os.ReadFile(cfgPath)
	if err != nil {
		slog.Error("read cfg err", "err", err.Error())
		panic(err.Error())
	}
	err = json.Unmarshal(b, rt)
	if err != nil {
		slog.Error("json decode err", "err", err.Error())
		panic(err.Error())
	}
	return rt
}

func LoadCfg[T any](cfgPath string) (*T, error) {
	rt := new(T)
	b, err := os.ReadFile(cfgPath)
	if err != nil {
		slog.Error(err.Error(), "path", cfgPath)
		return nil, err
	}

	err = json.Unmarshal(b, rt)
	if err != nil {
		slog.Error(err.Error(), "path", cfgPath)
		return nil, err
	}
	return rt, nil
}
