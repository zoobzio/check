package check

import (
	"encoding/base64"
	"encoding/json"
	"net"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

// Pre-compiled regular expressions for format validation.
var (
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9.!#$%&'*+/=?^_` + "`" + `{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$`)
	uuidRegex  = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[1-5][0-9a-fA-F]{3}-[89abAB][0-9a-fA-F]{3}-[0-9a-fA-F]{12}$`)
	uuid4Regex = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-4[0-9a-fA-F]{3}-[89abAB][0-9a-fA-F]{3}-[0-9a-fA-F]{12}$`)

	hexColorRegex     = regexp.MustCompile(`^#([0-9a-fA-F]{3}|[0-9a-fA-F]{6}|[0-9a-fA-F]{8})$`)
	hexColorFullRegex = regexp.MustCompile(`^#[0-9a-fA-F]{6}$`)

	semverRegex = regexp.MustCompile(`^v?(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(?:-((?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+([0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$`)

	e164Regex = regexp.MustCompile(`^\+[1-9]\d{1,14}$`)

	macRegex      = regexp.MustCompile(`^([0-9A-Fa-f]{2}:){5}[0-9A-Fa-f]{2}$`)
	macDashRegex  = regexp.MustCompile(`^([0-9A-Fa-f]{2}-){5}[0-9A-Fa-f]{2}$`)
	hostnameRegex = regexp.MustCompile(`^[a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(\.[a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$`)
)

// Email validates that a string is a valid email address.
func Email(v, field string) *Validation {
	var err error
	if !emailRegex.MatchString(v) {
		err = fieldErr(field, "must be a valid email address")
	}
	return validation(err, field, "email")
}

// URL validates that a string is a valid URL.
func URL(v, field string) *Validation {
	var err error
	u, parseErr := url.Parse(v)
	if parseErr != nil || u.Scheme == "" || u.Host == "" {
		err = fieldErr(field, "must be a valid URL")
	}
	return validation(err, field, "url")
}

// URLWithScheme validates that a string is a valid URL with one of the given schemes.
func URLWithScheme(v string, schemes []string, field string) *Validation {
	var err error
	u, parseErr := url.Parse(v)
	if parseErr != nil || u.Scheme == "" || u.Host == "" {
		err = fieldErr(field, "must be a valid URL")
	} else {
		found := false
		for _, s := range schemes {
			if strings.EqualFold(u.Scheme, s) {
				found = true
				break
			}
		}
		if !found {
			err = fieldErrf(field, "must have scheme: %s", strings.Join(schemes, ", "))
		}
	}
	return validation(err, field, "url")
}

// HTTPOrHTTPS validates that a string is a valid HTTP or HTTPS URL.
func HTTPOrHTTPS(v, field string) *Validation {
	return URLWithScheme(v, []string{"http", "https"}, field)
}

// UUID validates that a string is a valid UUID (versions 1-5).
func UUID(v, field string) *Validation {
	var err error
	if !uuidRegex.MatchString(v) {
		err = fieldErr(field, "must be a valid UUID")
	}
	return validation(err, field, "uuid")
}

// UUID4 validates that a string is a valid UUID version 4.
func UUID4(v, field string) *Validation {
	var err error
	if !uuid4Regex.MatchString(v) {
		err = fieldErr(field, "must be a valid UUID v4")
	}
	return validation(err, field, "uuid4")
}

// IP validates that a string is a valid IP address (v4 or v6).
func IP(v, field string) *Validation {
	var err error
	if net.ParseIP(v) == nil {
		err = fieldErr(field, "must be a valid IP address")
	}
	return validation(err, field, "ip")
}

// IPv4 validates that a string is a valid IPv4 address.
func IPv4(v, field string) *Validation {
	var err error
	ip := net.ParseIP(v)
	if ip == nil || ip.To4() == nil {
		err = fieldErr(field, "must be a valid IPv4 address")
	}
	return validation(err, field, "ipv4")
}

// IPv6 validates that a string is a valid IPv6 address.
func IPv6(v, field string) *Validation {
	var err error
	ip := net.ParseIP(v)
	if ip == nil || ip.To4() != nil {
		err = fieldErr(field, "must be a valid IPv6 address")
	}
	return validation(err, field, "ipv6")
}

// CIDR validates that a string is a valid CIDR notation.
func CIDR(v, field string) *Validation {
	var err error
	_, _, parseErr := net.ParseCIDR(v)
	if parseErr != nil {
		err = fieldErr(field, "must be a valid CIDR notation")
	}
	return validation(err, field, "cidr")
}

