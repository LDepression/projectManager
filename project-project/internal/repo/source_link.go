/**
 * @Author: lenovo
 * @Description:
 * @File:  source_link
 * @Version: 1.0.0
 * @Date: 2023/08/29 21:53
 */

package repo

import (
	"context"
	"projectManager/project-project/internal/data"
)

type SourceLinkRepo interface {
	Save(ctx context.Context, link *data.SourceLink) error
	FindByTaskCode(ctx context.Context, taskCode int64) (list []*data.SourceLink, err error)
}
