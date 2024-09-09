package models

import (
	"errors"
	"io"

	"github.com/suyashkumar/dicom"
	"github.com/suyashkumar/dicom/pkg/tag"
)

type DicomDTO struct {
	InstanceId string         `json:"instanceId"`
	Dataset    *dicom.Dataset `json:"-"`
	ImagePath  string         `json:"-"`
}

func NewDicomDTO(in io.Reader, bytesToRead int64) (*DicomDTO, error) {
	dataset, err := dicom.Parse(in, bytesToRead, nil)
	if err != nil {
		return nil, err
	}

	element, err := dataset.FindElementByTag(tag.SOPInstanceUID)
	if err != nil {
		return nil, errors.New("SOP Instance UID not found")
	}

	instanceId := element.Value.GetValue().([]string)[0]

	return &DicomDTO{
		InstanceId: instanceId,
		Dataset:    &dataset,
	}, nil
}
