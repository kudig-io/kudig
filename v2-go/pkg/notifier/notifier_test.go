package notifier

import (
	"testing"

	"github.com/kudig/kudig/pkg/types"
)

func TestNewConfigFromEnv(t *testing.T) {
	// Save and restore environment variables
	oldSlackURL := ""
	oldMinSev := ""
	
	// Set test environment variables
	t.Setenv("KUDIG_SLACK_WEBHOOK_URL", "https://hooks.slack.com/test")
	t.Setenv("KUDIG_NOTIFY_MIN_SEVERITY", "warning")

	config := NewConfigFromEnv()

	if config.SlackWebhookURL != "https://hooks.slack.com/test" {
		t.Errorf("Expected Slack URL to be set, got '%s'", config.SlackWebhookURL)
	}

	if config.MinSeverity != types.SeverityWarning {
		t.Errorf("Expected min severity warning, got %v", config.MinSeverity)
	}

	// Reset
	if oldSlackURL != "" {
		t.Setenv("KUDIG_SLACK_WEBHOOK_URL", oldSlackURL)
	} else {
		t.Setenv("KUDIG_SLACK_WEBHOOK_URL", "")
	}
	if oldMinSev != "" {
		t.Setenv("KUDIG_NOTIFY_MIN_SEVERITY", oldMinSev)
	} else {
		t.Setenv("KUDIG_NOTIFY_MIN_SEVERITY", "")
	}
}

func TestConfig_IsEnabled(t *testing.T) {
	tests := []struct {
		name     string
		config   Config
		expected bool
	}{
		{
			name:     "Empty config",
			config:   Config{},
			expected: false,
		},
		{
			name: "Slack enabled",
			config: Config{
				SlackWebhookURL: "https://hooks.slack.com/test",
			},
			expected: true,
		},
		{
			name: "DingTalk enabled",
			config: Config{
				DingTalkWebhookURL: "https://oapi.dingtalk.com/test",
			},
			expected: true,
		},
		{
			name: "WeChat enabled",
			config: Config{
				WeChatWebhookURL: "https://qyapi.weixin.qq.com/test",
			},
			expected: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.config.IsEnabled() != tc.expected {
				t.Errorf("Expected IsEnabled() = %v", tc.expected)
			}
		})
	}
}

func TestConfig_ShouldNotify(t *testing.T) {
	config := &Config{
		SlackWebhookURL: "https://hooks.slack.com/test",
		MinSeverity:     types.SeverityWarning,
	}

	tests := []struct {
		name     string
		issues   []types.Issue
		expected bool
	}{
		{
			name:     "No issues",
			issues:   []types.Issue{},
			expected: false,
		},
		{
			name: "Only info issues",
			issues: []types.Issue{
				{Severity: types.SeverityInfo},
			},
			expected: false,
		},
		{
			name: "Has warning issue",
			issues: []types.Issue{
				{Severity: types.SeverityWarning},
			},
			expected: true,
		},
		{
			name: "Has critical issue",
			issues: []types.Issue{
				{Severity: types.SeverityCritical},
			},
			expected: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if config.ShouldNotify(tc.issues) != tc.expected {
				t.Errorf("Expected ShouldNotify() = %v for %s", tc.expected, tc.name)
			}
		})
	}
}

func TestNewMultiNotifier(t *testing.T) {
	config := &Config{
		SlackWebhookURL:    "https://hooks.slack.com/test",
		DingTalkWebhookURL: "https://oapi.dingtalk.com/test",
		WeChatWebhookURL:   "https://qyapi.weixin.qq.com/test",
	}

	mn := NewMultiNotifier(config)

	if len(mn.Notifiers) != 3 {
		t.Errorf("Expected 3 notifiers, got %d", len(mn.Notifiers))
	}
}

func TestSlackNotifier_Name(t *testing.T) {
	n := NewSlackNotifier("https://hooks.slack.com/test")
	if n.Name() != "Slack" {
		t.Errorf("Expected name 'Slack', got '%s'", n.Name())
	}
}

func TestDingTalkNotifier_Name(t *testing.T) {
	n := NewDingTalkNotifier("https://oapi.dingtalk.com/test")
	if n.Name() != "DingTalk" {
		t.Errorf("Expected name 'DingTalk', got '%s'", n.Name())
	}
}

func TestWeChatNotifier_Name(t *testing.T) {
	n := NewWeChatNotifier("https://qyapi.weixin.qq.com/test")
	if n.Name() != "WeChat Work" {
		t.Errorf("Expected name 'WeChat Work', got '%s'", n.Name())
	}
}

func TestMultiNotifier_Send_NoNotifiers(t *testing.T) {
	mn := &MultiNotifier{}
	errors := mn.Send("Test", "Message", nil)

	if len(errors) != 0 {
		t.Errorf("Expected 0 errors with no notifiers, got %d", len(errors))
	}
}
