// Code generated by protoc-gen-govalidator. DO NOT EDIT.
// versions:
// 	protoc-gen-govalidator v2.3.0
// 	protoc                 v3.17.3
// source: valtest.proto

package valtest

import (
	fmt "fmt"
	proto "github.com/SafetyCulture/s12-proto/s12/protobuf/proto"
	transform "golang.org/x/text/transform"
	norm "golang.org/x/text/unicode/norm"
	strings "strings"
	utf8 "unicode/utf8"
)

func (m *ValTestMessage) Validate() error {
	if !proto.IsUUIDv4(m.Id) {
		isValidId := false
		if !isValidId {
			return fmt.Errorf(`id: value must be parsable as UUIDv4`)
		}
	}
	for _, item := range m.Ids {
		if !proto.IsUUIDv4(item) {
			isValidId := false
			if !isValidId {
				return fmt.Errorf(`ids: value must be parsable as UUIDv4`)
			}
		}
	}
	if m.MediaId != "" {
		if !proto.IsUUIDv4(m.MediaId) {
			isValidId := false
			if !isValidId {
				return fmt.Errorf(`media_id: value must be parsable as UUIDv4`)
			}
		}
	}
	if !proto.IsUUIDv4(m.LegacyId) {
		isValidId := false
		if proto.IsLegacyID(m.LegacyId) {
			isValidId = true
		}
		if !isValidId {
			return fmt.Errorf(`legacy_id: value must be parsable as UUIDv4 or legacy ID`)
		}
	}
	if m.InnerLegacyId != nil {
		if v, ok := interface{}(m.InnerLegacyId).(proto.Validator); ok {
			if err := v.Validate(); err != nil {
				return proto.FieldError("inner_legacy_id", err)
			}
		}
	}
	if !proto.IsValidEmail(m.Email, false) {
		return fmt.Errorf(`email: value must be parsable as an email address`)
	}
	if m.OptEmail != "" {
		if !proto.IsValidEmail(m.OptEmail, false) {
			return fmt.Errorf(`opt_email: value must be parsable as an email address`)
		}
	}
	if !norm.NFC.IsNormalString(m.Description) && norm.NFD.IsNormalString(m.Description) {
		// normalise NFD to NFC string
		var normErr error
		m.Description, _, normErr = transform.String(transform.Chain(norm.NFD, norm.NFC), m.Description)
		if normErr != nil {
			return fmt.Errorf(`description: value must must be normalisable to NFC`)
		}
	}
	if strings.ContainsRune(m.Description, utf8.RuneError) {
		return fmt.Errorf(`description: value must must have valid encoding`)
	} else if !utf8.ValidString(m.Description) {
		return fmt.Errorf(`description: value must must be a valid UTF-8-encoded string`)
	}
	var _len_ValTestMessage_Description = len(m.Description)
	if !(_len_ValTestMessage_Description >= 1 && _len_ValTestMessage_Description <= 750) {
		return fmt.Errorf(`description: value must have length between 1 and 750`)
	}
	if !_regex_d4db71516b8749dc594e5bf604c6a110.MatchString(m.Description) {
		return fmt.Errorf(`description: value must only have valid characters`)
	}
	if !norm.NFC.IsNormalString(m.Password) && norm.NFD.IsNormalString(m.Password) {
		// normalise NFD to NFC string
		var normErr error
		m.Password, _, normErr = transform.String(transform.Chain(norm.NFD, norm.NFC), m.Password)
		if normErr != nil {
			return fmt.Errorf(`password: value must must be normalisable to NFC`)
		}
	}
	if strings.ContainsRune(m.Password, utf8.RuneError) {
		return fmt.Errorf(`password: value must must have valid encoding`)
	} else if !utf8.ValidString(m.Password) {
		return fmt.Errorf(`password: value must must be a valid UTF-8-encoded string`)
	}
	var _len_ValTestMessage_Password = len(m.Password)
	if !(_len_ValTestMessage_Password >= 8 && _len_ValTestMessage_Password <= 130) {
		return fmt.Errorf(`password: value must have length between 8 and 130`)
	}
	if !_regex_51116fcfa477f1949f7055f2f1bf33db.MatchString(m.Password) {
		return fmt.Errorf(`password: value must only have valid characters`)
	}
	if !norm.NFC.IsNormalString(m.Title) && norm.NFD.IsNormalString(m.Title) {
		// normalise NFD to NFC string
		var normErr error
		m.Title, _, normErr = transform.String(transform.Chain(norm.NFD, norm.NFC), m.Title)
		if normErr != nil {
			return fmt.Errorf(`title: value must must be normalisable to NFC`)
		}
	}
	if strings.ContainsRune(m.Title, utf8.RuneError) {
		return fmt.Errorf(`title: value must must have valid encoding`)
	} else if !utf8.ValidString(m.Title) {
		return fmt.Errorf(`title: value must must be a valid UTF-8-encoded string`)
	}
	var _len_ValTestMessage_Title = len(m.Title)
	if !(_len_ValTestMessage_Title >= 3 && _len_ValTestMessage_Title <= 50) {
		return fmt.Errorf(`title: value must have length between 3 and 50`)
	}
	if !_regex_d4db71516b8749dc594e5bf604c6a110.MatchString(m.Title) {
		return fmt.Errorf(`title: value must only have valid characters`)
	}
	if !norm.NFC.IsNormalString(m.FixedString) && norm.NFD.IsNormalString(m.FixedString) {
		// normalise NFD to NFC string
		var normErr error
		m.FixedString, _, normErr = transform.String(transform.Chain(norm.NFD, norm.NFC), m.FixedString)
		if normErr != nil {
			return fmt.Errorf(`fixed_string: value must must be normalisable to NFC`)
		}
	}
	if strings.ContainsRune(m.FixedString, utf8.RuneError) {
		return fmt.Errorf(`fixed_string: value must must have valid encoding`)
	} else if !utf8.ValidString(m.FixedString) {
		return fmt.Errorf(`fixed_string: value must must be a valid UTF-8-encoded string`)
	}
	var _len_ValTestMessage_FixedString = len(m.FixedString)
	if !(_len_ValTestMessage_FixedString == 4) {
		return fmt.Errorf(`fixed_string: value must have length 4`)
	}
	if !_regex_d4db71516b8749dc594e5bf604c6a110.MatchString(m.FixedString) {
		return fmt.Errorf(`fixed_string: value must only have valid characters`)
	}
	if !norm.NFC.IsNormalString(m.RuneString) && norm.NFD.IsNormalString(m.RuneString) {
		// normalise NFD to NFC string
		var normErr error
		m.RuneString, _, normErr = transform.String(transform.Chain(norm.NFD, norm.NFC), m.RuneString)
		if normErr != nil {
			return fmt.Errorf(`rune_string: value must must be normalisable to NFC`)
		}
	}
	if strings.ContainsRune(m.RuneString, utf8.RuneError) {
		return fmt.Errorf(`rune_string: value must must have valid encoding`)
	} else if !utf8.ValidString(m.RuneString) {
		return fmt.Errorf(`rune_string: value must must be a valid UTF-8-encoded string`)
	}
	var _len_ValTestMessage_RuneString = utf8.RuneCountInString(m.RuneString)
	if !(_len_ValTestMessage_RuneString == 4) {
		return fmt.Errorf(`rune_string: value must have length 4`)
	}
	if !_regex_471b9980a6f67c9193bc7044ec96c4da.MatchString(m.RuneString) {
		return fmt.Errorf(`rune_string: value must only have valid characters`)
	}
	if !norm.NFC.IsNormalString(m.ReplaceString) && norm.NFD.IsNormalString(m.ReplaceString) {
		// normalise NFD to NFC string
		var normErr error
		m.ReplaceString, _, normErr = transform.String(transform.Chain(norm.NFD, norm.NFC), m.ReplaceString)
		if normErr != nil {
			return fmt.Errorf(`replace_string: value must must be normalisable to NFC`)
		}
	}
	if strings.ContainsRune(m.ReplaceString, utf8.RuneError) {
		return fmt.Errorf(`replace_string: value must must have valid encoding`)
	} else if !utf8.ValidString(m.ReplaceString) {
		return fmt.Errorf(`replace_string: value must must be a valid UTF-8-encoded string`)
	}
	m.ReplaceString = proto.UnsafeCharReplacer.Replace(m.ReplaceString)
	var _len_ValTestMessage_ReplaceString = len(m.ReplaceString)
	if !(_len_ValTestMessage_ReplaceString >= 1 && _len_ValTestMessage_ReplaceString <= 130) {
		return fmt.Errorf(`replace_string: value must have length between 1 and 130`)
	}
	if !_regex_7f420dacf785a9f7b630597663705292.MatchString(m.ReplaceString) {
		return fmt.Errorf(`replace_string: value must only have valid characters`)
	}
	if !norm.NFC.IsNormalString(m.NotReplaceString) && norm.NFD.IsNormalString(m.NotReplaceString) {
		// normalise NFD to NFC string
		var normErr error
		m.NotReplaceString, _, normErr = transform.String(transform.Chain(norm.NFD, norm.NFC), m.NotReplaceString)
		if normErr != nil {
			return fmt.Errorf(`not_replace_string: value must must be normalisable to NFC`)
		}
	}
	if strings.ContainsRune(m.NotReplaceString, utf8.RuneError) {
		return fmt.Errorf(`not_replace_string: value must must have valid encoding`)
	} else if !utf8.ValidString(m.NotReplaceString) {
		return fmt.Errorf(`not_replace_string: value must must be a valid UTF-8-encoded string`)
	}
	var _len_ValTestMessage_NotReplaceString = len(m.NotReplaceString)
	if !(_len_ValTestMessage_NotReplaceString >= 1 && _len_ValTestMessage_NotReplaceString <= 130) {
		return fmt.Errorf(`not_replace_string: value must have length between 1 and 130`)
	}
	if !_regex_d4db71516b8749dc594e5bf604c6a110.MatchString(m.NotReplaceString) {
		return fmt.Errorf(`not_replace_string: value must only have valid characters`)
	}
	if !norm.NFC.IsNormalString(m.AllowString) && norm.NFD.IsNormalString(m.AllowString) {
		// normalise NFD to NFC string
		var normErr error
		m.AllowString, _, normErr = transform.String(transform.Chain(norm.NFD, norm.NFC), m.AllowString)
		if normErr != nil {
			return fmt.Errorf(`allow_string: value must must be normalisable to NFC`)
		}
	}
	if strings.ContainsRune(m.AllowString, utf8.RuneError) {
		return fmt.Errorf(`allow_string: value must must have valid encoding`)
	} else if !utf8.ValidString(m.AllowString) {
		return fmt.Errorf(`allow_string: value must must be a valid UTF-8-encoded string`)
	}
	m.AllowString = strings.ReplaceAll(m.AllowString, "\u0023", "\u0020")
	var _len_ValTestMessage_AllowString = len(m.AllowString)
	if !(_len_ValTestMessage_AllowString >= 1 && _len_ValTestMessage_AllowString <= 130) {
		return fmt.Errorf(`allow_string: value must have length between 1 and 130`)
	}
	if !_regex_2d8966526d95f557bf61ba41d6e869ca.MatchString(m.AllowString) {
		return fmt.Errorf(`allow_string: value must only have valid characters`)
	}
	if !norm.NFC.IsNormalString(m.SymbolString) && norm.NFD.IsNormalString(m.SymbolString) {
		// normalise NFD to NFC string
		var normErr error
		m.SymbolString, _, normErr = transform.String(transform.Chain(norm.NFD, norm.NFC), m.SymbolString)
		if normErr != nil {
			return fmt.Errorf(`symbol_string: value must must be normalisable to NFC`)
		}
	}
	if strings.ContainsRune(m.SymbolString, utf8.RuneError) {
		return fmt.Errorf(`symbol_string: value must must have valid encoding`)
	} else if !utf8.ValidString(m.SymbolString) {
		return fmt.Errorf(`symbol_string: value must must be a valid UTF-8-encoded string`)
	}
	var _len_ValTestMessage_SymbolString = len(m.SymbolString)
	if !(_len_ValTestMessage_SymbolString >= 1 && _len_ValTestMessage_SymbolString <= 130) {
		return fmt.Errorf(`symbol_string: value must have length between 1 and 130`)
	}
	if !_regex_cfc4f3f27c72b7eff33e9065de4993b8.MatchString(m.SymbolString) {
		return fmt.Errorf(`symbol_string: value must only have valid characters`)
	}
	if !norm.NFC.IsNormalString(m.SymbolsString) && norm.NFD.IsNormalString(m.SymbolsString) {
		// normalise NFD to NFC string
		var normErr error
		m.SymbolsString, _, normErr = transform.String(transform.Chain(norm.NFD, norm.NFC), m.SymbolsString)
		if normErr != nil {
			return fmt.Errorf(`symbols_string: value must must be normalisable to NFC`)
		}
	}
	if strings.ContainsRune(m.SymbolsString, utf8.RuneError) {
		return fmt.Errorf(`symbols_string: value must must have valid encoding`)
	} else if !utf8.ValidString(m.SymbolsString) {
		return fmt.Errorf(`symbols_string: value must must be a valid UTF-8-encoded string`)
	}
	var _len_ValTestMessage_SymbolsString = len(m.SymbolsString)
	if !(_len_ValTestMessage_SymbolsString >= 1 && _len_ValTestMessage_SymbolsString <= 130) {
		return fmt.Errorf(`symbols_string: value must have length between 1 and 130`)
	}
	if !_regex_c5418abde4a1025792e46f9de3e163a8.MatchString(m.SymbolsString) {
		return fmt.Errorf(`symbols_string: value must only have valid characters`)
	}
	if !norm.NFC.IsNormalString(m.NewlineString) && norm.NFD.IsNormalString(m.NewlineString) {
		// normalise NFD to NFC string
		var normErr error
		m.NewlineString, _, normErr = transform.String(transform.Chain(norm.NFD, norm.NFC), m.NewlineString)
		if normErr != nil {
			return fmt.Errorf(`newline_string: value must must be normalisable to NFC`)
		}
	}
	if strings.ContainsRune(m.NewlineString, utf8.RuneError) {
		return fmt.Errorf(`newline_string: value must must have valid encoding`)
	} else if !utf8.ValidString(m.NewlineString) {
		return fmt.Errorf(`newline_string: value must must be a valid UTF-8-encoded string`)
	}
	m.NewlineString = strings.ReplaceAll(m.NewlineString, "\r", "")
	var _len_ValTestMessage_NewlineString = len(m.NewlineString)
	if !(_len_ValTestMessage_NewlineString >= 1 && _len_ValTestMessage_NewlineString <= 130) {
		return fmt.Errorf(`newline_string: value must have length between 1 and 130`)
	}
	if !_regex_b87d0bf989a4ccd3a8f851bcdcfadd5b.MatchString(m.NewlineString) {
		return fmt.Errorf(`newline_string: value must only have valid characters`)
	}
	if !norm.NFC.IsNormalString(m.InvalidEncodingString) && norm.NFD.IsNormalString(m.InvalidEncodingString) {
		// normalise NFD to NFC string
		var normErr error
		m.InvalidEncodingString, _, normErr = transform.String(transform.Chain(norm.NFD, norm.NFC), m.InvalidEncodingString)
		if normErr != nil {
			return fmt.Errorf(`invalid_encoding_string: value must must be normalisable to NFC`)
		}
	}
	var _len_ValTestMessage_InvalidEncodingString = len(m.InvalidEncodingString)
	if !(_len_ValTestMessage_InvalidEncodingString >= 1 && _len_ValTestMessage_InvalidEncodingString <= 130) {
		return fmt.Errorf(`invalid_encoding_string: value must have length between 1 and 130`)
	}
	if !_regex_721ec450fcf3c35a27a9e280064d4c50.MatchString(m.InvalidEncodingString) {
		return fmt.Errorf(`invalid_encoding_string: value must only have valid characters`)
	}
	if m.OptString != "" {
		if !norm.NFC.IsNormalString(m.OptString) && norm.NFD.IsNormalString(m.OptString) {
			// normalise NFD to NFC string
			var normErr error
			m.OptString, _, normErr = transform.String(transform.Chain(norm.NFD, norm.NFC), m.OptString)
			if normErr != nil {
				return fmt.Errorf(`opt_string: value must must be normalisable to NFC`)
			}
		}
		if strings.ContainsRune(m.OptString, utf8.RuneError) {
			return fmt.Errorf(`opt_string: value must must have valid encoding`)
		} else if !utf8.ValidString(m.OptString) {
			return fmt.Errorf(`opt_string: value must must be a valid UTF-8-encoded string`)
		}
		var _len_ValTestMessage_OptString = len(m.OptString)
		if !(_len_ValTestMessage_OptString >= 1 && _len_ValTestMessage_OptString <= 130) {
			return fmt.Errorf(`opt_string: value must have length between 1 and 130`)
		}
		if !_regex_d4db71516b8749dc594e5bf604c6a110.MatchString(m.OptString) {
			return fmt.Errorf(`opt_string: value must only have valid characters`)
		}
	}
	if !norm.NFC.IsNormalString(m.TrimString) && norm.NFD.IsNormalString(m.TrimString) {
		// normalise NFD to NFC string
		var normErr error
		m.TrimString, _, normErr = transform.String(transform.Chain(norm.NFD, norm.NFC), m.TrimString)
		if normErr != nil {
			return fmt.Errorf(`trim_string: value must must be normalisable to NFC`)
		}
	}
	if strings.ContainsRune(m.TrimString, utf8.RuneError) {
		return fmt.Errorf(`trim_string: value must must have valid encoding`)
	} else if !utf8.ValidString(m.TrimString) {
		return fmt.Errorf(`trim_string: value must must be a valid UTF-8-encoded string`)
	}
	m.TrimString = strings.TrimSpace(m.TrimString)
	var _len_ValTestMessage_TrimString = len(m.TrimString)
	if !(_len_ValTestMessage_TrimString >= 1 && _len_ValTestMessage_TrimString <= 130) {
		return fmt.Errorf(`trim_string: value must have length between 1 and 130`)
	}
	if !_regex_d4db71516b8749dc594e5bf604c6a110.MatchString(m.TrimString) {
		return fmt.Errorf(`trim_string: value must only have valid characters`)
	}
	if !norm.NFC.IsNormalString(m.AllString) && norm.NFD.IsNormalString(m.AllString) {
		// normalise NFD to NFC string
		var normErr error
		m.AllString, _, normErr = transform.String(transform.Chain(norm.NFD, norm.NFC), m.AllString)
		if normErr != nil {
			return fmt.Errorf(`all_string: value must must be normalisable to NFC`)
		}
	}
	m.AllString = strings.ReplaceAll(m.AllString, "\u003E", "\u02C3")
	m.AllString = proto.SymbolCharReplacer.Replace(m.AllString)
	var _len_ValTestMessage_AllString = len(m.AllString)
	if !(_len_ValTestMessage_AllString >= 1 && _len_ValTestMessage_AllString <= 130) {
		return fmt.Errorf(`all_string: value must have length between 1 and 130`)
	}
	if !_regex_d834caddc03154e20f75eb8af178a1ac.MatchString(m.AllString) {
		return fmt.Errorf(`all_string: value must only have valid characters`)
	}
	if m.Name != "" {
		if !norm.NFC.IsNormalString(m.Name) && norm.NFD.IsNormalString(m.Name) {
			// normalise NFD to NFC string
			var normErr error
			m.Name, _, normErr = transform.String(transform.Chain(norm.NFD, norm.NFC), m.Name)
			if normErr != nil {
				return fmt.Errorf(`name: value must must be normalisable to NFC`)
			}
		}
		if strings.ContainsRune(m.Name, utf8.RuneError) {
			return fmt.Errorf(`name: value must must have valid encoding`)
		} else if !utf8.ValidString(m.Name) {
			return fmt.Errorf(`name: value must must be a valid UTF-8-encoded string`)
		}
		m.Name = strings.ReplaceAll(m.Name, "\u0027", "\u2019")
		m.Name = strings.ReplaceAll(m.Name, "\u002D", "\u2212")
		var _len_ValTestMessage_Name = len(m.Name)
		if !(_len_ValTestMessage_Name >= 1 && _len_ValTestMessage_Name <= 50) {
			return fmt.Errorf(`name: value must have length between 1 and 50`)
		}
		if !_regex_d4a3ad03e6e76d647c361a8c16ac9395.MatchString(m.Name) {
			return fmt.Errorf(`name: value must only have valid characters`)
		}
	}
	if m.ScTitle != "" {
		if !norm.NFC.IsNormalString(m.ScTitle) && norm.NFD.IsNormalString(m.ScTitle) {
			// normalise NFD to NFC string
			var normErr error
			m.ScTitle, _, normErr = transform.String(transform.Chain(norm.NFD, norm.NFC), m.ScTitle)
			if normErr != nil {
				return fmt.Errorf(`sc_title: value must must be normalisable to NFC`)
			}
		}
		if strings.ContainsRune(m.ScTitle, utf8.RuneError) {
			return fmt.Errorf(`sc_title: value must must have valid encoding`)
		} else if !utf8.ValidString(m.ScTitle) {
			return fmt.Errorf(`sc_title: value must must be a valid UTF-8-encoded string`)
		}
		m.ScTitle = proto.UnsafeCharReplacer.Replace(m.ScTitle)
		var _len_ValTestMessage_ScTitle = len(m.ScTitle)
		if !(_len_ValTestMessage_ScTitle >= 1 && _len_ValTestMessage_ScTitle <= 500) {
			return fmt.Errorf(`sc_title: value must have length between 1 and 500`)
		}
		if !_regex_7f420dacf785a9f7b630597663705292.MatchString(m.ScTitle) {
			return fmt.Errorf(`sc_title: value must only have valid characters`)
		}
	}
	if m.ScPermissive != "" {
		if !norm.NFC.IsNormalString(m.ScPermissive) && norm.NFD.IsNormalString(m.ScPermissive) {
			// normalise NFD to NFC string
			var normErr error
			m.ScPermissive, _, normErr = transform.String(transform.Chain(norm.NFD, norm.NFC), m.ScPermissive)
			if normErr != nil {
				return fmt.Errorf(`sc_permissive: value must must be normalisable to NFC`)
			}
		}
		if strings.ContainsRune(m.ScPermissive, utf8.RuneError) {
			return fmt.Errorf(`sc_permissive: value must must have valid encoding`)
		} else if !utf8.ValidString(m.ScPermissive) {
			return fmt.Errorf(`sc_permissive: value must must be a valid UTF-8-encoded string`)
		}
		m.ScPermissive = proto.UnsafeCharReplacer.Replace(m.ScPermissive)
		m.ScPermissive = proto.SymbolCharReplacer.Replace(m.ScPermissive)
		m.ScPermissive = proto.RegexPua.ReplaceAllString(m.ScPermissive, "")
		var _len_ValTestMessage_ScPermissive = len(m.ScPermissive)
		if !(_len_ValTestMessage_ScPermissive >= 1 && _len_ValTestMessage_ScPermissive <= 1000) {
			return fmt.Errorf(`sc_permissive: value must have length between 1 and 1000`)
		}
		if !_regex_b15160b6f45559b046acfd0ffd80eb84.MatchString(m.ScPermissive) {
			return fmt.Errorf(`sc_permissive: value must only have valid characters`)
		}
	}
	if m.NotSanitisePua != "" {
		if !norm.NFC.IsNormalString(m.NotSanitisePua) && norm.NFD.IsNormalString(m.NotSanitisePua) {
			// normalise NFD to NFC string
			var normErr error
			m.NotSanitisePua, _, normErr = transform.String(transform.Chain(norm.NFD, norm.NFC), m.NotSanitisePua)
			if normErr != nil {
				return fmt.Errorf(`not_sanitise_pua: value must must be normalisable to NFC`)
			}
		}
		if strings.ContainsRune(m.NotSanitisePua, utf8.RuneError) {
			return fmt.Errorf(`not_sanitise_pua: value must must have valid encoding`)
		} else if !utf8.ValidString(m.NotSanitisePua) {
			return fmt.Errorf(`not_sanitise_pua: value must must be a valid UTF-8-encoded string`)
		}
		var _len_ValTestMessage_NotSanitisePua = len(m.NotSanitisePua)
		if !(_len_ValTestMessage_NotSanitisePua >= 1 && _len_ValTestMessage_NotSanitisePua <= 130) {
			return fmt.Errorf(`not_sanitise_pua: value must have length between 1 and 130`)
		}
		if !_regex_d4db71516b8749dc594e5bf604c6a110.MatchString(m.NotSanitisePua) {
			return fmt.Errorf(`not_sanitise_pua: value must only have valid characters`)
		}
	}
	if m.SanitisePua != "" {
		if !norm.NFC.IsNormalString(m.SanitisePua) && norm.NFD.IsNormalString(m.SanitisePua) {
			// normalise NFD to NFC string
			var normErr error
			m.SanitisePua, _, normErr = transform.String(transform.Chain(norm.NFD, norm.NFC), m.SanitisePua)
			if normErr != nil {
				return fmt.Errorf(`sanitise_pua: value must must be normalisable to NFC`)
			}
		}
		if strings.ContainsRune(m.SanitisePua, utf8.RuneError) {
			return fmt.Errorf(`sanitise_pua: value must must have valid encoding`)
		} else if !utf8.ValidString(m.SanitisePua) {
			return fmt.Errorf(`sanitise_pua: value must must be a valid UTF-8-encoded string`)
		}
		m.SanitisePua = proto.RegexPua.ReplaceAllString(m.SanitisePua, "")
		var _len_ValTestMessage_SanitisePua = len(m.SanitisePua)
		if !(_len_ValTestMessage_SanitisePua >= 1 && _len_ValTestMessage_SanitisePua <= 130) {
			return fmt.Errorf(`sanitise_pua: value must have length between 1 and 130`)
		}
		if !_regex_d4db71516b8749dc594e5bf604c6a110.MatchString(m.SanitisePua) {
			return fmt.Errorf(`sanitise_pua: value must only have valid characters`)
		}
	}
	if m.SanitiseLength != "" {
		if !norm.NFC.IsNormalString(m.SanitiseLength) && norm.NFD.IsNormalString(m.SanitiseLength) {
			// normalise NFD to NFC string
			var normErr error
			m.SanitiseLength, _, normErr = transform.String(transform.Chain(norm.NFD, norm.NFC), m.SanitiseLength)
			if normErr != nil {
				return fmt.Errorf(`sanitise_length: value must must be normalisable to NFC`)
			}
		}
		if strings.ContainsRune(m.SanitiseLength, utf8.RuneError) {
			return fmt.Errorf(`sanitise_length: value must must have valid encoding`)
		} else if !utf8.ValidString(m.SanitiseLength) {
			return fmt.Errorf(`sanitise_length: value must must be a valid UTF-8-encoded string`)
		}
		m.SanitiseLength = proto.RegexPua.ReplaceAllString(m.SanitiseLength, "")
		var _len_ValTestMessage_SanitiseLength = len(m.SanitiseLength)
		if !(_len_ValTestMessage_SanitiseLength == 2) {
			return fmt.Errorf(`sanitise_length: value must have length 2`)
		}
		if !_regex_d4db71516b8749dc594e5bf604c6a110.MatchString(m.SanitiseLength) {
			return fmt.Errorf(`sanitise_length: value must only have valid characters`)
		}
	}
	if x, ok := m.ContactOneof.(*ValTestMessage_Phone); ok {
		if !norm.NFC.IsNormalString(x.Phone) && norm.NFD.IsNormalString(x.Phone) {
			// normalise NFD to NFC string
			var normErr error
			x.Phone, _, normErr = transform.String(transform.Chain(norm.NFD, norm.NFC), x.Phone)
			if normErr != nil {
				return fmt.Errorf(`phone: value must must be normalisable to NFC`)
			}
		}
		if strings.ContainsRune(x.Phone, utf8.RuneError) {
			return fmt.Errorf(`phone: value must must have valid encoding`)
		} else if !utf8.ValidString(x.Phone) {
			return fmt.Errorf(`phone: value must must be a valid UTF-8-encoded string`)
		}
		var _len_ValTestMessage_Phone = len(x.Phone)
		if !(_len_ValTestMessage_Phone == 11) {
			return fmt.Errorf(`phone: value must have length 11`)
		}
		if !_regex_d4db71516b8749dc594e5bf604c6a110.MatchString(x.Phone) {
			return fmt.Errorf(`phone: value must only have valid characters`)
		}
	}
	if m.MsgRequired == nil {
		return fmt.Errorf("field msg_required is required")
	}
	if m.MsgRequired != nil {
		if v, ok := interface{}(m.MsgRequired).(proto.Validator); ok {
			if err := v.Validate(); err != nil {
				return proto.FieldError("msg_required", err)
			}
		}
	}
	if m.NestedMessage != nil {
		if v, ok := interface{}(m.NestedMessage).(proto.Validator); ok {
			if err := v.Validate(); err != nil {
				return proto.FieldError("nested_message", err)
			}
		}
	}
	// Validation of proto3 map<> fields is unsupported.
	if !(len(m.ContactsWithLengthConstraint) >= 1) {
		return fmt.Errorf(`contacts_with_length_constraint: length must be greater than or equal to 1`)
	}
	if !(len(m.ContactsWithLengthConstraint) <= 10) {
		return fmt.Errorf(`contacts_with_length_constraint: length must be lesser than or equal to 10`)
	}
	if len(m.ContactsWithLengthConstraint) > 0 {
		for _, item := range m.ContactsWithLengthConstraint {
			if item != nil {
				if v, ok := interface{}(item).(proto.Validator); ok {
					if err := v.Validate(); err != nil {
						return proto.FieldError("contacts_with_length_constraint", err)
					}
				}
			}
		}
	}
	if len(m.ContactsWithoutLengthConstraint) > 0 {
		for _, item := range m.ContactsWithoutLengthConstraint {
			if item != nil {
				if v, ok := interface{}(item).(proto.Validator); ok {
					if err := v.Validate(); err != nil {
						return proto.FieldError("contacts_without_length_constraint", err)
					}
				}
			}
		}
	}
	if !proto.IsUUIDv4(m.S12Id) {
		isValidId := false
		if !isValidId && proto.IsS12ID(m.S12Id) {
			isValidId = true
		}
		if !isValidId {
			return fmt.Errorf(`s12_id: value must be parsable as UUIDv4 or S12 ID`)
		}
	}
	if m.InnerS12Id != nil {
		if v, ok := interface{}(m.InnerS12Id).(proto.Validator); ok {
			if err := v.Validate(); err != nil {
				return proto.FieldError("inner_s12_id", err)
			}
		}
	}
	if m.AllId != "" {
		if !proto.IsUUIDv4(m.AllId) {
			isValidId := false
			if proto.IsLegacyID(m.AllId) {
				isValidId = true
			}
			if !isValidId && proto.IsS12ID(m.AllId) {
				isValidId = true
			}
			if !isValidId {
				return fmt.Errorf(`all_id: value must be parsable as UUIDv4 or legacy ID or S12 ID`)
			}
		}
	}
	_schemes_ValTestMessage_Url := []string{"https"}
	if _, err := proto.IsValidURL(m.Url, _schemes_ValTestMessage_Url, false); err != nil {
		return fmt.Errorf(`url: value must be parsable as a URL: %v (%v)`, err, m.Url)
	}
	if m.UrlAllOpts != "" {
		_schemes_ValTestMessage_UrlAllOpts := []string{"ftp", "ftps", "http"}
		if _, err := proto.IsValidURL(m.UrlAllOpts, _schemes_ValTestMessage_UrlAllOpts, true); err != nil {
			return fmt.Errorf(`url_all_opts: value must be parsable as a URL: %v (%v)`, err, m.UrlAllOpts)
		}
	}
	return nil
}

