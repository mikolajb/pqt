package pqtsql

import (
	"bytes"
	"errors"
	"fmt"
	"io"

	"github.com/piotrkowalczuk/pqt"
)

type generator struct{}

// Generator ...
func Generator() *generator {
	return &generator{}
}

func (g *generator) Generate(s *pqt.Schema) ([]byte, error) {
	code, err := g.generate(s)
	if err != nil {
		return nil, err
	}

	return code.Bytes(), nil
}

func (g *generator) GenerateTo(s *pqt.Schema, w io.Writer) error {
	code, err := g.generate(s)
	if err != nil {
		return err
	}

	_, err = code.WriteTo(w)
	return err
}

func (g *generator) generate(s *pqt.Schema) (*bytes.Buffer, error) {
	code := bytes.NewBufferString("-- do not modify, generated by pqt\n\n")
	if s.Name != "" {
		fmt.Fprintf(code, "CREATE SCHEMA IF NOT EXISTS %s; \n\n", s.Name)
	}
	for _, t := range s.Tables {
		if err := g.generateCreateTable(code, t); err != nil {
			return nil, err
		}
	}

	return code, nil
}

func (g *generator) generateCreateTable(buf *bytes.Buffer, t *pqt.Table) error {
	if t == nil {
		return nil
	}

	if t.Name == "" {
		return errors.New("pqt: missing table name")
	}
	if len(t.Columns) == 0 {
		return fmt.Errorf("pqt: table %s has no columns", t.Name)
	}

	constraints := tableConstraints(t)

	buf.WriteString("CREATE ")
	if t.Temporary {
		buf.WriteString("TEMPORARY ")
	}
	buf.WriteString("TABLE ")
	if t.IfNotExists {
		buf.WriteString("IF NOT EXISTS ")
	}
	if t.Schema != nil {
		buf.WriteString(t.Schema.Name)
		buf.WriteRune('.')
		buf.WriteString(t.Name)
	} else {
		buf.WriteString(t.Name)
	}
	buf.WriteString(" (\n")
	for i, c := range t.Columns {
		buf.WriteRune('	')
		buf.WriteString(c.Name)
		buf.WriteRune(' ')
		buf.WriteString(c.Type.String())
		if c.Collate != "" {
			buf.WriteRune(' ')
			buf.WriteString(c.Collate)
		}
		if c.Default != "" {
			buf.WriteString(" DEFAULT ")
			buf.WriteString(c.Default)
		}
		if c.NotNull {
			buf.WriteString(" NOT NULL")
		}
		if i < len(t.Columns)-1 || len(constraints) > 0 {
			buf.WriteRune(',')
		}
		buf.WriteRune('\n')
	}

	if len(constraints) > 0 {
		buf.WriteRune('\n')
	}

	for i, c := range constraints {
		buf.WriteString("	")
		err := g.generateConstraint(buf, c)
		if err != nil {
			return err
		}
		if i < len(constraints)-1 {
			buf.WriteRune(',')
		}
		buf.WriteRune('\n')
	}

	buf.WriteString(");\n\n")

	return nil
}

func (g *generator) generateConstraint(buf *bytes.Buffer, c *pqt.Constraint) error {
	switch c.Type {
	case pqt.ConstraintTypeUnique:
		uniqueConstraintQuery(buf, c)
	case pqt.ConstraintTypePrimaryKey:
		primaryKeyConstraintQuery(buf, c)
	case pqt.ConstraintTypeForeignKey:
		return foreignKeyConstraintQuery(buf, c)
	case pqt.ConstraintTypeCheck:
		checkConstraintQuery(buf, c)
	default:
		return fmt.Errorf("pqt: unknown constraint type: %s", c.Type)
	}

	return nil
}

func uniqueConstraintQuery(buf *bytes.Buffer, c *pqt.Constraint) {
	fmt.Fprintf(buf, `CONSTRAINT "%s" UNIQUE (%s)`, c.Name(), pqt.JoinColumns(c.Columns, ", "))
}

func primaryKeyConstraintQuery(buf *bytes.Buffer, c *pqt.Constraint) {
	fmt.Fprintf(buf, `CONSTRAINT "%s" PRIMARY KEY (%s)`, c.Name(), pqt.JoinColumns(c.Columns, ", "))
}

func foreignKeyConstraintQuery(buf *bytes.Buffer, c *pqt.Constraint) error {
	switch {
	case len(c.Columns) == 0:
		return errors.New("pqt: foreign key constraint require at least one column")
	case len(c.ReferenceColumns) == 0:
		return errors.New("pqt: foreign key constraint require at least one reference column")
	case c.ReferenceTable == nil:
		return errors.New("pqt: foreiqn key constraint missing reference table")
	}

	fmt.Fprintf(buf, `CONSTRAINT "%s" FOREIGN KEY (%s) REFERENCES %s (%s)`,
		c.Name(),
		pqt.JoinColumns(c.Columns, ", "),
		c.ReferenceTable.FullName(),
		pqt.JoinColumns(c.ReferenceColumns, ", "),
	)

	return nil
}

func checkConstraintQuery(buf *bytes.Buffer, c *pqt.Constraint) {
	fmt.Fprintf(buf, `CONSTRAINT "%s" CHECK (%s)`, c.Name(), c.Check)
}

func tableConstraints(t *pqt.Table) []*pqt.Constraint {
	constraints := make([]*pqt.Constraint, 0, len(t.Constraints)+len(t.Columns))
	for _, c := range t.Columns {
		constraints = append(constraints, c.Constraints()...)
	}

	return append(constraints, t.Constraints...)
}
