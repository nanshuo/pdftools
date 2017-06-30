package job

import (
	"os"
	"time"
)

const (
	ConvertJobStatusPutSuccess = iota
	ConvertJobStatusPutFailed
	ConvertJobStatusWait
	ConvertJobStatusIng
	ConvertJobStatusDoneSuccess
	ConvertJobStatusDoneFailed
	ConvertJobStatusUnknow
)

var (
	statusInfoMap = map[int]string{
		ConvertJobStatusPutSuccess:  "add job success",
		ConvertJobStatusPutFailed:   "add job failed",
		ConvertJobStatusWait:        "waiting in task queue",
		ConvertJobStatusIng:         "converting, please wait",
		ConvertJobStatusDoneSuccess: "convert success",
		ConvertJobStatusDoneFailed:  "convert failed",
	}
)

type ConvertJob struct {
	JobId      string `json:"job_id"`
	Status     int    `json:"status"`
	StatusInfo string `json:"status_info"`

	ConvertType string `json:"convert_type"`

	InputFileInfo  os.FileInfo `json:"input_file_info"`
	OutputFileInfo os.FileInfo `json:"output_file_info"`

	// time information
	AddTime   time.Duration `json:"add_time"`
	StartTime time.Duration `json:"start_time"`
	EndTime   time.Duration `json:"end_time"`
}

func NewConvertJob() *ConvertJob {
	j := &ConvertJob{}

	return j
}