func (m *ValTestMessage_NestedMessage) Validate() error {
	if !norm.NFC.IsNormalString(m.Val) && norm.NFD.IsNormalString(m.Val) {
		// normalise NFD to NFC string
		var normErr error
		m.Val, _, normErr = transform.String(transform.Chain(norm.NFD, norm.NFC), m.Val)
		if normErr != nil {
			return fmt.Errorf(`val: value must must be normalisable to NFC`)
		}
	}
	if strings.ContainsRune(m.Val, utf8.RuneError) {
		return fmt.Errorf(`val: value must must have valid encoding`)
	} else if !utf8.ValidString(m.Val) {
		return fmt.Errorf(`val: value must must be a valid UTF-8-encoded string`)
	}
	var _len_ValTestMessage_NestedMessage_Val = len(m.Val)
	if !(_len_ValTestMessage_NestedMessage_Val >= 1 && _len_ValTestMessage_NestedMessage_Val <= 100) {
		return fmt.Errorf(`val: value must have length between 1 and 100`)
	}
	if !_regex_d4db71516b8749dc594e5bf604c6a110.MatchString(m.Val) {
		return fmt.Errorf(`val: value must only have valid characters`)
	}
	if !proto.IsValidEmail(m.NestedEmail, false) {
		return fmt.Errorf(`nested_email: value must be parsable as an email address`)
	}
	if !(len(m.MemberEmails) >= 2) {
		return fmt.Errorf(`member_emails: length must be greater than or equal to 2`)
	}
	if !(len(m.MemberEmails) <= 5) {
		return fmt.Errorf(`member_emails: length must be lesser than or equal to 5`)
	}
	for _, item := range m.MemberEmails {
		if !proto.IsValidEmail(item, false) {
			return fmt.Errorf(`member_emails: value must be parsable as an email address`)
		}
	}
	return nil
}

