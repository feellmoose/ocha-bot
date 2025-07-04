package helper

import (
	"bytes"
	"log"
	"text/template"
)

type Template struct {
	inner *template.Template
}

func (t Template) String() string {
	if t.inner == nil || t.inner.Tree == nil || t.inner.Tree.Root == nil {
		return "<nil template>"
	}
	return t.inner.Tree.Root.String()
}

func (t Template) Execute(data any) (string, error) {
	var buf bytes.Buffer
	err := t.inner.Execute(&buf, data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

var templates = map[string]map[string]string{
	"en": {
		"stat.all.note":                   "Stat report:\n<blockquote expandable>Bot:\nID: {{.BotID}}\nName: {{.BotName}}\nVersion: {{.Version}}\nUpdate: {{.Update}}\n\nRepos:\nsize: {{.RepoSize}}\n{{.Repos}}\n{{.Mine}}\n\nAnalysis:\nTime: {{.Now}}</blockquote>",
		"stat.repo.note":                  "Repo: {{ .Name }}\n\t| type: {{ .Type }}\n\t| size: {{ .DataSize }}\n\t| objs: {{ .ObjsSize }}\n",
		"stat.game.mine.note":             "Mine-sweeper-game:\n\t| running: {{.Running}}\n\t| active: {{.Active}}\n\t| total: {{.Total}}",
		"lang.note":                       "@{{ .Username }}\nLanguage updated successfully",
		"lang.chat.note":                  "@{{ .Username }}\nThe default language for chat group {{ .ChatName }} has been successfully updated",
		"lang.menu.note":                  "@{{ .Username }}\nPlease click the button below to update your language setting saved in ocha. Your personal setting will take precedence over the chat group’s default language",
		"lang.chat.menu.note":             "@{{ .Username }}\nAdmins, please click the button below to update the default language setting for chat group {{ .ChatName }} saved in ocha. Personal language settings will take precedence over the group’s default setting",
		"lang.zh.button":                  "简体中文",
		"lang.en.button":                  "English",
		"lang.cxg.button":                 "nya大人",
		"mine.game.quit.note":             "@{{ .Username }}\nQuit game success",
		"menu.back.button":                "Back",
		"menu.cancel.button":              "Cancel",
		"mine.game.menu.note":             "@{{ .Username }}\nWelcome to the entertainment service provided by ocha. You can start a Minesweeper game using this menu.\nPlease click the button below to select a difficulty level.",
		"mine.game.menu.easy.button":      "Easy",
		"mine.game.menu.normal.button":    "Normal",
		"mine.game.menu.hard.button":      "Hard",
		"mine.game.menu.nightmare.button": "Nightmare",
		"mine.game.menu.random.button":    "Random Map",
		"mine.game.menu.rank.button":      "Leaderboard",
		"mine.game.menu.classic.button":   "Classic",
		"mine.game.rank.start.note":       "@{{ .Username }}\nWelcome to the entertainment service provided by ocha.  If you successfully complete this Minesweeper challenge, your result will be added to the leaderboard. You have started a new {{ .Width }} × {{ .Height }} Minesweeper map with {{ .Mines }} mines in total.",
		"mine.game.rank.win.note":         "@{{ .Username }}\nCongratulations! 🎉\nYou successfully completed the game in {{ .Seconds }} seconds.\nLeaderboard score\\rank: {{ .Score }}\\{{ .Rank }}\nMap size: {{ .Width }} × {{ .Height }}\nMine count: {{ .Mines }}\nUse this command to view the full leaderboard:\n/mine_rank@{{ .BotName }}",
		"mine.game.rank.lose.note":        "@{{ .Username }}\nBoom! 💣\nUnfortunately, this run did not qualify for the leaderboard.\nTime taken: {{ .Seconds }} seconds.\nMap size: {{ .Width }} × {{ .Height }}\nMine count: {{ .Mines }}\nUse this command to view the full leaderboard:\n/mine_rank@{{ .BotName }}",
		"mine.game.rank.res.note":         "@{{ .Username }}\nHere is the current Minesweeper leaderboard:\n<blockquote expandable>{{.RankLines}}</blockquote>\nLast updated: {{.Update}}",
		"mine.game.rank.line.note":        "Rank: {{.Index}}\n\t|User: {{.Username}}\n\t|Map size: {{ .Width }} × {{ .Height }}\n\t|Mines: {{ .Mines }}\n\t|Steps: {{ .Steps }}\n\t|Duration: {{.Duration}}\n\t|Score: {{.Score}}\n\n",
		"mine.game.start.note":            "@{{ .Username }}\nWelcome to the entertainment service provided by ocha. You have started a new {{ .Width }} × {{ .Height }} Minesweeper map.\nThere are {{ .Mines }} mines in total.",
		"mine.game.start.button":          "Click to Start",
		"mine.game.win.note":              "@{{ .Username }}\nCongratulations! 🎉\nYou successfully completed the game in {{ .Seconds }} seconds.\nMap size: {{ .Width }} × {{ .Height }}\nNumber of mines: {{ .Mines }}",
		"mine.game.win.button":            "Play Again",
		"mine.game.lose.note":             "@{{ .Username }}\nBoom! 💣\nTime taken: {{ .Seconds }} seconds.\nMap size: {{ .Width }} × {{ .Height }}\nNumber of mines: {{ .Mines }}",
		"mine.game.lose.button":           "Try Again",
		"mine.game.opt.quit":              "Exit",
		"mine.game.opt.flag":              "Flag",
		"mine.game.opt.click":             "Sweep",
		"cron.help.note":                  "@{{ .Username }}\nWelcome to the cron message service provided by ocha. \nYou can schedule message tasks here through Cron expressions to implement the function of sending messages at a scheduled time: <blockquote expandable>/cron * * * * * '<message>'\n- - - - -\n| | | | |\n| | | | +----- day of the week (0-6)\n\n| | | +------- month (1-12)\n| | +--------- day of the month (1-31)\n| +----------- hour (0-23)\n+------------- minute (0-59)\n\nE.g.\n/cron 0 * * * * 'hello'</blockquote>",
		"cron.list.note":                  "@{{ .Username }}\nHere is current tasks:\n<blockquote expandable>{{.TaskLines}}</blockquote>\nLast updated: {{.Update}}",
		"error":                           "@{{ .Username }} Oops! Something went wrong! {{.Message}}!",
		"help.note":                       "@{{ .Username }}\nWelcome to ocha!\nHere are some commands to help you get started:\n/mine\n/mine  &lt;width&gt; &lt;height&gt; &lt;mines&gt;\n/lang  [ zh | en | cxg ]\n/lang_chat  [ zh | en | cxg ]\n/help\n<blockquote expandable>{{.BotName}}\nAuthor: @feellmoose_dev\nVersion: {{.Version}}\nUpdated on: {{.Update}}\n</blockquote>",
	},
	"zh": {
		"stat.all.note":                   "Stat report:\n<blockquote expandable>Bot:\nID: {{.BotID}}\nName: {{.BotName}}\nVersion: {{.Version}}\nUpdate: {{.Update}}\n\nRepos:\nsize: {{.RepoSize}}\n{{.Repos}}\n{{.Mine}}\n\nAnalysis:\nTime: {{.Now}}</blockquote>",
		"stat.repo.note":                  "Repo: {{ .Name }}\n\t| type: {{ .Type }}\n\t| size: {{ .DataSize }}\n\t| objs: {{ .ObjsSize }}\n",
		"stat.game.mine.note":             "Mine-sweeper-game:\n\t| running: {{.Running}}\n\t| active: {{.Active}}\n\t| total: {{.Total}}",
		"lang.note":                       "@{{ .Username }}\n语言修改成功",
		"lang.chat.note":                  "@{{ .Username }}\n本聊天群组 {{ .ChatName }} 的默认语言修改成功",
		"lang.menu.note":                  "@{{ .Username }}\n请点击下方按钮修改您在 ocha 留存的语言设置，个人语言设置将优先于聊天群组的默认语言设置显示",
		"lang.chat.menu.note":             "@{{ .Username }}\n请管理员点击下方按钮修改本聊天群组 {{ .ChatName }} 在 ocha 留存的默认语言设置，个人语言设置将优先于聊天群组的默认语言设置显示",
		"lang.zh.button":                  "简体中文",
		"lang.en.button":                  "English",
		"lang.cxg.button":                 "nya大人",
		"mine.game.quit.note":             "@{{ .Username }}\n成功退出游戏",
		"menu.back.button":                "返回",
		"menu.cancel.button":              "取消",
		"mine.game.menu.note":             "@{{ .Username }}\n欢迎使用 ocha 为您提供的娱乐服务，您可以通过此菜单开始一个扫雷游戏。\n请点击下面的按钮选择难度",
		"mine.game.menu.easy.button":      "简单",
		"mine.game.menu.normal.button":    "普通",
		"mine.game.menu.hard.button":      "困难",
		"mine.game.menu.nightmare.button": "噩梦模式",
		"mine.game.menu.random.button":    "随机地图",
		"mine.game.menu.rank.button":      "天梯赛",
		"mine.game.menu.classic.button":   "经典模式",
		"mine.game.rank.start.note":       "@{{ .Username }}\n欢迎使用 ocha 为您提供的娱乐服务，若本次扫雷任务成功，则会被记录在天梯赛榜单内。您已开始一个新的 {{ .Width }} × {{ .Height }} 扫雷地图。\n共有 {{ .Mines }} 个地雷",
		"mine.game.rank.win.note":         "@{{ .Username }}\n恭喜！🎉\n您成功在 {{ .Seconds }} 秒内完成了游戏。\n天梯赛得分\\排位：{{ .Score }}\\{{ .Rank }}\n地图尺寸：{{ .Width }} × {{ .Height }}\n地雷数量：{{ .Mines }}\n使用指令查看详细榜单:\n/mine_rank@{{.BotName}}",
		"mine.game.rank.lose.note":        "@{{ .Username }}\n砰！💣\n很遗憾，此次记录未能加入天梯赛排位中。\n耗时：{{ .Seconds }} 秒。\n地图尺寸：{{ .Width }} × {{ .Height }}\n地雷数量：{{ .Mines }}\n使用指令查看详细榜单:\n/mine_rank@{{.BotName}}",
		"mine.game.rank.res.note":         "@{{.Username}}\n当前的扫雷天梯榜单如下：\n<blockquote expandable>{{.RankLines}}</blockquote>\n更新时间：{{.Update}}",
		"mine.game.rank.line.note":        "排行：{{.Index}}\n\t|用户：{{.Username}}\n\t|地图：{{ .Width }} × {{ .Height }}\n\t|雷数：{{ .Mines }}\n\t|步数：{{ .Steps }}\n\t|用时：{{.Duration}}\n\t|最终得分：{{.Score}}\n\n",
		"mine.game.start.note":            "@{{ .Username }}\n欢迎使用 ocha 为您提供的娱乐服务，您已开始一个新的 {{ .Width }} × {{ .Height }} 扫雷地图。\n共有 {{ .Mines }} 个地雷",
		"mine.game.start.button":          "点击开始",
		"mine.game.win.note":              "@{{ .Username }}\n恭喜！🎉\n您成功在 {{ .Seconds }} 秒内完成了游戏。\n地图尺寸：{{ .Width }} × {{ .Height }}\n地雷数量：{{ .Mines }}",
		"mine.game.win.button":            "再来一局",
		"mine.game.lose.note":             "@{{ .Username }}\n砰！💣\n耗时：{{ .Seconds }} 秒。\n地图尺寸：{{ .Width }} × {{ .Height }}\n地雷数量：{{ .Mines }}",
		"mine.game.lose.button":           "再试一次",
		"mine.game.opt.quit":              "退出",
		"mine.game.opt.flag":              "插旗",
		"mine.game.opt.click":             "扫雷",
		"cron.help.note":                  "@{{ .Username }}\n欢迎使用 ocha 为您提供的cron定时消息服务. \n以下是使用样例: <blockquote expandable>/cron * * * * * '<message>'\n- - - - -\n| | | | |\n| | | | +----- 周 (0-6)\n\n| | | +------- 月 (1-12)\n| | +--------- 日 (1-31)\n| +----------- 时 (0-23)\n+------------- 分 (0-59)\n\nE.g.\n/cron 0 * * * * 'hello'</blockquote>",
		"cron.list.note":                  "@{{ .Username }}\n活跃任务:\n<blockquote expandable>{{.TaskLines}}</blockquote>\n更新时间: {{.Update}}",
		"error":                           "@{{ .Username }} 哎呀！出了点问题！{{.Message}}！",
		"help.note":                       "@{{ .Username }}\n欢迎使用 ocha ！\n以下是一些帮助您入门的命令：\n/mine\n/mine  &lt; 宽 &gt; &lt; 高 &gt; &lt; 雷数 &gt;\n/lang  [ zh | en | cxg ]\n/lang_chat  [ zh | en | cxg ]\n/help\n<blockquote expandable>{{.BotName}}\n作者: @feellmoose_dev\n版本信息:{{.Version}}\n更新于:{{.Update}}\n</blockquote>",
	},
	"cxg": {
		"stat.all.note":                   "Stat report:\n<blockquote expandable>Bot:\nID: {{.BotID}}\nName: {{.BotName}}\nVersion: {{.Version}}\nUpdate: {{.Update}}\n\nRepos:\nsize: {{.RepoSize}}\n{{.Repos}}\n{{.Mine}}\n\nAnalysis:\nTime: {{.Now}}</blockquote>",
		"stat.repo.note":                  "Repo: {{ .Name }}\n\t| type: {{ .Type }}\n\t| size: {{ .DataSize }}\n\t| objs: {{ .ObjsSize }}\n",
		"stat.game.mine.note":             "Mine-sweeper-game:\n\t| running: {{.Running}}\n\t| active: {{.Active}}\n\t| total: {{.Total}}",
		"lang.note":                       "@{{ .Username }}\n哼哼！本nya大人已经优雅地把你的语言换好啦！快感谢我吧！",
		"lang.chat.note":                  "@{{ .Username }}\n哼哼！本nya大人已经优雅地把聊天群组 {{ .ChatName }} 的默认语言换好啦！快感谢我吧！",
		"lang.menu.note":                  "@{{ .Username }}\n快点自己选一个语言记录在nya大人的小本本上哦 ~ 不要让本喵亲自动手！咱才不会承认这个语言会比群组默认的那个要重要得多呢！哼！",
		"lang.chat.menu.note":             "@{{ .Username }}\n管理员大人！快点选一个聊天群组 {{ .ChatName }} 的默认语言，然后记录在nya大人的身体上 ~ 不要让本喵求您呜呜 ~ 没有自己设置语言的杂鱼都会被强制使用这个语言呢 ~ 嗯哼 ~",
		"lang.zh.button":                  "简体中文",
		"lang.en.button":                  "English",
		"lang.cxg.button":                 "nya大人",
		"mine.game.quit.note":             "@{{ .Username }}\n有笨蛋逃跑了呢~真是杂鱼！",
		"menu.back.button":                "返回喵",
		"menu.cancel.button":              "取消喵",
		"mine.game.menu.easy.button":      "杂鱼",
		"mine.game.menu.normal.button":    "一般",
		"mine.game.menu.hard.button":      "勉强",
		"mine.game.menu.nightmare.button": "找虐喵",
		"mine.game.menu.random.button":    "随本喵心意",
		"mine.game.start.button":          "扫雷~启动！",
		"mine.game.win.button":            "再战！",
		"mine.game.lose.button":           "不服？咱还要玩！",
		"mine.game.opt.quit":              "逃跑喵",
		"mine.game.opt.flag":              "插旗旗",
		"mine.game.opt.click":             "点爆它",
		"mine.game.menu.rank.button":      "最新最热最好的！天梯赛！",
		"mine.game.menu.classic.button":   "适合老年人的经典模式",
		"mine.game.rank.start.note":       "@{{ .Username }}\n喵喵喵~你的游戏开始啦~ 只要您这次扫雷挑战完成，成绩就会被记录到天梯赛榜单上哦~ 您已踏入全新 {{ .Width }} × {{ .Height }} 扫雷地图，埋伏了 {{ .Mines }} 颗地雷",
		"mine.game.rank.win.note":         "@{{ .Username }}\n你竟然赢了喵！？哼哼~你是不是偷偷作弊了？不然怎么可能在 {{ .Seconds }} 秒就通关。\n天梯赛得分\\排位：{{ .Score }}\\{{ .Rank }}\n地图尺寸：{{ .Width }} × {{ .Height }}\n地雷数量：{{ .Mines }}\n要看详细榜单，请键入咒语:\n/mine_rank@{{ .BotName }}",
		"mine.game.rank.lose.note":        "@{{ .Username }}\n砰！💣\n好可惜，这次记录没能挤进天梯赛排位里…\n耗时：{{ .Seconds }} 秒\n地图尺寸：{{ .Width }} × {{ .Height }}\n地雷数量：{{ .Mines }}\n要看详细榜单，请键入咒语:\n/mine_rank@{{ .BotName }}",
		"mine.game.rank.res.note":         "@{{.Username}}\n哦呀！这里是扫雷天梯赛的结果看板哦:\n<blockquote expandable>{{.RankLines}}</blockquote>\n更新时间: {{.Update}}",
		"mine.game.rank.line.note":        "杂鱼排行：{{.Index}}\n\t|杂鱼：{{.Username}}\n\t|地图：{{ .Width }} × {{ .Height }}\n\t|雷数：{{ .Mines }}\n\t|步数：{{ .Steps }}\n\t|用时：{{.Duration}}\n\t|杂鱼得分：{{.Score}}\n\n",
		"mine.game.menu.note":             "@{{ .Username }}\n欢迎来到本nya大人精心布置的雷之乐园~♡\n喵呼呼~快选个难度试试看你能撑几步喵？别怕爆炸哦，本nya大人会在一旁看好戏的~♪",
		"mine.game.start.note":            "@{{ .Username }}\n喵喵喵~你的游戏开始啦~ \n尺寸：{{ .Width }} × {{ .Height }}，地雷数：{{ .Mines }} 个。\n本nya大人已经布好雷，等你来踩爆~♡",
		"mine.game.win.note":              "@{{ .Username }}\n你竟然赢了喵！？哼哼~你是不是偷偷作弊了？不然怎么可能在 {{ .Seconds }} 秒就完成地图：{{ .Width }}×{{ .Height }}，地雷数：{{ .Mines }} 个！\n本nya大人才没那么容易认输呢~下次让你哭着投降！",
		"mine.game.lose.note":             "@{{ .Username }}\n砰～💣哇咔咔~你爆炸啦~本nya大人就知道你会踩雷喵！\n时间：{{ .Seconds }} 秒，地图：{{ .Width }}×{{ .Height }}，雷数：{{ .Mines }}。\n可怜兮兮的小笨蛋，要不要本nya大人抱抱呀~？嘻嘻~",
		"cron.help.note":                  "@{{ .Username }}\n迷路的小猫咪要找帮助吗？本nya大人大发慈悲告诉你一点线索喵~\ncron是这样用的喵: <blockquote expandable>/cron * * * * * '<message>'\n- - - - -\n| | | | |\n| | | | +----- 周 (0-6)\n\n| | | +------- 月 (1-12)\n| | +--------- 日 (1-31)\n| +----------- 时 (0-23)\n+------------- 分 (0-59)\n\nE.g.\n/cron 0 * * * * 'hello'</blockquote>",
		"cron.list.note":                  "@{{ .Username }}\n目前在线的任务喵:\n<blockquote expandable>{{.TaskLines}}</blockquote>\n更新时间: {{.Update}}",
		"error":                           "@{{ .Username }} 哎呀出错了喵~ 你果然不行呢~连 {{ .Message }} 都搞不清楚~要不要本nya大人教教你啊？喵呼呼~",
		"help.note":                       "@{{ .Username }}\n迷路的小猫咪要找帮助吗？本nya大人大发慈悲告诉你一点线索喵~\n/mine\n/mine  &lt; 宽 &gt; &lt; 高 &gt; &lt; 雷数 &gt;\n/cron * * * * * '<message>'\n/lang  [ zh | en | cxg ]\n/lang_chat  [ zh | en | cxg ]\n/help\n<blockquote expandable>{{.BotName}}\n作者: @feellmoose_dev\n版本：{{.Version}}\n更新时间：{{.Update}}\n</blockquote>",
	},
}

func convert(templates map[string]map[string]string) map[string]map[string]Template {
	m := map[string]map[string]Template{}
	for lang, ts := range templates {
		msg := map[string]Template{}
		for key, text := range ts {
			t, err := template.New(lang + "_" + key).Parse(text)
			if err != nil {
				log.Panic(err)
			}
			msg[key] = Template{inner: t}
		}
		m[lang] = msg
	}
	return m
}

var Messages = convert(templates)
