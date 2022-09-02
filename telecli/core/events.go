package core

type CustomEventTypeEnum string
type CustomEvent map[string]interface{}

const (
	ChatSelectedEventType                 CustomEventTypeEnum = "chatSelectedEventType"
	AuthorizationStateWaitPhoneNumberType CustomEventTypeEnum = "authorizationStateWaitPhoneNumber"
	AuthorizationStateWaitCodeType        CustomEventTypeEnum = "authorizationStateWaitCodeType"
	AuthorizationStateReadyType           CustomEventTypeEnum = "authorizationStateReadyType"
	UpdateNewMessageTextType              CustomEventTypeEnum = "updateNewMessageTextType"
)

func NewChatSelectedEvent(chatId int64) CustomEvent {
	var event = make(CustomEvent)
	event["@type"] = string(ChatSelectedEventType)
	event["chatId"] = chatId
	return event
}

func NewSimpleCustomEvent(state CustomEventTypeEnum) CustomEvent {
	var event = make(CustomEvent)
	event["@type"] = string(state)
	return event
}

func NewUpdateNewMessageTextEvent(chatId int64, chatTitle, text string, timeStamp int64) CustomEvent {
	var event = make(CustomEvent)
	event["@type"] = string(UpdateNewMessageTextType)
	event["chatId"] = chatId
	event["chatTitle"] = chatTitle
	event["text"] = text
	event["timeStamp"] = timeStamp
	return event
}

var events = make(chan CustomEvent)

func PublishEvents(eventsToPublish ...CustomEvent) {
	for _, event := range eventsToPublish {
		events <- event
	}
}

func ReadEvent() CustomEvent {
	return <-events
}