func (m *ValTestMessage_NestedMessage_InnerNestedMessage) Validate() error {
	if !norm.NFC.IsNormalString(m.InnerVal) && norm.NFD.IsNormalString(m.InnerVal) {
		// normalise NFD to NFC string
		var normErr error
		m.InnerVal, _, normErr = transform.String(transform.Chain(norm.NFD, norm.NFC), m.InnerVal)
		if normErr != nil {
			return fmt.Errorf(`inner_val: value must must be normalisable to NFC`)
		}
	}
	if strings.ContainsRune(m.InnerVal, utf8.RuneError) {
		return fmt.Errorf(`inner_val: value must must have valid encoding`)
	} else if !utf8.ValidString(m.InnerVal) {
		return fmt.Errorf(`inner_val: value must must be a valid UTF-8-encoded string`)
	}
	var _len_ValTestMessage_NestedMessage_InnerNestedMessage_InnerVal = len(m.InnerVal)
	if !(_len_ValTestMessage_NestedMessage_InnerNestedMessage_InnerVal >= 1 && _len_ValTestMessage_NestedMessage_InnerNestedMessage_InnerVal <= 100) {
		return fmt.Errorf(`inner_val: value must have length between 1 and 100`)
	}
	if !_regex_d4db71516b8749dc594e5bf604c6a110.MatchString(m.InnerVal) {
		return fmt.Errorf(`inner_val: value must only have valid characters`)
	}
	return nil
}

