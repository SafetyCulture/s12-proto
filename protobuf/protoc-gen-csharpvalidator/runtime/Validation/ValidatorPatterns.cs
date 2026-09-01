// Copyright (c) 2026 SafetyCulture Pty Ltd. All Rights Reserved.

using System.Text.RegularExpressions;

namespace S12.Protobuf.Validation;

/// <summary>
/// The fixed patterns behind the format rules, translated from the RE2 sources in s12-proto.
/// </summary>
/// <remarks>
/// Translation from RE2 to the .NET dialect changes four things: <c>\x{N}</c> becomes <c>\uNNNN</c>,
/// <c>\d</c> becomes <c>[0-9]</c> because .NET reads it as every Unicode digit, <c>^</c> and <c>$</c>
/// become <c>\A</c> and <c>\z</c> because .NET <c>$</c> also matches before a trailing newline, and
/// <c>\p{Latin}</c> becomes an explicit block union because .NET has no Unicode script classes, and
/// an escaped word character such as <c>\_</c> loses its backslash because .NET rejects escapes it
/// does not define where RE2 allows any punctuation to be escaped.
/// Every pattern uses the non-backtracking engine, which runs in time linear in the input and
/// accepts exactly the constructs RE2 allows.
/// </remarks>
internal static class ValidatorPatterns
{
    private const RegexOptions Options = RegexOptions.NonBacktracking | RegexOptions.CultureInvariant;

    /// <summary>UUID version 3, 4 or 5, with or without hyphens.</summary>
    internal static readonly Regex Uuid = new(@"\A[0-9a-f]{8}-?[0-9a-f]{4}-?[0-9a-f]{4}-?[0-9a-f]{4}-?[0-9a-f]{12}\z", Options);

    /// <summary>UUID version 4.</summary>
    internal static readonly Regex UuidV4 = new(@"\A[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}\z", Options);

    /// <summary>Legacy identifier: a UUID optionally followed by a shard suffix.</summary>
    internal static readonly Regex LegacyId = new(@"(?i)\A[0-9a-f]{8}-?[0-9a-f]{4}-?[0-9a-f]{4}-?[0-9a-f]{4}-?[0-9a-f]{12}(-?[0-9a-f]{2,5}-?[0-9a-f]{1,16})?\z", Options);

    /// <summary>Legacy identifier, lowercase hexadecimal only.</summary>
    internal static readonly Regex LegacyIdLowercase = new(@"\A[0-9a-f]{8}-?[0-9a-f]{4}-?[0-9a-f]{4}-?[0-9a-f]{4}-?[0-9a-f]{12}(-?[0-9a-f]{2,5}-?[0-9a-f]{1,16})?\z", Options);

    /// <summary>Prefixed S12 identifier, for example <c>user_0123…</c>.</summary>
    internal static readonly Regex S12Id = new(@"\A(audit|template|action|user|ntfmsg|evidence|role|location|responseset|response|preference|heads_up|subscription|folder|scheduleitem)_([0-9a-fA-F]){32}\z", Options);

    /// <summary>Prefixed S12 identifier, lowercase hexadecimal only.</summary>
    internal static readonly Regex S12IdLowercase = new(@"\A(audit|template|action|user|ntfmsg|evidence|role|location|responseset|response|preference|heads_up|subscription|folder|scheduleitem)_([0-9a-f]){32}\z", Options);

    /// <summary>Audit or template identifier with a 50-53 character uppercase body.</summary>
    internal static readonly Regex LongPrefixedLegacyId = new(@"\A(audit|template)_([0-9A-F]){50,53}\z", Options);

    /// <summary>RFC 5322 address.</summary>
    internal static readonly Regex Email = new(@"\A((((([a-zA-Z]|[0-9]|[!#\$%&'\*\+\-\/=\?\^_`{\|}~]|[\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF])+(\.([a-zA-Z]|[0-9]|[!#\$%&'\*\+\-\/=\?\^_`{\|}~]|[\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF])+)*)|((\x22)((((\x20|\x09)*(\x0d\x0a))?(\x20|\x09)+)?(([\x01-\x08\x0b\x0c\x0e-\x1f\x7f]|\x21|[\x23-\x5b]|[\x5d-\x7e]|[\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF])|(\([\x01-\x09\x0b\x0c\x0d-\x7f]|[\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF]))))*(((\x20|\x09)*(\x0d\x0a))?(\x20|\x09)+)?(\x22)))){1,14}@((([a-zA-Z]|[0-9]|[\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF])|(([a-zA-Z]|[0-9]|[\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF])([a-zA-Z]|[0-9]|-|\.|_|~|[\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF])*([a-zA-Z]|[0-9]|[\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF])))\.)+(([a-zA-Z]|[\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF])|(([a-zA-Z]|[\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF])([a-zA-Z]|[0-9]|-|_|~|[\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF])*([a-zA-Z]|[\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF])))\z", Options);

    /// <summary>Username of Latin letters and ASCII digits, with dot, underscore or hyphen inside.</summary>
    internal static readonly Regex NonEmailUsername = new(@"\A[a-zA-Z\u00C0-\u00D6\u00D8-\u00F6\u00F8-\u024F\u1E00-\u1EFF0-9][a-zA-Z\u00C0-\u00D6\u00D8-\u00F6\u00F8-\u024F\u1E00-\u1EFF0-9._-]*[a-zA-Z\u00C0-\u00D6\u00D8-\u00F6\u00F8-\u024F\u1E00-\u1EFF0-9]\z", Options);

    /// <summary>Characters permitted in a URL. Does not check structure.</summary>
    internal static readonly Regex Url = new(@"\A[@\:\/\?#\.\-_\%\;\=\~\&\+a-zA-Z0-9]+\z", Options);

    /// <summary>A complete URL appearing anywhere in the value.</summary>
    internal static readonly Regex RejectUrl = new(@"https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{2,256}\.[a-z]{2,4}([-a-zA-Z0-9@:%_\+.~#?&//=]*)", Options);

    /// <summary>Private-use code point in the Basic Multilingual Plane.</summary>
    internal static readonly Regex PrivateUseArea = new(@"[\uE000-\uF8FF]", Options);

}
