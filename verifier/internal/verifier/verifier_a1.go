package verifier

import (
	"github.com/game-of-nfts/gon-toolbox/verifier/internal/chain"
)

type A1Params struct {
	ChainID string
	TxHash  string
	ClassID string
}

type A1Verifier struct {
	r *chain.Registry
}

func (v A1Verifier) Do(req Request, res chan<- *Respone) {
	result := &Respone{
		TaskNo:   req.TaskNo,
		TeamName: req.User.TeamName,
	}

	params, ok := req.Params.(A1Params)
	if ok {
		result.Reason = "非法参数"
		res <- result
		return
	}

	if len(params.ChainID) == 0 {
		result.Reason = "chainID不能为空"
		res <- result
		return
	}

	if len(params.TxHash) == 0 {
		result.Reason = "txHash不能为空"
		res <- result
		return
	}

	chain := v.r.GetChain(params.ChainID)
	tx, err := chain.GetTx(params.TxHash)
	if err != nil {
		result.Reason = err.Error()
		res <- result
		return
	}

	if req.User.Address[params.ChainID] != tx.Sender {
		result.Reason = "交易地址非用户注册地址"
		res <- result
		return
	}

	class, err := chain.GetClass(params.ClassID)
	if err != nil {
		result.Reason = err.Error()
		res <- result
		return
	}

	if len(class.Uri) == 0 {
		result.Reason = "不符合规则,Uri不能为空"
		res <- result
		return
	}
	if len(class.Data) == 0 {
		result.Reason = "不符合规则,Data不能为空"
		res <- result
		return
	}
	//TODO

	result.Point = PointMap[req.TaskNo]
	res <- result
}