func (m *ValTestMessage_Contact) Validate() error {
	if m.Phone != "" {
	}
	if !proto.IsValidEmail(m.Email, false) {
		return fmt.Errorf(`email: value must be parsable as an email address`)
	}
	return nil
}

func (m *OuterMessageUsingNestedMessage) Validate() error {
	if m.SomeMessage != nil {
		if v, ok := interface{}(m.SomeMessage).(proto.Validator); ok {
			if err := v.Validate(); err != nil {
				return proto.FieldError("some_message", err)
			}
		}
	}
	return nil
}

func (m *InnerMessage) Validate() error {
	if !proto.IsUUID(m.Id) {
		return fmt.Errorf(`id: value must be parsable as a UUID`)
	}
	return nil
}

func (m *InnerMessageWithLegacyId) Validate() error {
	if !proto.IsUUIDv4(m.Id) {
		isValidId := false
		if proto.IsLegacyID(m.Id) {
			isValidId = true
		}
		if !isValidId {
			return fmt.Errorf(`id: value must be parsable as UUIDv4 or legacy ID`)
		}
	}
	return nil
}

func (m *InnerMessageWithS12Id) Validate() error {
	if !proto.IsUUIDv4(m.Id) {
		isValidId := false
		if !isValidId && proto.IsS12ID(m.Id) {
			isValidId = true
		}
		if !isValidId {
			return fmt.Errorf(`id: value must be parsable as UUIDv4 or S12 ID`)
		}
	}
	return nil
}

