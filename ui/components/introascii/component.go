package components

import "github.com/rivo/tview"

type Component struct{}

func NewComponent() *Component {
	return &Component{}
}

func (c *Component) IntroASCII() *tview.TextView {
	bannerText := `
 $$$$$$\  $$$$$$$\  $$$$$$\ $$\   $$\ 
$$  __$$\ $$  __$$\ \_$$  _|$$$\  $$ |
$$ /  $$ |$$ |  $$ |  $$ |  $$$$\ $$ |
$$ |  $$ |$$ |  $$ |  $$ |  $$ $$\$$ |
$$ |  $$ |$$ |  $$ |  $$ |  $$ \$$$$ |
$$ |  $$ |$$ |  $$ |  $$ |  $$ |\$$$ |
 $$$$$$  |$$$$$$$  |$$$$$$\ $$ | \$$ |
 \______/ \_______/ \______|\__|  \__|
	`

	return tview.NewTextView().
		SetText(bannerText).
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(false)
}
