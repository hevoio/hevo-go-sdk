package pipelineobjects

// ObjectSelection selects or deselects one object and optionally changes its load mode and
// fields. Build it with the client's NewObjectSelection.
type ObjectSelection struct {
	name          *string
	selected      *bool
	loadMode      *string
	includeFields []string
	excludeFields []string
}

type objectSelectionRequest struct {
	Name          *string  `json:"name,omitempty"`
	Selected      *bool    `json:"selected,omitempty"`
	LoadMode      *string  `json:"load_mode,omitempty"`
	IncludeFields []string `json:"include_fields,omitempty"`
	ExcludeFields []string `json:"exclude_fields,omitempty"`
}

// Name sets the object, as its name appears in the list of objects.
func (o *ObjectSelection) Name(value string) *ObjectSelection {
	o.name = &value

	return o
}

// Selected starts (true) or stops (false) replicating the object.
func (o *ObjectSelection) Selected(value bool) *ObjectSelection {
	o.selected = &value

	return o
}

// LoadMode sets MERGE, APPEND or HISTORY_TRACKING for this object.
func (o *ObjectSelection) LoadMode(value string) *ObjectSelection {
	o.loadMode = &value

	return o
}

// IncludeFields starts replicating these fields.
func (o *ObjectSelection) IncludeFields(values ...string) *ObjectSelection {
	o.includeFields = values

	return o
}

// ExcludeFields stops replicating these fields.
func (o *ObjectSelection) ExcludeFields(values ...string) *ObjectSelection {
	o.excludeFields = values

	return o
}

// request returns the selection as the API expects it.
func (o *ObjectSelection) request() *objectSelectionRequest {
	if o == nil {
		return nil
	}

	return &objectSelectionRequest{
		Name:          o.name,
		Selected:      o.selected,
		LoadMode:      o.loadMode,
		IncludeFields: o.includeFields,
		ExcludeFields: o.excludeFields,
	}
}
