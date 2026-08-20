// Copyright (c) 2026 SafetyCulture Pty Ltd. All Rights Reserved.

using System.Buffers;
using System.Collections.Frozen;
using System.Text;

namespace S12.Protobuf.Validation;

/// <summary>The rewrites applied to a field value before its length and character class are checked.</summary>
public static class StringMutators
{
    /// <summary>
    /// Converts a fully decomposed value to composed form. Returns <c>false</c> when the value cannot
    /// be normalised at all, such as when it carries an unpaired surrogate, leaving
    /// <paramref name="normalized"/> unchanged.
    /// </summary>
    public static bool TryNormalizeToNfc(string value, out string normalized)
    {
        normalized = value;
        try
        {
            if (value.IsNormalized(NormalizationForm.FormC) || !value.IsNormalized(NormalizationForm.FormD))
            {
                return true;
            }

            normalized = value.Normalize(NormalizationForm.FormC);
            return true;
        }
        catch (ArgumentException)
        {
            return false;
        }
    }

    /// <summary>Reports whether the value carries the Unicode replacement character.</summary>
    public static bool ContainsReplacementCharacter(string value) => value.Contains('\uFFFD');

    /// <summary>Removes leading and trailing Unicode whitespace.</summary>
    public static string TrimSpace(string value) => value.Trim();

    /// <summary>Removes Basic Multilingual Plane private-use code points.</summary>
    public static string StripPrivateUseArea(string value) =>
        ValidatorPatterns.PrivateUseArea.Replace(value, string.Empty);

    /// <summary>Substitutes look-alikes for the characters that carry injection risk.</summary>
    public static string ReplaceUnsafeCharacters(string value) => Substitute(value, UnsafeTargets, UnsafeMap);

    /// <summary>Folds exotic spacing, dashes and fullwidth punctuation to their plain equivalents.</summary>
    public static string ReplaceSymbolCharacters(string value) => Substitute(value, SymbolTargets, SymbolMap);

    /// <summary>
    /// Folds exotic spacing, dashes and fullwidth punctuation to their plain equivalents, leaving line
    /// feeds in place.
    /// </summary>
    public static string ReplaceSymbolCharactersMultiline(string value) =>
        Substitute(value, SymbolMultilineTargets, SymbolMultilineMap);

    /// <summary>
    /// Inserts a space after a full stop that directly joins two words, so the value does not read as a
    /// link once rendered.
    /// </summary>
    public static string BreakPartialUrls(string value)
    {
        if (value.Length < 2 || value.IndexOf('.', 1) < 0)
        {
            return value;
        }

        var builder = new StringBuilder(value.Length + 8);
        for (var i = 0; i < value.Length; i++)
        {
            if (value[i] == '.' && i > 0 && IsAsciiWord(value[i - 1]) &&
                i + 1 < value.Length && !IsBreakStop(value[i + 1]))
            {
                var next = i + 1;
                var width = char.IsHighSurrogate(value[next]) && next + 1 < value.Length &&
                            char.IsLowSurrogate(value[next + 1])
                    ? 2
                    : 1;
                builder.Append(". ").Append(value.AsSpan(next, width));
                i += width;
                continue;
            }

            builder.Append(value[i]);
        }

        return builder.ToString();
    }

    private static bool IsAsciiWord(char c) =>
        c is >= '0' and <= '9' or >= 'A' and <= 'Z' or >= 'a' and <= 'z' or '_';

    private static bool IsBreakStop(char c) =>
        c is '\t' or '\n' or '\f' or '\r' or ' ' or (>= '0' and <= '9');

    private static string Substitute(string value, SearchValues<char> targets, FrozenDictionary<char, string> map)
    {
        var first = value.AsSpan().IndexOfAny(targets);
        if (first < 0)
        {
            return value;
        }

        var builder = new StringBuilder(value.Length);
        builder.Append(value.AsSpan(0, first));
        for (var i = first; i < value.Length; i++)
        {
            if (map.TryGetValue(value[i], out var replacement))
            {
                builder.Append(replacement);
            }
            else
            {
                builder.Append(value[i]);
            }
        }

        return builder.ToString();
    }