func (m *NestedLevel3Message) Validate() error {
	if !norm.NFC.IsNormalString(m.OrgId5) && norm.NFD.IsNormalString(m.OrgId5) {
		// normalise NFD to NFC string
		var normErr error
		m.OrgId5, _, normErr = transform.String(transform.Chain(norm.NFD, norm.NFC), m.OrgId5)
		if normErr != nil {
			return fmt.Errorf(`org_id5: value must must be normalisable to NFC`)
		}
	}
	if strings.ContainsRune(m.OrgId5, utf8.RuneError) {
		return fmt.Errorf(`org_id5: value must must have valid encoding`)
	} else if !utf8.ValidString(m.OrgId5) {
		return fmt.Errorf(`org_id5: value must must be a valid UTF-8-encoded string`)
	}
	var _len_NestedLevel3Message_OrgId5 = len(m.OrgId5)
	if !(_len_NestedLevel3Message_OrgId5 == 5) {
		return fmt.Errorf(`org_id5: value must have length 5`)
	}
	if !_regex_d4db71516b8749dc594e5bf604c6a110.MatchString(m.OrgId5) {
		return fmt.Errorf(`org_id5: value must only have valid characters`)
	}
	return nil
}

func (m *NestedLevel2Message) Validate() error {
	if !norm.NFC.IsNormalString(m.OrgId4) && norm.NFD.IsNormalString(m.OrgId4) {
		// normalise NFD to NFC string
		var normErr error
		m.OrgId4, _, normErr = transform.String(transform.Chain(norm.NFD, norm.NFC), m.OrgId4)
		if normErr != nil {
			return fmt.Errorf(`org_id4: value must must be normalisable to NFC`)
		}
	}
	if strings.ContainsRune(m.OrgId4, utf8.RuneError) {
		return fmt.Errorf(`org_id4: value must must have valid encoding`)
	} else if !utf8.ValidString(m.OrgId4) {
		return fmt.Errorf(`org_id4: value must must be a valid UTF-8-encoded string`)
	}
	var _len_NestedLevel2Message_OrgId4 = len(m.OrgId4)
	if !(_len_NestedLevel2Message_OrgId4 == 4) {
		return fmt.Errorf(`org_id4: value must have length 4`)
	}
	if !_regex_d4db71516b8749dc594e5bf604c6a110.MatchString(m.OrgId4) {
		return fmt.Errorf(`org_id4: value must only have valid characters`)
	}
	if m.OrgNested != nil {
		if v, ok := interface{}(m.OrgNested).(proto.Validator); ok {
			if err := v.Validate(); err != nil {
				return proto.FieldError("org_nested", err)
			}
		}
	}
	return nil
}

