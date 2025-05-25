package helpers

import "github.com/fatih/color"

func Information(format string, a ...interface{}) {
	c := color.New(color.FgWhite, color.Bold)
	c.Printf("[INFO] "+format+"\n", a...)
}

func Warning(format string, a ...interface{}) {
	c := color.New(color.FgYellow, color.Bold)
	c.Printf("[WARNING] "+format+"\n", a...)
}

func Error(format string, a ...interface{}) {
	c := color.New(color.FgRed, color.Bold)
	c.Printf("[ERROR] "+format+"\n", a...)
}

func Success(format string, a ...interface{}) {
	c := color.New(color.FgGreen, color.Bold)
	c.Printf("[SUCCESS] "+format+"\n", a...)
}

func Debug(format string, a ...interface{}) {
	c := color.New(color.FgCyan, color.Bold)
	c.Printf("[DEBUG] "+format+"\n", a...)
}

func Help(format string, a ...interface{}) {
	c := color.New(color.FgMagenta, color.Bold)
	c.Printf(""+format+"\n", a...)
}
