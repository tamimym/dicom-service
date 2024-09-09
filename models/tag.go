package models

import (
	"errors"
	"log/slog"
	"strconv"
	"strings"

	"github.com/suyashkumar/dicom/pkg/tag"
)

func ParseTag(tagValue string) (tag.Tag, error) {
	parts := strings.Split(strings.Trim(tagValue, "()"), ",")
	if len(parts) != 2 {
		slog.Error("Tag value is invalid", slog.String("tag", tagValue))
		return tag.Tag{}, errors.New("invalid tag")
	}

	group, err := strconv.ParseInt(parts[0], 16, 0)
	if err != nil {
		return tag.Tag{}, err
	}
	elem, err := strconv.ParseInt(parts[1], 16, 0)
	if err != nil {
		return tag.Tag{}, err
	}
	return tag.Tag{Group: uint16(group), Element: uint16(elem)}, nil
}
