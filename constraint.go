package pqt

import (
	"fmt"
	"regexp"
	"strings"
)

const (
	// ConstraintTypeUnknown ...
	ConstraintTypeUnknown ConstraintType = "unknown"
	// ConstraintTypePrimaryKey ...
	ConstraintTypePrimaryKey ConstraintType = "pkey"
	// ConstraintTypeCheck ...
	ConstraintTypeCheck ConstraintType = "check"
	// ConstraintTypeUnique ...
	ConstraintTypeUnique ConstraintType = "key"
	// ConstraintTypeIndex ...
	ConstraintTypeIndex ConstraintType = "idx"
	// ConstraintTypeForeignKey ...
	ConstraintTypeForeignKey ConstraintType = "fkey"
	// ConstraintTypeExclusion ...
	ConstraintTypeExclusion ConstraintType = "excl"
	// ConstraintTypeUniqueIndex ...
	ConstraintTypeUniqueIndex ConstraintType = "uidx"
)

type ConstraintType string

// ConstraintOption ...
type ConstraintOption func(*Constraint)

// Constraint ...
type Constraint struct {
	Type                                                                 ConstraintType
	Where, Check                                                         string
	PrimaryTable, Table                                                  *Table
	PrimaryColumns, Columns                                              Columns
	Attribute                                                            []*Attribute
	Match, OnDelete, OnUpdate                                            int32
	NoInherit, DeferrableInitiallyDeferred, DeferrableInitiallyImmediate bool
}

// Name ...
func (c *Constraint) Name() string {
	var schema string

	switch {
	case c.PrimaryTable == nil:
		return "<missing table>"
	case c.PrimaryTable.Schema == nil || c.PrimaryTable.Schema.Name == "":
		schema = "public"
	default:
		schema = c.PrimaryTable.Schema.Name
	}

	if len(c.PrimaryColumns) == 0 {
		return fmt.Sprintf("%s.%s_%s", schema, c.PrimaryTable.ShortName, c.Type)
	}
	tmp := make([]string, 0, len(c.PrimaryColumns))
	for _, col := range c.PrimaryColumns {
		if col.ShortName != "" {
			tmp = append(tmp, col.ShortName)
			continue
		}
		tmp = append(tmp, col.Name)
	}

	if len(c.Where) > 0 {
		tmp = append(tmp, "WHERE")
		reg := regexp.MustCompile("([a-zA-Z0-9]+)")
		tmp = append(tmp, reg.FindAllString(c.Where, 16)...)
	}

	return fmt.Sprintf("%s.%s_%s_%s", schema, c.PrimaryTable.ShortName, strings.Join(tmp, "_"), c.Type)
}

// Unique constraint ensure that the data contained in a column or a group of columns is unique with respect to all the rows in the table.
func Unique(table *Table, columns ...*Column) *Constraint {
	return &Constraint{
		Type:           ConstraintTypeUnique,
		PrimaryTable:   table,
		PrimaryColumns: columns,
	}
}

// PrimaryKey constraint is simply a combination of a unique constraint and a not-null constraint.
func PrimaryKey(table *Table, columns ...*Column) *Constraint {
	return &Constraint{
		Type:           ConstraintTypePrimaryKey,
		PrimaryTable:   table,
		PrimaryColumns: columns,
	}
}

// Check ...
func Check(table *Table, check string, columns ...*Column) *Constraint {
	return &Constraint{
		Type:           ConstraintTypeCheck,
		PrimaryTable:   table,
		PrimaryColumns: columns,
		Check:          check,
	}
}

//
//// Exclusion constraint ensure that if any two rows are compared on the specified columns
//// or expressions using the specified operators,
//// at least one of these operator comparisons will return false or null.
//func Exclusion(table *Table, exclude Exclude, columns ...*Column) *Constraint {
//	return &Constraint{
//		Type:    ConstraintTypeExclusion,
//		Table:   table,
//		Exclude: exclude,
//		Columns: columns,
//	}
//}

// Reference ...
type Reference struct {
	From, To *Column
}

// ForeignKey constraint specifies that the values in a column (or a group of columns)
// must match the values appearing in some row of another table.
// We say this maintains the referential integrity between two related tables.
func ForeignKey(primaryColumns, referenceColumns Columns, opts ...ConstraintOption) *Constraint {
	if len(referenceColumns) == 0 {
		panic("foreign key expects at least one reference column")
	}
	for _, c := range primaryColumns {
		if c.Table != primaryColumns[0].Table {
			panic("column tables inconsistency")
		}
	}
	for _, r := range referenceColumns {
		if r.Table != referenceColumns[0].Table {
			panic("reference column tables inconsistency")
		}
	}
	fk := &Constraint{
		Type:           ConstraintTypeForeignKey,
		PrimaryTable:   primaryColumns[0].Table,
		PrimaryColumns: primaryColumns,
		Table:          referenceColumns[0].Table,
		Columns:        referenceColumns,
	}

	for _, o := range opts {
		o(fk)
	}

	return fk
}

// Index ...
func Index(table *Table, columns ...*Column) *Constraint {
	return &Constraint{
		Type:           ConstraintTypeIndex,
		PrimaryTable:   table,
		PrimaryColumns: columns,
	}
}

// UniqueIndex ...
func UniqueIndex(table *Table, columns ...*Column) *Constraint {
	return &Constraint{
		Type:           ConstraintTypeUniqueIndex,
		PrimaryTable:   table,
		PrimaryColumns: columns,
	}
}

// PartialUniqueIndex ...
func PartialUniqueIndex(table *Table, where string, columns ...*Column) *Constraint {
	return &Constraint{
		Type:           ConstraintTypeUniqueIndex,
		PrimaryTable:   table,
		PrimaryColumns: columns,
		Where:          where,
	}
}

// String implements Stringer interface.
func (c *Constraint) String() string {
	return c.Name()
}

// IsForeignKey returns true if string has suffix "_fkey".
func IsForeignKey(c string) bool {
	return strings.HasSuffix(c, string(ConstraintTypeForeignKey))
}

// IsUnique returns true if string has suffix "_key".
func IsUnique(c string) bool {
	return strings.HasSuffix(c, string(ConstraintTypeUnique))
}

// IsPrimaryKey returns true if string has suffix "_pkey".
func IsPrimaryKey(c string) bool {
	return strings.HasSuffix(c, string(ConstraintTypePrimaryKey))
}

// IsCheck returns true if string has suffix "_check".
func IsCheck(c string) bool {
	return strings.HasSuffix(c, string(ConstraintTypeCheck))
}

//// IsExclusion returns true if string has suffix "_excl".
//func IsExclusion(c string) bool {
//	return strings.HasSuffix(c, string(ConstraintTypeExclusion))
//}

// IsIndex returns true if string has suffix "_idx".
func IsIndex(c string) bool {
	return strings.HasSuffix(c, string(ConstraintTypeIndex))
}

type Constraints []*Constraint

// CountOf returns number of constraints of given type.
// If nothing is given return length of entire slice.
func (c Constraints) CountOf(types ...ConstraintType) int {
	if len(types) == 0 {
		return len(c)
	}
	var count int
OuterLoop:
	for _, cc := range c {
		for _, t := range types {
			if cc.Type == t {
				count++
				continue OuterLoop
			}
		}
	}
	return count
}