// MAC validates that a string is a valid MAC address.
func MAC(v, field string) *Validation {
	var err error
	if !macRegex.MatchString(v) && !macDashRegex.MatchString(v) {
		err = fieldErr(field, "must be a valid MAC address")
	}
	return validation(err, field, "mac")
}

// Hostname validates that a string is a valid hostname.
func Hostname(v, field string) *Validation {
	var err error
	if len(v) > 253 || !hostnameRegex.MatchString(v) {
		err = fieldErr(field, "must be a valid hostname")
	}
	return validation(err, field, "hostname")
}

// Port validates that a string is a valid port number (1-65535).
func Port(v, field string) *Validation {
	var err error
	p, parseErr := strconv.Atoi(v)
	if parseErr != nil || p < 1 || p > 65535 {
		err = fieldErr(field, "must be a valid port number (1-65535)")
	}
	return validation(err, field, "port")
}

// HostPort validates that a string is a valid host:port combination.
func HostPort(v, field string) *Validation {
	var err error
	host, port, splitErr := net.SplitHostPort(v)
	switch {
	case splitErr != nil:
		err = fieldErr(field, "must be a valid host:port")
	case host == "":
		err = fieldErr(field, "host must not be empty")
	default:
		p, parseErr := strconv.Atoi(port)
		if parseErr != nil || p < 1 || p > 65535 {
			err = fieldErr(field, "port must be a valid number (1-65535)")
		}
	}
	return validation(err, field, "hostport")
}

// HexColor validates that a string is a valid hex color (#RGB, #RRGGBB, or #RRGGBBAA).
func HexColor(v, field string) *Validation {
	var err error
	if !hexColorRegex.MatchString(v) {
		err = fieldErr(field, "must be a valid hex color")
	}
	return validation(err, field, "hexcolor")
}

// HexColorFull validates that a string is a valid 6-digit hex color (#RRGGBB).
func HexColorFull(v, field string) *Validation {
	var err error
	if !hexColorFullRegex.MatchString(v) {
		err = fieldErr(field, "must be a valid hex color (#RRGGBB)")
	}
	return validation(err, field, "hexcolor")
}

// Base64 validates that a string is valid base64.
func Base64(v, field string) *Validation {
	var err error
	_, decodeErr := base64.StdEncoding.DecodeString(v)
	if decodeErr != nil {
		err = fieldErr(field, "must be valid base64")
	}
	return validation(err, field, "base64")
}

// Base64URL validates that a string is valid URL-safe base64.
func Base64URL(v, field string) *Validation {
	var err error
	_, decodeErr := base64.URLEncoding.DecodeString(v)
	if decodeErr != nil {
		err = fieldErr(field, "must be valid URL-safe base64")
	}
	return validation(err, field, "base64url")
}

// JSON validates that a string is valid JSON.
func JSON(v, field string) *Validation {
	var err error
	if !json.Valid([]byte(v)) {
		err = fieldErr(field, "must be valid JSON")
	}
	return validation(err, field, "json")
}

// Semver validates that a string is a valid semantic version.
func Semver(v, field string) *Validation {
	var err error
	if !semverRegex.MatchString(v) {
		err = fieldErr(field, "must be a valid semantic version")
	}
	return validation(err, field, "semver")
}

// E164 validates that a string is a valid E.164 phone number.
func E164(v, field string) *Validation {
	var err error
	if !e164Regex.MatchString(v) {
		err = fieldErr(field, "must be a valid E.164 phone number")
	}
	return validation(err, field, "e164")
}

// CreditCard validates that a string is a valid credit card number using the Luhn algorithm.
func CreditCard(v, field string) *Validation {
	var err error
	// Remove spaces and dashes
	clean := strings.ReplaceAll(strings.ReplaceAll(v, " ", ""), "-", "")
	if len(clean) < 13 || len(clean) > 19 {
		err = fieldErr(field, "must be a valid credit card number")
	} else {
		valid := true
		// Check all digits
		for _, r := range clean {
			if r < '0' || r > '9' {
				valid = false
				break
			}
		}
		if valid {
			// Luhn algorithm
			sum := 0
			nDigits := len(clean)
			parity := nDigits % 2
			for i := 0; i < nDigits; i++ {
				digit := int(clean[i] - '0')
				if i%2 == parity {
					digit *= 2
					if digit > 9 {
						digit -= 9
					}
				}
				sum += digit
			}
			if sum%10 != 0 {
				valid = false
			}
		}
		if !valid {
			err = fieldErr(field, "must be a valid credit card number")
		}
	}
	return validation(err, field, "creditcard")
}