    private static SearchValues<char> TargetsOf((char From, string To)[] pairs) =>
        SearchValues.Create(Array.ConvertAll(pairs, p => p.From));

    private static FrozenDictionary<char, string> MapOf((char From, string To)[] pairs) =>
        pairs.ToFrozenDictionary(p => p.From, p => p.To);

    private static readonly (char From, string To)[] UnsafePairs =
    [
        ('\u0021', "\uFF01"),
        ('\u0022', "\u201D"),
        ('\u0023', "\u0020"),
        ('\u0025', "\u2052"),
        ('\u0026', "\u0020"),
        ('\u0027', "\u2019"),
        ('\u002A', "\u2217"),
        ('\u002B', "\uFF0B"),
        ('\u002D', "\u2212"),
        ('\u002F', "\u2215"),
        ('\u003B', "\u037E"),
        ('\u003C', "\u02C2"),
        ('\u003D', "\u2E40"),
        ('\u003E', "\u02C3"),
        ('\u005C', "\uFF3C"),
        ('\u0060', "\u2019"),
        ('\u00B4', "\u2019"),
        ('\u007C', "\uFFE8"),
        ('\u3164', "\u0020"),
    ];

    private static readonly (char From, string To)[] SymbolPairs =
    [
        ('\u00A0', "\u0020"),
        ('\u1680', "\u0020"),
        ('\u2000', "\u0020"),
        ('\u2001', "\u0020"),
        ('\u2002', "\u0020"),
        ('\u2003', "\u0020"),
        ('\u2004', "\u0020"),
        ('\u2005', "\u0020"),
        ('\u2006', "\u0020"),
        ('\u2007', "\u0020"),
        ('\u2008', "\u0020"),
        ('\u2009', "\u0020"),
        ('\u200A', "\u0020"),
        ('\u202F', "\u0020"),
        ('\u205F', "\u0020"),
        ('\u3000', "\u0020"),
        ('\u200B', string.Empty),
        ('\u200C', string.Empty),
        ('\u200D', string.Empty),
        ('\uFEFF', string.Empty),
        ('\u2014', "\u2013"),
        ('\u2018', "\u2019"),
        ('\u3002', "\u002E"),
        ('\uFF0C', "\u002C"),
        ('\uFF1A', "\u003A"),
        ('\u0009', "\u0020"),
        ('\u000D', "\u0020"),
        ('\u000A', "\u0020"),
    ];

    private static readonly (char From, string To)[] SymbolMultilinePairs =
    [
        ('\u00A0', "\u0020"),
        ('\u1680', "\u0020"),
        ('\u2000', "\u0020"),
        ('\u2001', "\u0020"),
        ('\u2002', "\u0020"),
        ('\u2003', "\u0020"),
        ('\u2004', "\u0020"),
        ('\u2005', "\u0020"),
        ('\u2006', "\u0020"),
        ('\u2007', "\u0020"),
        ('\u2008', "\u0020"),
        ('\u2009', "\u0020"),
        ('\u200A', "\u0020"),
        ('\u202F', "\u0020"),
        ('\u205F', "\u0020"),
        ('\u3000', "\u0020"),
        ('\u200B', string.Empty),
        ('\u200C', string.Empty),
        ('\u200D', string.Empty),
        ('\uFEFF', string.Empty),
        ('\u2014', "\u2013"),
        ('\u2018', "\u2019"),
        ('\u3002', "\u002E"),
        ('\uFF0C', "\u002C"),
        ('\uFF1A', "\u003A"),
        ('\u0009', "\u0020"),
        ('\u000D', "\u0020"),
    ];

    private static readonly SearchValues<char> UnsafeTargets = TargetsOf(UnsafePairs);
    private static readonly SearchValues<char> SymbolTargets = TargetsOf(SymbolPairs);
    private static readonly SearchValues<char> SymbolMultilineTargets = TargetsOf(SymbolMultilinePairs);

    private static readonly FrozenDictionary<char, string> UnsafeMap = MapOf(UnsafePairs);
    private static readonly FrozenDictionary<char, string> SymbolMap = MapOf(SymbolPairs);
    private static readonly FrozenDictionary<char, string> SymbolMultilineMap = MapOf(SymbolMultilinePairs);
}