func (m *NestedLevel1Message) Validate() error {
	if !norm.NFC.IsNormalString(m.OrgId3) && norm.NFD.IsNormalString(m.OrgId3) {
		// normalise NFD to NFC string
		var normErr error
		m.OrgId3, _, normErr = transform.String(transform.Chain(norm.NFD, norm.NFC), m.OrgId3)
		if normErr != nil {
			return fmt.Errorf(`org_id3: value must must be normalisable to NFC`)
		}
	}
	if strings.ContainsRune(m.OrgId3, utf8.RuneError) {
		return fmt.Errorf(`org_id3: value must must have valid encoding`)
	} else if !utf8.ValidString(m.OrgId3) {
		return fmt.Errorf(`org_id3: value must must be a valid UTF-8-encoded string`)
	}
	var _len_NestedLevel1Message_OrgId3 = len(m.OrgId3)
	if !(_len_NestedLevel1Message_OrgId3 == 3) {
		return fmt.Errorf(`org_id3: value must have length 3`)
	}
	if !_regex_d4db71516b8749dc594e5bf604c6a110.MatchString(m.OrgId3) {
		return fmt.Errorf(`org_id3: value must only have valid characters`)
	}
	if m.OrgNested != nil {
		if v, ok := interface{}(m.OrgNested).(proto.Validator); ok {
			if err := v.Validate(); err != nil {
				return proto.FieldError("org_nested", err)
			}
		}
	}
	return nil
}

