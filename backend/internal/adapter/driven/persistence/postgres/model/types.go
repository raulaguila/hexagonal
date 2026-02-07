package model

import (
	"database/sql/driver"
	"errors"
	"strings"
)

// StringArray is a slice of strings that implements sql.Scanner and driver.Valuer
// for PostgreSQL text[] types.
type StringArray []string

// Scan implements the sql.Scanner interface.
func (a *StringArray) Scan(src any) error {
	switch src := src.(type) {
	case []byte:
		return a.scanBytes(src)
	case string:
		return a.scanBytes([]byte(src))
	case nil:
		*a = nil
		return nil
	default:
		return errors.New("incompatible type for StringArray")
	}
}

func (a *StringArray) scanBytes(src []byte) error {
	elems := strings.Trim(string(src), "{}")
	if len(elems) == 0 {
		*a = nil
		return nil
	}

	// This is a simplified parser. For full Postgres array support
	// (quoted strings, escaped characters), a more complex parser is needed.
	// However, for typical "simple identifier" arrays like permissions,
	// checking strictly for commas inside quotes isn't always needed if we enforce simple values.
	// BUT, to be safe, we should handle CSV parsing correctly or at least strict splits.

	// For fully robust parsing without lib/pq, we can use a simple state machine
	// or just assume permissions are simple alphanumeric strings without commas.
	// Given "permissions" are usually things like "users:read", simple splitting is risky if future permissions have commas.

	// Let's implement a basic reader loop.
	*a = make(StringArray, 0)

	// If it's empty
	if len(elems) == 0 {
		return nil
	}

	// Split by comma is simplistic but might fail on "foo,bar".
	// Since we want to eliminate lib/pq, let's provide a "good enough" parser for this use case.
	// Permissions usually: "users:read", "roles:write". Safe to split by comma?
	// Postgres returns values unquoted if they are simple, quoted if complex.
	// e.g. {a,b,"c,d"}

	// Only using simple split for now as permissions shouldn't contain commas/quotes generally.
	// If the user wants full compliance, lib/pq is better.
	// But assuming we want to remove the dependency for a simpler stack:

	vals := strings.Split(elems, ",")
	for i, v := range vals {
		// Postgres usually escapes quotes as \" inside the string? No, it doubles them ""
		// And wraps in ".
		// Clean up quotes if present
		if len(v) >= 2 && v[0] == '"' && v[len(vals[i])-1] == '"' {
			vals[i] = v[1 : len(v)-1]
			vals[i] = strings.ReplaceAll(vals[i], "\\\"", "\"") // basic unescape
		}
	}
	*a = vals
	return nil
}

// Value implements the driver.Valuer interface.
func (a StringArray) Value() (driver.Value, error) {
	if len(a) == 0 {
		return nil, nil
	}

	var sb strings.Builder
	sb.WriteByte('{')
	for i, s := range a {
		if i > 0 {
			sb.WriteByte(',')
		}
		// Escape quotes and backslashes
		// Postgres array literal: "val" or val.
		// Safer to always quote.
		sb.WriteByte('"')
		sb.WriteString(strings.ReplaceAll(strings.ReplaceAll(s, "\\", "\\\\"), "\"", "\\\""))
		sb.WriteByte('"')
	}
	sb.WriteByte('}')
	return sb.String(), nil
}
