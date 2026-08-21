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

    /// <summary>Reports whether the value names a zone in the IANA time zone database.</summary>
    public static bool IsValidTimeZone(string value) =>
        value.Length != 0 && TimeZoneInfo.TryFindSystemTimeZoneById(value, out _);

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

        if (!Uri.TryCreate(value, UriKind.Absolute, out var uri) || HasHashInAuthority(value))
        {
            failure = "Invalid URL format";
            return false;
        }

        // Uri rewrites a single-slash authority such as https:/host into a host, leaving nothing for the
        // host check below to reject.
        if (!value.Contains("://", StringComparison.Ordinal) ||
            uri.Host.Length <= MinHostLength || uri.Host.StartsWith('.'))
        {
            failure = "Invalid host format";
            return false;
        }

        if (uri.Scheme.Length == 0)
        {
            failure = "Missing scheme";
            return false;
        }

        var schemeAllowed = false;
        foreach (var scheme in schemes)
        {
            if (string.Equals(uri.Scheme, scheme, StringComparison.Ordinal))
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
    /// Reports whether a hash sits in the authority, where it terminates the host and leaves the value
    /// unparsable. A hash after the authority opens a fragment and parses.
    /// </summary>
    private static bool HasHashInAuthority(string value)
    {
        var separator = value.IndexOf("://", StringComparison.Ordinal);
        if (separator < 0)
        {
            return false;
        }

        var start = separator + 3;
        var rest = value.AsSpan(start);
        var end = rest.IndexOfAny('/', '?');
        return (end < 0 ? rest : rest[..end]).Contains('#');
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
