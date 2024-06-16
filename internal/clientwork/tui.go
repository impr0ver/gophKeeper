package clientwork

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"path"
	"regexp"

	"github.com/impr0ver/gophKeeper/internal/handlers"
	"github.com/impr0ver/gophKeeper/internal/storage"
	"github.com/impr0ver/gophKeeper/internal/userdata"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	log "github.com/sirupsen/logrus"
	"golang.design/x/clipboard"
)

const (
	picture = `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AAAAAjMzNAAAAAACMzM0AP/AABEIANwBEgMBIgACEQEDEQH/xAAfAAABBQEBAQEBAQAAAAAAAAAAAQIDBAUGBwgJCgv/xAC1EAACAQMDAgQDBQUEBAAAAX0BAgMABBEFEiExQQYTUWEHInEUMoGRoQgjQrHBFVLR8CQzYnKCCQoWFxgZGiUmJygpKjQ1Njc4OTpDREVGR0hJSlNUVVZXWFlaY2RlZmdoaWpzdHV2d3h5eoOEhYaHiImKkpOUlZaXmJmaoqOkpaanqKmqsrO0tba3uLm6wsPExcbHyMnK0tPU1dbX2Nna4eLj5OXm5+jp6vHy8/T19vf4+fr/xAAfAQADAQEBAQEBAQEBAAAAAAAAAQIDBAUGBwgJCgv/xAC1EQACAQIEBAMEBwUEBAABAncAAQIDEQQFITEGEkFRB2FxEyIygQgUQpGhscEJIzNS8BVictEKFiQ04SXxFxgZGiYnKCkqNTY3ODk6Q0RFRkdISUpTVFVWV1hZWmNkZWZnaGlqc3R1dnd4eXqCg4SFhoeIiYqSk5SVlpeYmZqio6Slpqeoqaqys7S1tre4ubrCw8TFxsfIycrS09TV1tfY2dri4+Tl5ufo6ery8/T19vf4+fr/2wBDAAEBAQEBAQIBAQIDAgICAwQDAwMDBAYEBAQEBAYHBgYGBgYGBwcHBwcHBwcICAgICAgJCQkJCQsLCwsLCwsLCwv/2wBDAQICAgMDAwUDAwULCAYICwsLCwsLCwsLCwsLCwsLCwsLCwsLCwsLCwsLCwsLCwsLCwsLCwsLCwsLCwsLCwsLCwv/3QAEABL/2gAMAwEAAhEDEQA/AP4Z6KKKDMKKKKACiiigAooooAKKKKACiiigAooooAKKKKACiiigAooooAKKKKACiiigAooooAKKKKACiiigAooooAKKKKACiiigD//Q/hnooooMwooooAKKKKACiiigAooooAKKKKACiiigAooooAKKKKACiiigAooooAKKKKACiiigAooooAKKKKACiiigAooooAKKKKAP/9H+GeiiigzCiiigAooooAKKKKACiiigAooooAKKKKACiiigAooooAKKKKACiiigAooooAKKKKACiiigAooooAKKKKACiiigAooooA//0v4Z6KKKDMKKKKACiiigAooooAKKKKACiiigAooooAKKKKACiiigAooooAKK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`
)

// TUI is a struct for terminal user interface.
type TUI struct {
	*tview.Application
	pages       *tview.Pages
	client      handlers.ClientHandlers
	maxFileSize int64
}

// NewTUI gets new terminal user interface for client.
func NewTUI(client handlers.ClientHandlers, fileSize int64) *TUI {
	application := tview.NewApplication()
	pages := tview.NewPages()

	if err := clipboard.Init(); err != nil {
		log.Fatalln("Failed init clipboard:", err)
	}

	application.SetRoot(pages, true).EnableMouse(true)

	tui := &TUI{
		Application: application,
		client:      client,
		pages:       pages,
		maxFileSize: fileSize,
	}

	tui.authPage("Please set login & password for continue =>>>")

	return tui
}

