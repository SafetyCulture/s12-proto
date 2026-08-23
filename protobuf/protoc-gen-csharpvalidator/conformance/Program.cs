// Copyright (c) 2026 SafetyCulture Pty Ltd. All Rights Reserved.

using System.Collections;
using System.Reflection;
using System.Text;
using System.Text.Json;
using Google.Protobuf;
using Google.Protobuf.Reflection;
using S12.Protobuf.Validation;
using Valtest;

namespace S12.Protobuf.Validation.Conformance;

/// <summary>
/// Runs the generated validators and the runtime over the committed vectors and reports every
/// result that differs from the one Go recorded.
/// </summary>
internal static class Program
{
    private static int Main(string[] args)
    {
        var directory = args.Length > 0 ? args[0] : DefaultDataDirectory();
        if (!Directory.Exists(directory))
        {
            Console.Error.WriteLine($"no vector directory at {directory}");
            return 2;
        }

        var report = new Report();
        var corpus = Corpus.Read(Path.Combine(directory, "inputs.txt"));

        CheckTables(directory, report);
        CheckHelpers(directory, corpus, report);
        CheckFields(directory, corpus, report);
        CheckNestedFields(directory, corpus, report);
        CheckDeepNesting(directory, corpus, report);

        return report.Summarise(corpus);
    }

    /// <summary>
    /// Checks the character substitutions each replacer makes against the ones Go makes, found the
    /// same way: by running every code point through it. This is what says the hand-written tables
    /// in the runtime are complete.
    /// </summary>
    private static void CheckTables(string directory, Report report)
    {
        var replacers = new (string Name, Func<string, string> Apply)[]
        {
            ("unsafe", StringMutators.ReplaceUnsafeCharacters),
            ("other", StringMutators.ReplaceSymbolCharacters),
            ("otherMultiline", StringMutators.ReplaceSymbolCharactersMultiline),
        };

        var produced = new List<(string Replacer, string From, string To)>();
        foreach (var (name, apply) in replacers)
        {
            for (var codePoint = 0; codePoint <= 0x10FFFF; codePoint++)
            {
                if (codePoint is >= 0xD800 and <= 0xDFFF)
                {
                    continue; // a surrogate half is not a code point on its own
                }

                var input = char.ConvertFromUtf32(codePoint);
                var output = apply(input);
                if (output != input)
                {
                    produced.Add((name, codePoint.ToString("X4"), Base64(output)));
                }
            }
        }

        var expected = Vectors.Read(Path.Combine(directory, "tables.jsonl"));
        if (produced.Count != expected.Count)
        {
            report.Add("tables.jsonl", "-",
                $"{expected.Count} substitutions", $"{produced.Count} substitutions");
            return;
        }

        for (var i = 0; i < expected.Count; i++)
        {
            var want = expected[i];
            var (replacer, from, to) = produced[i];
            var key = $"{replacer} U+{from}";
            report.Compare("tables.jsonl", key, "replacer", want.String("replacer"), replacer);
            report.Compare("tables.jsonl", key, "from", want.String("from"), from);
            report.Compare("tables.jsonl", key, "to", want.String("to"), to);
        }

        report.Counted("tables.jsonl", produced.Count);
    }

