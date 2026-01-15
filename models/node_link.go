package models

// NodeLink 节点关联（仅用于阶段递进关联）
type NodeLink struct {
	BaseModel
	// 源节点ID（前一阶段）
	SourceNodeID uint `gorm:"index;not null" json:"source_node_id"`
	// 目标节点ID（后一阶段）
	TargetNodeID uint `gorm:"index;not null" json:"target_node_id"`
	// 关联类型：stage_progression-阶段递进
	LinkType string `gorm:"size:50;not null;default:'stage_progression'" json:"link_type"`
	// 创建人ID
	CreatorID uint `gorm:"index;not null" json:"creator_id"`

	// 关联
	SourceNode *PlanNode `gorm:"foreignKey:SourceNodeID" json:"source_node,omitempty"`
	TargetNode *PlanNode `gorm:"foreignKey:TargetNodeID" json:"target_node,omitempty"`
	Creator    *User     `gorm:"foreignKey:CreatorID" json:"creator,omitempty"`
}

// TableName 指定表名
func (NodeLink) TableName() string {
	return "node_links"
}

// 节点关联类型常量
const (
	NodeLinkTypeStageProgression = "stage_progression" // 阶段递进
)
