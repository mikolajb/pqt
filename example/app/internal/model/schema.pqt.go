// Code generated by pqt.
// source: cmd/appg/main.go
// DO NOT EDIT!
package model

import (
	"bytes"
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/lib/pq"
	"github.com/piotrkowalczuk/pqt/pqtgo"
)

func ScanCategoryRows(rows *sql.Rows) (entities []*CategoryEntity, err error) {

	for rows.Next() {
		var ent CategoryEntity
		err = rows.Scan(
			&ent.Content,
			&ent.CreatedAt,
			&ent.ID,
			&ent.Name,
			&ent.ParentID,
			&ent.UpdatedAt,
		)
		if err != nil {
			return
		}

		entities = append(entities, &ent)
	}
	if err = rows.Err(); err != nil {
		return
	}

	return
}

// CategoryEntity ...
type CategoryEntity struct {
	// Content ...
	Content string
	// CreatedAt ...
	CreatedAt time.Time
	// ID ...
	ID int64
	// Name ...
	Name string
	// ParentID ...
	ParentID sql.NullInt64
	// UpdatedAt ...
	UpdatedAt pq.NullTime
	// ChildCategory ...
	ChildCategory []*CategoryEntity
	// ParentCategory ...
	ParentCategory *CategoryEntity
	// Packages ...
	Packages []*PackageEntity
}

type CategoryCriteria struct {
	Offset, Limit int64
	Sort          map[string]bool
	Content       sql.NullString
	CreatedAt     pq.NullTime
	ID            sql.NullInt64
	Name          sql.NullString
	ParentID      sql.NullInt64
	UpdatedAt     pq.NullTime
}

const (
	TableCategory                             = "example.category"
	TableCategoryColumnContent                = "content"
	TableCategoryColumnCreatedAt              = "created_at"
	TableCategoryColumnID                     = "id"
	TableCategoryColumnName                   = "name"
	TableCategoryColumnParentID               = "parent_id"
	TableCategoryColumnUpdatedAt              = "updated_at"
	TableCategoryConstraintPrimaryKey         = "example.category_id_pkey"
	TableCategoryConstraintParentIDForeignKey = "example.category_parent_id_fkey"
)

var (
	TableCategoryColumns = []string{
		TableCategoryColumnContent,
		TableCategoryColumnCreatedAt,
		TableCategoryColumnID,
		TableCategoryColumnName,
		TableCategoryColumnParentID,
		TableCategoryColumnUpdatedAt,
	}
)

type CategoryRepositoryBase struct {
	Table   string
	Columns []string
	DB      *sql.DB
	Debug   bool
	Log     log.Logger
}

