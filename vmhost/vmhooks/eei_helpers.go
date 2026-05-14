package vmhooks

const esdtRoleLocalMint = "ESDTRoleLocalMint"
const esdtRoleLocalBurn = "ESDTRoleLocalBurn"
const esdtRoleNFTCreate = "ESDTRoleNFTCreate"
const esdtRoleNFTAddQuantity = "ESDTRoleNFTAddQuantity"
const esdtRoleNFTBurn = "ESDTRoleNFTBurn"
const esdtRoleNFTUpdateAttributes = "ESDTRoleNFTUpdateAttributes"
const esdtRoleNFTAddURI = "ESDTRoleNFTAddURI"
const esdtRoleNFTRecreate = "ESDTRoleNFTRecreate"
const esdtRoleModifyCreator = "ESDTRoleModifyCreator"
const esdtRoleModifyRoyalties = "ESDTRoleModifyRoyalties"
const esdtRoleSetNewURI = "ESDTRoleSetNewURI"

const tickerMinLength = 3
const tickerMaxLength = 10
const additionalRandomCharsLength = 6
const identifierMinLength = tickerMinLength + additionalRandomCharsLength + 1
const identifierMaxLength = tickerMaxLength + additionalRandomCharsLength + 1

// constants defining roles values
const (
	RoleMint = 1 << iota
	RoleBurn
	RoleNFTCreate
	RoleNFTAddQuantity
	RoleNFTBurn
	RoleNFTUpdateAttributes
	RoleNFTAddURI
	RoleNFTRecreate
	RoleModifyCreator
	RoleModifyRoyalties
	RoleSetNewURI
)

func roleFromByteArray(bytes []byte) int64 {
	stringValue := string(bytes)
	switch stringValue {
	case esdtRoleLocalMint:
		return RoleMint
	case esdtRoleLocalBurn:
		return RoleBurn
	case esdtRoleNFTCreate:
		return RoleNFTCreate
	case esdtRoleNFTAddQuantity:
		return RoleNFTAddQuantity
	case esdtRoleNFTBurn:
		return RoleNFTBurn
	default:
		return 0
	}
}

func roleFromByteArrayV2(bytes []byte) int64 {
	stringValue := string(bytes)
	switch stringValue {
	case esdtRoleLocalMint:
		return RoleMint
	case esdtRoleLocalBurn:
		return RoleBurn
	case esdtRoleNFTCreate:
		return RoleNFTCreate
	case esdtRoleNFTAddQuantity:
		return RoleNFTAddQuantity
	case esdtRoleNFTBurn:
		return RoleNFTBurn
	case esdtRoleNFTUpdateAttributes:
		return RoleNFTUpdateAttributes
	case esdtRoleNFTAddURI:
		return RoleNFTAddURI
	case esdtRoleNFTRecreate:
		return RoleNFTRecreate
	case esdtRoleModifyCreator:
		return RoleModifyCreator
	case esdtRoleModifyRoyalties:
		return RoleModifyRoyalties
	case esdtRoleSetNewURI:
		return RoleSetNewURI
	default:
		return 0
	}
}

// getESDTRoles parses a serialized ESDT roles record.
//
// ISSUE-084: each iteration was previously vulnerable to slice-out-of-range
// panic on malformed input. Three sites were unprotected:
//
//  1. After the `currentIndex += 1` to skip the \n delimiter, the loop
//     read `dataBuffer[currentIndex]` without verifying currentIndex was
//     still in range. A buffer that ended on a \n delimiter (i.e., \n is
//     the last byte) would over-read by one.
//  2. `endIndex := currentIndex + roleLen` could exceed valueLen.
//  3. `dataBuffer[currentIndex:endIndex]` then panics if endIndex > valueLen.
//
// FAIL-CLOSED policy: on ANY bounds violation we return `0` (= no roles
// granted). This is the safe default for AUTHORIZATION data — a malformed
// buffer must NEVER grant a privilege that a well-formed parse would not
// grant. An earlier version of the fix returned the roles parsed so far
// before the malformed section, which was fail-permissive: a buffer that
// began with a high-privilege role byte and then trailed off into garbage
// would have granted that role. Returning 0 closes that vector.
//
// The signature returns int64 with no error path; the only graceful
// signal the parser has is the bitmask itself. The caller treats a
// non-zero return as "these specific role bits are granted." Returning
// 0 is equivalent to "no roles granted" — the appropriate fail-closed
// outcome for an authorization parser.
//
// Threat model note: dataBuffer originates from blockchain state via the
// BlockchainHook (set by built-in functions, not directly by user input),
// so direct attacker control is normally absent. The fail-closed policy
// is defense in depth against state-corruption / built-in-function bugs
// / future code that exposes a more permissive path.
func getESDTRoles(dataBuffer []byte, cryptoOpcodesV2Enabled bool) int64 {
	result := int64(0)
	currentIndex := 0
	valueLen := len(dataBuffer)

	for currentIndex < valueLen {
		// first character before each role is a \n, so we skip it
		currentIndex += 1

		// ISSUE-084 fail-closed: any bounds violation discards all roles
		// parsed so far and returns 0.
		if currentIndex >= valueLen {
			return 0
		}

		// next is the length of the role as string
		roleLen := int(dataBuffer[currentIndex])
		currentIndex += 1

		// ISSUE-084 fail-closed: roleLen exceeding remaining buffer is a
		// malformed-record signal. Return 0 rather than partial roles.
		endIndex := currentIndex + roleLen
		if endIndex > valueLen {
			return 0
		}

		// next is role's ASCII string representation
		roleName := dataBuffer[currentIndex:endIndex]
		currentIndex = endIndex

		if cryptoOpcodesV2Enabled {
			result |= roleFromByteArrayV2(roleName)
		} else {
			result |= roleFromByteArray(roleName)
		}
	}
	return result
}

// ValidateToken - validates the token ID
func ValidateToken(tokenID []byte) bool {
	tokenIDLen := len(tokenID)
	if tokenIDLen < identifierMinLength || tokenIDLen > identifierMaxLength {
		return false
	}

	tickerLen := tokenIDLen - additionalRandomCharsLength

	if !isTickerValid(tokenID[0 : tickerLen-1]) {
		return false
	}

	// dash char between the random chars and the ticker
	if tokenID[tickerLen-1] != '-' {
		return false
	}

	if !randomCharsAreValid(tokenID[tickerLen:tokenIDLen]) {
		return false
	}

	return true
}

// ticker must be all uppercase alphanumeric
func isTickerValid(tickerName []byte) bool {
	if len(tickerName) < tickerMinLength || len(tickerName) > tickerMaxLength {
		return false
	}
	for _, ch := range tickerName {
		isBigCharacter := ch >= 'A' && ch <= 'Z'
		isNumber := ch >= '0' && ch <= '9'
		isReadable := isBigCharacter || isNumber
		if !isReadable {
			return false
		}
	}

	return true
}

// random chars are alphanumeric lowercase
func randomCharsAreValid(chars []byte) bool {
	if len(chars) != additionalRandomCharsLength {
		return false
	}
	for _, ch := range chars {
		isSmallCharacter := ch >= 'a' && ch <= 'f'
		isNumber := ch >= '0' && ch <= '9'
		isReadable := isSmallCharacter || isNumber
		if !isReadable {
			return false
		}
	}

	return true
}
