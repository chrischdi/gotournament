package frontend

var MenuLinks []*MenuLink

func init() {
	MenuLinks = []*MenuLink{}

	AddMenuLink("/matchplan", "Matchplan", "")
	AddMenuLink("/table", "Table", "")
	AddMenuLink("/fulltable", "Complete Table", "")
	AddMenuLink("/ko", "KO", "")
	AddMenuLink("/setup", "Setup", "")
	AddMenuLink("/save/tofile", "Save", "")
	AddMenuLink("/load/fromfile", "Load", "uploadDB")
}

type MenuLink struct {
	Uri  string
	Name string
	ID   string
}

func AddMenuLink(uri, name, id string) {
	MenuLinks = append(MenuLinks, &MenuLink{uri, name, id})
}