func (r *CategoryRepositoryBase) InsertQuery(e *CategoryEntity) (string, []interface{}, error) {
	ins := pqtgo.NewComposer(6)
	buf := bytes.NewBufferString("INSERT INTO " + r.Table)
	col := bytes.NewBuffer(nil)
	if col.Len() > 0 {
		col.WriteString(", ")
	}
	if _, err := col.WriteString(TableCategoryColumnContent); err != nil {
		return "", nil, err
	}
	if ins.Dirty {
		ins.WriteString(", ")
	}
	if err := ins.WritePlaceholder(); err != nil {
		return "", nil, err
	}
	ins.Add(e.Content)
	ins.Dirty = true
	if col.Len() > 0 {
		col.WriteString(", ")
	}
	if _, err := col.WriteString(TableCategoryColumnCreatedAt); err != nil {
		return "", nil, err
	}
	if ins.Dirty {
		ins.WriteString(", ")
	}
	if err := ins.WritePlaceholder(); err != nil {
		return "", nil, err
	}
	ins.Add(e.CreatedAt)
	ins.Dirty = true
	if col.Len() > 0 {
		col.WriteString(", ")
	}
	if _, err := col.WriteString(TableCategoryColumnName); err != nil {
		return "", nil, err
	}
	if ins.Dirty {
		ins.WriteString(", ")
	}
	if err := ins.WritePlaceholder(); err != nil {
		return "", nil, err
	}
	ins.Add(e.Name)
	ins.Dirty = true
	if e.ParentID.Valid {
		if col.Len() > 0 {
			col.WriteString(", ")
		}
		if _, err := col.WriteString(TableCategoryColumnParentID); err != nil {
			return "", nil, err
		}
		if ins.Dirty {
			ins.WriteString(", ")
		}
		if err := ins.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		ins.Add(e.ParentID)
		ins.Dirty = true
	}
	if e.UpdatedAt.Valid {
		if col.Len() > 0 {
			col.WriteString(", ")
		}
		if _, err := col.WriteString(TableCategoryColumnUpdatedAt); err != nil {
			return "", nil, err
		}
		if ins.Dirty {
			ins.WriteString(", ")
		}
		if err := ins.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		ins.Add(e.UpdatedAt)
		ins.Dirty = true
	}
	if col.Len() > 0 {
		buf.WriteString(" (")
		buf.ReadFrom(col)
		buf.WriteString(") VALUES (")
		buf.ReadFrom(ins)
		buf.WriteString(") ")
		if len(r.Columns) > 0 {
			buf.WriteString("RETURNING ")
			buf.WriteString(strings.Join(r.Columns, ", "))
		}
	}
	return buf.String(), ins.Args(), nil
}
func (r *CategoryRepositoryBase) Insert(ctx context.Context, e *CategoryEntity) (*CategoryEntity, error) {
	query, args, err := r.InsertQuery(e)
	if err != nil {
		return nil, err
	}
	if err := r.DB.QueryRowContext(ctx, query, args...).Scan(&e.Content,
		&e.CreatedAt,
		&e.ID,
		&e.Name,
		&e.ParentID,
		&e.UpdatedAt,
	); err != nil {
		if r.Debug {
			r.Log.Log("level", "error", "timestamp", time.Now().Format(time.RFC3339), "msg", "insert query failure", "query", query, "table", r.Table, "error", err.Error())
		}
		return nil, err
	}
	if r.Debug {
		r.Log.Log("level", "debug", "timestamp", time.Now().Format(time.RFC3339), "msg", "insert query success", "query", query, "table", r.Table)
	}
	return e, nil
}
func (r *CategoryRepositoryBase) FindQuery(s []string, c *CategoryCriteria) (string, []interface{}, error) {
	where := pqtgo.NewComposer(6)
	buf := bytes.NewBufferString("SELECT ")
	buf.WriteString(strings.Join(s, ", "))
	buf.WriteString(" FROM ")
	buf.WriteString(r.Table)
	buf.WriteString(" ")
	if c.Content.Valid {
		if where.Dirty {
			where.WriteString(" AND ")
		}
		if _, err := where.WriteString(TableCategoryColumnContent); err != nil {
			return "", nil, err
		}
		if _, err := where.WriteString("="); err != nil {
			return "", nil, err
		}
		if err := where.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		where.Add(c.Content)
		where.Dirty = true
	}
	if c.CreatedAt.Valid {
		if where.Dirty {
			where.WriteString(" AND ")
		}
		if _, err := where.WriteString(TableCategoryColumnCreatedAt); err != nil {
			return "", nil, err
		}
		if _, err := where.WriteString("="); err != nil {
			return "", nil, err
		}
		if err := where.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		where.Add(c.CreatedAt)
		where.Dirty = true
	}
	if c.ID.Valid {
		if where.Dirty {
			where.WriteString(" AND ")
		}
		if _, err := where.WriteString(TableCategoryColumnID); err != nil {
			return "", nil, err
		}
		if _, err := where.WriteString("="); err != nil {
			return "", nil, err
		}
		if err := where.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		where.Add(c.ID)
		where.Dirty = true
	}
	if c.Name.Valid {
		if where.Dirty {
			where.WriteString(" AND ")
		}
		if _, err := where.WriteString(TableCategoryColumnName); err != nil {
			return "", nil, err
		}
		if _, err := where.WriteString("="); err != nil {
			return "", nil, err
		}
		if err := where.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		where.Add(c.Name)
		where.Dirty = true
	}
	if c.ParentID.Valid {
		if where.Dirty {
			where.WriteString(" AND ")
		}
		if _, err := where.WriteString(TableCategoryColumnParentID); err != nil {
			return "", nil, err
		}
		if _, err := where.WriteString("="); err != nil {
			return "", nil, err
		}
		if err := where.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		where.Add(c.ParentID)
		where.Dirty = true
	}
	if c.UpdatedAt.Valid {
		if where.Dirty {
			where.WriteString(" AND ")
		}
		if _, err := where.WriteString(TableCategoryColumnUpdatedAt); err != nil {
			return "", nil, err
		}
		if _, err := where.WriteString("="); err != nil {
			return "", nil, err
		}
		if err := where.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		where.Add(c.UpdatedAt)
		where.Dirty = true
	}

	if where.Dirty {
		buf.WriteString("WHERE ")
		buf.ReadFrom(where)
	}
	return buf.String(), where.Args(), nil
}
func (r *CategoryRepositoryBase) Find(ctx context.Context, c *CategoryCriteria) ([]*CategoryEntity, error) {
	query, args, err := r.FindQuery(r.Columns, c)
	if err != nil {
		return nil, err
	}
	rows, err := r.DB.QueryContext(ctx, query, args...)
	if err != nil {
		if r.Debug {
			r.Log.Log("level", "error", "timestamp", time.Now().Format(time.RFC3339), "msg", "insert query failure", "query", query, "table", r.Table, "error", err.Error())
		}
		return nil, err
	}
	defer rows.Close()

	if r.Debug {
		r.Log.Log("level", "debug", "timestamp", time.Now().Format(time.RFC3339), "msg", "insert query success", "query", query, "table", r.Table)
	}

	return ScanCategoryRows(rows)
}
func (r *CategoryRepositoryBase) Count(ctx context.Context, c *CategoryCriteria) (int64, error) {
	query, args, err := r.FindQuery([]string{"COUNT(*)"}, c)
	if err != nil {
		return 0, err
	}
	var count int64
	if err := r.DB.QueryRowContext(ctx, query, args...).Scan(&count); err != nil {
		if r.Debug {
			r.Log.Log("level", "error", "timestamp", time.Now().Format(time.RFC3339), "msg", "insert query failure", "query", query, "table", r.Table, "error", err.Error())
		}
		return 0, err
	}

	if r.Debug {
		r.Log.Log("level", "debug", "timestamp", time.Now().Format(time.RFC3339), "msg", "insert query success", "query", query, "table", r.Table)
	}

	return count, nil
}
func ScanPackageRows(rows *sql.Rows) (entities []*PackageEntity, err error) {

	for rows.Next() {
		var ent PackageEntity
		err = rows.Scan(
			&ent.Break,
			&ent.CategoryID,
			&ent.CreatedAt,
			&ent.ID,
			&ent.UpdatedAt,
		)
		if err != nil {
			return
		}

		entities = append(entities, &ent)
	}
	if err = rows.Err(); err != nil {
		return
	}

	return
}

