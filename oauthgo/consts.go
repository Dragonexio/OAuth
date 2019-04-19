package oauthgo

// scopes
const (
	ScopeLogin       = 1
	ScopeUserInfo    = 2
	ScopeWithdraw    = 3
	ScopeTrade       = 4
	ScopePayUser2App = 5
	ScopePayApp2User = 6
	ScopeAssetInfo   = 7
	ScopeOpenApi     = 8
	ScopeAutoAuth    = 9
)

// transfer direction between user and app
const (
	OrderDirectionUser2App = 1
	OrderDirectionApp2User = 2
)

// dragonex pay order status
const (
	OrderStatusSucceed      = 1
	OrderStatusFailed       = 2
	OrderStatusTransferring = 3
)
