package vo

import "shiyu-admin-backend/internal/model/entity"

// MenuVO represents menu info for API responses.
type MenuVO struct {
	MenuCode   string    `json:"menu_code"`
	ParentCode string    `json:"parent_code"`
	MenuType   string    `json:"menu_type"`
	MenuName   string    `json:"menu_name"`
	Perms      string    `json:"perms"`
	Path       string    `json:"path"`
	Component  string    `json:"component"`
	Status     int       `json:"status"`
	Children   []*MenuVO `json:"children,omitempty"`
}

// BuildMenuVO converts entity to VO.
func BuildMenuVO(m *entity.Menu) *MenuVO {
	if m == nil {
		return nil
	}
	return &MenuVO{
		MenuCode:   m.MenuCode,
		ParentCode: m.ParentCode,
		MenuType:   m.MenuType,
		MenuName:   m.MenuName,
		Perms:      m.Perms,
		Path:       m.Path,
		Component:  m.Component,
		Status:     m.Status,
	}
}

// BuildMenuTree builds a tree structure from flat menu list.
func BuildMenuTree(menus []*entity.Menu) []*MenuVO {
	menuMap := make(map[string]*MenuVO)
	var roots []*MenuVO

	// First pass: create all menu VOs
	for _, m := range menus {
		vo := BuildMenuVO(m)
		menuMap[m.MenuCode] = vo
	}

	// Second pass: build tree
	for _, m := range menus {
		vo := menuMap[m.MenuCode]
		if m.ParentCode == "" {
			roots = append(roots, vo)
		} else {
			parent := menuMap[m.ParentCode]
			if parent != nil {
				if parent.Children == nil {
					parent.Children = []*MenuVO{}
				}
				parent.Children = append(parent.Children, vo)
			}
		}
	}

	return roots
}

