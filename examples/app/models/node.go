package models

import (
	_activity "github.com/gtantech/pdm/activity"
)

type Node struct {
	_activity.Activity[Activity]
}

func NewNode(activity Activity) *Node {
	return &Node{Activity: _activity.New(activity)}
}
