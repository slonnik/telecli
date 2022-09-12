package core

type CustomEventTypeEnum string
type CustomEvent map[string]interface{}

const (
	ChatSelectedEventType CustomEventTypeEnum = "chatSelectedEventType"
	UpdateScreenEventType CustomEventTypeEnum = "updateScreenEventType"
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

func NewUpdateScreenEvent() CustomEvent {
	return NewSimpleCustomEvent(UpdateScreenEventType)
}

var CoreEvents = make(chan CustomEvent)