// PackageEntity ...
type PackageEntity struct {
	// Break ...
	Break sql.NullString
	// CategoryID ...
	CategoryID sql.NullInt64
	// CreatedAt ...
	CreatedAt time.Time
	// ID ...
	ID int64
	// UpdatedAt ...
	UpdatedAt pq.NullTime
	// Category ...
	Category *CategoryEntity
}

type PackageCriteria struct {
	Offset, Limit int64
	Sort          map[string]bool
	Break         sql.NullString
	CategoryID    sql.NullInt64
	CreatedAt     pq.NullTime
	ID            sql.NullInt64
	UpdatedAt     pq.NullTime
}

const (
	TablePackage                               = "example.package"
	TablePackageColumnBreak                    = "break"
	TablePackageColumnCategoryID               = "category_id"
	TablePackageColumnCreatedAt                = "created_at"
	TablePackageColumnID                       = "id"
	TablePackageColumnUpdatedAt                = "updated_at"
	TablePackageConstraintCategoryIDForeignKey = "example.package_category_id_fkey"
	TablePackageConstraintPrimaryKey           = "example.package_id_pkey"
)

var (
	TablePackageColumns = []string{
		TablePackageColumnBreak,
		TablePackageColumnCategoryID,
		TablePackageColumnCreatedAt,
		TablePackageColumnID,
		TablePackageColumnUpdatedAt,
	}
)

type PackageRepositoryBase struct {
	Table   string
	Columns []string
	DB      *sql.DB
	Debug   bool
	Log     log.Logger
}

