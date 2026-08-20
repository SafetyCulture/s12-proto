// Copyright (c) 2026 SafetyCulture Pty Ltd. All Rights Reserved.

using System.Globalization;
using System.Text;
using System.Text.RegularExpressions;

namespace S12.Protobuf.Validation;

/// <summary>Unicode general categories, as a set.</summary>
[Flags]
public enum UnicodeCategories : uint
{
    /// <summary>No category.</summary>
    None = 0,

    /// <summary>UppercaseLetter.</summary>
    Lu = 1u << 0,

    /// <summary>LowercaseLetter.</summary>
    Ll = 1u << 1,

    /// <summary>TitlecaseLetter.</summary>
    Lt = 1u << 2,

    /// <summary>ModifierLetter.</summary>
    Lm = 1u << 3,

    /// <summary>OtherLetter.</summary>
    Lo = 1u << 4,

    /// <summary>NonSpacingMark.</summary>
    Mn = 1u << 5,

    /// <summary>SpacingCombiningMark.</summary>
    Mc = 1u << 6,

    /// <summary>EnclosingMark.</summary>
    Me = 1u << 7,

    /// <summary>DecimalDigitNumber.</summary>
    Nd = 1u << 8,

    /// <summary>LetterNumber.</summary>
    Nl = 1u << 9,

    /// <summary>OtherNumber.</summary>
    No = 1u << 10,

    /// <summary>SpaceSeparator.</summary>
    Zs = 1u << 11,

    /// <summary>LineSeparator.</summary>
    Zl = 1u << 12,

    /// <summary>ParagraphSeparator.</summary>
    Zp = 1u << 13,

    /// <summary>Control.</summary>
    Cc = 1u << 14,

    /// <summary>Format.</summary>
    Cf = 1u << 15,

    /// <summary>Surrogate.</summary>
    Cs = 1u << 16,

    /// <summary>PrivateUse.</summary>
    Co = 1u << 17,

    /// <summary>ConnectorPunctuation.</summary>
    Pc = 1u << 18,

    /// <summary>DashPunctuation.</summary>
    Pd = 1u << 19,

    /// <summary>OpenPunctuation.</summary>
    Ps = 1u << 20,

    /// <summary>ClosePunctuation.</summary>
    Pe = 1u << 21,

    /// <summary>InitialQuotePunctuation.</summary>
    Pi = 1u << 22,

    /// <summary>FinalQuotePunctuation.</summary>
    Pf = 1u << 23,

    /// <summary>OtherPunctuation.</summary>
    Po = 1u << 24,

    /// <summary>MathSymbol.</summary>
    Sm = 1u << 25,

    /// <summary>CurrencySymbol.</summary>
    Sc = 1u << 26,

    /// <summary>ModifierSymbol.</summary>
    Sk = 1u << 27,

    /// <summary>OtherSymbol.</summary>
    So = 1u << 28,

    /// <summary>OtherNotAssigned.</summary>
    Cn = 1u << 29,

    /// <summary>Every L category.</summary>
    L = Lu | Ll | Lt | Lm | Lo,

    /// <summary>Every M category.</summary>
    M = Mn | Mc | Me,

    /// <summary>Every N category.</summary>
    N = Nd | Nl | No,

    /// <summary>Every Z category.</summary>
    Z = Zs | Zl | Zp,

    /// <summary>Every C category.</summary>
    C = Cc | Cf | Cs | Co | Cn,

    /// <summary>Every P category.</summary>
    P = Pc | Pd | Ps | Pe | Pi | Pf | Po,

    /// <summary>Every S category.</summary>
    S = Sm | Sc | Sk | So,
}

/// <summary>Matches a value against the character class declared for a field.</summary>
public static class CharacterClass
{
    /// <summary>
    /// Reports whether every character of <paramref name="value"/> is permitted by
    /// <paramref name="pattern"/>, treating a code point above the Basic Multilingual Plane as
    /// permitted when its general category is in <paramref name="astralAllowed"/>.
    /// </summary>
    /// <param name="astralAllowed">
    /// The general categories the pattern admits, expanded from its <c>\p{…}</c> tokens by the
    /// generator. .NET matches over UTF-16 code units, so a pattern that admits a category still
    /// rejects a code point above the Basic Multilingual Plane in that category: each half of the
    /// surrogate pair carries category Surrogate rather than the category of the pair. Those code
    /// points are therefore removed before matching rather than matched.
    /// </param>
    public static bool IsMatch(Regex pattern, string value, UnicodeCategories astralAllowed)
    {
        if (astralAllowed == UnicodeCategories.None)
        {
            return pattern.IsMatch(value);
        }

        var remainder = RemoveAllowedAstral(value, astralAllowed);
        if (remainder.Length == 0 && value.Length != 0)
        {
            return true;
        }

        return pattern.IsMatch(remainder);
    }

    private static string RemoveAllowedAstral(string value, UnicodeCategories allowed)
    {
        StringBuilder? builder = null;
        for (var i = 0; i < value.Length; i++)
        {
            if (!char.IsHighSurrogate(value[i]) || i + 1 >= value.Length || !char.IsLowSurrogate(value[i + 1]))
            {
                builder?.Append(value[i]);
                continue;
            }

            if ((allowed & Of(CharUnicodeInfo.GetUnicodeCategory(value, i))) != UnicodeCategories.None)
            {
                builder ??= new StringBuilder(value.Length).Append(value.AsSpan(0, i));
                i++;
                continue;
            }

            builder?.Append(value[i]).Append(value[i + 1]);
            i++;
        }

        return builder?.ToString() ?? value;
    }

    private static UnicodeCategories Of(UnicodeCategory category) => (UnicodeCategories)(1u << (int)category);
}
