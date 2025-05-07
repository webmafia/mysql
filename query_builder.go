package mysql

import (
	"slices"
	"sync"

	"github.com/webmafia/fast/buffer"
)

type queryBuilder struct {
	buf buffer.Pool
	arg sync.Pool
}

func (qb *queryBuilder) getArgs(cap int) *[]any {
	if args, ok := qb.arg.Get().(*[]any); ok {
		*args = slices.Grow(*args, cap)
		return args
	}

	args := make([]any, 0, cap)
	return &args
}

func (qb *queryBuilder) putArgs(args *[]any) {
	if args == nil {
		return
	}

	for i := range *args {
		(*args)[i] = nil
	}

	*args = (*args)[:0]

	qb.arg.Put(args)
}

func (qb *queryBuilder) buildQuery(dstQuery *buffer.Buffer, dstArgs *[]any, query string, args []any) (err error) {
	return dstQuery.Str().WritefCb(query, args, func(b *buffer.Buffer, c byte, arg any) (err error) {
		writeAny(b, dstArgs, arg)
		return
	})
}