// Latitude validates that a string is a valid latitude (-90 to 90).
func Latitude(v, field string) *Validation {
	var err error
	lat, parseErr := strconv.ParseFloat(v, 64)
	if parseErr != nil || lat < -90 || lat > 90 {
		err = fieldErr(field, "must be a valid latitude (-90 to 90)")
	}
	return validation(err, field, "latitude")
}

// Longitude validates that a string is a valid longitude (-180 to 180).
func Longitude(v, field string) *Validation {
	var err error
	lon, parseErr := strconv.ParseFloat(v, 64)
	if parseErr != nil || lon < -180 || lon > 180 {
		err = fieldErr(field, "must be a valid longitude (-180 to 180)")
	}
	return validation(err, field, "longitude")
}

// CountryCode2 validates that a string is a valid ISO 3166-1 alpha-2 country code.
func CountryCode2(v, field string) *Validation {
	var err error
	if len(v) != 2 {
		err = fieldErr(field, "must be a valid ISO 3166-1 alpha-2 country code")
	} else {
		for _, r := range v {
			if r < 'A' || r > 'Z' {
				err = fieldErr(field, "must be a valid ISO 3166-1 alpha-2 country code")
				break
			}
		}
	}
	return validation(err, field, "iso3166_1_alpha2")
}

// CountryCode3 validates that a string is a valid ISO 3166-1 alpha-3 country code.
func CountryCode3(v, field string) *Validation {
	var err error
	if len(v) != 3 {
		err = fieldErr(field, "must be a valid ISO 3166-1 alpha-3 country code")
	} else {
		for _, r := range v {
			if r < 'A' || r > 'Z' {
				err = fieldErr(field, "must be a valid ISO 3166-1 alpha-3 country code")
				break
			}
		}
	}
	return validation(err, field, "iso3166_1_alpha3")
}

// LanguageCode validates that a string is a valid ISO 639-1 language code.
func LanguageCode(v, field string) *Validation {
	var err error
	if len(v) != 2 {
		err = fieldErr(field, "must be a valid ISO 639-1 language code")
	} else {
		for _, r := range v {
			if r < 'a' || r > 'z' {
				err = fieldErr(field, "must be a valid ISO 639-1 language code")
				break
			}
		}
	}
	return validation(err, field, "iso639_1")
}

// CurrencyCode validates that a string is a valid ISO 4217 currency code.
func CurrencyCode(v, field string) *Validation {
	var err error
	if len(v) != 3 {
		err = fieldErr(field, "must be a valid ISO 4217 currency code")
	} else {
		for _, r := range v {
			if r < 'A' || r > 'Z' {
				err = fieldErr(field, "must be a valid ISO 4217 currency code")
				break
			}
		}
	}
	return validation(err, field, "iso4217")
}

// Hex validates that a string contains only hexadecimal characters.
func Hex(v, field string) *Validation {
	var err error
	for _, r := range v {
		if !isHexDigit(r) {
			err = fieldErr(field, "must be a valid hexadecimal string")
			break
		}
	}
	return validation(err, field, "hex")
}

func isHexDigit(r rune) bool {
	return (r >= '0' && r <= '9') || (r >= 'a' && r <= 'f') || (r >= 'A' && r <= 'F')
}

// DataURI validates that a string is a valid data URI.
func DataURI(v, field string) *Validation {
	var err error
	if !strings.HasPrefix(v, "data:") {
		err = fieldErr(field, "must be a valid data URI")
	} else {
		commaIdx := strings.Index(v, ",")
		if commaIdx == -1 {
			err = fieldErr(field, "must be a valid data URI")
		}
	}
	return validation(err, field, "datauri")
}

// FilePath validates that a string looks like a file path (contains path separators).
func FilePath(v, field string) *Validation {
	var err error
	if v == "" || strings.ContainsRune(v, 0) {
		err = fieldErr(field, "must be a valid file path")
	}
	return validation(err, field, "filepath")
}

// UnixPath validates that a string is a valid Unix-style path.
func UnixPath(v, field string) *Validation {
	var err error
	if v == "" || strings.ContainsRune(v, 0) {
		err = fieldErr(field, "must be a valid Unix path")
	}
	return validation(err, field, "unixpath")
}