// getImage get image from base64 string
func (app *TUI) getImage() image.Image {
	b, err := base64.StdEncoding.DecodeString(picture)
	if err != nil {
		log.Info("base64 string error")
	}
	photo, err := jpeg.Decode(bytes.NewReader(b))
	if err != nil {
		log.Info("jpeg decode string error")
	}
	return photo
}

// authPage switches to authentication page, where user can log in or register.
// login page
func (app *TUI) authPage(message string) {
	credentials := userdata.UserCredentials{}

	form := tview.NewForm()
	form.SetBorder(true)
	form.SetBorderColor(tcell.ColorDarkGrey)
	form.SetFieldTextColor(tcell.ColorBlack)
	form.SetFieldBackgroundColor(tcell.ColorWhite)
	form.SetButtonBackgroundColor(tcell.ColorGray)
	form.SetLabelColor(tcell.ColorLightGreen)

	form.AddImage("\"GophKeeper client\"", app.getImage(), 0, 14, 16777216)
	form.AddInputField("Login", "", 35, nil, func(login string) {
		credentials.Login = login
	})
	form.AddPasswordField("Password", "", 35, '*', func(password string) {
		credentials.Password = password
	})
	form.AddPasswordField("AES Key", "", 35, '*', func(masterKey string) {
		credentials.AESKey = masterKey
	})

	form.AddButton("Login", func() {
		err := app.client.Login(credentials)

		if errors.Is(err, storage.ErrWrongCredentials) {
			log.Infoln(storage.ErrWrongCredentials)

			app.authPage("[red]Wrong credentials. Please try again.[white]")
			return
		}
		if errors.Is(err, handlers.ErrEmptyField) {
			log.Infoln(handlers.ErrEmptyField)

			app.authPage("[red]Some fields are empty.[white]")
			return
		}
		if errors.Is(err, storage.ErrUnknown) || err != nil {
			log.Infoln(storage.ErrUnknown)

			app.authPage("[red]Something is wrong. ;([white]")
			return
		}

		app.recordsInfoPage("[green]Login successfully![white]")
	})

	form.AddButton("Register", func() {
		err := app.client.Register(credentials)

		if errors.Is(err, storage.ErrLoginExists) {
			log.Infoln(storage.ErrLoginExists)

			app.authPage("[red]Login exists. Please try again.[white]")
			return
		}
		if errors.Is(err, handlers.ErrEmptyField) {
			log.Infoln(handlers.ErrEmptyField)

			app.authPage("[red]Some fields are empty.[white]")
			return
		}
		if errors.Is(err, storage.ErrUnknown) || err != nil {
			log.Infoln(storage.ErrUnknown)

			app.authPage("[red]Something is wrong. ;([white]")
			return
		}

		app.recordsInfoPage("[green]Registered successfully.[white]")
	})

	frame := tview.NewFrame(form).SetBorders(0, 0, 0, 1, 4, 4).
		AddText(
			"TAB - switch fields / Enter - choose option",
			false,
			tview.AlignLeft,
			tcell.ColorWhite,
		).
		AddText(
			"Ctrl+C - exit",
			false,
			tview.AlignLeft,
			tcell.ColorWhite,
		).
		AddText(
			message,
			false,
			tview.AlignRight,
			tcell.ColorWhite,
		)

	app.pages.AddPage("authentication", frame, true, true)
	app.pages.SwitchToPage("authentication")
}

