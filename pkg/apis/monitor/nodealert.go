// Copyright 2019 Yunion
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package monitor

import (
	"time"

	"yunion.io/x/onecloud/pkg/apis"
)

const (
	NodeAlertTypeGuest = "guest"
	NodeAlertTypeHost  = "host"
)

type ResourceAlertV1CreateInput struct {
	*AlertCreateInput

	// 查询指标周期
	Period string `json:"period"`
	// 每隔多久查询一次
	Window string `json:"window"`
	// 比较运算符, 比如: >, <, >=, <=
	Comparator string `json:"comparator"`
	// 报警阀值
	Threshold float64 `json:"threshold"`
	// 报警级别
	Level string `json:"level"`
	// 通知方式, 比如: email, mobile
	Channel string `json:"channel"`
	// 通知接受者
	Recipients string `json:"recipients"`
}

type NodeAlertCreateInput struct {
	ResourceAlertV1CreateInput

	// 监控指标名称
	Metric string `json:"metric"`
	// 监控资源类型, 比如: guest, host
	Type string `json:"type"`
	// 监控资源名称
	NodeName string `json:"node_name"`
	// 监控资源 Id
	NodeId string `json:"node_id"`
}

func (input NodeAlertCreateInput) ToAlertCreateInput(
	name string,
	field string,
	measurement string,
	db string,
	notifications []string) AlertCreateInput {
	freq, _ := time.ParseDuration(input.Window)
	ret := AlertCreateInput{
		Name:      name,
		Frequency: int64(freq / time.Second),
		Settings: AlertSetting{
			Level:         input.Level,
			Notifications: notifications,
			Conditions: []AlertCondition{
				{
					Type:     "query",
					Operator: "and",
					Query: AlertQuery{
						Model: input.GetQuery(field, measurement, db),
						From:  input.Period,
						To:    "now",
					},
					Evaluator: input.GetEvaluator(),
					Reducer: Condition{
						Type: "avg",
					},
				},
			},
		},
	}
	return ret
}

func (input NodeAlertCreateInput) GetQuery(field, measurement, db string) MetricQuery {
	return GetNodeAlertQuery(input.Type, field, measurement, db, input.NodeId)
}

func GetNodeAlertQuery(typ, field, measurement, db, nodeId string) MetricQuery {
	var idField string
	switch typ {
	case NodeAlertTypeGuest:
		idField = "vm_id"
	case NodeAlertTypeHost:
		idField = "host_id"
	}
	sels := make([]MetricQuerySelect, 0)
	sels = append(sels, NewMetricQuerySelect(MetricQueryPart{Type: "field", Params: []string{field}}))
	return MetricQuery{
		Selects: sels,
		Tags: []MetricQueryTag{
			{
				Key:   idField,
				Value: nodeId,
			},
		},
		GroupBy: []MetricQueryPart{
			{
				Type:   "field",
				Params: []string{"*"},
			},
		},
		Measurement: measurement,
		Database:    db,
	}
}

func (input NodeAlertCreateInput) GetEvaluator() Condition {
	return GetNodeAlertEvaluator(input.Comparator, input.Threshold)
}

func GetNodeAlertEvaluator(comparator string, threshold float64) Condition {
	typ := "gt"
	switch comparator {
	case ">=", ">":
		typ = "gt"
	case "<=", "<":
		typ = "lt"
	}
	return Condition{
		Type:   typ,
		Params: []float64{threshold},
	}
}

type NodeAlertListInput struct {
	apis.VirtualResourceListInput

	// 监控指标名称
	Metric string `json:"metric"`
	// 监控资源类型, 比如: guest, host
	Type string `json:"type"`
	// 监控资源名称
	NodeName string `json:"node_name"`
	// 监控资源 Id
	NodeId string `json:"node_id"`
}

func (input NodeAlertListInput) ToAlertListInput() AlertListInput {
	return AlertListInput{
		VirtualResourceListInput: input.VirtualResourceListInput,
		Metric:                   input.Metric,
	}
}

type AlertV1Details struct {
	apis.VirtualResourceDetails

	Name        string  `json:"name"`
	Period      string  `json:"period"`
	Window      string  `json:"window"`
	Comparator  string  `json:"comparator"`
	Threshold   float64 `json:"threshold"`
	Recipients  string  `json:"recipients"`
	Level       string  `json:"level"`
	Channel     string  `json:"channel"`
	DB          string  `json:"db"`
	Measurement string  `json:"measurement"`
	Field       string  `json:"field"`
	NotifierId  string  `json:"notifier_id"`
}

type NodeAlertDetails struct {
	AlertV1Details

	Type     string `json:"type"`
	Metric   string `json:"metric"`
	NodeId   string `json:"node_id"`
	NodeName string `json:"node_name"`
}

type NodeAlertUpdateInput struct {
	// 监控指标名称
	Metric *string `json:"metric"`
	// 监控资源类型, 比如: guest, host
	Type *string `json:"type"`
	// 监控资源名称
	NodeName *string `json:"node_name"`
	// 监控资源 Id
	NodeId *string `json:"node_id"`
	// 查询指标周期
	Period *string `json:"period"`
	// 每隔多久查询一次
	Window *string `json:"window"`
	// 比较运算符, 比如: >, <, >=, <=
	Comparator *string `json:"comparator"`
	// 报警阀值
	Threshold *float64 `json:"threshold"`
	// 报警级别
	Level *string `json:"level"`
	// 通知方式, 比如: email, mobile
	Channel *string `json:"channel"`
	// 通知接受者
	Recipients *string `json:"recipients"`
}
