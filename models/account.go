package models

import "github.com/ServiceWeaver/weaver"

// note: most of these fields will be changed to better fit criteria. this is
// simply a prototype for now.
type Account struct {
	weaver.AutoMarshal
	ID       int
	Alias    string //account username
	Name     string //acount holders name
	Email    string //account email
	Password string //account passphrase
	//Recovery_key  *string   //secret key generate for recovery, is deleted once shown
	//Profile_limit *uint8    //number of profiles the account can maintain, depends on tier
	Tier int //account tier level- default: 1
	//TwoFactorAuth *bool     //if 2FA is enabled or not on the account
	//IP            *[]string //list of users devices IP address

	//optional fields
	//Secondary_email *string //secondary email for verification & other purposes
	//Phone_number    *string //acount owners phone number
}