// recordInfoPage switches to page, where are all records shown.
// main page
func (app *TUI) recordsInfoPage(message string) {

	records, err := app.client.GetRecordsInfo()

	if errors.Is(err, storage.ErrUnauthenticated) {
		log.Infoln(storage.ErrUnauthenticated)

		app.authPage("[red]Session expired. Please login again.[white]")
		return
	}
	if errors.Is(err, storage.ErrUnknown) || err != nil {
		log.Infoln(storage.ErrUnknown)

		message = "[red]Something is wrong. ;([white]"
		return
	}

	list := tview.NewList()
	list.SetBorder(true)
	list.SetBorderColor(tcell.ColorDarkGrey)
	list.SetMainTextColor(tcell.ColorGrey)
	list.SetSecondaryTextColor(tcell.ColorLightGreen)
	list.SetShortcutColor(tcell.ColorLightGreen)
	for _, record := range records {

		f := func(record userdata.Record) func() {
			return func() {
				app.recordPage(record.ID, "")
			}
		}(record)

		if record.Metadata == "" {
			record.Metadata = "no metadata"
		}

		list.AddItem(record.ID, "Type: "+record.Type.String()+" | Metadata: "+record.Metadata+" | AES key hint: "+record.KeyHint, '⏺', f)
	}

	listFrame := tview.NewFrame(list).SetBorders(0, 0, 0, 1, 4, 4).
		AddText(
			"↑ or ↓ - switch records / Enter - choose option",
			false,
			tview.AlignLeft,
			tcell.ColorWhite,
		).
		AddText(
			"Ctrl+N - create new record / Ctrl+R - refresh page",
			false,
			tview.AlignLeft,
			tcell.ColorWhite,
		).
		AddText(
			"Ctrl+K - change AES key",
			false,
			tview.AlignLeft,
			tcell.ColorWhite,
		).
		AddText(
			"ESC - return to the login page",
			false,
			tview.AlignLeft,
			tcell.ColorWhite,
		).
		AddText(
			message, false,
			tview.AlignRight,
			tcell.ColorWhite,
		)

	listFrame.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlN {
			app.createRecordPage("")
		}
		if event.Key() == tcell.KeyCtrlR {
			app.recordsInfoPage("[green]Refreshed.[white]")
		}
		if event.Key() == tcell.KeyCtrlK {
			app.setNewAESKey("Change AES Key")
		}
		if event.Key() == tcell.KeyESC {
			app.authPage("Logget out")
		}
		return event
	})

	app.pages.AddPage("records", listFrame, true, true)
	app.pages.SwitchToPage("records")
}

// setNewAESKey set new AES key for decrypt records
func (app *TUI) setNewAESKey(message string) {
	var newAESKey string

	form := tview.NewForm()

	form.SetBorder(true)
	form.SetBorderColor(tcell.ColorDarkGrey)
	form.SetFieldTextColor(tcell.ColorBlack)
	form.SetFieldBackgroundColor(tcell.ColorWhite)
	form.SetButtonBackgroundColor(tcell.ColorGray)
	form.SetLabelColor(tcell.ColorLightGreen)

	form.AddPasswordField("AES Key", "", 35, '*', func(masterKey string) {
		newAESKey = masterKey
	})

	form.AddButton("Set AES key", func() {
		err := app.client.SetAESKey(newAESKey)
		if errors.Is(err, handlers.ErrEmptyField) {
			log.Infoln(handlers.ErrEmptyField)

			app.authPage("[red]AES field is empty.[white]")
			return
		}
		if errors.Is(err, storage.ErrUnknown) || err != nil {
			log.Infoln(storage.ErrUnknown)

			app.authPage("[red]Something is wrong. ;([white]")
			return
		}

		app.recordsInfoPage("[green]AES key change successfully![white]")
	})

	frame := tview.NewFrame(form).SetBorders(0, 0, 0, 1, 4, 4).
		AddText(
			"Enter - choose option / ESC - return to the menu",
			false,
			tview.AlignLeft,
			tcell.ColorWhite,
		).
		AddText(
			message,
			false,
			tview.AlignRight,
			tcell.ColorWhite,
		)

	app.pages.AddPage("setaeskey", frame, true, true)
	app.pages.SwitchToPage("setaeskey")

	frame.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyESC {
			app.recordsInfoPage("Returned to menu.")
		}
		return event
	})
}

