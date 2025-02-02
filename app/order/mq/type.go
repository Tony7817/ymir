package mq

import (
	"encoding/json"

	"github.com/hibiken/asynq"
)

type OrderTask struct {
	OrderId int64
	UserId  int64
	Success bool
}

func NewOrderTask(orderId int64, userId int64) (*asynq.Task, error) {
	payload, err := json.Marshal(OrderTask{
		OrderId: orderId,
		UserId:  userId,
	})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask("payment", payload), nil
}
