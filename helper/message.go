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
		"mine.game.start.level.note": `@{{ .Username }}
Hey there! 👋 Thanks for choosing Mine Sweeper Bot Plus!
Ready to play? Just follow the steps below to start a new game.
<pre>
level     w × h  mines
Easy      6 × 6      5
Normal    8 × 8     10
Hard      8 × 8     13
</pre>`,
		"mine.game.start.level.easy":   `Easy`,
		"mine.game.start.level.normal": `Normal`,
		"mine.game.start.level.hard":   `Hard`,
		"mine.game.menu.guide.note": `@{{ .Username }}
Hey there! 👋 Thanks for choosing Mine Sweeper Bot Plus!
Ready to play? Just follow the steps below to start a new game.`,
		"mine.game.menu.classic.button": `Classic`,
		"mine.game.menu.level.button":   `Level`,
		"mine.game.menu.random.button":  `Random a map`,
		"mine.game.menu.back.button":    `Back`,
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
/mine_cxg [&lt;width&gt; &lt;height&gt; &lt;mine&gt;]
/mine_random
/mine_level [level]
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
		"mine.game.start.level.note": `@{{ .Username }}
哈啰！👋 感谢您使用 Mine Sweeper Bot Plus！
请按照以下步骤开始扫雷。
<pre>
难度     宽 × 高   地雷数
简单     6 × 6        5
普通     8 × 8       10
困难     8 × 8       13
</pre>`,
		"mine.game.start.level.easy":   `简单`,
		"mine.game.start.level.normal": `普通`,
		"mine.game.start.level.hard":   `困难`,
		"mine.game.menu.guide.note": `@{{ .Username }}
哈啰！👋 感谢您使用 Mine Sweeper Bot Plus！
请按照以下步骤开始扫雷。`,
		"mine.game.menu.classic.button": `经典模式`,
		"mine.game.menu.level.button":   `选择难度`,
		"mine.game.menu.random.button":  `随机地图`,
		"mine.game.menu.back.button":    `返回`,
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
/mine_cxg [&lt;宽度&gt; &lt;高度&gt; &lt;地雷数&gt;]
/mine_random
/mine_level [难度等级]
/help

<blockquote>
Mine Sweeper Bot Plus
作者: @feellmoose_dev
版本信息:{{.Version}}
更新于:{{.Update}}
</blockquote>`,
	},
	"zh_CN_CiXioGui": {
		"mine.game.start.note": `@{{ .Username }}
哼哼~又来挑战我啦？👋 欢迎使用 Mine Sweeper Bot Plus！
你已经开始了一个新的 {{ .Width }} × {{ .Height }} 地图
共有 {{ .Mines }} 个地雷
区区杂鱼~杂鱼~♡肯定是要踩到地雷的呢~桀桀桀~♡`,
		"mine.game.start.button": `点我开始吧~`,
		"mine.game.start.level.note": `@{{ .Username }}
哼哼~又来挑战我啦~👋
欢迎使用 Mine Sweeper Bot Plus！
按照下面的步骤开始扫雷吧~
<pre>
难度      宽 × 高    地雷数
小杂鱼     6 × 6        5
一般的     8 × 8       10
厉害的     8 × 8       13
</pre>`,
		"mine.game.start.level.easy":   `小杂鱼`,
		"mine.game.start.level.normal": `一般的杂鱼`,
		"mine.game.start.level.hard":   `厉害的杂鱼~♡`,
		"mine.game.menu.guide.note": `@{{ .Username }}
哼哼~又来挑战我啦~👋
欢迎使用 Mine Sweeper Bot Plus！
按照下面的步骤开始扫雷吧~`,
		"mine.game.menu.classic.button": `经典模式~`,
		"mine.game.menu.level.button":   `选择难度~`,
		"mine.game.menu.random.button":  `随机地图~`,
		"mine.game.menu.back.button":    `返回~`,
		"mine.game.win.note": `@{{ .Username }}
什么！你居然赢了！
可恶的杂鱼居然~ 把那种地方（脸红）
全都摸了个遍呢~ 哈啊~ 明明是区区杂鱼！
在 {{ .Seconds }} 秒内完成了游戏
地图尺寸：{{ .Width }} × {{ .Height }}
地雷数量：{{ .Mines }}
杂鱼~杂鱼~我是不会承认的♡`,
		"mine.game.win.button": `再来一次？再来也是杂鱼~♡`,
		"mine.game.lose.note": `@{{ .Username }}
哦呀哦呀~踩到地雷了呢！💣
区区杂鱼不过如此嘛~
用时：{{ .Seconds }} 秒
地图尺寸：{{ .Width }} × {{ .Height }}
地雷数量：{{ .Mines }}
杂鱼杂鱼~这就输了吗~
笨蛋杂鱼就只能眼睁睁看着咱却不会玩~
自己弄爆炸掉了呢~♡ 真是杂鱼~杂鱼~♡`,
		"mine.game.lose.button": `再来一次？再来也是杂鱼~♡`,
		"mine.game.opt.quit":    `承认自己是杂鱼`,
		"mine.game.opt.flag":    `插旗子`,
		"mine.game.opt.click":   `摸地板`,
		"mine.game.error":       `@{{ .Username }} 哎呀~ 出错惹呜呜~  {{.Message}}`,
		"mine.help": `@{{ .Username }}
欢迎使用 Mine Sweeper Bot Plus！
以下是一些帮助你入门的命令：

/mine [&lt;宽度&gt; &lt;高度&gt; &lt;地雷数&gt;]
/mine_cxg [&lt;宽度&gt; &lt;高度&gt; &lt;地雷数&gt;]
/mine_random
/mine_level [难度等级]
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
