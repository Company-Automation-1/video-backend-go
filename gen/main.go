// Package main GORM Gen 代码生成器入口
package main

import (
	"github.com/Company-Automation-1/video-backend-go/src/models"
	"gorm.io/gen"
)

func main() {
	// 创建生成器（基于模型生成）
	g := gen.NewGenerator(gen.Config{
		OutPath:        "./src/query", // 生成代码的输出目录
		OutFile:        "gen.go",      // 所有模型共享一个主文件
		Mode:           gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface,
		FieldNullable:  true, // 字段可为空
		FieldCoverable: true, // 字段可覆盖
	})

	// 定义需要生成查询代码的模型列表
	// 后续添加新模型时，只需在此处添加即可
	modelsToGenerate := []interface{}{
		models.User{},
		// 后续添加新模型示例：
		// models.Article{},
		// models.Comment{},
	}

	// 为所有模型生成基础 DAO API（所有模型共享 Query 结构）
	for _, model := range modelsToGenerate {
		g.ApplyBasic(model)
	}

	// 执行生成（一次生成所有模型）
	g.Execute()
}
