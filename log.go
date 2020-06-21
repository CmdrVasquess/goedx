package edgx

import (
	"git.fractalqb.de/fractalqb/c4hgol"
	"git.fractalqb.de/fractalqb/qbsllm"
	"github.com/CmdrVasquess/watched"
)

var (
	log    = qbsllm.New(qbsllm.Lnormal, "edgx", nil, nil)
	LogCfg = c4hgol.Config(qbsllm.NewConfig(log), watched.LogCfg)
)
