package dbeaver

import "encoding/json"

// dataSources mirrors the top-level structure of DBeaver's data-sources.json.
// json.RawMessage is used for connection and folder values to preserve unknown fields.
type dataSources struct {
	Folders     map[string]json.RawMessage `json:"folders"`
	Connections map[string]json.RawMessage `json:"connections"`
}

// MergeDataSources merges incoming data-sources.json bytes into existing ones.
//
// Rules (connection and folder maps):
//   - Existing IDs are always preserved (existing wins on conflict).
//   - New IDs from incoming are added.
//   - IDs present only in existing are never removed.
//   - Result is idempotent: MergeDataSources(MergeDataSources(e,i), i) == MergeDataSources(e,i).
func MergeDataSources(existing, incoming []byte) ([]byte, error) {
	var ex, inc dataSources
	if err := json.Unmarshal(existing, &ex); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(incoming, &inc); err != nil {
		return nil, err
	}

	if ex.Connections == nil {
		ex.Connections = map[string]json.RawMessage{}
	}
	if ex.Folders == nil {
		ex.Folders = map[string]json.RawMessage{}
	}

	// Add new connections from incoming (existing wins on ID conflict).
	for id, v := range inc.Connections {
		if _, exists := ex.Connections[id]; !exists {
			ex.Connections[id] = v
		}
	}

	// Add new folders from incoming (existing wins on ID conflict).
	for id, v := range inc.Folders {
		if _, exists := ex.Folders[id]; !exists {
			ex.Folders[id] = v
		}
	}

	return json.MarshalIndent(ex, "", "  ")
}
