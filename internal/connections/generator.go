package connections

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strings"
	"text/template"
)

// OutFile is an output artifact: a relative path and its content.
// Mirrors output.File but avoids a circular import when tests reference this package.
type OutFile struct {
	Path    string
	Content []byte
}

// ConnConfig describes a single database connection for generation purposes.
// The CLI layer maps config.ConnConfig → connections.ConnConfig before calling Generate.
type ConnConfig struct {
	Name     string
	Type     string
	Host     string
	Port     int
	Database string
}

// ConnSpec holds everything needed to generate DBeaver connection files.
type ConnSpec struct {
	Project      string
	Environments []string
	Connections  []ConnConfig
}

// Generate produces data-sources.json and data-sources-config.json as output files.
// data-sources.json follows the DBeaver 26.0 workspace format.
// data-sources-config.json is always emitted as {}.
// save-password is always set to false — no exceptions.
func Generate(spec ConnSpec) ([]OutFile, error) {
	folders := map[string]folderEntry{}
	conns := map[string]connEntry{}

	for _, env := range spec.Environments {
		folders[env] = folderEntry{Description: env + " connections"}
	}

	for _, env := range spec.Environments {
		for _, cc := range spec.Connections {
			driver, err := Lookup(cc.Type)
			if err != nil {
				return nil, fmt.Errorf("connections.Generate: %w", err)
			}

			renderedName, err := renderTemplate(cc.Name, struct {
				Project  string
				Env      string
				DBEngine string
			}{
				Project:  spec.Project,
				Env:      env,
				DBEngine: cc.Type,
			})
			if err != nil {
				return nil, fmt.Errorf("connections.Generate: render name %q: %w", cc.Name, err)
			}

			id := connectionID(driver.Provider, driver.Driver, cc.Name, env)

			conns[id] = connEntry{
				Provider:          driver.Provider,
				Driver:            driver.Driver,
				Name:              renderedName,
				SavePassword:      false,
				ShowSystemObjects: true,
				ReadOnly:          false,
				Folder:            env,
				Configuration: connConfiguration{
					Host:      cc.Host,
					Port:      fmt.Sprintf("%d", cc.Port),
					Database:  cc.Database,
					URL:       buildJDBCURL(cc.Type, cc.Host, cc.Port, cc.Database),
					Type:      strings.ToLower(env),
					AuthModel: "native",
				},
			}
		}
	}

	dsDoc := map[string]interface{}{
		"folders":     folders,
		"connections": conns,
	}

	dsJSON, err := json.MarshalIndent(dsDoc, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("connections.Generate: marshal data-sources.json: %w", err)
	}

	return []OutFile{
		{Path: "data-sources.json", Content: dsJSON},
		{Path: "data-sources-config.json", Content: []byte("{}")},
	}, nil
}

// ── helpers ──────────────────────────────────────────────────────────────────

type folderEntry struct {
	Description string `json:"description"`
}

type connEntry struct {
	Provider          string            `json:"provider"`
	Driver            string            `json:"driver"`
	Name              string            `json:"name"`
	SavePassword      bool              `json:"save-password"`
	ShowSystemObjects bool              `json:"show-system-objects"`
	ReadOnly          bool              `json:"read-only"`
	Folder            string            `json:"folder"`
	Configuration     connConfiguration `json:"configuration"`
}

type connConfiguration struct {
	Host      string `json:"host"`
	Port      string `json:"port"`
	Database  string `json:"database"`
	URL       string `json:"url"`
	Type      string `json:"type"`
	AuthModel string `json:"auth-model"`
}

// connectionID returns the first 16 hex chars of SHA-256(provider|driver|name|env).
func connectionID(provider, driver, name, env string) string {
	h := sha256.Sum256([]byte(provider + "|" + driver + "|" + name + "|" + env))
	return fmt.Sprintf("%x", h)[:16]
}

// renderTemplate executes a Go text/template string with the given data.
func renderTemplate(tmpl string, data interface{}) (string, error) {
	t, err := template.New("").Parse(tmpl)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// buildJDBCURL constructs a JDBC connection URL from its components.
func buildJDBCURL(engineType, host string, port int, database string) string {
	switch engineType {
	case "postgresql":
		return fmt.Sprintf("jdbc:postgresql://%s:%d/%s", host, port, database)
	case "sqlserver":
		return fmt.Sprintf("jdbc:sqlserver://%s:%d;databaseName=%s", host, port, database)
	case "mysql":
		return fmt.Sprintf("jdbc:mysql://%s:%d/%s", host, port, database)
	case "oracle":
		return fmt.Sprintf("jdbc:oracle:thin:@%s:%d:%s", host, port, database)
	case "mongodb":
		return fmt.Sprintf("mongodb://%s:%d/%s", host, port, database)
	default:
		return fmt.Sprintf("jdbc:%s://%s:%d/%s", engineType, host, port, database)
	}
}
