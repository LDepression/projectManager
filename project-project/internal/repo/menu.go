/**
 * @Author: lenovo
 * @Description:
 * @File:  menu
 * @Version: 1.0.0
 * @Date: 2023/07/21 16:15
 */

package repo

import (
	"context"
	"projectManager/project-project/internal/data/menu"
)

type MenuRepo interface {
	FindMenus(ctx context.Context) ([]*menu.ProjectMenu, error)
}
