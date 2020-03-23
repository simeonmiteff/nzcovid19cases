package nzcovid19cases

import (
	"encoding/json"
	"fmt"
)

type AlertLevel struct {
	Level int
	LevelName string
}

func RenderLevels(levelInt int, levelString, viewType string) (string, error) {
	if viewType != "json" {
		return "", InvalidUsageError{fmt.Sprintf("unknown view type: %v", viewType)}
	}
	level := AlertLevel{
		Level:     levelInt,
		LevelName: levelString,
	}
	b, err := json.MarshalIndent(level, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}