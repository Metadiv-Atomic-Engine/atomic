package atomic

type IJwt interface {
	GetUUID() string
	SetUUID(uuid string)

	GetUserID() uint
	SetUserID(userID uint)

	GetWorkspaceID() uint
	SetWorkspaceID(workspaceID uint)

	GetIP() string
	SetIP(ip string)

	GetAgent() string
	SetAgent(agent string)
}

/*
This is the basic structure of a JWT.
*/
type Jwt struct {
	UUID        string `json:"uuid"`
	UserID      uint   `json:"user_id"`
	WorkspaceID uint   `json:"workspace_id"`
	IP          string `json:"ip"`
	Agent       string `json:"agent"`
}

func (j *Jwt) GetUUID() string {
	return j.UUID
}

func (j *Jwt) SetUUID(uuid string) {
	j.UUID = uuid
}

func (j *Jwt) GetUserID() uint {
	return j.UserID
}

func (j *Jwt) SetUserID(userID uint) {
	j.UserID = userID
}

func (j *Jwt) GetWorkspaceID() uint {
	return j.WorkspaceID
}

func (j *Jwt) SetWorkspaceID(workspaceID uint) {
	j.WorkspaceID = workspaceID
}

func (j *Jwt) GetIP() string {
	return j.IP
}

func (j *Jwt) SetIP(ip string) {
	j.IP = ip
}

func (j *Jwt) GetAgent() string {
	return j.Agent
}

func (j *Jwt) SetAgent(agent string) {
	j.Agent = agent
}