func (r *PackageRepositoryBase) InsertQuery(e *PackageEntity) (string, []interface{}, error) {
	ins := pqtgo.NewComposer(5)
	buf := bytes.NewBufferString("INSERT INTO " + r.Table)
	col := bytes.NewBuffer(nil)
	if e.Break.Valid {
		if col.Len() > 0 {
			col.WriteString(", ")
		}
		if _, err := col.WriteString(TablePackageColumnBreak); err != nil {
			return "", nil, err
		}
		if ins.Dirty {
			ins.WriteString(", ")
		}
		if err := ins.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		ins.Add(e.Break)
		ins.Dirty = true
	}
	if e.CategoryID.Valid {
		if col.Len() > 0 {
			col.WriteString(", ")
		}
		if _, err := col.WriteString(TablePackageColumnCategoryID); err != nil {
			return "", nil, err
		}
		if ins.Dirty {
			ins.WriteString(", ")
		}
		if err := ins.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		ins.Add(e.CategoryID)
		ins.Dirty = true
	}
	if col.Len() > 0 {
		col.WriteString(", ")
	}
	if _, err := col.WriteString(TablePackageColumnCreatedAt); err != nil {
		return "", nil, err
	}
	if ins.Dirty {
		ins.WriteString(", ")
	}
	if err := ins.WritePlaceholder(); err != nil {
		return "", nil, err
	}
	ins.Add(e.CreatedAt)
	ins.Dirty = true
	if e.UpdatedAt.Valid {
		if col.Len() > 0 {
			col.WriteString(", ")
		}
		if _, err := col.WriteString(TablePackageColumnUpdatedAt); err != nil {
			return "", nil, err
		}
		if ins.Dirty {
			ins.WriteString(", ")
		}
		if err := ins.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		ins.Add(e.UpdatedAt)
		ins.Dirty = true
	}
	if col.Len() > 0 {
		buf.WriteString(" (")
		buf.ReadFrom(col)
		buf.WriteString(") VALUES (")
		buf.ReadFrom(ins)
		buf.WriteString(") ")
		if len(r.Columns) > 0 {
			buf.WriteString("RETURNING ")
			buf.WriteString(strings.Join(r.Columns, ", "))
		}
	}
	return buf.String(), ins.Args(), nil
}
func (r *PackageRepositoryBase) Insert(ctx context.Context, e *PackageEntity) (*PackageEntity, error) {
	query, args, err := r.InsertQuery(e)
	if err != nil {
		return nil, err
	}
	if err := r.DB.QueryRowContext(ctx, query, args...).Scan(&e.Break,
		&e.CategoryID,
		&e.CreatedAt,
		&e.ID,
		&e.UpdatedAt,
	); err != nil {
		if r.Debug {
			r.Log.Log("level", "error", "timestamp", time.Now().Format(time.RFC3339), "msg", "insert query failure", "query", query, "table", r.Table, "error", err.Error())
		}
		return nil, err
	}
	if r.Debug {
		r.Log.Log("level", "debug", "timestamp", time.Now().Format(time.RFC3339), "msg", "insert query success", "query", query, "table", r.Table)
	}
	return e, nil
}
func (r *PackageRepositoryBase) FindQuery(s []string, c *PackageCriteria) (string, []interface{}, error) {
	where := pqtgo.NewComposer(5)
	buf := bytes.NewBufferString("SELECT ")
	buf.WriteString(strings.Join(s, ", "))
	buf.WriteString(" FROM ")
	buf.WriteString(r.Table)
	buf.WriteString(" ")
	if c.Break.Valid {
		if where.Dirty {
			where.WriteString(" AND ")
		}
		if _, err := where.WriteString(TablePackageColumnBreak); err != nil {
			return "", nil, err
		}
		if _, err := where.WriteString("="); err != nil {
			return "", nil, err
		}
		if err := where.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		where.Add(c.Break)
		where.Dirty = true
	}
	if c.CategoryID.Valid {
		if where.Dirty {
			where.WriteString(" AND ")
		}
		if _, err := where.WriteString(TablePackageColumnCategoryID); err != nil {
			return "", nil, err
		}
		if _, err := where.WriteString("="); err != nil {
			return "", nil, err
		}
		if err := where.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		where.Add(c.CategoryID)
		where.Dirty = true
	}
	if c.CreatedAt.Valid {
		if where.Dirty {
			where.WriteString(" AND ")
		}
		if _, err := where.WriteString(TablePackageColumnCreatedAt); err != nil {
			return "", nil, err
		}
		if _, err := where.WriteString("="); err != nil {
			return "", nil, err
		}
		if err := where.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		where.Add(c.CreatedAt)
		where.Dirty = true
	}
	if c.ID.Valid {
		if where.Dirty {
			where.WriteString(" AND ")
		}
		if _, err := where.WriteString(TablePackageColumnID); err != nil {
			return "", nil, err
		}
		if _, err := where.WriteString("="); err != nil {
			return "", nil, err
		}
		if err := where.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		where.Add(c.ID)
		where.Dirty = true
	}
	if c.UpdatedAt.Valid {
		if where.Dirty {
			where.WriteString(" AND ")
		}
		if _, err := where.WriteString(TablePackageColumnUpdatedAt); err != nil {
			return "", nil, err
		}
		if _, err := where.WriteString("="); err != nil {
			return "", nil, err
		}
		if err := where.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		where.Add(c.UpdatedAt)
		where.Dirty = true
	}

	if where.Dirty {
		buf.WriteString("WHERE ")
		buf.ReadFrom(where)
	}
	return buf.String(), where.Args(), nil
}
func (r *PackageRepositoryBase) Find(ctx context.Context, c *PackageCriteria) ([]*PackageEntity, error) {
	query, args, err := r.FindQuery(r.Columns, c)
	if err != nil {
		return nil, err
	}
	rows, err := r.DB.QueryContext(ctx, query, args...)
	if err != nil {
		if r.Debug {
			r.Log.Log("level", "error", "timestamp", time.Now().Format(time.RFC3339), "msg", "insert query failure", "query", query, "table", r.Table, "error", err.Error())
		}
		return nil, err
	}
	defer rows.Close()

	if r.Debug {
		r.Log.Log("level", "debug", "timestamp", time.Now().Format(time.RFC3339), "msg", "insert query success", "query", query, "table", r.Table)
	}

	return ScanPackageRows(rows)
}
func (r *PackageRepositoryBase) Count(ctx context.Context, c *PackageCriteria) (int64, error) {
	query, args, err := r.FindQuery([]string{"COUNT(*)"}, c)
	if err != nil {
		return 0, err
	}
	var count int64
	if err := r.DB.QueryRowContext(ctx, query, args...).Scan(&count); err != nil {
		if r.Debug {
			r.Log.Log("level", "error", "timestamp", time.Now().Format(time.RFC3339), "msg", "insert query failure", "query", query, "table", r.Table, "error", err.Error())
		}
		return 0, err
	}

	if r.Debug {
		r.Log.Log("level", "debug", "timestamp", time.Now().Format(time.RFC3339), "msg", "insert query success", "query", query, "table", r.Table)
	}

	return count, nil
}
func ScanNewsRows(rows *sql.Rows) (entities []*NewsEntity, err error) {

	for rows.Next() {
		var ent NewsEntity
		err = rows.Scan(
			&ent.Content,
			&ent.Continue,
			&ent.CreatedAt,
			&ent.ID,
			&ent.Lead,
			&ent.Title,
			&ent.UpdatedAt,
		)
		if err != nil {
			return
		}

		entities = append(entities, &ent)
	}
	if err = rows.Err(); err != nil {
		return
	}

	return
}

// NewsEntity ...
type NewsEntity struct {
	// Content ...
	Content string
	// Continue ...
	Continue bool
	// CreatedAt ...
	CreatedAt time.Time
	// ID ...
	ID int64
	// Lead ...
	Lead sql.NullString
	// Title ...
	Title string
	// UpdatedAt ...
	UpdatedAt pq.NullTime
	// CommentsByNewsTitle ...
	CommentsByNewsTitle []*CommentEntity
	// Comments ...
	Comments []*CommentEntity
}

type NewsCriteria struct {
	Offset, Limit int64
	Sort          map[string]bool
	Content       sql.NullString
	Continue      sql.NullBool
	CreatedAt     pq.NullTime
	ID            sql.NullInt64
	Lead          sql.NullString
	Title         sql.NullString
	UpdatedAt     pq.NullTime
}

const (
	TableNews                          = "example.news"
	TableNewsColumnContent             = "content"
	TableNewsColumnContinue            = "continue"
	TableNewsColumnCreatedAt           = "created_at"
	TableNewsColumnID                  = "id"
	TableNewsColumnLead                = "lead"
	TableNewsColumnTitle               = "title"
	TableNewsColumnUpdatedAt           = "updated_at"
	TableNewsConstraintPrimaryKey      = "example.news_id_pkey"
	TableNewsConstraintTitleUnique     = "example.news_title_key"
	TableNewsConstraintTitleLeadUnique = "example.news_title_lead_key"
)

