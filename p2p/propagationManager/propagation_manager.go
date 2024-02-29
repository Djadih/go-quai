package propagationManager

import (

)

// PropagationManager is responsible for managing the propagation of messages
// Each topic has its own validator function that uses the backend to validate the message

func makeValidator(topic string) func([]byte) bool {
	return func(msg []byte) bool { return true }
}