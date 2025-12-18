// Package tools 查询构建器
package tools

import (
	"gorm.io/gen"
	"gorm.io/gen/field"
)

// ConditionBuilder 条件构建器（流式API）
// 使用示例：
//
//	conditions := NewConditionBuilder().
//	    EqUint(query.User.ID, req.ID).
//	    EqString(query.User.Username, req.Username).
//	    Like(query.User.Email, req.EmailLike).
//	    GteInt64(query.User.CreatedAt, req.CreatedAtMin).
//	    LteInt64(query.User.CreatedAt, req.CreatedAtMax).
//	    Build()
type ConditionBuilder struct {
	conditions []gen.Condition
}

// NewConditionBuilder 创建条件构建器
func NewConditionBuilder() *ConditionBuilder {
	return &ConditionBuilder{conditions: make([]gen.Condition, 0, 8)}
}

// Build 构建并返回条件数组
func (b *ConditionBuilder) Build() []gen.Condition {
	return b.conditions
}

// --- 字符串条件 ---

// EqString 字符串精确匹配
func (b *ConditionBuilder) EqString(f *field.String, value string) *ConditionBuilder {
	if value != "" {
		b.conditions = append(b.conditions, f.Eq(value))
	}
	return b
}

// Like 字符串模糊匹配（LIKE %value%）
func (b *ConditionBuilder) Like(f *field.String, value string) *ConditionBuilder {
	if value != "" {
		b.conditions = append(b.conditions, f.Like("%"+value+"%"))
	}
	return b
}

// InStrings 字符串 IN 查询
func (b *ConditionBuilder) InStrings(f *field.String, values []string) *ConditionBuilder {
	if len(values) > 0 {
		b.conditions = append(b.conditions, f.In(values...))
	}
	return b
}

// --- 整数条件 ---

// EqUint Uint 精确匹配
func (b *ConditionBuilder) EqUint(f *field.Uint, value *uint) *ConditionBuilder {
	if value != nil {
		b.conditions = append(b.conditions, f.Eq(*value))
	}
	return b
}

// EqInt Int 精确匹配
func (b *ConditionBuilder) EqInt(f *field.Int, value int) *ConditionBuilder {
	if value != 0 {
		b.conditions = append(b.conditions, f.Eq(value))
	}
	return b
}

// GteInt Int 大于等于
func (b *ConditionBuilder) GteInt(f *field.Int, value *int) *ConditionBuilder {
	if value != nil {
		b.conditions = append(b.conditions, f.Gte(*value))
	}
	return b
}

// LteInt Int 小于等于
func (b *ConditionBuilder) LteInt(f *field.Int, value *int) *ConditionBuilder {
	if value != nil {
		b.conditions = append(b.conditions, f.Lte(*value))
	}
	return b
}

// InInts Int IN 查询
func (b *ConditionBuilder) InInts(f *field.Int, values []int) *ConditionBuilder {
	if len(values) > 0 {
		b.conditions = append(b.conditions, f.In(values...))
	}
	return b
}

// BetweenInt Int 范围查询（BETWEEN minVal AND maxVal）
func (b *ConditionBuilder) BetweenInt(f *field.Int, minVal, maxVal int) *ConditionBuilder {
	if minVal != 0 && maxVal != 0 {
		b.conditions = append(b.conditions, f.Between(minVal, maxVal))
	}
	return b
}

// --- Int64条件 ---

// EqInt64 Int64 精确匹配
func (b *ConditionBuilder) EqInt64(f *field.Int64, value int64) *ConditionBuilder {
	if value != 0 {
		b.conditions = append(b.conditions, f.Eq(value))
	}
	return b
}

// GteInt64 Int64 大于等于
func (b *ConditionBuilder) GteInt64(f *field.Int64, value *int64) *ConditionBuilder {
	if value != nil {
		b.conditions = append(b.conditions, f.Gte(*value))
	}
	return b
}

// LteInt64 Int64 小于等于
func (b *ConditionBuilder) LteInt64(f *field.Int64, value *int64) *ConditionBuilder {
	if value != nil {
		b.conditions = append(b.conditions, f.Lte(*value))
	}
	return b
}

// --- 布尔条件 ---

// EqBool Bool 精确匹配
func (b *ConditionBuilder) EqBool(f *field.Bool, value *bool) *ConditionBuilder {
	if value != nil {
		b.conditions = append(b.conditions, f.Is(*value))
	}
	return b
}

// --- 条件版本方法（带条件判断）---

// EqStringIf 条件字符串精确匹配
func (b *ConditionBuilder) EqStringIf(f *field.String, value string, condition bool) *ConditionBuilder {
	if condition && value != "" {
		b.conditions = append(b.conditions, f.Eq(value))
	}
	return b
}

// LikeIf 条件字符串模糊匹配
func (b *ConditionBuilder) LikeIf(f *field.String, value string, condition bool) *ConditionBuilder {
	if condition && value != "" {
		b.conditions = append(b.conditions, f.Like("%"+value+"%"))
	}
	return b
}

// GteIntIf 条件 Int 大于等于
func (b *ConditionBuilder) GteIntIf(f *field.Int, value *int, condition bool) *ConditionBuilder {
	if condition && value != nil {
		b.conditions = append(b.conditions, f.Gte(*value))
	}
	return b
}

// --- 空值条件 ---

// IsNull 空值判断（isNull=true 时查询 NULL，false 时查询 NOT NULL）
func (b *ConditionBuilder) IsNull(f *field.Field, isNull bool) *ConditionBuilder {
	if isNull {
		b.conditions = append(b.conditions, f.IsNull())
	} else {
		b.conditions = append(b.conditions, f.IsNotNull())
	}
	return b
}

// --- 逻辑组合 ---

// And 添加 AND 条件（默认就是 AND，直接追加）
func (b *ConditionBuilder) And(conds ...gen.Condition) *ConditionBuilder {
	b.conditions = append(b.conditions, conds...)
	return b
}
