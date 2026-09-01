// Copyright (c) 2026 SafetyCulture Pty Ltd. All Rights Reserved.

using System.Text.Json;

namespace S12.Protobuf.Validation.Conformance;

/// <summary>One recorded result, read by column name.</summary>
internal readonly struct Vector(JsonElement element)
{
    /// <summary>The column, or the empty string when it was omitted as empty.</summary>
    public string String(string name) =>
        element.TryGetProperty(name, out var property) ? property.GetString() ?? "" : "";

    /// <summary>The column, or false when it was omitted as empty.</summary>
    public bool Bool(string name) =>
        element.TryGetProperty(name, out var property) && property.GetBoolean();
}

internal static class Vectors
{
    /// <summary>Reads a vector file, one record per line.</summary>
    public static IReadOnlyList<Vector> Read(string path)
    {
        var records = new List<Vector>();
        foreach (var line in File.ReadLines(path))
        {
            if (line.Length != 0)
            {
                records.Add(new Vector(JsonDocument.Parse(line).RootElement));
            }
        }
        return records;
    }
}
