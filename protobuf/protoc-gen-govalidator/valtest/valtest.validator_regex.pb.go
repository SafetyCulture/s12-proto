// Code generated by protoc-gen-govalidator. DO NOT EDIT.
// versions:
// 	protoc-gen-govalidator v2.6.0
// 	protoc                 v5.28.0
// source: valtest.proto

package valtest

import (
	_ "github.com/SafetyCulture/s12-proto/s12/protobuf/proto"
	regexp "regexp"
)

// Pattern for ValTestMessage_Description
const _regex_val_d4db71516b8749dc594e5bf604c6a110 = `^[\pL\pN\x{0020}\x{0028}\x{0029}\x{002C}\x{002E}\x{003A}\x{003F}\x{0040}\x{005B}\x{005D}\x{005F}\x{00BF}\x{2013}]+$`

var _regex_d4db71516b8749dc594e5bf604c6a110 = regexp.MustCompile(_regex_val_d4db71516b8749dc594e5bf604c6a110)

// Pattern for ValTestMessage_Password
const _regex_val_51116fcfa477f1949f7055f2f1bf33db = `^[\pL\pN\x{0020}\x{0021}\x{0022}\x{0023}\x{0024}\x{0025}\x{0026}\x{0027}\x{0028}\x{0029}\x{002A}\x{002B}\x{002C}\x{002D}\x{002E}\x{002F}\x{003A}\x{003B}\x{003C}\x{003D}\x{003E}\x{003F}\x{0040}\x{005B}\x{005C}\x{005D}\x{005E}\x{005F}\x{0060}\x{007B}\x{007C}\x{007D}\x{007E}\x{00BF}\x{2013}]+$`

var _regex_51116fcfa477f1949f7055f2f1bf33db = regexp.MustCompile(_regex_val_51116fcfa477f1949f7055f2f1bf33db)

// ValTestMessage_Title is using regex d4db71516b8749dc594e5bf604c6a110
// ValTestMessage_FixedString is using regex d4db71516b8749dc594e5bf604c6a110
// Pattern for ValTestMessage_RuneString
const _regex_val_471b9980a6f67c9193bc7044ec96c4da = `^[\pL\pN\pSo\x{0020}\x{0028}\x{0029}\x{002C}\x{002E}\x{003A}\x{003F}\x{0040}\x{005B}\x{005D}\x{005F}\x{00BF}\x{200D}\x{2013}]+$`

var _regex_471b9980a6f67c9193bc7044ec96c4da = regexp.MustCompile(_regex_val_471b9980a6f67c9193bc7044ec96c4da)

// Pattern for ValTestMessage_ReplaceString
const _regex_val_7f420dacf785a9f7b630597663705292 = `^[\pL\pN\x{0020}\x{0028}\x{0029}\x{002C}\x{002E}\x{003A}\x{003F}\x{0040}\x{005B}\x{005D}\x{005F}\x{00BF}\x{02C2}\x{02C3}\x{037E}\x{2013}\x{2019}\x{201D}\x{2052}\x{2212}\x{2215}\x{2217}\x{2E40}\x{FF01}\x{FF0B}\x{FF3C}\x{FFE8}]+$`

var _regex_7f420dacf785a9f7b630597663705292 = regexp.MustCompile(_regex_val_7f420dacf785a9f7b630597663705292)

// ValTestMessage_NotReplaceString is using regex d4db71516b8749dc594e5bf604c6a110
// Pattern for ValTestMessage_AllowString
const _regex_val_2d8966526d95f557bf61ba41d6e869ca = `^[\pL\pN\x{0020}\x{0028}\x{0029}\x{002C}\x{002E}\x{003A}\x{003F}\x{0040}\x{005B}\x{005D}\x{005F}\x{007E}\x{00BF}\x{2013}]+$`

var _regex_2d8966526d95f557bf61ba41d6e869ca = regexp.MustCompile(_regex_val_2d8966526d95f557bf61ba41d6e869ca)