func (m *MyReqMessage) Validate() error {
	if !norm.NFC.IsNormalString(m.UserId) && norm.NFD.IsNormalString(m.UserId) {
		// normalise NFD to NFC string
		var normErr error
		m.UserId, _, normErr = transform.String(transform.Chain(norm.NFD, norm.NFC), m.UserId)
		if normErr != nil {
			return fmt.Errorf(`user_id: value must must be normalisable to NFC`)
		}
	}
	if strings.ContainsRune(m.UserId, utf8.RuneError) {
		return fmt.Errorf(`user_id: value must must have valid encoding`)
	} else if !utf8.ValidString(m.UserId) {
		return fmt.Errorf(`user_id: value must must be a valid UTF-8-encoded string`)
	}
	var _len_MyReqMessage_UserId = len(m.UserId)
	if !(_len_MyReqMessage_UserId == 2) {
		return fmt.Errorf(`user_id: value must have length 2`)
	}
	if !_regex_d4db71516b8749dc594e5bf604c6a110.MatchString(m.UserId) {
		return fmt.Errorf(`user_id: value must only have valid characters`)
	}
	if m.OrgNested != nil {
		if v, ok := interface{}(m.OrgNested).(proto.Validator); ok {
			if err := v.Validate(); err != nil {
				return proto.FieldError("org_nested", err)
			}
		}
	}
	return nil
}