var (
	TableNewsColumns = []string{
		TableNewsColumnContent,
		TableNewsColumnContinue,
		TableNewsColumnCreatedAt,
		TableNewsColumnID,
		TableNewsColumnLead,
		TableNewsColumnTitle,
		TableNewsColumnUpdatedAt,
	}
)

type NewsRepositoryBase struct {
	Table   string
	Columns []string
	DB      *sql.DB
	Debug   bool
	Log     log.Logger
}

func (r *NewsRepositoryBase) InsertQuery(e *NewsEntity) (string, []interface{}, error) {
	ins := pqtgo.NewComposer(7)
	buf := bytes.NewBufferString("INSERT INTO " + r.Table)
	col := bytes.NewBuffer(nil)
	if col.Len() > 0 {
		col.WriteString(", ")
	}
	if _, err := col.WriteString(TableNewsColumnContent); err != nil {
		return "", nil, err
	}
	if ins.Dirty {
		ins.WriteString(", ")
	}
	if err := ins.WritePlaceholder(); err != nil {
		return "", nil, err
	}
	ins.Add(e.Content)
	ins.Dirty = true
	if col.Len() > 0 {
		col.WriteString(", ")
	}
	if _, err := col.WriteString(TableNewsColumnContinue); err != nil {
		return "", nil, err
	}
	if ins.Dirty {
		ins.WriteString(", ")
	}
	if err := ins.WritePlaceholder(); err != nil {
		return "", nil, err
	}
	ins.Add(e.Continue)
	ins.Dirty = true
	if col.Len() > 0 {
		col.WriteString(", ")
	}
	if _, err := col.WriteString(TableNewsColumnCreatedAt); err != nil {
		return "", nil, err
	}
	if ins.Dirty {
		ins.WriteString(", ")
	}
	if err := ins.WritePlaceholder(); err != nil {
		return "", nil, err
	}
	ins.Add(e.CreatedAt)
	ins.Dirty = true
	if e.Lead.Valid {
		if col.Len() > 0 {
			col.WriteString(", ")
		}
		if _, err := col.WriteString(TableNewsColumnLead); err != nil {
			return "", nil, err
		}
		if ins.Dirty {
			ins.WriteString(", ")
		}
		if err := ins.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		ins.Add(e.Lead)
		ins.Dirty = true
	}
	if col.Len() > 0 {
		col.WriteString(", ")
	}
	if _, err := col.WriteString(TableNewsColumnTitle); err != nil {
		return "", nil, err
	}
	if ins.Dirty {
		ins.WriteString(", ")
	}
	if err := ins.WritePlaceholder(); err != nil {
		return "", nil, err
	}
	ins.Add(e.Title)
	ins.Dirty = true
	if e.UpdatedAt.Valid {
		if col.Len() > 0 {
			col.WriteString(", ")
		}
		if _, err := col.WriteString(TableNewsColumnUpdatedAt); err != nil {
			return "", nil, err
		}
		if ins.Dirty {
			ins.WriteString(", ")
		}
		if err := ins.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		ins.Add(e.UpdatedAt)
		ins.Dirty = true
	}
	if col.Len() > 0 {
		buf.WriteString(" (")
		buf.ReadFrom(col)
		buf.WriteString(") VALUES (")
		buf.ReadFrom(ins)
		buf.WriteString(") ")
		if len(r.Columns) > 0 {
			buf.WriteString("RETURNING ")
			buf.WriteString(strings.Join(r.Columns, ", "))
		}
	}
	return buf.String(), ins.Args(), nil
}
func (r *NewsRepositoryBase) Insert(ctx context.Context, e *NewsEntity) (*NewsEntity, error) {
	query, args, err := r.InsertQuery(e)
	if err != nil {
		return nil, err
	}
	if err := r.DB.QueryRowContext(ctx, query, args...).Scan(&e.Content,
		&e.Continue,
		&e.CreatedAt,
		&e.ID,
		&e.Lead,
		&e.Title,
		&e.UpdatedAt,
	); err != nil {
		if r.Debug {
			r.Log.Log("level", "error", "timestamp", time.Now().Format(time.RFC3339), "msg", "insert query failure", "query", query, "table", r.Table, "error", err.Error())
		}
		return nil, err
	}
	if r.Debug {
		r.Log.Log("level", "debug", "timestamp", time.Now().Format(time.RFC3339), "msg", "insert query success", "query", query, "table", r.Table)
	}
	return e, nil
}
func (r *NewsRepositoryBase) FindQuery(s []string, c *NewsCriteria) (string, []interface{}, error) {
	where := pqtgo.NewComposer(7)
	buf := bytes.NewBufferString("SELECT ")
	buf.WriteString(strings.Join(s, ", "))
	buf.WriteString(" FROM ")
	buf.WriteString(r.Table)
	buf.WriteString(" ")
	if c.Content.Valid {
		if where.Dirty {
			where.WriteString(" AND ")
		}
		if _, err := where.WriteString(TableNewsColumnContent); err != nil {
			return "", nil, err
		}
		if _, err := where.WriteString("="); err != nil {
			return "", nil, err
		}
		if err := where.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		where.Add(c.Content)
		where.Dirty = true
	}
	if c.Continue.Valid {
		if where.Dirty {
			where.WriteString(" AND ")
		}
		if _, err := where.WriteString(TableNewsColumnContinue); err != nil {
			return "", nil, err
		}
		if _, err := where.WriteString("="); err != nil {
			return "", nil, err
		}
		if err := where.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		where.Add(c.Continue)
		where.Dirty = true
	}
	if c.CreatedAt.Valid {
		if where.Dirty {
			where.WriteString(" AND ")
		}
		if _, err := where.WriteString(TableNewsColumnCreatedAt); err != nil {
			return "", nil, err
		}
		if _, err := where.WriteString("="); err != nil {
			return "", nil, err
		}
		if err := where.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		where.Add(c.CreatedAt)
		where.Dirty = true
	}
	if c.ID.Valid {
		if where.Dirty {
			where.WriteString(" AND ")
		}
		if _, err := where.WriteString(TableNewsColumnID); err != nil {
			return "", nil, err
		}
		if _, err := where.WriteString("="); err != nil {
			return "", nil, err
		}
		if err := where.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		where.Add(c.ID)
		where.Dirty = true
	}
	if c.Lead.Valid {
		if where.Dirty {
			where.WriteString(" AND ")
		}
		if _, err := where.WriteString(TableNewsColumnLead); err != nil {
			return "", nil, err
		}
		if _, err := where.WriteString("="); err != nil {
			return "", nil, err
		}
		if err := where.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		where.Add(c.Lead)
		where.Dirty = true
	}
	if c.Title.Valid {
		if where.Dirty {
			where.WriteString(" AND ")
		}
		if _, err := where.WriteString(TableNewsColumnTitle); err != nil {
			return "", nil, err
		}
		if _, err := where.WriteString("="); err != nil {
			return "", nil, err
		}
		if err := where.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		where.Add(c.Title)
		where.Dirty = true
	}
	if c.UpdatedAt.Valid {
		if where.Dirty {
			where.WriteString(" AND ")
		}
		if _, err := where.WriteString(TableNewsColumnUpdatedAt); err != nil {
			return "", nil, err
		}
		if _, err := where.WriteString("="); err != nil {
			return "", nil, err
		}
		if err := where.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		where.Add(c.UpdatedAt)
		where.Dirty = true
	}

	if where.Dirty {
		buf.WriteString("WHERE ")
		buf.ReadFrom(where)
	}
	return buf.String(), where.Args(), nil
}
func (r *NewsRepositoryBase) Find(ctx context.Context, c *NewsCriteria) ([]*NewsEntity, error) {
	query, args, err := r.FindQuery(r.Columns, c)
	if err != nil {
		return nil, err
	}
	rows, err := r.DB.QueryContext(ctx, query, args...)
	if err != nil {
		if r.Debug {
			r.Log.Log("level", "error", "timestamp", time.Now().Format(time.RFC3339), "msg", "insert query failure", "query", query, "table", r.Table, "error", err.Error())
		}
		return nil, err
	}
	defer rows.Close()

	if r.Debug {
		r.Log.Log("level", "debug", "timestamp", time.Now().Format(time.RFC3339), "msg", "insert query success", "query", query, "table", r.Table)
	}

	return ScanNewsRows(rows)
}
func (r *NewsRepositoryBase) Count(ctx context.Context, c *NewsCriteria) (int64, error) {
	query, args, err := r.FindQuery([]string{"COUNT(*)"}, c)
	if err != nil {
		return 0, err
	}
	var count int64
	if err := r.DB.QueryRowContext(ctx, query, args...).Scan(&count); err != nil {
		if r.Debug {
			r.Log.Log("level", "error", "timestamp", time.Now().Format(time.RFC3339), "msg", "insert query failure", "query", query, "table", r.Table, "error", err.Error())
		}
		return 0, err
	}

	if r.Debug {
		r.Log.Log("level", "debug", "timestamp", time.Now().Format(time.RFC3339), "msg", "insert query success", "query", query, "table", r.Table)
	}

	return count, nil
}
func ScanCommentRows(rows *sql.Rows) (entities []*CommentEntity, err error) {

	for rows.Next() {
		var ent CommentEntity
		err = rows.Scan(
			&ent.Content,
			&ent.CreatedAt,
			&ent.ID,
			&ent.NewsID,
			&ent.NewsTitle,
			&ent.UpdatedAt,
		)
		if err != nil {
			return
		}

		entities = append(entities, &ent)
	}
	if err = rows.Err(); err != nil {
		return
	}

	return
}

