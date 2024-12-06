package bot

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/mingyueyu/myeasygo/bot/botTool"
)

func RefreshSetting(data []byte) {
	err := json.Unmarshal([]byte(data), &botTool.Bot)
	if err != nil {
		fmt.Println("更新bot配置失败:", err)
	}else {
		fmt.Println("更新bot配置成功")
	}
}

func Bot(content string) error {
	return BotWithName("", content)
}

func BotAlert(content string) error {
	return BotAlertWithName("", content)
}

func BotErr(code int64, errString string) error {
	return BotErrWithName("", code, errString)
}

func BotWithName(name string, content string) error {
	if len(botTool.Bot.Bots) == 0 {
		return errors.New("bot is not initialized")
	}
	return botTool.Note("", content, botTool.Default)
}

func BotAlertWithName(name string, content string) error {
	if len(botTool.Bot.Bots) == 0 {
		return errors.New("bot is not initialized")
	}
	return botTool.Note("", content, botTool.Alert)
}

func BotErrWithName(name string, code int64, errString string) error {
	if len(botTool.Bot.Bots) == 0 {
		return errors.New("bot is not initialized")
	}
	return botTool.Note("", fmt.Sprintf("(code %d)%s", code, errString), botTool.Error)
}


