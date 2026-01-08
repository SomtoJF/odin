package components

type Component struct{}

func NewComponent() *Component {
	return &Component{}
}

func (c *Component) IntroASCII() string {
	return `
 $$$$$$\  $$$$$$$\  $$$$$$\ $$\   $$\ 
$$  __$$\ $$  __$$\ \_$$  _|$$$\  $$ |
$$ /  $$ |$$ |  $$ |  $$ |  $$$$\ $$ |
$$ |  $$ |$$ |  $$ |  $$ |  $$ $$\$$ |
$$ |  $$ |$$ |  $$ |  $$ |  $$ \$$$$ |
$$ |  $$ |$$ |  $$ |  $$ |  $$ |\$$$ |
 $$$$$$  |$$$$$$$  |$$$$$$\ $$ | \$$ |
 \______/ \_______/ \______|\__|  \__|
	`
}