// CommentEntity ...
type CommentEntity struct {
	// Content ...
	Content string
	// CreatedAt ...
	CreatedAt time.Time
	// ID ...
	ID sql.NullInt64
	// NewsID ...
	NewsID int64
	// NewsTitle ...
	NewsTitle string
	// UpdatedAt ...
	UpdatedAt pq.NullTime
	// NewsByTitle ...
	NewsByTitle *NewsEntity
	// NewsByID ...
	NewsByID *NewsEntity
}

type CommentCriteria struct {
	Offset, Limit int64
	Sort          map[string]bool
	Content       sql.NullString
	CreatedAt     pq.NullTime
	ID            sql.NullInt64
	NewsID        sql.NullInt64
	NewsTitle     sql.NullString
	UpdatedAt     pq.NullTime
}

const (
	TableComment                              = "example.comment"
	TableCommentColumnContent                 = "content"
	TableCommentColumnCreatedAt               = "created_at"
	TableCommentColumnID                      = "id"
	TableCommentColumnNewsID                  = "news_id"
	TableCommentColumnNewsTitle               = "news_title"
	TableCommentColumnUpdatedAt               = "updated_at"
	TableCommentConstraintNewsIDForeignKey    = "example.comment_news_id_fkey"
	TableCommentConstraintNewsTitleForeignKey = "example.comment_news_title_fkey"
)

