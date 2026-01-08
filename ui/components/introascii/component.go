package components

import "github.com/rivo/tview"

type Component struct{}

func NewComponent() *Component {
	return &Component{}
}

func (c *Component) IntroASCII() *tview.TextView {
	bannerText := `
  $$$$$$\  $$$$$$$\  $$$$$$\ $$\   $$\        $$$$$$\   $$$$$$\  $$$$$$$\  $$$$$$$$\ 
$$  __$$\ $$  __$$\ \_$$  _|$$$\  $$ |      $$  __$$\ $$  __$$\ $$  __$$\ $$  _____|
$$ /  $$ |$$ |  $$ |  $$ |  $$$$\ $$ |      $$ /  \__|$$ /  $$ |$$ |  $$ |$$ |      
$$ |  $$ |$$ |  $$ |  $$ |  $$ $$\$$ |      $$ |      $$ |  $$ |$$ |  $$ |$$$$$\    
$$ |  $$ |$$ |  $$ |  $$ |  $$ \$$$$ |      $$ |      $$ |  $$ |$$ |  $$ |$$  __|   
$$ |  $$ |$$ |  $$ |  $$ |  $$ |\$$$ |      $$ |  $$\ $$ |  $$ |$$ |  $$ |$$ |      
 $$$$$$  |$$$$$$$  |$$$$$$\ $$ | \$$ |      \$$$$$$  | $$$$$$  |$$$$$$$  |$$$$$$$$\ 
 \______/ \_______/ \______|\__|  \__|       \______/  \______/ \_______/ \________|
	`

	return tview.NewTextView().
		SetText(bannerText).
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(false)
}
