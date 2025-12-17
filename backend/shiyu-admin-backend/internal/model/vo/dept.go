package vo

import "shiyu-admin-backend/internal/model/entity"

// DeptVO represents department info for API responses.
type DeptVO struct {
	DeptCode   string    `json:"dept_code"`
	ParentCode string    `json:"parent_code"`
	DeptName   string    `json:"dept_name"`
	Status     int       `json:"status"`
	Children   []*DeptVO `json:"children,omitempty"`
}

// BuildDeptVO converts entity to VO.
func BuildDeptVO(d *entity.Dept) *DeptVO {
	if d == nil {
		return nil
	}
	return &DeptVO{
		DeptCode:   d.DeptCode,
		ParentCode: d.ParentCode,
		DeptName:   d.DeptName,
		Status:     d.Status,
	}
}

// BuildDeptTree builds a tree structure from flat dept list.
func BuildDeptTree(depts []*entity.Dept) []*DeptVO {
	deptMap := make(map[string]*DeptVO)
	var roots []*DeptVO

	// First pass: create all dept VOs
	for _, d := range depts {
		vo := BuildDeptVO(d)
		deptMap[d.DeptCode] = vo
	}

	// Second pass: build tree
	for _, d := range depts {
		vo := deptMap[d.DeptCode]
		if d.ParentCode == "" {
			roots = append(roots, vo)
		} else {
			parent := deptMap[d.ParentCode]
			if parent != nil {
				if parent.Children == nil {
					parent.Children = []*DeptVO{}
				}
				parent.Children = append(parent.Children, vo)
			}
		}
	}

	return roots
}

