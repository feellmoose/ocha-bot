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
		"menu.back.button":                "Back",
		"menu.cancel.button":              "Cancel",
		"mine.game.menu.note":             "@{{ .Username }}\nWelcome to the entertainment service provided by ocha. You can start a Minesweeper game using this menu.\nPlease click the button below to select a difficulty level.",
		"mine.game.menu.easy.button":      "Easy",
		"mine.game.menu.normal.button":    "Normal",
		"mine.game.menu.hard.button":      "Hard",
		"mine.game.menu.nightmare.button": "Nightmare",
		"mine.game.menu.random.button":    "Random Map",
		"mine.game.start.note":            "@{{ .Username }}\nWelcome to the entertainment service provided by ocha. You have started a new {{ .Width }} × {{ .Height }} Minesweeper map.\nThere are {{ .Mines }} mines in total.",
		"mine.game.start.button":          "Click to Start",
		"mine.game.win.note":              "@{{ .Username }}\nCongratulations! 🎉\nYou successfully completed the game in {{ .Seconds }} seconds.\nMap size: {{ .Width }} × {{ .Height }}\nNumber of mines: {{ .Mines }}",
		"mine.game.win.button":            "Play Again",
		"mine.game.lose.note":             "@{{ .Username }}\nBoom! 💣\nTime taken: {{ .Seconds }} seconds.\nMap size: {{ .Width }} × {{ .Height }}\nNumber of mines: {{ .Mines }}",
		"mine.game.lose.button":           "Try Again",
		"mine.game.opt.quit":              "Exit",
		"mine.game.opt.flag":              "Flag",
		"mine.game.opt.click":             "Sweep",
		"mine.game.error":                 "@{{ .Username }} Oops! Something went wrong! {{.Message}}!",
		"mine.help":                       "@{{ .Username }}\nWelcome to ocha!\nHere are some commands to help you get started:\n\n/mine <width> <height> <number of mines>\n/help\n\n<blockquote>\nocha bot\nAuthor: @feellmoose_dev\nVersion: {{.Version}}\nUpdated on: {{.Update}}\n</blockquote>",
	},
	"zh": {
		"menu.back.button":                "返回",
		"menu.cancel.button":              "取消",
		"mine.game.menu.note":             "@{{ .Username }}\n欢迎使用 ocha 为您提供的娱乐服务，您可以通过此菜单开始一个扫雷游戏。\n请点击下面的按钮选择难度",
		"mine.game.menu.easy.button":      "简单",
		"mine.game.menu.normal.button":    "普通",
		"mine.game.menu.hard.button":      "困难",
		"mine.game.menu.nightmare.button": "噩梦模式",
		"mine.game.menu.random.button":    "随机地图",
		"mine.game.start.note":            "@{{ .Username }}\n欢迎使用 ocha 为您提供的娱乐服务，您已开始一个新的 {{ .Width }} × {{ .Height }} 扫雷地图。\n共有 {{ .Mines }} 个地雷",
		"mine.game.start.button":          "点击开始",
		"mine.game.win.note":              "@{{ .Username }}\n恭喜！🎉\n您成功在 {{ .Seconds }} 秒内完成了游戏。\n地图尺寸：{{ .Width }} × {{ .Height }}\n地雷数量：{{ .Mines }}",
		"mine.game.win.button":            "再来一局",
		"mine.game.lose.note":             "@{{ .Username }}\n砰！💣\n耗时：{{ .Seconds }} 秒。\n地图尺寸：{{ .Width }} × {{ .Height }}\n地雷数量：{{ .Mines }}",
		"mine.game.lose.button":           "再试一次",
		"mine.game.opt.quit":              "退出",
		"mine.game.opt.flag":              "插旗",
		"mine.game.opt.click":             "扫雷",
		"mine.game.error":                 "@{{ .Username }} 哎呀！出了点问题！{{.Message}}！`,\n\t\t\"mine.help\": `@{{ .Username }}\n欢迎使用 ocha ！\n以下是一些帮助您入门的命令：\n\n/mine [&lt;宽度&gt; &lt;高度&gt; &lt;地雷数&gt;]\n/help\n\n<blockquote>\nocha bot\n作者: @feellmoose_dev\n版本信息:{{.Version}}\n更新于:{{.Update}}\n</blockquote>",
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
