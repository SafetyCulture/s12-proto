// Copyright (c) 2026 SafetyCulture Pty Ltd. All Rights Reserved.

namespace S12.Protobuf.Validation.Conformance;

/// <summary>Collects the differences found and prints them once the run is over.</summary>
internal sealed class Report
{
    private const int MaxPrinted = 20;

    private readonly List<string> _differences = [];
    private readonly Dictionary<string, int> _counts = [];

    public void Compare(string file, string key, string column, string expected, string actual)
    {
        if (!string.Equals(expected, actual, StringComparison.Ordinal))
        {
            Add(file, $"{key} [{column}]", expected, actual);
        }
    }

    public void Compare(string file, string key, string column, bool expected, bool actual)
    {
        if (expected != actual)
        {
            Add(file, $"{key} [{column}]", expected.ToString(), actual.ToString());
        }
    }

    public void Add(string file, string key, string expected, string actual) =>
        _differences.Add($"{file}  {key}\n    go: {expected}\n    cs: {actual}");

    public void Counted(string file, int records) => _counts[file] = records;

    public int Summarise(Corpus corpus)
    {
        foreach (var (file, records) in _counts)
        {
            Console.WriteLine($"{file}: {records} records checked");
        }

        var skipped = corpus.Unrepresentable.ToList();
        if (skipped.Count > 0)
        {
            Console.WriteLine($"{skipped.Count} of {corpus.Total} inputs have no C# spelling and were stepped over:");
            foreach (var input in skipped)
            {
                Console.WriteLine($"    {input}");
            }
        }

        if (_differences.Count == 0)
        {
            Console.WriteLine("C# matches Go on every checked record.");
            return 0;
        }

        Console.Error.WriteLine($"\n{_differences.Count} differences from Go:");
        foreach (var difference in _differences.Take(MaxPrinted))
        {
            Console.Error.WriteLine(difference);
        }
        if (_differences.Count > MaxPrinted)
        {
            Console.Error.WriteLine($"... and {_differences.Count - MaxPrinted} more");
        }
        return 1;
    }
}
