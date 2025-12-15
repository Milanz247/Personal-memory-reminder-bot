package command

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// HelpCommand handles the /help command
type HelpCommand struct{}

// NewHelpCommand creates a new help command
func NewHelpCommand() *HelpCommand {
	return &HelpCommand{}
}

// Name returns the command name
func (c *HelpCommand) Name() string {
	return "help"
}

// Description returns the command description
func (c *HelpCommand) Description() string {
	return "Show help"
}

// Execute executes the help command
func (c *HelpCommand) Execute(ctx context.Context, bot BotAPI, message *tgbotapi.Message) error {
	helpText := `📚 *Memory Bot - Complete Command Guide*
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

*💾 SAVING MEMORIES:*
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

*Method 1:* ` + "`/save [text]`" + ` - Quick save
*Method 2:* ` + "`/save`" + ` - Interactive template

*🔑 Recommended Template:*
` + "`[What I did] [How I felt] #tag1 #tag2`" + `

*📚 Examples:*
• ` + "`/save Amazing project breakthrough! Felt excited and proud. #work #achievement`" + `
• ` + "`/save Had wonderful conversation with family today. #personal #happy`" + `
• ` + "`/save Completed database migration successfully. #tech #project`" + `

*🧠 What Gets Analyzed:*
• Emotional words → Weight (0-100%)
• Time & day → Context encoding
• Tags → Organization & search
• Priority → Review scheduling

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
*🔍 SEARCHING MEMORIES:*
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

*Basic:* ` + "`/search keyword`" + `
*Tags:* ` + "`/search #work`" + `
*Multiple:* ` + "`/search project meeting`" + `
*Context:* ` + "`/search Monday`" + ` or ` + "`/search morning`" + `

*🎯 Smart Features:*
• Wildcard matching (` + "`tele*`" + ` finds telegram, telephone)
• Context detection (Monday, morning, etc.)
• Tag filtering (` + "`#work`, `#health`" + `)
• Relevance ranking (BM25 algorithm)

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
*📋 OTHER COMMANDS:*
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

` + "`/recent`" + ` - View latest 10 memories
` + "`/stats`" + ` - Memory statistics & insights
` + "`/start`" + ` - Welcome & feature overview
` + "`/help`" + ` - This guide

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
*💡 PRO TIPS:*
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

✅ *Use emotional words* → Better Amygdala tagging
✅ *Add #hashtags* → Easy organization
✅ *Be specific* → More context = Better recall
✅ *Include feelings* → Emotions strengthen memory
✅ *Regular reviews* → Spaced repetition works!

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
*🧠 BIOLOGICAL FEATURES:*
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

🔬 Amygdala - Emotional tagging (0-100%)
🧩 Hippocampus - Context encoding
💤 Sleep Consolidation - Priority boost
🔄 LTP - Spaced repetition (1,3,7,14,30 days)
📉 Forgetting Curve - Smart scheduling`

	msg := tgbotapi.NewMessage(message.Chat.ID, helpText)
	msg.ParseMode = "Markdown"

	// Add action buttons
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📝 Save Memory", "cmd_save"),
			tgbotapi.NewInlineKeyboardButtonData("🔍 Search", "cmd_search"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📚 Recent", "cmd_recent"),
			tgbotapi.NewInlineKeyboardButtonData("📊 Stats", "cmd_stats"),
		),
	)
	msg.ReplyMarkup = keyboard

	_, err := bot.Send(msg)
	return err
}
