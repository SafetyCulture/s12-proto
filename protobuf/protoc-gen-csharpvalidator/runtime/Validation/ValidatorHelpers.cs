// Copyright (c) 2026 SafetyCulture Pty Ltd. All Rights Reserved.

using System.Text;

namespace S12.Protobuf.Validation;

/// <summary>The format checks and measurements behind the field rules.</summary>
public static class ValidatorHelpers
{
    private const int MaxEmailByteLength = 321;
    private const int MinUsernameRuneCount = 2;
    private const int MaxUsernameRuneCount = 64;
    private const int MaxUrlRuneCount = 2048;
    private const int MinUrlByteLength = 3;
    private const int MinHostLength = 3;

    /// <summary>Length of the value in UTF-8 bytes.</summary>
    public static int Utf8Length(string value) => Encoding.UTF8.GetByteCount(value);

    /// <summary>Length of the value in Unicode code points.</summary>
    public static int RuneCount(string value)
    {
        var count = 0;
        for (var i = 0; i < value.Length; i++, count++)
        {
            if (char.IsHighSurrogate(value[i]) && i + 1 < value.Length && char.IsLowSurrogate(value[i + 1]))
            {
                i++;
            }
        }

        return count;
    }

    /// <summary>Reports whether the value is a UUID of version 3, 4 or 5.</summary>
    public static bool IsUuid(string value) => ValidatorPatterns.Uuid.IsMatch(value);

    /// <summary>Reports whether the value is a version 4 UUID.</summary>
    public static bool IsUuidV4(string value) => ValidatorPatterns.UuidV4.IsMatch(value);

    /// <summary>Reports whether the value is a legacy identifier, which carries no document prefix.</summary>
    public static bool IsLegacyId(string value, bool lowerCaseOnly) =>
        (lowerCaseOnly ? ValidatorPatterns.LegacyIdLowercase : ValidatorPatterns.LegacyId).IsMatch(value);

    /// <summary>Reports whether the value is a prefixed S12 identifier.</summary>
    public static bool IsS12Id(string value, bool lowerCaseOnly) =>
        (lowerCaseOnly ? ValidatorPatterns.S12IdLowercase : ValidatorPatterns.S12Id).IsMatch(value);

    /// <summary>Reports whether the value is an audit or template identifier with a long body.</summary>
    public static bool IsLongPrefixedLegacyId(string value) =>
        ValidatorPatterns.LongPrefixedLegacyId.IsMatch(value);

    /// <summary>Reports whether the value is an RFC 5322 address.</summary>
    public static bool IsValidEmail(string value) =>
        Utf8Length(value) <= MaxEmailByteLength && ValidatorPatterns.Email.IsMatch(value);

    /// <summary>
    /// Reports whether the value is a username: 2 to 64 code points of lowercase Latin letters and
    /// ASCII digits, starting and ending on a letter or digit, with dot, underscore or hyphen inside
    /// and no consecutive dots. Accented letters count only in composed form.
    /// </summary>
    public static bool IsValidNonEmail(string value)
    {
        var runes = RuneCount(value);
        if (runes < MinUsernameRuneCount || runes > MaxUsernameRuneCount)
        {
            return false;
        }

        if (!string.Equals(value, value.ToLowerInvariant(), StringComparison.Ordinal))
        {
            return false;
        }

        if (value.Contains("..", StringComparison.Ordinal))
        {
            return false;
        }

        return ValidatorPatterns.NonEmailUsername.IsMatch(value);
    }

    /// <summary>
    /// Reports whether every code unit takes part in a well-formed code point, which is what it
    /// means for the value to be encodable as UTF-8. A surrogate standing on its own is not.
    /// </summary>
    public static bool IsWellFormedUtf16(string value)
    {
        for (var i = 0; i < value.Length; i++)
        {
            if (!char.IsSurrogate(value[i]))
            {
                continue;
            }

            if (!char.IsHighSurrogate(value[i]) || i + 1 >= value.Length || !char.IsLowSurrogate(value[i + 1]))
            {
                return false;
            }

            i++;
        }

        return true;
    }

    /// <summary>Reports whether a complete URL appears anywhere in the value.</summary>
    public static bool ContainsUrl(string value) => ValidatorPatterns.RejectUrl.IsMatch(value);

