package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/fatih/color"
	"github.com/himidori/pm/db"
	"github.com/himidori/pm/utils"
	"github.com/ogier/pflag"
)

var (
	show        bool
	name        string
	group       string
	new         bool
	link        string
	user        string
	comment     string
	pass        string
	length      int
	remove      bool
	id          int
	open        bool
	interactive bool
	menu        bool
	rofi        bool
	table       bool
)

func printUsage() {
	fmt.Println(`Simple password manager written in Go

-s                      show password
-n [Name of resource]   name of resource
-g [Group name]         group name
-o                      open link
-t                      show passwords as table
-w                      store new password
-I                      interactive mode for adding new password
-l [Link]               link to resource
-u                      username
-c                      comment
-p [Password]           password
                        (if password is omitted PM will
                        generate a secure password)
-L [Length]             length of generated password
-r                      remove password
-i                      password ID
-m                      show dmenu
-R                      show rofi
-h                      show help`)
}

func initArgs() {
	pflag.BoolVarP(&show, "show", "s", false, "show password")
	pflag.StringVarP(&name, "name", "n", "", "name of the resource")
	pflag.StringVarP(&group, "group", "g", "", "name of the group")
	pflag.BoolVarP(&new, "write", "w", false, "add new password")
	pflag.StringVarP(&link, "link", "l", "", "link to the resource")
	pflag.StringVarP(&user, "user", "u", "", "username of the resource")
	pflag.StringVarP(&comment, "comment", "c", "", "comment")
	pflag.StringVarP(&pass, "password", "p", "", "password")
	pflag.IntVarP(&length, "length", "L", 16, "length of generated password")
	pflag.BoolVarP(&remove, "remove", "r", false, "remove password")
	pflag.IntVarP(&id, "id", "i", -1, "password id")
	pflag.BoolVarP(&open, "open", "o", false, "open link in browser")
	pflag.BoolVarP(&interactive, "interactive", "I", false, "interactive mode")
	pflag.BoolVarP(&menu, "menu", "m", false, "show dmenu")
	pflag.BoolVarP(&rofi, "rofi", "R", false, "show rofi")
	pflag.BoolVarP(&table, "table", "t", false, "print passwords table")
	pflag.Usage = printUsage

	pflag.Parse()
}