// recordPage - record page (decrypted record data, copy or delete record).
func (app *TUI) recordPage(recordID string, message string) {
	record, err := app.client.GetRecord(recordID)

	if errors.Is(err, storage.ErrUnauthenticated) {
		log.Infoln(storage.ErrUnauthenticated)

		app.authPage("[red]Session expired. Please login again.[white]")
		return
	}
	if errors.Is(err, storage.ErrNotFound) {
		log.Infoln(storage.ErrNotFound)

		app.recordsInfoPage("[red]Not found this record.[white]")
		return
	}

	if err != nil {
		log.Infoln("Failed get record. Wrong AES key???")

		app.recordsInfoPage("[red]Failed get record. Wrong AES key???[white]")
		return
	}

	if record.Metadata == "" {
		record.Metadata = "no metadata"
	}

	frame := tview.NewFrame(
		tview.NewTextView().
			SetText(string(record.Data)).
			SetTextColor(tcell.ColorLightGreen).
			SetDisabled(true)).
		SetBorders(0, 0, 0, 1, 4, 4).
		AddText(
			record.Metadata+" | "+record.Type.String(),
			true,
			tview.AlignCenter,
			tcell.ColorLightGreen,
		).
		AddText(
			"Ctrl+K - copy / Ctrl+D - delete / ESC - return to the menu",
			false,
			tview.AlignLeft,
			tcell.ColorWhite,
		).
		AddText(
			message,
			false,
			tview.AlignRight,
			tcell.ColorWhite,
		)

	frame.SetBorder(true)
	frame.SetBorderColor(tcell.ColorDarkGrey)

	frame.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyESC:
			app.recordsInfoPage("Returned to menu.")
		case tcell.KeyCtrlK:
			app.recordPage(recordID, "[green]Copied successfully.[white]")
			clipboard.Write(clipboard.FmtText, record.Data)
		case tcell.KeyCtrlD:
			err := app.client.DeleteRecord(recordID)

			if errors.Is(err, storage.ErrUnauthenticated) {
				log.Infoln(storage.ErrUnauthenticated)

				app.authPage("[red]Session expired. Please login again.[white]")
				return event
			}
			if errors.Is(err, handlers.ErrWrongAESKey) {
				log.Infoln(handlers.ErrWrongAESKey)

				app.authPage("[red]Wrong AES key. Please login again.[white]")
				return event
			}
			if errors.Is(err, storage.ErrNotFound) {
				log.Infoln(storage.ErrNotFound)

				app.recordsInfoPage("[red]Failed to delete. Not found record.[white]")
				return event
			}
			if errors.Is(err, storage.ErrUnknown) || err != nil {
				log.Infoln(storage.ErrUnknown)

				app.recordPage(recordID, "[red]Something is wrong. ;([white]")
				return event
			}

			app.recordsInfoPage("[green]Deleted successfully.[white]")
		}

		return event
	})

	app.pages.AddPage("record", frame, true, true)
	app.pages.SwitchToPage("record")
}

