package core

import "github.com/Arman92/go-tdlib"

type CustomEventTypeEnum string
type CustomEvent map[string]interface{}

const (
	ChatSelectedEventType                 CustomEventTypeEnum = "chatSelectedEventType"
	AuthorizationStateWaitPhoneNumberType CustomEventTypeEnum = "authorizationStateWaitPhoneNumber"
	AuthorizationStateWaitCodeType        CustomEventTypeEnum = "authorizationStateWaitCodeType"
	AuthorizationStateReadyType           CustomEventTypeEnum = "authorizationStateReadyType"
	UpdateNewMessageTextType              CustomEventTypeEnum = "updateNewMessageTextType"
	UpdateScreenEventType                 CustomEventTypeEnum = "updateScreenEventType"
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

func NewUpdateNewMessageTextEvent(message tdlib.Message) CustomEvent {
	var event = make(CustomEvent)
	event["@type"] = string(UpdateNewMessageTextType)
	event["message"] = message
	return event
}

func NewUpdateScreenEvent() CustomEvent {
	return NewSimpleCustomEvent(UpdateScreenEventType)
}

var events = make(chan CustomEvent)

func PublishEvents(eventsToPublish ...CustomEvent) {
	go func() {
		for _, event := range eventsToPublish {
			events <- event
		}
	}()
}

func ReadEvent() CustomEvent {
	return <-events
}
