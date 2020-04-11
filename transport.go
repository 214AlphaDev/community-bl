package community_bl

type Transport interface {
	SendConfirmationCode(confirmationCode ConfirmationCode) error
}
