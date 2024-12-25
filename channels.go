package main

import (
	"encoding/json"
	"os"
)

type Trigger struct {
	Type     string  `json:"type"`
	Time     float64 `json:"time"`
	Unit     string  `json:"unit"`
	Distance float64 `json:"distance"`
	Velocity float64 `json:"velocity"`
	Topic    string  `json:"topic"`
}

type Location struct {
	Name          string  `json:"name"`
	Description   string  `json:"description"`
	Trigger       Trigger `json:"trigger"`
	TopicPrevious string  `json:"topicprevious"`
	TopicEntry    string  `json:"topicentry"`
	TopicExit     string  `json:"topicexit"`
}

type Tag struct {
	TagName      string `json:"tagname"`
	RawTopic     string `json:"rawtopic"`
	VariableType string `json:"variabletype"`
	Location     string `json:"location"`
	Unit         string `json:"unit"`
}

type Channels struct {
	ModelName    string     `json:"modelname"`
	ModelId      string     `json:"modelid"`
	ModelVersion string     `json:"modelversion"`
	Tags         []Tag      `json:"tags"`
	Locations    []Location `json:"locations"`
}

var channels Channels

func ReadChannels(filename string) (Channels, error) {
	contents, err := os.ReadFile(filename)
	if err != nil {
		return Channels{}, err
	}
	var channels Channels
	err = json.Unmarshal(contents, &channels)
	if err != nil {
		return Channels{}, err
	}
	return channels, nil
}

func GetKnownChannelTopics() []string {
	knownChannelTopics := []string{}
	for _, variable := range channels.Tags {
		knownChannelTopics = append(knownChannelTopics, variable.RawTopic)
	}
	return knownChannelTopics
}

func IsAKnownTopic(topic string) bool {
	for _, knownTopic := range GetKnownChannelTopics() {
		if knownTopic == topic {
			return true
		}
	}
	return false
}

func GetTriggerStartTopic() string {
	for _, c := range channels.Tags {
		if c.VariableType == "TriggerStart" {
			return c.RawTopic
		}
	}
	return ""
}
