package base

type ModelWorkspace struct {
	WorkspaceID uint `json:"workspace_id" csv:"workspace_id"` // 0 means global, 1 means default workspace.
}

func (m *ModelWorkspace) GetWorkspaceID() uint {
	return m.WorkspaceID
}

func (m *ModelWorkspace) SetWorkspaceID(workspaceID uint) {
	m.WorkspaceID = workspaceID
}
