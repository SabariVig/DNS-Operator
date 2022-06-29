package providers

import (
	"context"

	v1 "github.com/SabariVig/DNS-Operator/api/v1"
)

type Providers interface {
	AddRecord(context.Context, *v1.Record) error
	DeleteRecord(context.Context, *v1.Record) error
}
