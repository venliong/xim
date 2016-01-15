/*
 * 消息路由
 */

package logic

import (
	"github.com/liuhengloveyou/nodenet"
	"github.com/liuhengloveyou/xim/common"
)

func init() {
	nodenet.RegisterWorker("ForwardMessage", common.MessageForward{}, ForwardMessage)
}

func ForwardMessage(data interface{}) (result interface{}, err error) {

	return nil, nil
}