// createTextRecord creates text record.
func (app *TUI) createTextRecord() {
	record := userdata.Record{Type: userdata.TypeText}
	form := tview.NewForm()

	form.SetBorder(true)
	form.SetBorderColor(tcell.ColorDarkGrey)
	form.SetFieldTextColor(tcell.ColorBlack)
	form.SetFieldBackgroundColor(tcell.ColorWhite)
	form.SetButtonBackgroundColor(tcell.ColorGray)
	form.SetLabelColor(tcell.ColorLightGreen)

	form.AddTextArea("Text", "", 50, 6, 0, func(text string) {
		record.Data = []byte(text)
	})
	form.AddInputField("Metadata", "", 20, nil, func(text string) {
		record.Metadata = text
	})
	form.AddButton("OK", func() {
		err := app.client.CreateRecord(record)

		if errors.Is(err, storage.ErrUnauthenticated) {
			log.Infoln(storage.ErrUnauthenticated)

			app.authPage("[red]Session expired. Please login again.[white]")
			return
		}
		if errors.Is(err, storage.ErrUnknown) {
			log.Infoln(storage.ErrUnknown)

			app.recordsInfoPage("[red]Something is wrong. ;([white]")
			return
		}
		if errors.Is(err, handlers.ErrWrongAESKey) {
			app.authPage("[red]Wrong AES key. Please login again.[white]")
			return
		}

		app.recordsInfoPage("[green]Created record successfully.[white]")
	})

	frame := tview.NewFrame(form).SetBorders(0, 0, 0, 1, 4, 4).
		AddText(
			"TAB - switch fields / Enter - choose option",
			false,
			tview.AlignLeft,
			tcell.ColorWhite,
		).
		AddText("ESC - return to the menu.", false, tview.AlignLeft, tcell.ColorWhite)

	frame.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyESC {
			app.recordsInfoPage("Returned to menu.")
		}
		return event
	})

	app.pages.AddPage("createTextRecord", frame, true, true)
	app.pages.SwitchToPage("createTextRecord")
}

// createCredentialsRecord creates credentials (login and password) record.
func (app *TUI) createCredentialsRecord() {
	record := userdata.Record{
		Type: userdata.TypeLoginAndPassword,
	}

	loginAndPassword := userdata.LoginAndPassword{}
	form := tview.NewForm()

	form.SetBorder(true)
	form.SetBorderColor(tcell.ColorDarkGrey)
	form.SetFieldTextColor(tcell.ColorBlack)
	form.SetFieldBackgroundColor(tcell.ColorWhite)
	form.SetButtonBackgroundColor(tcell.ColorGray)
	form.SetLabelColor(tcell.ColorLightGreen)

	form.AddInputField("Login", "", 20, nil, func(text string) {
		loginAndPassword.Login = text
	})
	form.AddInputField("Password", "", 20, nil, func(text string) {
		loginAndPassword.Password = text
	})
	form.AddInputField("Metadata", "", 20, nil, func(text string) {
		record.Metadata = text
	})
	form.AddButton("OK", func() {
		record.Data, _ = loginAndPassword.Bytes()
		err := app.client.CreateRecord(record)

		if errors.Is(err, storage.ErrUnauthenticated) {
			app.authPage("[red]Session expired. Please login again.[white]")
			return
		}
		if errors.Is(err, storage.ErrUnknown) {
			app.recordsInfoPage("[red]Something is wrong. ;([white]")
			return
		}
		if errors.Is(err, handlers.ErrWrongAESKey) {
			app.authPage("[red]Wrong AES key. Please login again.[white]")
			return
		}

		app.recordsInfoPage("[green]Created record successfully.[white]")
	})

	frame := tview.NewFrame(form).SetBorders(0, 0, 0, 1, 4, 4).
		AddText("TAB - switch fields / Enter - choose option", false, tview.AlignLeft, tcell.ColorWhite).
		AddText("ESC - exit to all records.", false, tview.AlignLeft, tcell.ColorWhite)

	frame.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyESC {
			app.recordsInfoPage("Returned to menu.")
		}
		return event
	})

	app.pages.AddPage("createTextRecord", frame, true, true)
	app.pages.SwitchToPage("createTextRecord")
}

