package core

type CustomEventTypeEnum string
type CustomEvent map[string]interface{}

const (
	ChatSelectedEventType                 CustomEventTypeEnum = "chatSelectedEventType"
	AuthorizationStateWaitPhoneNumberType CustomEventTypeEnum = "authorizationStateWaitPhoneNumber"
	AuthorizationStateWaitCodeType        CustomEventTypeEnum = "authorizationStateWaitCodeType"
	AuthorizationStateReadyType           CustomEventTypeEnum = "authorizationStateReadyType"
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

var events = make(chan CustomEvent)

func PublishEvent(event CustomEvent) {
	events <- event
}

func ReadEvent() CustomEvent {
	return <-events
}