func parseArgs() {
	if menu {
		ok, err := utils.IsIntalled("dmenu")
		if err != nil {
			fmt.Println("failed to check dmenu installation:", err)
			return
		}
		if !ok {
			fmt.Println("dmenu is not installed")
			return
		}

		passwords, err := db.SelectAll()
		if err != nil {
			fmt.Println("failed to fetch passwords:", err)
			return
		}
		if passwords == nil {
			fmt.Println("no passwords found")
			return
		}

		str := ""
		for _, p := range passwords {
			str += p.Name + "|"
			if p.Group != "" {
				str += p.Group + "|"
			}
			str += p.Resource + "\n"
		}
		res, err := utils.ShowMenu("dmenu", str)
		if err != nil {
			fmt.Println("failed to show menu:", err)
			return
		}
		if res == "" {
			return
		}

		data := strings.Split(res, "|")
		n := strings.Split(res, "|")[0]
		g := ""
		if len(data) == 3 {
			g = data[1]
		}
		fmt.Println(g)
		for _, p := range passwords {
			if p.Name == n && p.Group == g {
				err = clipboard.WriteAll(p.Password)
				if err != nil {
					utils.Notify(p.Name, "failed to copy password to the clipboard")
				}
				utils.Notify(p.Name, "copied password to the clipboard!")

				return
			}
		}

	}

	if rofi {
		ok, err := utils.IsIntalled("rofi")
		if err != nil {
			fmt.Println("failed to check rofi installation:", err)
			return
		}
		if !ok {
			fmt.Println("rofi is not installed")
			return
		}

		passwords, err := db.SelectAll()
		if err != nil {
			fmt.Println("failed to fetch passwords:", err)
			return
		}

		if passwords == nil {
			fmt.Println("no passwords found")
			return
		}

		str := ""
		longestName := getLongestNameField(passwords)
		longestGroup := getLongestGroupField(passwords)
		for _, p := range passwords {
			nameSpaces := longestName - len(p.Name) + 1
			groupSpaces := longestGroup - len(p.Group) + 1
			str += p.Name +
				strings.Repeat(" ", nameSpaces) +
				p.Group +
				strings.Repeat(" ", groupSpaces) +
				p.Resource + "\n"
		}
		res, err := utils.ShowMenu("rofi", str)
		if err != nil {
			fmt.Println("failed to show menu:", err)
			return
		}
		if res == "" {
			return
		}

		fields := strings.Fields(res)
		n := fields[0]
		g := ""
		if len(fields) == 3 {
			g = fields[1]
		}

		for _, p := range passwords {
			if p.Name == n && p.Group == g {
				err = clipboard.WriteAll(p.Password)
				if err != nil {
					utils.Notify(p.Name, "failed to copy password to the clipboard")
				}
				utils.Notify(p.Name, "copied password to the clipboard!")

				return
			}
		}
	}

	if !show && !new && !remove {
		printUsage()
		return
	}

	if show {
		if name == "" && group == "" {
			printUsage()
			return
		}

		if name != "" && group == "" {
			passwd, err := db.SelectByName(name)
			if err != nil {
				fmt.Println("failed to get password:", err)
				return
			}
			if passwd == nil {
				fmt.Println("no passwords found for name", name)
				return
			}

			if len(passwd) > 1 || name == "all " {
				if table {
					printTable(passwd)
				} else {
					printPasswords(passwd)
				}
				return
			}

			err = clipboard.WriteAll(passwd[0].Password)
			if err != nil {
				fmt.Println("failed to copy password to the clipboard")
			} else {
				fmt.Println("password was copied to the clipboard!")
			}

			fmt.Print("URL: ")
			color.Blue(passwd[0].Resource)
			fmt.Print("User: ")
			color.Yellow(passwd[0].Username)
			if passwd[0].Group != "" {
				fmt.Print("Group: ")
				color.Magenta(passwd[0].Group)
			}

			if open {
				utils.OpenURL(passwd[0].Resource)
			}
		}

		if name == "" && group != "" {
			passwords, err := db.SelectByGroup(group)
			if err != nil {
				fmt.Println("failed to get passwords:", err)
				return
			}

			if passwords == nil {
				fmt.Println("no passwords found for group", group)
				return
			}

			if table {
				printTable(passwords)
			} else {
				fmt.Print("Group: ")
				color.Magenta(group)
				printPasswords(passwords)
			}
		}

		if name != "" && group != "" {
			passwords, err := db.SelectByGroupAndName(name, group)
			if err != nil {
				fmt.Println("failed to get passwords:", err)
				return
			}

			if passwords == nil {
				fmt.Println("no password found")
				return
			}

			err = clipboard.WriteAll(passwords[0].Password)
			if err != nil {
				fmt.Println("failed to copy password to the clipboard")
			} else {
				fmt.Println("password was copied to the clipboard!")
			}

			fmt.Print("URL: ")
			color.Blue(passwords[0].Resource)
			fmt.Print("User: ")
			color.Yellow(passwords[0].Username)
			if passwords[0].Group != "" {
				fmt.Print("Group: ")
				color.Magenta(passwords[0].Group)
			}

			if open {
				utils.OpenURL(passwords[0].Resource)
			}
		}
	}

	if remove {
		if id == -1 {
			printUsage()
			return
		}

		err := db.RemovePassword(id)
		if err != nil {
			fmt.Println("failed to remove password:", err)
			return
		}

		fmt.Println("successfuly removed password with id", id)
	}

	if new {
		if interactive {
			addInteractive()
			return
		}

		if name == "" || link == "" {
			printUsage()
			return
		}

		passwd := ""
		var err error

		if pass != "" {
			passwd = pass
		} else {
			passwd, err = db.GeneratePassword(length)
			if err != nil {
				fmt.Println("failed to generate password:", err)
				return
			}
		}

		err = db.AddPassword(&db.Password{
			Name:     name,
			Resource: link,
			Password: passwd,
			Username: user,
			Comment:  comment,
			Group:    group,
		})

		if err != nil {
			fmt.Println("failed to add password:", err)
			return
		}

		fmt.Println("successfuly added new password!")
	}
}