// createCardRecord creates credit card record (card number, expiration date, cvc).
func (app *TUI) createCardRecord() {
	record := userdata.Record{
		Type: userdata.TypeCreditCard,
	}

	creditCard := userdata.CreditCard{}
	form := tview.NewForm()

	form.SetBorder(true)
	form.SetBorderColor(tcell.ColorDarkGrey)
	form.SetFieldTextColor(tcell.ColorBlack)
	form.SetFieldBackgroundColor(tcell.ColorWhite)
	form.SetButtonBackgroundColor(tcell.ColorGray)
	form.SetLabelColor(tcell.ColorLightGreen)

	form.AddInputField("Card number", "", 20, nil, func(text string) {
		creditCard.Number = text
	})
	form.AddInputField("Expiration", "", 20, nil, func(text string) {
		creditCard.ExpirationDate = text
	})
	form.AddInputField("CVC", "", 4, nil, func(text string) {
		creditCard.CVC = text
	})
	form.AddInputField("Metadata", "", 20, nil, func(text string) {
		record.Metadata = text
	})

	form.AddButton("OK", func() {
		// Regex for credit card - exp date from regex101.com
		result, err := regexp.Match(
			`^(0[1-9]|1[0-2])[|/]?([0-9]{4}|[0-9]{2})$`,
			[]byte(creditCard.ExpirationDate),
		)
		if err != nil {
			log.Fatalln("Failed parse regex for check card expiration date.")
		}
		if !result {
			app.recordsInfoPage("[red]Incorrect expiration date.[white]")
			return
		}

		// Regex for credit card - card number from regex101.com
		result, err = regexp.Match(
			`^(?:4[0-9]{12}(?:[0-9]{3})?|[25][1-7][0-9]{14}|6(?:011|5[0-9][0-9])[0-9]{12}|3[47][0-9]{13}|3(?:0[0-5]|[68][0-9])[0-9]{11}|(?:2131|1800|35\d{3})\d{11})$`,
			[]byte(creditCard.Number),
		)
		if err != nil {
			log.Fatalln("Failed parse regex for check card expiration date.")
			return
		}
		if !result {
			app.recordsInfoPage("[red]Incorrect card number.[white]")
			return
		}

		// Regex for credit card - card CVC from regex101.com
		result, err = regexp.Match(`\d{3}`, []byte(creditCard.CVC))
		if err != nil {
			log.Fatalln("Failed parse regex for check card expiration date.")
		}
		if !result {
			app.recordsInfoPage("[red]Incorrect CVC code.[white]")
			return
		}

		record.Data, _ = creditCard.Bytes()
		err = app.client.CreateRecord(record)

		if errors.Is(err, storage.ErrUnauthenticated) {
			app.authPage("[red]Session expired. Please login again.[white]")
			return
		}
		if errors.Is(err, storage.ErrUnknown) {
			app.recordsInfoPage("[red]Something is wrong. ;([white]")
			return
		}
		if errors.Is(err, handlers.ErrWrongAESKey) {
			app.authPage("[red]Wrong AES key. Please login again.[white]")
			return
		}

		app.recordsInfoPage("[green]Created record successfully.[white]")
	})

	frame := tview.NewFrame(form).SetBorders(0, 0, 0, 1, 4, 4).
		AddText(
			"[yellow]Attention! Regex used![white]",
			false,
			tview.AlignCenter,
			tcell.ColorLightGreen,
		).
		AddText(
			"TAB - switch fields / Enter - choose option",
			false,
			tview.AlignLeft,
			tcell.ColorWhite,
		).
		AddText(
			"ESC - return to the menu.",
			false,
			tview.AlignLeft,
			tcell.ColorWhite,
		)

	frame.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyESC {
			app.recordsInfoPage("Returned to menu.")
		}
		return event
	})

	app.pages.AddPage("createCardRecord", frame, true, true)
	app.pages.SwitchToPage("createCardRecord")
}

