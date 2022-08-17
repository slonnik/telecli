package main

import (
	"github.com/Arman92/go-tdlib"
	"github.com/rivo/tview"
)

const (
	startPageLabel = "Start"
	phonePageLabel = "Phone"
	codePageLabel  = "Code"
	mainPageLabel  = "Main"
)

func main() {
	tdlib.SetLogVerbosityLevel(1)
	tdlib.SetFilePath("./errors.txt")

	// Create new instance of client
	client := tdlib.NewClient(tdlib.Config{
		APIID:               "187786",
		APIHash:             "e782045df67ba48e441ccb105da8fc85",
		SystemLanguageCode:  "en",
		DeviceModel:         "Server",
		SystemVersion:       "1.0.0",
		ApplicationVersion:  "1.0.0",
		UseMessageDatabase:  true,
		UseFileDatabase:     true,
		UseChatInfoDatabase: true,
		UseTestDataCenter:   false,
		DatabaseDirectory:   "./tdlib-db",
		FileDirectory:       "./tdlib-files",
		IgnoreFileNames:     false,
	})

	app := tview.NewApplication()
	pages := tview.NewPages()

	phoneNumber := ""
	codeNumber := ""

	pages.
		AddPage(startPageLabel, tview.
			NewForm(), true, true).
		AddPage(phonePageLabel, tview.
			NewForm().
			AddInputField("Phone", "", 20, func(textToCheck string, lastChar rune) bool {
				return true
			}, func(text string) { phoneNumber = text }).
			AddButton("Cancel", func() { app.Stop() }).
			AddButton("Next", func() { client.SendPhoneNumber(phoneNumber) }), true, true).
		AddPage(codePageLabel, tview.
			NewForm().
			AddInputField("code", "", 20, func(textToCheck string, lastChar rune) bool {
				return true
			}, func(text string) { codeNumber = text }).
			AddButton("Cancel", func() { app.Stop() }).
			AddButton("Previous", func() { pages.SwitchToPage(phonePageLabel) }).
			AddButton("Done", func() { client.SendAuthCode(codeNumber) }), true, true).
		AddPage(mainPageLabel,
			NewTeleList(), true, true)

	pages.SwitchToPage(startPageLabel)

	authStates := make(chan tdlib.AuthorizationStateEnum)
	go func() {
		for {
			state := <-authStates
			if state == tdlib.AuthorizationStateWaitPhoneNumberType {
				app.QueueUpdate(func() {
					pages.SwitchToPage(phonePageLabel)
				})
				app.Draw()
			} else if state == tdlib.AuthorizationStateWaitCodeType {
				app.QueueUpdate(func() {
					pages.SwitchToPage(codePageLabel)
				})
				app.Draw()
			} else if state == tdlib.AuthorizationStateWaitPasswordType {

			}
			if state == tdlib.AuthorizationStateReadyType {
				app.QueueUpdate(func() {
					pages.SwitchToPage(mainPageLabel)
				})
				app.Draw()
				break
			}
		}
	}()

	updates := make(chan tdlib.UpdateData)
	go func() {
		for {
			update := <-updates
			switch tdlib.UpdateEnum(update["@type"].(string)) {
			case tdlib.UpdateNewMessageType:

				//
				message := update["message"].(map[string]interface{})
				messageType := message["content"].(map[string]interface{})["@type"]
				chat, _ := client.GetChat(int64(message["chat_id"].(float64)))
				switch tdlib.MessageContentEnum(messageType.(string)) {
				case tdlib.MessageTextType:
					messageText := message["content"].(map[string]interface{})["text"].(map[string]interface{})["text"]
					app.QueueUpdate(func() {
						_, page := pages.GetFrontPage()
						list := page.(*TeleList)
						list.AddItem(chat.Title, messageText.(string))
					})
					app.Draw()
				}

			}
		}
	}()

	go func() {

		var previousStateEnum tdlib.AuthorizationStateEnum

		for {
			currentState, _ := client.Authorize()
			if currentState.GetAuthorizationStateEnum() != previousStateEnum {
				previousStateEnum = currentState.GetAuthorizationStateEnum()
				authStates <- previousStateEnum
			}

			if currentState.GetAuthorizationStateEnum() == tdlib.AuthorizationStateReadyType {
				//fmt.Println("Authorization Ready! Let's rock")
				break
			}
		}

		// Main loop
		rawUpdates := client.GetRawUpdatesChannel(100)
		for update := range rawUpdates {
			updates <- update.Data
		}
	}()

	err := app.SetRoot(pages, true).EnableMouse(true).Run()

	if err != nil {
		panic(err)
	}
}
