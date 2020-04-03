/*
Copyright 2020 Rafael Fernández López <ereslibre@ereslibre.es>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package conditions

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	commonv1alpha1 "github.com/oneinfra/oneinfra/apis/common/v1alpha1"
)

// ConditionType represents a condition type
type ConditionType string

// ConditionStatus represents a condition status
type ConditionStatus string

const (
	// ConditionTrue is a condition that holds true
	ConditionTrue ConditionStatus = "True"
	// ConditionFalse is a condition that holds false
	ConditionFalse ConditionStatus = "False"
	// ConditionUnknown is a condition in an unknown state
	ConditionUnknown ConditionStatus = "Unknown"
)

// Condition represents a condition
type Condition struct {
	Type               ConditionType
	Status             ConditionStatus
	LastTransitionTime metav1.Time
	Reason             string
	Message            string
}

// ConditionList represents a list of conditions
type ConditionList []Condition

// NewConditionListFromv1alpha1 returns an internal condition list
// from a versioned condition list
func NewConditionListFromv1alpha1(conditionList commonv1alpha1.ConditionList) ConditionList {
	res := ConditionList{}
	for _, condition := range conditionList {
		res = append(res, Condition{
			Type:               ConditionType(condition.Type),
			Status:             ConditionStatus(condition.Status),
			LastTransitionTime: condition.LastTransitionTime,
			Reason:             condition.Reason,
			Message:            condition.Message,
		})
	}
	return res
}

// IsCondition checks whether a condition type matches a condition
// status
func (conditionList ConditionList) IsCondition(conditionType ConditionType, conditionStatus ConditionStatus) bool {
	for _, condition := range conditionList {
		if condition.Type == conditionType {
			return condition.Status == conditionStatus
		}
	}
	return false
}

// SetCondition sets the condition type to the condition status
func (conditionList *ConditionList) SetCondition(conditionType ConditionType, conditionStatus ConditionStatus) {
	newConditionList := ConditionList{}
	for _, condition := range *conditionList {
		if condition.Type == conditionType {
			if condition.Status == conditionStatus {
				return
			}
		} else {
			newConditionList = append(newConditionList, condition)
		}
	}
	newConditionList = append(newConditionList, Condition{
		Type:               conditionType,
		Status:             conditionStatus,
		LastTransitionTime: metav1.Now(),
	})
	*conditionList = newConditionList
}

// Export exports the internal condition list to a versioned condition
// list
func (conditionList ConditionList) Export() commonv1alpha1.ConditionList {
	res := commonv1alpha1.ConditionList{}
	for _, condition := range conditionList {
		res = append(res, commonv1alpha1.Condition{
			Type:               commonv1alpha1.ConditionType(condition.Type),
			Status:             commonv1alpha1.ConditionStatus(condition.Status),
			LastTransitionTime: condition.LastTransitionTime,
			Reason:             condition.Reason,
			Message:            condition.Message,
		})
	}
	return res
}