    /// <summary>
    /// Reports whether the value names a zone in the IANA time zone database. The empty string and
    /// "Local" name a zone too, as UTC and the host's own zone.
    /// </summary>
    public static bool IsValidTimeZone(string value) =>
        value is "" or "UTC" or "Local" || TimeZoneInfo.TryFindSystemTimeZoneById(value, out _);

    /// <summary>
    /// Reports whether the value is an absolute URL on one of <paramref name="schemes"/>, and names the
    /// unmet requirement in <paramref name="failure"/> when it is not.
    /// </summary>
    public static bool TryValidateUrl(string value, ReadOnlySpan<string> schemes, bool allowFragment,
        out string? failure)
    {
        failure = null;

        if (value.Length == 0 || RuneCount(value) >= MaxUrlRuneCount ||
            Utf8Length(value) <= MinUrlByteLength || value.StartsWith('.'))
        {
            failure = "Length too long";
            return false;
        }

        if (!allowFragment && value.Contains('#'))
        {
            failure = "Fragment not allowed";
            return false;
        }

        if (!IsRequestUri(value))
        {
            failure = "Invalid URL format";
            return false;
        }

        // An absolute URL with no authority, such as https:/host, carries the host in its path and so
        // has none. A host of at most three characters cannot be a real one, e.g. a.nl is the shortest.
        var host = Authority(value);
        if (host.Length <= MinHostLength || host.StartsWith('.'))
        {
            failure = "Invalid host format";
            return false;
        }

        var scheme = Scheme(value);
        if (scheme.Length == 0)
        {
            failure = "Missing scheme";
            return false;
        }

        var schemeAllowed = false;
        foreach (var allowed in schemes)
        {
            if (string.Equals(scheme, allowed, StringComparison.Ordinal))
            {
                schemeAllowed = true;
                break;
            }
        }

        if (!schemeAllowed)
        {
            failure = "Invalid scheme";
            return false;
        }

        if (!ValidatorPatterns.Url.IsMatch(value))
        {
            failure = "Invalid URL characters";
            return false;
        }

        return true;
    }

    /// <summary>
    /// Reports whether the value has the shape of a URL received in a request: an absolute path, or a
    /// scheme followed by everything else, carrying no space, no control character, no malformed
    /// percent escape and no hash in the host.
    /// </summary>
    private static bool IsRequestUri(string value)
    {
        for (var i = 0; i < value.Length; i++)
        {
            if (value[i] <= ' ' || value[i] == '\u007F')
            {
                return false;
            }

            if (value[i] == '%' &&
                (i + 2 >= value.Length || !char.IsAsciiHexDigit(value[i + 1]) || !char.IsAsciiHexDigit(value[i + 2])))
            {
                return false;
            }
        }

        // A hash in the authority is a hash in the host name, which no host may carry.
        if (Authority(value).Contains('#'))
        {
            return false;
        }

        return value.StartsWith('/') || Scheme(value).Length != 0;
    }

    /// <summary>The scheme of an absolute URL, or the empty string when it declares none.</summary>
    private static string Scheme(string value)
    {
        var colon = value.IndexOf(':');
        if (colon <= 0 || !char.IsAsciiLetter(value[0]))
        {
            return "";
        }

        for (var i = 1; i < colon; i++)
        {
            if (!char.IsAsciiLetterOrDigit(value[i]) && value[i] is not ('+' or '-' or '.'))
            {
                return "";
            }
        }

        return value[..colon];
    }

    /// <summary>
    /// The host and port of an absolute URL, or the empty string when it declares no authority.
    /// Any userinfo ahead of the host is left out.
    /// </summary>
    private static string Authority(string value)
    {
        var separator = value.IndexOf("://", StringComparison.Ordinal);
        if (separator < 0)
        {
            return "";
        }

        var rest = value.AsSpan(separator + 3);
        var end = rest.IndexOfAny('/', '?');
        var authority = end < 0 ? rest : rest[..end];

        var userinfo = authority.LastIndexOf('@');
        return (userinfo < 0 ? authority : authority[(userinfo + 1)..]).ToString();
    }

    /// <summary>Base64 of at most <paramref name="maxLength"/> leading characters of the value.</summary>
    public static string TruncateAndEncode(string value, int maxLength)
    {
        var length = Math.Min(value.Length, maxLength);
        if (length > 0 && length < value.Length && char.IsHighSurrogate(value[length - 1]))
        {
            length--;
        }

        return Convert.ToBase64String(Encoding.UTF8.GetBytes(value[..length]));
    }
}