var (
	TableCommentColumns = []string{
		TableCommentColumnContent,
		TableCommentColumnCreatedAt,
		TableCommentColumnID,
		TableCommentColumnNewsID,
		TableCommentColumnNewsTitle,
		TableCommentColumnUpdatedAt,
	}
)

type CommentRepositoryBase struct {
	Table   string
	Columns []string
	DB      *sql.DB
	Debug   bool
	Log     log.Logger
}

func (r *CommentRepositoryBase) InsertQuery(e *CommentEntity) (string, []interface{}, error) {
	ins := pqtgo.NewComposer(6)
	buf := bytes.NewBufferString("INSERT INTO " + r.Table)
	col := bytes.NewBuffer(nil)
	if col.Len() > 0 {
		col.WriteString(", ")
	}
	if _, err := col.WriteString(TableCommentColumnContent); err != nil {
		return "", nil, err
	}
	if ins.Dirty {
		ins.WriteString(", ")
	}
	if err := ins.WritePlaceholder(); err != nil {
		return "", nil, err
	}
	ins.Add(e.Content)
	ins.Dirty = true
	if col.Len() > 0 {
		col.WriteString(", ")
	}
	if _, err := col.WriteString(TableCommentColumnCreatedAt); err != nil {
		return "", nil, err
	}
	if ins.Dirty {
		ins.WriteString(", ")
	}
	if err := ins.WritePlaceholder(); err != nil {
		return "", nil, err
	}
	ins.Add(e.CreatedAt)
	ins.Dirty = true
	if col.Len() > 0 {
		col.WriteString(", ")
	}
	if _, err := col.WriteString(TableCommentColumnNewsID); err != nil {
		return "", nil, err
	}
	if ins.Dirty {
		ins.WriteString(", ")
	}
	if err := ins.WritePlaceholder(); err != nil {
		return "", nil, err
	}
	ins.Add(e.NewsID)
	ins.Dirty = true
	if col.Len() > 0 {
		col.WriteString(", ")
	}
	if _, err := col.WriteString(TableCommentColumnNewsTitle); err != nil {
		return "", nil, err
	}
	if ins.Dirty {
		ins.WriteString(", ")
	}
	if err := ins.WritePlaceholder(); err != nil {
		return "", nil, err
	}
	ins.Add(e.NewsTitle)
	ins.Dirty = true
	if e.UpdatedAt.Valid {
		if col.Len() > 0 {
			col.WriteString(", ")
		}
		if _, err := col.WriteString(TableCommentColumnUpdatedAt); err != nil {
			return "", nil, err
		}
		if ins.Dirty {
			ins.WriteString(", ")
		}
		if err := ins.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		ins.Add(e.UpdatedAt)
		ins.Dirty = true
	}
	if col.Len() > 0 {
		buf.WriteString(" (")
		buf.ReadFrom(col)
		buf.WriteString(") VALUES (")
		buf.ReadFrom(ins)
		buf.WriteString(") ")
		if len(r.Columns) > 0 {
			buf.WriteString("RETURNING ")
			buf.WriteString(strings.Join(r.Columns, ", "))
		}
	}
	return buf.String(), ins.Args(), nil
}
func (r *CommentRepositoryBase) Insert(ctx context.Context, e *CommentEntity) (*CommentEntity, error) {
	query, args, err := r.InsertQuery(e)
	if err != nil {
		return nil, err
	}
	if err := r.DB.QueryRowContext(ctx, query, args...).Scan(&e.Content,
		&e.CreatedAt,
		&e.ID,
		&e.NewsID,
		&e.NewsTitle,
		&e.UpdatedAt,
	); err != nil {
		if r.Debug {
			r.Log.Log("level", "error", "timestamp", time.Now().Format(time.RFC3339), "msg", "insert query failure", "query", query, "table", r.Table, "error", err.Error())
		}
		return nil, err
	}
	if r.Debug {
		r.Log.Log("level", "debug", "timestamp", time.Now().Format(time.RFC3339), "msg", "insert query success", "query", query, "table", r.Table)
	}
	return e, nil
}
func (r *CommentRepositoryBase) FindQuery(s []string, c *CommentCriteria) (string, []interface{}, error) {
	where := pqtgo.NewComposer(6)
	buf := bytes.NewBufferString("SELECT ")
	buf.WriteString(strings.Join(s, ", "))
	buf.WriteString(" FROM ")
	buf.WriteString(r.Table)
	buf.WriteString(" ")
	if c.Content.Valid {
		if where.Dirty {
			where.WriteString(" AND ")
		}
		if _, err := where.WriteString(TableCommentColumnContent); err != nil {
			return "", nil, err
		}
		if _, err := where.WriteString("="); err != nil {
			return "", nil, err
		}
		if err := where.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		where.Add(c.Content)
		where.Dirty = true
	}
	if c.CreatedAt.Valid {
		if where.Dirty {
			where.WriteString(" AND ")
		}
		if _, err := where.WriteString(TableCommentColumnCreatedAt); err != nil {
			return "", nil, err
		}
		if _, err := where.WriteString("="); err != nil {
			return "", nil, err
		}
		if err := where.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		where.Add(c.CreatedAt)
		where.Dirty = true
	}
	if c.ID.Valid {
		if where.Dirty {
			where.WriteString(" AND ")
		}
		if _, err := where.WriteString(TableCommentColumnID); err != nil {
			return "", nil, err
		}
		if _, err := where.WriteString("="); err != nil {
			return "", nil, err
		}
		if err := where.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		where.Add(c.ID)
		where.Dirty = true
	}
	if c.NewsID.Valid {
		if where.Dirty {
			where.WriteString(" AND ")
		}
		if _, err := where.WriteString(TableCommentColumnNewsID); err != nil {
			return "", nil, err
		}
		if _, err := where.WriteString("="); err != nil {
			return "", nil, err
		}
		if err := where.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		where.Add(c.NewsID)
		where.Dirty = true
	}
	if c.NewsTitle.Valid {
		if where.Dirty {
			where.WriteString(" AND ")
		}
		if _, err := where.WriteString(TableCommentColumnNewsTitle); err != nil {
			return "", nil, err
		}
		if _, err := where.WriteString("="); err != nil {
			return "", nil, err
		}
		if err := where.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		where.Add(c.NewsTitle)
		where.Dirty = true
	}
	if c.UpdatedAt.Valid {
		if where.Dirty {
			where.WriteString(" AND ")
		}
		if _, err := where.WriteString(TableCommentColumnUpdatedAt); err != nil {
			return "", nil, err
		}
		if _, err := where.WriteString("="); err != nil {
			return "", nil, err
		}
		if err := where.WritePlaceholder(); err != nil {
			return "", nil, err
		}
		where.Add(c.UpdatedAt)
		where.Dirty = true
	}

	if where.Dirty {
		buf.WriteString("WHERE ")
		buf.ReadFrom(where)
	}
	return buf.String(), where.Args(), nil
}
func (r *CommentRepositoryBase) Find(ctx context.Context, c *CommentCriteria) ([]*CommentEntity, error) {
	query, args, err := r.FindQuery(r.Columns, c)
	if err != nil {
		return nil, err
	}
	rows, err := r.DB.QueryContext(ctx, query, args...)
	if err != nil {
		if r.Debug {
			r.Log.Log("level", "error", "timestamp", time.Now().Format(time.RFC3339), "msg", "insert query failure", "query", query, "table", r.Table, "error", err.Error())
		}
		return nil, err
	}
	defer rows.Close()

	if r.Debug {
		r.Log.Log("level", "debug", "timestamp", time.Now().Format(time.RFC3339), "msg", "insert query success", "query", query, "table", r.Table)
	}

	return ScanCommentRows(rows)
}
func (r *CommentRepositoryBase) Count(ctx context.Context, c *CommentCriteria) (int64, error) {
	query, args, err := r.FindQuery([]string{"COUNT(*)"}, c)
	if err != nil {
		return 0, err
	}
	var count int64
	if err := r.DB.QueryRowContext(ctx, query, args...).Scan(&count); err != nil {
		if r.Debug {
			r.Log.Log("level", "error", "timestamp", time.Now().Format(time.RFC3339), "msg", "insert query failure", "query", query, "table", r.Table, "error", err.Error())
		}
		return 0, err
	}

	if r.Debug {
		r.Log.Log("level", "debug", "timestamp", time.Now().Format(time.RFC3339), "msg", "insert query success", "query", query, "table", r.Table)
	}

	return count, nil
}

