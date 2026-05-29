// Package dbeavergoMCP is the root module package hosting embedded filesystem assets.
// By placing embed directives here, go:embed paths can reference the templates/
// directory which sits at the same level as this file.
package dbeavergoMCP

import "embed"

// SQLTemplates embeds all SQL script templates under templates/sql/.
//
//go:embed templates/sql
var SQLTemplates embed.FS

// SnippetTemplates embeds all SQL snippet templates under templates/snippets/.
//
//go:embed templates/snippets
var SnippetTemplates embed.FS

// ProjectTemplates embeds built-in project templates under templates/projects/.
//
//go:embed templates/projects
var ProjectTemplates embed.FS