func addInteractive() {
	fmt.Print("name: ")
	name, err := utils.ReadLine()
	if err != nil {
		fmt.Println("failed to read line:", err)
		return
	}

	fmt.Print("resource: ")
	resource, err := utils.ReadLine()
	if err != nil {
		fmt.Println("failed to read line:", err)
		return
	}

	fmt.Print("password (leave empty to generate): ")
	passwd, err := utils.ReadLine()
	if err != nil {
		fmt.Println("failed to read line:", err)
		return
	}

	if passwd == "" {
		fmt.Print("length of generated password: ")
		le, err := utils.ReadLine()
		if err != nil {
			fmt.Println("failed to read line:", err)
			return
		}
		length, err = strconv.Atoi(le)
		if err != nil {
			fmt.Println("invalid input:", err)
			return
		}
	}

	fmt.Print("username: ")
	username, err := utils.ReadLine()
	if err != nil {
		fmt.Println("failed to read line:", err)
		return
	}

	fmt.Print("comment: ")
	comment, err := utils.ReadLine()
	if err != nil {
		fmt.Println("failed to read line:", err)
		return
	}

	fmt.Print("group: ")
	grp, err := utils.ReadLine()
	if err != nil {
		fmt.Println("failed to read line:", err)
	}

	if passwd == "" {
		passwd, err = db.GeneratePassword(length)
		if err != nil {
			fmt.Println("failed to generate password:", err)
			return
		}
	}

	err = db.AddPassword(&db.Password{
		Name:     name,
		Resource: resource,
		Password: passwd,
		Username: username,
		Comment:  comment,
		Group:    grp,
	})

	if err != nil {
		fmt.Println("failed to add password:", err)
		return
	}

	fmt.Println("successfuly added password to the database!")
}

// getLongestNameField return int with length of name
func getLongestNameField(passwords []*db.Password) int {
	counter := 0
	for i := range passwords {
		if len(passwords[i].Name) > counter {
			counter = len(passwords[i].Name)
		}
	}
	return counter
}

// getLongestGroupField return int with length of group
func getLongestGroupField(passwords []*db.Password) int {
	counter := 0
	for i := range passwords {
		if len(passwords[i].Group) > counter {
			counter = len(passwords[i].Group)
		}
	}
	return counter
}

// getLongestResourceField return int with length of resource
func getLongestResourceField(passwords []*db.Password) int {
	counter := 0
	for i := range passwords {
		if len(passwords[i].Resource) > counter {
			counter = len(passwords[i].Resource)
		}
	}
	return counter
}

func getLongestUsernameField(passwords []*db.Password) int {
	counter := 0
	for i := range passwords {
		if len(passwords[i].Username) > counter {
			counter = len(passwords[i].Username)
		}
	}
	return counter
}

func getLongestCommentField(passwords []*db.Password) int {
	counter := 0
	for i := range passwords {
		if len(passwords[i].Comment) > counter {
			counter = len(passwords[i].Comment)
		}
	}
	return counter
}

func getLongestIdField(passwords []*db.Password) int {
	counter := 0
	for i := range passwords {
		length := len(strconv.Itoa(passwords[i].Id))
		if length > counter {
			counter = length
		}
	}
	return counter
}