// Pattern for ValTestMessage_SymbolString
const _regex_val_cfc4f3f27c72b7eff33e9065de4993b8 = `^[\pL\pN\pSc\x{0020}\x{0028}\x{0029}\x{002C}\x{002E}\x{003A}\x{003F}\x{0040}\x{005B}\x{005D}\x{005F}\x{00BF}\x{2013}]+$`

var _regex_cfc4f3f27c72b7eff33e9065de4993b8 = regexp.MustCompile(_regex_val_cfc4f3f27c72b7eff33e9065de4993b8)

// Pattern for ValTestMessage_SymbolsString
const _regex_val_c5418abde4a1025792e46f9de3e163a8 = `^[\pL\pM\pN\pSo\x{0020}\x{0028}\x{0029}\x{002C}\x{002E}\x{003A}\x{003F}\x{0040}\x{005B}\x{005D}\x{005F}\x{00BF}\x{200D}\x{2013}]+$`

var _regex_c5418abde4a1025792e46f9de3e163a8 = regexp.MustCompile(_regex_val_c5418abde4a1025792e46f9de3e163a8)

// Pattern for ValTestMessage_NewlineString
const _regex_val_b87d0bf989a4ccd3a8f851bcdcfadd5b = `^[\pL\pN\x{000A}\x{0020}\x{0028}\x{0029}\x{002C}\x{002E}\x{003A}\x{003F}\x{0040}\x{005B}\x{005D}\x{005F}\x{00BF}\x{2013}]+$`

var _regex_b87d0bf989a4ccd3a8f851bcdcfadd5b = regexp.MustCompile(_regex_val_b87d0bf989a4ccd3a8f851bcdcfadd5b)

// Pattern for ValTestMessage_InvalidEncodingString
const _regex_val_721ec450fcf3c35a27a9e280064d4c50 = `^[\pL\pN\x{0020}\x{0028}\x{0029}\x{002C}\x{002E}\x{003A}\x{003F}\x{0040}\x{005B}\x{005D}\x{005F}\x{00BF}\x{2013}\x{FFFD}]+$`

var _regex_721ec450fcf3c35a27a9e280064d4c50 = regexp.MustCompile(_regex_val_721ec450fcf3c35a27a9e280064d4c50)

// ValTestMessage_OptString is using regex d4db71516b8749dc594e5bf604c6a110
// ValTestMessage_TrimString is using regex d4db71516b8749dc594e5bf604c6a110
// Pattern for ValTestMessage_AllString
const _regex_val_d834caddc03154e20f75eb8af178a1ac = `^[\pL\pM\pN\pPo\pSc\pSk\pSm\pSo\x{0020}\x{0028}\x{0029}\x{002C}\x{002E}\x{003A}\x{003F}\x{0040}\x{005B}\x{005D}\x{005F}\x{007E}\x{00BF}\x{02C3}\x{104B}\x{200D}\x{2013}\x{2018}\x{2019}\x{201C}\x{201D}\x{2022}\x{FFFD}}]+$`

var _regex_d834caddc03154e20f75eb8af178a1ac = regexp.MustCompile(_regex_val_d834caddc03154e20f75eb8af178a1ac)

// Pattern for ValTestMessage_Name
const _regex_val_d4a3ad03e6e76d647c361a8c16ac9395 = `^[\pL\pN\x{0020}\x{0028}\x{0029}\x{002C}\x{002E}\x{003A}\x{003F}\x{0040}\x{005B}\x{005D}\x{005F}\x{00BF}\x{2013}\x{2019}\x{2212}]+$`

var _regex_d4a3ad03e6e76d647c361a8c16ac9395 = regexp.MustCompile(_regex_val_d4a3ad03e6e76d647c361a8c16ac9395)

// ValTestMessage_ScTitle is using regex 7f420dacf785a9f7b630597663705292
// Pattern for ValTestMessage_ScPermissive
const _regex_val_b15160b6f45559b046acfd0ffd80eb84 = `^[\pL\pM\pN\pPo\pSc\pSk\pSm\pSo\x{0020}\x{0028}\x{0029}\x{002C}\x{002E}\x{003A}\x{003F}\x{0040}\x{005B}\x{005D}\x{005F}\x{007E}\x{00BF}\x{02C2}\x{02C3}\x{037E}\x{104B}\x{200D}\x{2013}\x{2018}\x{2019}\x{201C}\x{201D}\x{2022}\x{2052}\x{2212}\x{2215}\x{2217}\x{2E40}\x{FF01}\x{FF0B}\x{FF3C}\x{FFE8}}]+$`