    /// <summary>
    /// Checks each helper directly, rather than only through the fields that happen to call it.
    /// </summary>
    private static void CheckHelpers(string directory, Corpus corpus, Report report)
    {
        var expected = Vectors.Read(Path.Combine(directory, "helpers.jsonl"));
        if (expected.Count != corpus.Total)
        {
            report.Add("helpers.jsonl", "-", $"{corpus.Total} records", $"{expected.Count} records");
            return;
        }

        string[] https = ["https"];
        string[] multi = ["https", "http"];
        var compared = 0;

        for (var i = 0; i < corpus.Total; i++)
        {
            if (corpus.Values[i] is not { } value)
            {
                continue;
            }

            var want = expected[i];
            var key = corpus.Encoded[i];

            report.Compare("helpers.jsonl", key, "in", want.String("in"), Base64(value));
            report.Compare("helpers.jsonl", key, "uuid", want.Bool("uuid"), ValidatorHelpers.IsUuid(value));
            report.Compare("helpers.jsonl", key, "uuidV4", want.Bool("uuidV4"), ValidatorHelpers.IsUuidV4(value));
            report.Compare("helpers.jsonl", key, "legacyId", want.Bool("legacyId"), ValidatorHelpers.IsLegacyId(value, false));
            report.Compare("helpers.jsonl", key, "legacyIdLower", want.Bool("legacyIdLower"), ValidatorHelpers.IsLegacyId(value, true));
            report.Compare("helpers.jsonl", key, "s12Id", want.Bool("s12Id"), ValidatorHelpers.IsS12Id(value, false));
            report.Compare("helpers.jsonl", key, "s12IdLower", want.Bool("s12IdLower"), ValidatorHelpers.IsS12Id(value, true));
            report.Compare("helpers.jsonl", key, "longPrefixedLegacyId", want.Bool("longPrefixedLegacyId"), ValidatorHelpers.IsLongPrefixedLegacyId(value));
            report.Compare("helpers.jsonl", key, "email", want.Bool("email"), ValidatorHelpers.IsValidEmail(value));
            report.Compare("helpers.jsonl", key, "nonEmail", want.Bool("nonEmail"), ValidatorHelpers.IsValidNonEmail(value));
            report.Compare("helpers.jsonl", key, "containsUrl", want.Bool("containsUrl"), ValidatorHelpers.ContainsUrl(value));
            report.Compare("helpers.jsonl", key, "breakPartialUrl", want.String("breakPartialUrl"), Base64(StringMutators.BreakPartialUrls(value)));
            report.Compare("helpers.jsonl", key, "stripPua", want.String("stripPua"), Base64(StringMutators.StripPrivateUseArea(value)));
            report.Compare("helpers.jsonl", key, "replaceUnsafe", want.String("replaceUnsafe"), Base64(StringMutators.ReplaceUnsafeCharacters(value)));
            report.Compare("helpers.jsonl", key, "replaceOther", want.String("replaceOther"), Base64(StringMutators.ReplaceSymbolCharacters(value)));
            report.Compare("helpers.jsonl", key, "replaceOtherMultiline", want.String("replaceOtherMultiline"), Base64(StringMutators.ReplaceSymbolCharactersMultiline(value)));

            var okHttps = ValidatorHelpers.TryValidateUrl(value, https, false, out var httpsFailure);
            report.Compare("helpers.jsonl", key, "urlHttps", want.Bool("urlHttps"), okHttps);
            report.Compare("helpers.jsonl", key, "urlHttpsErr", want.String("urlHttpsErr"), httpsFailure ?? "");

            var okMulti = ValidatorHelpers.TryValidateUrl(value, multi, true, out var multiFailure);
            report.Compare("helpers.jsonl", key, "urlMulti", want.Bool("urlMulti"), okMulti);
            report.Compare("helpers.jsonl", key, "urlMultiErr", want.String("urlMultiErr"), multiFailure ?? "");

            compared++;
        }

        report.Counted("helpers.jsonl", compared);
    }

    /// <summary>
    /// Checks the generated validators, one record per string field per input. A verdict alone would
    /// not be enough: the sanitising rules rewrite the field, so the value left behind is compared too.
    /// </summary>
    private static void CheckFields(string directory, Corpus corpus, Report report)
    {
        var baseMessage = BaseMessage(directory);

        var expected = Vectors.Read(Path.Combine(directory, "fields.jsonl"));
        var fields = ValTestMessage.Descriptor.Fields.InDeclarationOrder()
            .Where(f => f.FieldType == FieldType.String)
            .ToList();

        if (expected.Count != fields.Count * corpus.Total)
        {
            report.Add("fields.jsonl", "-",
                $"{fields.Count * corpus.Total} records", $"{expected.Count} records");
            return;
        }

        var compared = 0;
        var line = 0;
        foreach (var field in fields)
        {
            foreach (var (value, encoded) in corpus.Values.Zip(corpus.Encoded))
            {
                var want = expected[line++];
                if (value is null)
                {
                    continue;
                }

                var message = baseMessage.Clone();
                if (field.IsRepeated)
                {
                    var list = (IList)field.Accessor.GetValue(message);
                    list.Clear();
                    list.Add(value);
                }
                else
                {
                    field.Accessor.SetValue(message, value);
                }

                var error = message.Validate();
                var key = $"{field.Name} {encoded}";
                report.Compare("fields.jsonl", key, "field", want.String("field"), field.Name);
                report.Compare("fields.jsonl", key, "in", want.String("in"), Base64(value));
                report.Compare("fields.jsonl", key, "ok", want.Bool("ok"), error is null);
                report.Compare("fields.jsonl", key, "err", want.String("err"), error?.Message ?? "");

                var after = field.IsRepeated
                    ? ((IList)field.Accessor.GetValue(message)) is { Count: > 0 } list2 ? (string)list2[0]! : ""
                    : (string)field.Accessor.GetValue(message);
                report.Compare("fields.jsonl", key, "out", want.String("out"), Base64(after));

                compared++;
            }
        }

        report.Counted("fields.jsonl", compared);
    }

