// Package infrastructure 基础设施层：邮件服务
package infrastructure

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"time"

	"github.com/Company-Automation-1/video-backend-go/src/config"
	"github.com/Company-Automation-1/video-backend-go/src/tools"
	"gopkg.in/gomail.v2"
)

// Email 邮件服务
type Email struct {
	cfg    *config.EmailConfig
	dialer *gomail.Dialer
}

// NewEmail 创建邮件服务
func NewEmail(cfg *config.Config) *Email {
	dialer := gomail.NewDialer(cfg.Email.Host, cfg.Email.Port, cfg.Email.User, cfg.Email.Password)
	dialer.SSL = cfg.Email.SSL

	return &Email{
		cfg:    &cfg.Email,
		dialer: dialer,
	}
}

// SendCaptchaEmail 发送验证码邮件
func (e *Email) SendCaptchaEmail(to, code string) error {
	if e == nil || e.dialer == nil {
		tools.Logf("邮件服务未配置")
		return errors.New("邮件服务未配置")
	}

	templatePath := filepath.Join(".", "email.html")
	tmplContent, err := os.ReadFile(templatePath) //nolint:gosec // 模板文件路径固定
	if err != nil {
		tools.Logf("读取邮件模板失败: %v", err)
		return errors.New("读取邮件模板失败")
	}

	tmpl, err := template.New("email").Parse(string(tmplContent))
	if err != nil {
		tools.Logf("解析邮件模板失败: %v", err)
		return errors.New("解析邮件模板失败")
	}

	var body bytes.Buffer
	data := map[string]string{
		"CODE": code,
	}
	if err := tmpl.Execute(&body, data); err != nil {
		tools.Logf("渲染邮件模板失败: %v", err)
		return errors.New("渲染邮件模板失败")
	}

	msg := gomail.NewMessage()
	msg.SetHeader("From", e.cfg.From)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", "邮箱验证码")
	msg.SetBody("text/html", body.String())

	// 发送邮件（使用 context 控制超时）
	timeout := e.cfg.Timeout

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	done := make(chan error, 1)
	go func() {
		done <- e.dialer.DialAndSend(msg)
	}()

	select {
	case err := <-done:
		if err != nil {
			tools.Logf("发送邮件失败: %v", err)
			return errors.New("发送邮件失败")
		}
		return nil
	case <-ctx.Done():
		tools.Logf("发送邮件超时 (%d秒)", timeout)
		return fmt.Errorf("发送邮件超时 (%d秒)", timeout)
	}
}
