// Copyright (c) 2026 SafetyCulture Pty Ltd. All Rights Reserved.

using System.Text;

namespace S12.Protobuf.Validation.Conformance;

/// <summary>
/// The shared input corpus, in the order the vectors record it.
/// </summary>
/// <remarks>
/// A Go string holds bytes and a C# string holds UTF-16 code units, so a few inputs exist on one
/// side only. An input whose bytes are not valid UTF-8 has no C# spelling — decoding it substitutes
/// the replacement character, which is a different value — so it is carried as <c>null</c> and the
/// checks step over it rather than comparing something the two languages never agreed on.
/// </remarks>
internal sealed class Corpus
{
    private Corpus(IReadOnlyList<string?> values, IReadOnlyList<string> encoded)
    {
        Values = values;
        Encoded = encoded;
    }

    public IReadOnlyList<string?> Values { get; }

    /// <summary>Base64 of each input's bytes, which names a record in a report.</summary>
    public IReadOnlyList<string> Encoded { get; }

    public int Total => Values.Count;

    public IEnumerable<string> Unrepresentable =>
        Values.Select((value, i) => (value, i)).Where(p => p.value is null).Select(p => Encoded[p.i]);

    public static Corpus Read(string path)
    {
        var strict = new UTF8Encoding(encoderShouldEmitUTF8Identifier: false, throwOnInvalidBytes: true);
        var values = new List<string?>();
        var encoded = new List<string>();

        foreach (var line in File.ReadLines(path))
        {
            if (line.Length == 0 || line.StartsWith('#'))
            {
                continue;
            }

            var bytes = line.StartsWith("b64:", StringComparison.Ordinal)
                ? Convert.FromBase64String(line[4..])
                : Encoding.UTF8.GetBytes(line);

            encoded.Add(Convert.ToBase64String(bytes));
            try
            {
                values.Add(strict.GetString(bytes));
            }
            catch (DecoderFallbackException)
            {
                values.Add(null);
            }
        }

        if (values.Count == 0)
        {
            throw new InvalidOperationException($"{path} holds no inputs");
        }

        return new Corpus(values, encoded);
    }
}
