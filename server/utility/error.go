package utility

// GenericErrorMessage can represent any custom error with a message
type GenericErrorMessage struct {
	Message string
}

// Error interface implementation
func (genericErrorMessage GenericErrorMessage) Error() string {
	return genericErrorMessage.Message
}
