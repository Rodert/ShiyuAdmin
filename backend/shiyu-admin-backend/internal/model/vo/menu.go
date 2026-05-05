package vo

import (
	"sort"

	"shiyu-admin-backend/internal/model/entity"
)

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
	SortOrder  int       `json:"sort_order"`
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
		SortOrder:  m.SortOrder,
	}
}

func sortMenuVONodes(nodes []*MenuVO) {
	if len(nodes) == 0 {
		return
	}
	if len(nodes) < 2 {
		for _, n := range nodes {
			if n.Children != nil {
				sortMenuVONodes(n.Children)
			}
		}
		return
	}
	sort.SliceStable(nodes, func(i, j int) bool {
		if nodes[i].SortOrder != nodes[j].SortOrder {
			return nodes[i].SortOrder < nodes[j].SortOrder
		}
		return nodes[i].MenuCode < nodes[j].MenuCode
	})
	for _, n := range nodes {
		if n.Children != nil {
			sortMenuVONodes(n.Children)
		}
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

	sortMenuVONodes(roots)
	return roots
}

