package newsbuilder

import (
	newsfeed "i2pgit.org/idk/newsgo/builder/feed"
)

type NewsBuilder struct {
	Nodes []newsfeed.Node
	File  string
}

func (n *NewsBuilder) LoadFeed() {
	n.Nodes = newsfeed.XMLData(n.File)
}