// createFileRecord creates file record.
func (app *TUI) createFileRecord() {
	record := userdata.Record{Type: userdata.TypeFile}
	file := userdata.BinaryFile{}
	form := tview.NewForm()

	form.SetBorder(true)
	form.SetBorderColor(tcell.ColorDarkGrey)
	form.SetFieldTextColor(tcell.ColorBlack)
	form.SetFieldBackgroundColor(tcell.ColorWhite)
	form.SetButtonBackgroundColor(tcell.ColorGray)
	form.SetLabelColor(tcell.ColorLightGreen)

	form.AddInputField("Please, enter filepath:", "", 70, nil, func(text string) {
		file.FilePath = text
	})

	form.AddButton("OK", func() {
		// Get filename from path
		filename := path.Base(file.FilePath)
		record.Metadata = filename

		data, err := file.Bytes()
		if err != nil {
			app.recordsInfoPage("[red]Failed opened file.[white]")
			return
		}

		// Get file size
		dataSize, err := file.Size()
		if err != nil {
			app.recordsInfoPage("[red]Failed get file size.[white]")
			return
		}

		// Check file in large, and size > cfg.MaxFileSize
		if dataSize > app.maxFileSize {
			app.recordsInfoPage("[red]File size is very big, please choose another file[white]")
			return
		}

		record.Data = data
		err = app.client.CreateRecord(record)

		if errors.Is(err, storage.ErrUnauthenticated) {
			app.authPage("[red]Session expired. Please login again.[white]")
			return
		}
		if errors.Is(err, storage.ErrUnknown) {
			app.recordsInfoPage("[red]Something is wrong. ;([white]")
			return
		}
		if errors.Is(err, handlers.ErrWrongAESKey) {
			app.authPage("[red]Wrong AES key. Please login again.[white]")
			return
		}

		app.recordsInfoPage("[green]Created record successfully.[white]")
	})

	frame := tview.NewFrame(form).SetBorders(0, 0, 0, 1, 4, 4).
		AddText(
			fmt.Sprintf("[yellow]Attention! File size must bee <= %d Bytes![white]", app.maxFileSize),
			false,
			tview.AlignCenter,
			tcell.ColorLightGreen,
		).
		AddText(
			"TAB - switch fields / Enter - choose option",
			false,
			tview.AlignLeft,
			tcell.ColorWhite,
		).
		AddText(
			"ESC - return to the menu.",
			false,
			tview.AlignLeft,
			tcell.ColorWhite,
		)

	frame.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyESC {
			app.recordsInfoPage("Returned to menu.")
		}
		return event
	})

	app.pages.AddPage("createFileRecord", frame, true, true)
	app.pages.SwitchToPage("createFileRecord")
}

// createRecordPage creates page, where can choose record type.
func (app *TUI) createRecordPage(message string) {
	form := tview.NewForm()

	form.SetBorder(true)
	form.SetBorderColor(tcell.ColorDarkGrey)
	form.SetFieldTextColor(tcell.ColorBlack)
	form.SetFieldBackgroundColor(tcell.ColorWhite)
	form.SetButtonBackgroundColor(tcell.ColorGray)
	form.SetLabelColor(tcell.ColorLightGreen)

	form.AddDropDown(
		"Please, select type data: ",
		[]string{"Text", "Login&password", "Credit card", "Binary file"}, 1,
		func(option string, optionIndex int) {
			switch option {
			case "Text":
				app.createTextRecord()
			case "Login&password":
				app.createCredentialsRecord()
			case "Credit card":
				app.createCardRecord()
			case "Binary file":
				app.createFileRecord()
			}
		})

	frame := tview.NewFrame(form).SetBorders(0, 0, 0, 1, 4, 4).
		AddText(
			"↑ or ↓  - switch fields / Enter - choose option",
			false,
			tview.AlignLeft,
			tcell.ColorWhite,
		).
		AddText(
			"ESC - exit to all records.",
			false,
			tview.AlignLeft,
			tcell.ColorWhite,
		).
		AddText(
			message,
			false,
			tview.AlignRight,
			tcell.ColorWhite,
		)

	frame.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyESC {
			app.recordsInfoPage("Returned to menu.")
		}
		return event
	})

	app.pages.AddPage("create", frame, true, true)
	app.pages.SwitchToPage("create")
}