    /// <summary>
    /// Checks a string field one level down, where the failure is re-rooted under the field holding
    /// the message. The re-rooted text is its own shape, not the leaf text with a prefix.
    /// </summary>
    private static void CheckNestedFields(string directory, Corpus corpus, Report report)
    {
        var baseMessage = BaseMessage(directory);

        var expected = Vectors.Read(Path.Combine(directory, "nested.jsonl"));
        var pairs = new List<(FieldDescriptor Holder, FieldDescriptor Leaf)>();
        foreach (var holder in ValTestMessage.Descriptor.Fields.InDeclarationOrder())
        {
            if (holder.FieldType != FieldType.Message || holder.IsRepeated || holder.IsMap ||
                holder.Accessor.GetValue(baseMessage) is not IMessage)
            {
                continue;
            }

            foreach (var leaf in holder.MessageType.Fields.InDeclarationOrder())
            {
                if (leaf.FieldType == FieldType.String)
                {
                    pairs.Add((holder, leaf));
                }
            }
        }

        if (expected.Count != pairs.Count * corpus.Total)
        {
            report.Add("nested.jsonl", "-",
                $"{pairs.Count * corpus.Total} records", $"{expected.Count} records");
            return;
        }

        var compared = 0;
        var line = 0;
        foreach (var (holder, leaf) in pairs)
        {
            foreach (var (value, encoded) in corpus.Values.Zip(corpus.Encoded))
            {
                var want = expected[line++];
                if (value is null)
                {
                    continue;
                }

                var message = baseMessage.Clone();
                var inner = (IMessage)holder.Accessor.GetValue(message);
                if (leaf.IsRepeated)
                {
                    var list = (IList)leaf.Accessor.GetValue(inner);
                    list.Clear();
                    list.Add(value);
                }
                else
                {
                    leaf.Accessor.SetValue(inner, value);
                }

                var error = message.Validate();
                var key = $"{holder.Name}.{leaf.Name} {encoded}";
                report.Compare("nested.jsonl", key, "field", want.String("field"), $"{holder.Name}.{leaf.Name}");
                report.Compare("nested.jsonl", key, "in", want.String("in"), Base64(value));
                report.Compare("nested.jsonl", key, "ok", want.Bool("ok"), error is null);
                report.Compare("nested.jsonl", key, "err", want.String("err"), error?.Message ?? "");

                compared++;
            }
        }

        report.Counted("nested.jsonl", compared);
    }

    /// <summary>
    /// Checks a failure at each depth of a three-level message chain. Re-rooting accumulates, so a
    /// port can get one level right and still repeat the prefix or lose a segment further down.
    /// </summary>
    private static void CheckDeepNesting(string directory, Corpus corpus, Report report)
    {
        var baseMessage = MyReqMessage.Parser.ParseFrom(
            Convert.FromBase64String(File.ReadAllText(Path.Combine(directory, "deep_base.b64")).Trim()));

        var expected = Vectors.Read(Path.Combine(directory, "deep.jsonl"));
        var paths = StringPaths(baseMessage, []);

        if (expected.Count != paths.Count * corpus.Total)
        {
            report.Add("deep.jsonl", "-",
                $"{paths.Count * corpus.Total} records", $"{expected.Count} records");
            return;
        }

        var compared = 0;
        var line = 0;
        foreach (var path in paths)
        {
            foreach (var (value, encoded) in corpus.Values.Zip(corpus.Encoded))
            {
                var want = expected[line++];
                if (value is null)
                {
                    continue;
                }

                var message = baseMessage.Clone();
                IMessage holder = message;
                for (var i = 0; i < path.Count - 1; i++)
                {
                    holder = (IMessage)FieldByName(holder, path[i]).Accessor.GetValue(holder);
                }
                FieldByName(holder, path[^1]).Accessor.SetValue(holder, value);

                var error = message.Validate();
                var name = string.Join(".", path);
                var key = $"{name} {encoded}";
                report.Compare("deep.jsonl", key, "field", want.String("field"), name);
                report.Compare("deep.jsonl", key, "in", want.String("in"), Base64(value));
                report.Compare("deep.jsonl", key, "ok", want.Bool("ok"), error is null);
                report.Compare("deep.jsonl", key, "err", want.String("err"), error?.Message ?? "");

                compared++;
            }
        }

        report.Counted("deep.jsonl", compared);
    }

    /// <summary>
    /// Every string field reachable from the message, descending into the message fields it already
    /// carries.
    /// </summary>
    private static List<List<string>> StringPaths(IMessage message, List<string> prefix)
    {
        var paths = new List<List<string>>();
        foreach (var field in message.Descriptor.Fields.InDeclarationOrder())
        {
            if (field.IsRepeated || field.IsMap)
            {
                continue;
            }

            var here = new List<string>(prefix) { field.Name };
            if (field.FieldType == FieldType.String)
            {
                paths.Add(here);
            }
            else if (field.FieldType == FieldType.Message &&
                     field.Accessor.GetValue(message) is IMessage nested)
            {
                paths.AddRange(StringPaths(nested, here));
            }
        }

        return paths;
    }

    private static FieldDescriptor FieldByName(IMessage message, string name) =>
        message.Descriptor.Fields.InDeclarationOrder().First(f => f.Name == name);

    /// <summary>The message every field record starts from.</summary>
    private static ValTestMessage BaseMessage(string directory) =>
        ValTestMessage.Parser.ParseFrom(
            Convert.FromBase64String(File.ReadAllText(Path.Combine(directory, "base.b64")).Trim()));

    private static string Base64(string value) => Convert.ToBase64String(Encoding.UTF8.GetBytes(value));

    /// <summary>The vector directory recorded at build time, so the runner works from any directory.</summary>
    private static string DefaultDataDirectory()
    {
        var recorded = Assembly.GetExecutingAssembly()
            .GetCustomAttributes<AssemblyMetadataAttribute>()
            .First(a => a.Key == "ConformanceData").Value!;
        return Path.GetFullPath(recorded);
    }
}