var _regex_b15160b6f45559b046acfd0ffd80eb84 = regexp.MustCompile(_regex_val_b15160b6f45559b046acfd0ffd80eb84)

// ValTestMessage_NotSanitisePua is using regex d4db71516b8749dc594e5bf604c6a110
// ValTestMessage_SanitisePua is using regex d4db71516b8749dc594e5bf604c6a110
// ValTestMessage_SanitiseLength is using regex d4db71516b8749dc594e5bf604c6a110
// ValTestMessage_OptionalString is using regex d4db71516b8749dc594e5bf604c6a110
// ValTestMessage_Phone is using regex d4db71516b8749dc594e5bf604c6a110
// Pattern for ValTestMessage_LongString
const _regex_val_1bb1e0a2437db38577f49b4f31ccfca2 = `^[\pL\pN\x{0020}\x{0021}\x{0022}\x{0023}\x{0025}\x{0026}\x{0027}\x{0028}\x{0029}\x{002A}\x{002C}\x{002D}\x{002E}\x{002F}\x{003A}\x{003F}\x{0040}\x{005B}\x{005C}\x{005D}\x{005F}\x{00BF}\x{2013}]+$`

var _regex_1bb1e0a2437db38577f49b4f31ccfca2 = regexp.MustCompile(_regex_val_1bb1e0a2437db38577f49b4f31ccfca2)

// ValTestMessage_NestedMessage_Val is using regex d4db71516b8749dc594e5bf604c6a110
// ValTestMessage_NestedMessage_InnerNestedMessage_InnerVal is using regex d4db71516b8749dc594e5bf604c6a110
// Pattern for LogOnlyValidationMessage_Title
const _regex_val_7d80f76b8ba61834c54bce3793157935 = `^[\pL\pN\x{000A}\x{0020}\x{0021}\x{0022}\x{0023}\x{0025}\x{0026}\x{0027}\x{0028}\x{0029}\x{002A}\x{002C}\x{002D}\x{002E}\x{002F}\x{003A}\x{003F}\x{0040}\x{005B}\x{005C}\x{005D}\x{005F}\x{00BF}\x{2013}]+$`

var _regex_7d80f76b8ba61834c54bce3793157935 = regexp.MustCompile(_regex_val_7d80f76b8ba61834c54bce3793157935)

// Pattern for LogOnlyValidationMessage_Name
const _regex_val_b3f79e2470927c095fff6ea841e2a650 = `^[\pL\pM\pN\pPo\pSc\pSk\pSm\pSo\x{0020}\x{0021}\x{0022}\x{0023}\x{0025}\x{0026}\x{0027}\x{0028}\x{0029}\x{002A}\x{002C}\x{002D}\x{002E}\x{002F}\x{003A}\x{003F}\x{0040}\x{005B}\x{005C}\x{005D}\x{005F}\x{007E}\x{00BF}\x{104B}\x{200D}\x{2013}\x{2018}\x{2019}\x{201C}\x{201D}\x{2022}]+$`

var _regex_b3f79e2470927c095fff6ea841e2a650 = regexp.MustCompile(_regex_val_b3f79e2470927c095fff6ea841e2a650)

// NestedLevel3Message_OrgId5 is using regex d4db71516b8749dc594e5bf604c6a110
// NestedLevel2Message_OrgId4 is using regex d4db71516b8749dc594e5bf604c6a110
// NestedLevel1Message_OrgId3 is using regex d4db71516b8749dc594e5bf604c6a110
// MyReqMessage_UserId is using regex d4db71516b8749dc594e5bf604c6a110
// NonUrlMessage_RejectPartialUrlTest is using regex 1bb1e0a2437db38577f49b4f31ccfca2
// NonUrlMessage_BreakPartialUrlTest is using regex 1bb1e0a2437db38577f49b4f31ccfca2
