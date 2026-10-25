// Package pipelineobjects reads and selects the Source objects a Pipeline replicates.
package pipelineobjects

import (
	"fmt"
	"net/url"

	"github.com/hevoio/hevo-go/common"
)

// PipelineObject is one Source object of a Pipeline, such as a table or an API resource.
type PipelineObject struct {
	FullyQualifiedName     string                    `json:"fully_qualified_name"`
	Name                   string                    `json:"name"`
	Schema                 string                    `json:"schema"`
	Database               string                    `json:"database"`
	Destination            PipelineObjectDestination `json:"destination"`
	Status                 string                    `json:"status"`
	Selected               bool                      `json:"selected"`
	LoadMode               string                    `json:"load_mode"`
	Replication            PipelineObjectReplication `json:"replication"`
	Position               *string                   `json:"position"`
	CreatedAt              int64                     `json:"created_at"`
	ParentObject           string                    `json:"parent_object"`
	AllowLoadModeSelection bool                      `json:"allow_load_mode_selection"`
	Fields                 []PipelineObjectField     `json:"fields"`
}

// PipelineObjectDestination is the table the object loads into.
type PipelineObjectDestination struct {
	Table    string `json:"table"`
	Schema   string `json:"schema"`
	Database string `json:"database"`
}

// PipelineObjectReplication is the object's replication status.
type PipelineObjectReplication struct {
	Status string  `json:"status"`
	Reason *string `json:"reason"`
}

// PipelineObjectField is one field of the object and its Destination column.
type PipelineObjectField struct {
	SourceName      string `json:"source_name"`
	SourceType      string `json:"source_type"`
	DestinationName string `json:"destination_name"`
	DestinationType string `json:"destination_type"`
	Status          string `json:"status"`
	PrimaryKey      bool   `json:"primary_key"`
	ForeignKey      bool   `json:"foreign_key"`
	Selected        bool   `json:"selected"`
	Selectable      bool   `json:"selectable"`
	Hashed          bool   `json:"hashed"`
}

// PipelineObjectResponse is one object.
type PipelineObjectResponse struct {
	common.CommonResponse
	Data PipelineObject `json:"data"`
}

// PipelineObjectsListResponse is one page of objects.
type PipelineObjectsListResponse struct {
	common.CommonResponse
	common.ListResponse
	Data []PipelineObject `json:"data"`
}

// Selection modes: what happens to the objects an update does not name.
const (
	SelectionModeIncludeAll     = "INCLUDE_ALL"
	SelectionModeExcludeAll     = "EXCLUDE_ALL"
	SelectionModeRetainExisting = "RETAIN_EXISTING"
)

type pipelineObjectsUpdateRequest struct {
	SelectionMode *string                   `json:"selection_mode,omitempty"`
	Objects       []*objectSelectionRequest `json:"objects,omitempty"`
}

type pipelineObjectsResyncRequest struct {
	ObjectNames []string `json:"object_names"`
	DropAndLoad *bool    `json:"drop_and_load,omitempty"`
}

func objectsPath(pipelineID int64, suffix string) string {
	return fmt.Sprintf("/api/v2/pipelines/%d/objects%s", pipelineID, suffix)
}

func objectPath(pipelineID int64, objectName string) string {
	return objectsPath(pipelineID, "/"+url.PathEscape(objectName))
}
