package funcs

type ListPermissionType int

const LPTRead ListPermissionType = 4
const LPTEdit ListPermissionType = 2
const LPTDelete ListPermissionType = 1

type ListPermissionStatus int

// 设置读权限
func (l *ListPermissionStatus) SetReadPermission() {
	*l = ListPermissionStatus(LPTRead)
}

// 设置编辑权限
func (l *ListPermissionStatus) SetEditPermission() {
	*l = ListPermissionStatus(LPTEdit)
}

// 设置删除权限
func (l *ListPermissionStatus) SetDeletePermission() {
	*l = ListPermissionStatus(LPTDelete)
}

// 设置读、编辑权限
func (l *ListPermissionStatus) SetREPermission() {
	*l = ListPermissionStatus(LPTRead) + ListPermissionStatus(LPTEdit)
}

// 设置读、删除权限
func (l *ListPermissionStatus) SetRDPermission() {
	*l = ListPermissionStatus(LPTRead) + ListPermissionStatus(LPTDelete)
}

// 设置编辑、删除权限
func (l *ListPermissionStatus) SetEDPermission() {
	*l = ListPermissionStatus(LPTEdit) + ListPermissionStatus(LPTDelete)
}

// 设置读、编辑、删除权限
func (l *ListPermissionStatus) SetREDPermission() {
	*l = ListPermissionStatus(LPTRead) + ListPermissionStatus(LPTEdit) + ListPermissionStatus(LPTDelete)
}

// 有读权限
func (l ListPermissionStatus) HaveRead() bool {
	status := int(l)
	if status == int(LPTRead) || status == int(LPTRead+LPTEdit) ||
		status == int(LPTRead+LPTDelete) || status == int(LPTRead+LPTEdit+LPTDelete) {
		return true
	}

	return false
}

// 有编辑权限
func (l ListPermissionStatus) HaveEdit() bool {
	status := int(l)
	if status == int(LPTEdit) || status == int(LPTRead+LPTEdit) ||
		status == int(LPTEdit+LPTDelete) || status == int(LPTRead+LPTEdit+LPTDelete) {
		return true
	}

	return false
}

// 有删除权限
func (l ListPermissionStatus) HaveDelete() bool {
	status := int(l)
	if status == int(LPTDelete) || status == int(LPTRead+LPTDelete) ||
		status == int(LPTEdit+LPTDelete) || status == int(LPTRead+LPTEdit+LPTDelete) {
		return true
	}

	return false
}

// 有读、编辑权限
func (l ListPermissionStatus) HaveRE() bool {
	if l.HaveRead() && l.HaveEdit() {
		return true
	}
	return false
}

// 有读、删除权限
func (l ListPermissionStatus) HaveRD() bool {
	if l.HaveRead() && l.HaveDelete() {
		return true
	}
	return false
}

// 有编辑、删除权限
func (l ListPermissionStatus) HaveED() bool {
	if l.HaveEdit() && l.HaveDelete() {
		return true
	}
	return false
}

// 有读、编辑、删除权限
func (l ListPermissionStatus) HaveRED() bool {
	if l.HaveEdit() && l.HaveEdit() && l.HaveDelete() {
		return true
	}
	return false
}

type ListPermission struct {
	ID               int
	ListID           string
	UserID           string
	PermissionStatus ListPermissionStatus
}
