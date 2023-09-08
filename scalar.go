// Copyright (c) Mondoo, Inc.
// SPDX-License-Identifier: BUSL-1.1

package mondoogql

import (
	"time"

	"github.com/shurcooL/graphql"
)

type (
	// Boolean represents true or false values.
	Boolean graphql.Boolean

	// Date is an ISO-8601 encoded UTC date.
	Date struct{ time.Time }

	// DateTime is an ISO-8601 encoded UTC date and time.
	DateTime struct{ time.Time }

	// Float represents signed double-precision fractional values as
	// specified by IEEE 754.
	Float graphql.Float

	// ID represents a unique identifier that is Base64 obfuscated. It is
	// often used to refetch an object or as key for a cache. The ID type
	// appears in a JSON response as a String; however, it is not intended
	// to be human-readable. When expected as an input type, any string
	// (such as "VXNlci0xMA==") or integer (such as 4) input value will be
	// accepted as an ID.
	ID graphql.ID

	// Int represents non-fractional signed whole numeric values.
	Int graphql.Int

	// String represents textual data as UTF-8 character sequences.
	String graphql.String

	// Map represents a generic map of string keys to values.
	Map map[string]interface{}
)

// NewBooleanPtr creates a new *Boolean.
func NewBooleanPtr(v Boolean) *Boolean { return &v }

// NewDatePtr creates a new *Date.
func NewDatePtr(v Date) *Date { return &v }

// NewDateTimePtr creates a new *DateTime.
func NewDateTimePtr(v DateTime) *DateTime { return &v }

// NewFloatPtr creates a new *Float.
func NewFloatPtr(v Float) *Float { return &v }

// NewIDPtr creates a new *ID.
func NewIDPtr(v ID) *ID { return &v }

// NewIntPtr creates a new *Int.
func NewIntPtr(v Int) *Int { return &v }

// NewStringPtr creates a new *String.
func NewStringPtr(v String) *String { return &v }
