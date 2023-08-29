package repo

import (
	"context"
	"projectManager/project-project/internal/data"
)

type FileRepo interface {
	Save(ctx context.Context, file *data.File) error
	FindByIds(back23ground context.Context, ids []int64) (list []*data.File, err error)
}
