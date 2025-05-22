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
		"mine.game.start.note": `@{{ .Username }}
Hey there! 👋 Thanks for choosing Mine Sweeper Bot Plus!
You have started a new {{ .Width }} × {{ .Height }} game with {{ .Mines }} mines.`,
		"mine.game.start.button": `Click to Start`,
		"mine.game.win.note": `@{{ .Username }}
Congratulations! 🎉
You've successfully completed the game in {{ .Seconds }} seconds.
Map Dimensions: {{ .Width }} × {{ .Height }}
Number of Mines: {{ .Mines }}
Well done on your achievement!`,
		"mine.game.win.button": `Try again?`,
		"mine.game.lose.note": `@{{ .Username }}
Boom! 💣
Unfortunately, you hit a mine and the game has ended.
Time Elapsed: {{ .Seconds }} seconds.
Map Dimensions: {{ .Width }} × {{ .Height }}
Number of Mines: {{ .Mines }}
Better luck next time!`,
		"mine.game.lose.button": `Try again?`,
		"mine.game.opt.quit":    `Quit`,
		"mine.game.opt.flag":    `Flag`,
		"mine.game.opt.click":   `Click`,
		"mine.game.error":       `@{{ .Username }} Oops! Something went wrong! {{.Message}}!`,
		"mine.help": `@{{ .Username }}
Hey there! 👋 Thanks for choosing Mine Sweeper Bot Plus!
Here's a list of commands to get you started:

/mine [&lt;width&gt; &lt;height&gt; &lt;mine&gt;]
/help

<blockquote>
Mine Sweeper Bot Plus created by @feellmoose_dev
Version {{.Version}}
Last update at {{.Update}}
</blockquote>`,
	},
	"zh": {
		"mine.game.start.note": `@{{ .Username }}
哈啰！👋 感谢您使用 Mine Sweeper Bot Plus！
您已开始一个新的 {{ .Width }} × {{ .Height }} 地图
共有 {{ .Mines }} 个地雷`,
		"mine.game.start.button": `点击开始`,
		"mine.game.win.note": `@{{ .Username }}
恭喜！🎉
您成功在 {{ .Seconds }} 秒内完成了游戏。
地图尺寸：{{ .Width }} × {{ .Height }}
地雷数量：{{ .Mines }}
干得漂亮！`,
		"mine.game.win.button": `再来一局`,
		"mine.game.lose.note": `@{{ .Username }}
砰！💣
很遗憾，您踩到了地雷，游戏结束。
耗时：{{ .Seconds }} 秒。
地图尺寸：{{ .Width }} × {{ .Height }}
地雷数量：{{ .Mines }}
祝您下次好运！`,
		"mine.game.lose.button": `再试一次`,
		"mine.game.opt.quit":    `退出`,
		"mine.game.opt.flag":    `切换\插旗`,
		"mine.game.opt.click":   `切换\扫雷`,
		"mine.game.error":       `@{{ .Username }} 哎呀！出了点问题！{{.Message}}！`,
		"mine.help": `@{{ .Username }}
感谢您使用 Mine Sweeper Bot Plus ！
以下是一些帮助您入门的命令：

/mine [&lt;宽度&gt; &lt;高度&gt; &lt;地雷数&gt;]
/help

<blockquote>
Mine Sweeper Bot Plus
作者: @feellmoose_dev
版本信息:{{.Version}}
更新于:{{.Update}}
</blockquote>`,
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
