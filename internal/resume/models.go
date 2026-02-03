package resume

// This file is intentionally minimal. We treat resumes as raw YAML strings
// and let Claude handle the structure, only doing basic validation.

// RenderCV has complex nested structures that vary by template.
// Rather than defining rigid Go structs, we pass YAML as strings
// and validate it can be parsed as valid YAML.
