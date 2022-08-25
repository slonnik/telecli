package main

import (
	"github.com/Arman92/go-tdlib"
	"github.com/rivo/tview"
	"math"
	"slonnik.ru/telecli/core"
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
			tview.NewFlex().
				AddItem(
					tview.NewFlex().
						SetDirection(tview.FlexRow).
						AddItem(core.NewTeleList().SetBorder(true), 0, 3, false).
						AddItem(tview.NewBox().SetBorder(true), 3, 1, false), 0, 2, false).
				AddItem(core.NewChatList().SetBorder(true).SetTitle(" Chats "), 20, 1, false),
			true, true)

	pages.SwitchToPage(startPageLabel)

	go func() {

		for {

			currentState, _ := client.Authorize()

			switch currentState.GetAuthorizationStateEnum() {
			case tdlib.AuthorizationStateWaitPhoneNumberType:
				{
					core.PublishEvent(core.NewSimpleCustomEvent(core.AuthorizationStateWaitPhoneNumberType))
				}
			case tdlib.AuthorizationStateWaitCodeType:
				{
					core.PublishEvent(core.NewSimpleCustomEvent(core.AuthorizationStateWaitCodeType))
				}
			case tdlib.AuthorizationStateReadyType:
				{
					core.PublishEvent(core.NewSimpleCustomEvent(core.AuthorizationStateReadyType))
					goto MAIN_LOOP
				}
			}

		}

	MAIN_LOOP: // Main loop
		rawUpdates := client.GetRawUpdatesChannel(100)
		for update := range rawUpdates {

			switch tdlib.UpdateEnum(update.Data["@type"].(string)) {
			case tdlib.UpdateNewMessageType:

				message := update.Data["message"].(map[string]interface{})
				messageType := message["content"].(map[string]interface{})["@type"]
				chat, _ := client.GetChat(int64(message["chat_id"].(float64)))
				switch tdlib.MessageContentEnum(messageType.(string)) {
				case tdlib.MessageTextType:
					messageText := message["content"].(map[string]interface{})["text"].(map[string]interface{})["text"]
					event := core.NewUpdateNewMessageTextEvent(chat.ID, chat.Title, messageText.(string))
					core.PublishEvent(event)

				}
			}
		}
	}()

	go func() {
		for {
			event := core.ReadEvent()
			switch core.CustomEventTypeEnum(event["@type"].(string)) {
			case core.ChatSelectedEventType:
				{
					chatId := event["chatId"].(int64)
					chat, _ := client.GetChat(chatId)
					messages, _ := client.GetChatHistory(chatId, chat.LastReadInboxMessageID, -10, 10, false)

					app.QueueUpdate(func() {
						_, page := pages.GetFrontPage()
						list := page.(*tview.Flex).GetItem(0).(*tview.Flex).GetItem(0).(*core.TeleList)
						list.ClearItems()
						for _, message := range messages.Messages {
							switch message.Content.GetMessageContentEnum() {
							case tdlib.MessageTextType:
								{
									list.AddItem(chat.Title, message.Content.(*tdlib.MessageText).Text.Text)

								}
							}
						}
					})
					app.Draw()
				}
			case core.AuthorizationStateWaitPhoneNumberType:
				{
					app.QueueUpdate(func() {
						pages.SwitchToPage(phonePageLabel)
					})
					app.Draw()
				}
			case core.AuthorizationStateWaitCodeType:
				{
					app.QueueUpdate(func() {
						pages.SwitchToPage(codePageLabel)
					})
					app.Draw()
				}
			case core.AuthorizationStateReadyType:
				{
					app.QueueUpdate(func() {
						pages.SwitchToPage(mainPageLabel)
						_, page := pages.GetFrontPage()
						chats := page.(*tview.Flex).GetItem(1).(*core.ChatList)

						chatList, _ := getChatList(client, 100)

						for _, chat := range chatList {
							chats.AddChat(chat.Title, chat.ID)
						}
						chats.SelectChat(7)
						chats.SetFocus()
					})
					app.Draw()
				}

			case core.UpdateNewMessageTextType:
				{
					chatTitle := event["chatTitle"]
					messageText := event["text"]
					chatId := event["chatId"]
					app.QueueUpdate(func() {
						_, page := pages.GetFrontPage()
						chatList := page.(*tview.Flex).GetItem(1).(*core.ChatList)
						if chatList.GetSelectedChatId() == chatId {
							mainList := page.(*tview.Flex).GetItem(0).(*tview.Flex).GetItem(0).(*core.TeleList)
							mainList.AddItem(chatTitle.(string), messageText.(string))
						}

					})
					app.Draw()
				}

			}

		}
	}()

	err := app.SetRoot(pages, true).EnableMouse(true).Run()

	if err != nil {
		panic(err)
	}
}

func getChatList(client *tdlib.Client, limit int) ([]*tdlib.Chat, error) {
	var allChats []*tdlib.Chat
	var offsetOrder = int64(math.MaxInt64)
	var offsetChatID = int64(0)
	var chatList = tdlib.NewChatListMain()
	var lastChat *tdlib.Chat

	for len(allChats) < limit {
		if len(allChats) > 0 {
			lastChat = allChats[len(allChats)-1]
			for i := 0; i < len(lastChat.Positions); i++ {
				//Find the main chat list
				if lastChat.Positions[i].List.GetChatListEnum() == tdlib.ChatListMainType {
					offsetOrder = int64(lastChat.Positions[i].Order)
				}
			}
			offsetChatID = lastChat.ID
		}

		// get chats (ids) from tdlib
		var chats, getChatsErr = client.GetChats(chatList, tdlib.JSONInt64(offsetOrder),
			offsetChatID, int32(limit-len(allChats)))
		if getChatsErr != nil {
			return nil, getChatsErr
		}
		if len(chats.ChatIDs) == 0 {
			return allChats, nil
		}

		for _, chatID := range chats.ChatIDs {
			// get chat info from tdlib
			var chat, getChatErr = client.GetChat(chatID)
			if getChatErr == nil {
				allChats = append(allChats, chat)
			} else {
				return nil, getChatErr
			}
		}
	}

	return allChats, nil
}