func printPasswords(passwords []*db.Password) {
	fmt.Println()

	for _, p := range passwords {
		c := color.New(color.FgYellow)
		c.Println("id: ", p.Id)
		c = color.New(color.FgRed)
		c.Println("name: ", p.Name)
		c = color.New(color.FgGreen)
		c.Println("resource: ", p.Resource)
		c = color.New(color.FgBlue)
		c.Println("username: ", p.Username)
		c = color.New(color.FgCyan)
		c.Println("comment:", p.Comment)
		c = color.New(color.FgMagenta)
		c.Println("group: ", p.Group)
		fmt.Println()
	}
}

func printTable(passwords []*db.Password) {
	longestId := getLongestIdField(passwords)
	idAddSpaces := 0
	longestName := getLongestNameField(passwords)
	nameAddSpaces := 0
	longestResource := getLongestResourceField(passwords)
	resourceAddSpaces := 0
	longestUsername := getLongestUsernameField(passwords)
	usernameAddSpaces := 0
	longestComment := getLongestCommentField(passwords)
	commentAddSpaces := 0
	longestGroup := getLongestGroupField(passwords)
	groupAddSpaces := 0
	totalSpaces := 0

	fmt.Println()

	c := color.New(color.FgYellow)

	if longestId < 2 {
		idAddSpaces = 1
		c.Printf("id ")
		totalSpaces += 3
	} else {
		c.Printf("id" + strings.Repeat(" ", (longestId-2)+1))
		totalSpaces += longestId + 1
	}

	c = color.New(color.FgRed)
	if longestName < 4 {
		nameAddSpaces = 4 - longestName
		totalSpaces += 5
		c.Printf("name ")
	} else {
		c.Printf("name" + strings.Repeat(" ", (longestName-4)+1))
		totalSpaces += longestName + 1
	}

	c = color.New(color.FgGreen)
	if longestResource < 8 {
		resourceAddSpaces = 8 - longestResource
		totalSpaces += 9
		c.Printf("resource ")
	} else {
		c.Printf("resource" + strings.Repeat(" ", (longestResource-8)+1))
		totalSpaces += longestResource + 1
	}

	c = color.New(color.FgBlue)
	if longestUsername < 8 {
		usernameAddSpaces = 8 - longestUsername
		totalSpaces += 9
		c.Printf("username ")
	} else {
		c.Printf("username" + strings.Repeat(" ", (longestUsername-8)+1))
		totalSpaces += longestUsername + 1
	}

	c = color.New(color.FgCyan)
	if longestComment < 7 {
		commentAddSpaces = 7 - longestComment
		totalSpaces += 8
		c.Printf("comment ")
	} else {
		c.Printf("comment" + strings.Repeat(" ", (longestComment-7)+1))
		totalSpaces += longestComment + 1
	}

	c = color.New(color.FgMagenta)
	if longestGroup < 5 {
		groupAddSpaces = 5 - longestGroup
		totalSpaces += 6
		c.Printf("group ")
	} else {
		c.Printf("group" + strings.Repeat(" ", (longestGroup-5)+1))
		totalSpaces += longestGroup
	}

	fmt.Println()
	fmt.Println(strings.Repeat("-", totalSpaces))

	for _, p := range passwords {
		idStr := strconv.Itoa(p.Id)
		idSpaces := (longestId - len(idStr)) + idAddSpaces + 1
		nameSpaces := (longestName - len(p.Name)) + nameAddSpaces + 1
		resourceSpaces := (longestResource - len(p.Resource)) + resourceAddSpaces + 1
		usernameSpaces := (longestUsername - len(p.Username)) + usernameAddSpaces + 1
		commentSpaces := (longestComment - len(p.Comment)) + commentAddSpaces + 1
		groupSpaces := (longestGroup - len(p.Group)) + groupAddSpaces + 1

		fmt.Println(
			idStr + strings.Repeat(" ", idSpaces) +
				p.Name + strings.Repeat(" ", nameSpaces) +
				p.Resource + strings.Repeat(" ", resourceSpaces) +
				p.Username + strings.Repeat(" ", usernameSpaces) +
				p.Comment + strings.Repeat(" ", commentSpaces) +
				p.Group + strings.Repeat(" ", groupSpaces),
		)
	}

	fmt.Println()
}
