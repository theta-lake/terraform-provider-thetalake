package thetalake

type apiSupervisionSpace struct {
	AllParticipants          bool                                  `json:"all_participants"`
	AllUsers                 bool                                  `json:"all_users"`
	CanDelete                bool                                  `json:"can_delete"`
	CanEnableAllParticipants bool                                  `json:"can_enable_all_participants"`
	CompiledParticipantList  []apiCompiledParticipant              `json:"compiled_participant_list"`
	CompiledUserList         []apiCompiledUser                     `json:"compiled_user_list"`
	CreatedAt                string                                `json:"created_at"`
	Description              string                                `json:"description"`
	DirectoryGroups          []apiSupervisionSpaceDirectoryGroup   `json:"directory_groups"`
	Disabled                 bool                                  `json:"disabled"`
	EntryStrategiesCount     int                                   `json:"entry_strategies_count"`
	ExternalId               string                                `json:"external_id"`
	HardEnforce              bool                                  `json:"hard_enforce"`
	ID                       int                                   `json:"id"`
	Integrations             []apiSupervisionSpaceIntegration      `json:"integrations"`
	MediaTypes               []apiSupervisionSpaceMediaType        `json:"media_types"`
	Name                     string                                `json:"name"`
	ParticipantCount         int                                   `json:"participant_count"`
	Participants             []apiSupervisionSpaceParticipant      `json:"participants"`
	RetentionLibraries       []apiSupervisionSpaceRetentionLibrary `json:"retention_libraries"`
	ReviewerCount            int                                   `json:"reviewer_count"`
	SupervisionSpacePriority int                                   `json:"supervision_space_priority"`
	UpdatedAt                string                                `json:"updated_at"`
	UserGroups               []apiSupervisionSpaceUserGroup        `json:"user_groups"`
	Users                    []apiSupervisionSpaceUser             `json:"users"`
}

type apiCompiledParticipant struct {
}

type apiCompiledUser struct {
}

type apiSupervisionSpaceDirectoryGroup struct {
	ID int `json:"id"`
}

type apiSupervisionSpaceIntegration struct {
	ID int `json:"id"`
}

type apiSupervisionSpaceMediaType struct {
	ID int `json:"id"`
}

type apiSupervisionSpaceParticipant struct {
}

type apiSupervisionSpaceRetentionLibrary struct {
	ID int `json:"id"`
}

type apiSupervisionSpaceUser struct {
}
type apiSupervisionSpaceUserGroup struct {
}

type apiCreateSupervisionSpaceRequest struct {
	AllParticipants          bool   `json:"all_participants"`
	AllUsers                 bool   `json:"all_users"`
	Description              string `json:"description"`
	DirectoryGroupIds        []int  `json:"directory_group_ids"`
	ExternalId               string `json:"external_id"`
	HardEnforce              bool   `json:"hard_enforce"`
	IntegrationIds           []int  `json:"integration_ids"`
	MediaTypeIds             []int  `json:"media_type_ids"`
	Name                     string `json:"name"`
	RetentionLibraryIds      []int  `json:"retention_library_ids"`
	SupervisionSpacePriority int    `json:"supervision_space_priority"`
}