func (m *ScimEmail) Validate() error {
	if !proto.IsValidEmail(m.Value, false) {
		return fmt.Errorf(`value: value must be parsable as an email address`)
	}
	return nil
}

func (m *ScimUser) Validate() error {
	if len(m.Emails) > 0 {
		for _, item := range m.Emails {
			if item != nil {
				if v, ok := interface{}(item).(proto.Validator); ok {
					if err := v.Validate(); err != nil {
						return proto.FieldError("emails", err)
					}
				}
			}
		}
	}
	return nil
}

func (m *MyMessageWithEnum) Validate() error {
	if int(m.Enum) == 0 {
		return fmt.Errorf("field enum must be specified and a non-zero value")
	}
	for _, item := range m.Enums {
		if int(item) == 0 {
			return fmt.Errorf("field enums must be specified and a non-zero value")
		}
	}
	return nil
}

func (m *MyMessageWithRepeatedEnum) Validate() error {
	for _, item := range m.Enums {
		if int(item) == 0 {
			return fmt.Errorf("field enums must be specified and a non-zero value")
		}
	}
	return nil
}

func (m *MyMessageWithRepeatedField) Validate() error {
	if !(len(m.MyInt) <= 5) {
		return fmt.Errorf(`my_int: length must be lesser than or equal to 5`)
	}
	return nil
}

func (m *MyOneOfMsg) Validate() error {
	if x, ok := m.MyField.(*MyOneOfMsg_MyFirstField); ok {
		if x.MyFirstField != nil {
			if v, ok := interface{}(x.MyFirstField).(proto.Validator); ok {
				if err := v.Validate(); err != nil {
					return proto.FieldError("my_first_field", err)
				}
			}
		}
	}
	if x, ok := m.MyField.(*MyOneOfMsg_MySecondField); ok {
		if x.MySecondField != nil {
			if v, ok := interface{}(x.MySecondField).(proto.Validator); ok {
				if err := v.Validate(); err != nil {
					return proto.FieldError("my_second_field", err)
				}
			}
		}
	}
	if x, ok := m.MyField.(*MyOneOfMsg_MyThirdField); ok {
		if !(len(x.MyThirdField) >= 1) {
			return fmt.Errorf(`my_third_field: value must have length greater than or equal to 1`)
		}
	}
	return nil
}

func (m *MyOneOfMsg_FirstType) Validate() error {
	if !(m.Value > 1) {
		return fmt.Errorf(`value: value must be greater than 1`)
	}
	return nil
}

func (m *MyOneOfMsg_SecondType) Validate() error {
	if !(m.Value > 2) {
		return fmt.Errorf(`value: value must be greater than 2`)
	}
	return nil
}