/// SQL ...
const SQL = `
-- do not modify, generated by pqt

CREATE SCHEMA IF NOT EXISTS example; 

CREATE TABLE IF NOT EXISTS example.category (
	content TEXT NOT NULL,
	created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
	id BIGSERIAL,
	name TEXT NOT NULL,
	parent_id BIGINT,
	updated_at TIMESTAMPTZ,

	CONSTRAINT "example.category_id_pkey" PRIMARY KEY (id),
	CONSTRAINT "example.category_parent_id_fkey" FOREIGN KEY (parent_id) REFERENCES example.category (id)
);

CREATE TABLE IF NOT EXISTS example.package (
	break TEXT,
	category_id BIGINT,
	created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
	id BIGSERIAL,
	updated_at TIMESTAMPTZ,

	CONSTRAINT "example.package_category_id_fkey" FOREIGN KEY (category_id) REFERENCES example.category (id),
	CONSTRAINT "example.package_id_pkey" PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS example.news (
	content TEXT NOT NULL,
	continue BOOL DEFAULT false NOT NULL,
	created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
	id BIGSERIAL,
	lead TEXT,
	title TEXT NOT NULL,
	updated_at TIMESTAMPTZ,

	CONSTRAINT "example.news_id_pkey" PRIMARY KEY (id),
	CONSTRAINT "example.news_title_key" UNIQUE (title),
	CONSTRAINT "example.news_title_lead_key" UNIQUE (title, lead)
);

CREATE TABLE IF NOT EXISTS example.comment (
	content TEXT NOT NULL,
	created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
	id BIGSERIAL,
	news_id BIGINT NOT NULL,
	news_title TEXT NOT NULL,
	updated_at TIMESTAMPTZ,

	CONSTRAINT "example.comment_news_id_fkey" FOREIGN KEY (news_id) REFERENCES example.news (id),
	CONSTRAINT "example.comment_news_title_fkey" FOREIGN KEY (news_title) REFERENCES example.news (title)
);

`
