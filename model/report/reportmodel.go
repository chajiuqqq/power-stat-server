package report

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ ReportModel = (*customReportModel)(nil)

type (
	// ReportModel is an interface to be customized, add more methods here,
	// and implement the added methods in customReportModel.
	ReportModel interface {
		reportModel
		FindLatest(ctx context.Context) (*Report, error)
	}

	customReportModel struct {
		*defaultReportModel
	}
)

// NewReportModel returns a model for the database table.
func NewReportModel(conn sqlx.SqlConn) ReportModel {
	return &customReportModel{
		defaultReportModel: newReportModel(conn),
	}
}

func (m *defaultReportModel) FindLatest(ctx context.Context) (*Report, error) {
	query := fmt.Sprintf("select %s from %s order by `id` desc limit 1", reportRows, m.table)
	var resp Report
	err := m.conn.QueryRowCtx(ctx, &resp, query)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}
