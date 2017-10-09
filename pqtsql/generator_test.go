package pqtsql_test

import (
	"testing"

	"github.com/piotrkowalczuk/pqt"
	"github.com/piotrkowalczuk/pqt/pqtsql"
)

func TestGenerator_Generate(t *testing.T) {
	success := []struct {
		expected string
		given    *pqt.Table
		version  float64
	}{
		{
			expected: `-- do not modify, generated by pqt

CREATE TEMPORARY TABLE schema.user (
	created_at TIMESTAMPTZ,
	password TEXT,
	username TEXT NOT NULL
);
CREATE INDEX IF NOT EXISTS "schema.user_username_idx" ON schema.user (username);

`,
			given: func() *pqt.Table {
				return pqt.NewTable("user", pqt.WithTemporary()).
					SetSchema(pqt.NewSchema("schema")).
					AddColumn(&pqt.Column{Name: "username", Type: pqt.TypeText(), NotNull: true, Index: true}).
					AddColumn(&pqt.Column{Name: "password", Type: pqt.TypeText()}).
					AddColumn(&pqt.Column{Name: "created_at", Type: pqt.TypeTimestampTZ()})
			}(),
			version: 9.5,
		},
		{
			expected: `-- do not modify, generated by pqt

CREATE TABLE IF NOT EXISTS table_name (
	created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
	created_by INTEGER NOT NULL,
	enabled BOOL,
	end_at TIMESTAMPTZ NOT NULL,
	id SERIAL,
	name TEXT,
	one BIGINT,
	price DECIMAL(10,1),
	rel_id INTEGER,
	slug TEXT NOT NULL,
	start_at TIMESTAMPTZ NOT NULL,
	two BIGINT,
	updated_at TIMESTAMPTZ,
	updated_by INTEGER,

	CONSTRAINT "public.table_name_id_pkey" PRIMARY KEY (id),
	CONSTRAINT "public.table_name_rel_id_fkey" FOREIGN KEY (rel_id) REFERENCES related_table (id),
	CONSTRAINT "public.table_name_name_key" UNIQUE (name),
	CONSTRAINT "public.table_name_slug_key" UNIQUE (slug),
	CONSTRAINT "public.table_name_start_at_end_at_check" CHECK ((start_at IS NULL AND end_at IS NULL) OR start_at < end_at)
);
CREATE INDEX "public.table_name_start_at_idx" ON table_name (start_at);
CREATE UNIQUE INDEX "public.table_name_one_two_V1tbrhqs_uidx" ON table_name (one,two) WHERE one IS NOT NULL AND two IS NULL AND one > 2;

`,
			given: func() *pqt.Table {
				id := pqt.Column{Name: "id", Type: pqt.TypeSerial()}
				startAt := &pqt.Column{Name: "start_at", Type: pqt.TypeTimestampTZ(), NotNull: true}
				endAt := &pqt.Column{Name: "end_at", Type: pqt.TypeTimestampTZ(), NotNull: true}

				one := pqt.NewColumn("one", pqt.TypeIntegerBig())
				two := pqt.NewColumn("two", pqt.TypeIntegerBig())

				_ = pqt.NewTable("related_table").
					AddColumn(&id)

				return pqt.NewTable("table_name", pqt.WithTableIfNotExists()).
					AddColumn(&pqt.Column{Name: "id", Type: pqt.TypeSerial(), PrimaryKey: true}).
					AddColumn(&pqt.Column{Name: "rel_id", Type: pqt.TypeInteger(), Reference: &id}).
					AddColumn(&pqt.Column{Name: "name", Type: pqt.TypeText(), Unique: true}).
					AddColumn(&pqt.Column{Name: "enabled", Type: pqt.TypeBool()}).
					AddColumn(&pqt.Column{Name: "price", Type: pqt.TypeDecimal(10, 1)}).
					AddColumn(startAt).
					AddIndex(startAt).
					AddColumn(endAt).
					AddColumn(pqt.NewColumn("created_at", pqt.TypeTimestampTZ(), pqt.WithNotNull(), pqt.WithDefault("NOW()"))).
					AddColumn(&pqt.Column{Name: "created_by", Type: pqt.TypeInteger(), NotNull: true}).
					AddColumn(&pqt.Column{Name: "updated_at", Type: pqt.TypeTimestampTZ()}).
					AddColumn(&pqt.Column{Name: "updated_by", Type: pqt.TypeInteger()}).
					AddColumn(&pqt.Column{Name: "slug", Type: pqt.TypeText(), NotNull: true, Unique: true}).
					AddCheck("(start_at IS NULL AND end_at IS NULL) OR start_at < end_at", startAt, endAt).
					AddColumn(one).
					AddColumn(two).
					AddPartialUniqueIndex("one IS NOT NULL AND two IS NULL AND one > 2", one, two)
			}(),
		},
	}

	for i, data := range success {
		g := &pqtsql.Generator{
			Version: data.version,
		}
		q, err := g.Generate(&pqt.Schema{
			Tables: []*pqt.Table{data.given},
		})

		if err != nil {
			t.Errorf("unexpected error for schema #%d : %s", i, err.Error())
			continue
		}

		if string(q) != data.expected {
			t.Errorf("wrong query, expected:\n'%s'\nbut got:\n'%s'", data.expected, q)
		}
	}
}
