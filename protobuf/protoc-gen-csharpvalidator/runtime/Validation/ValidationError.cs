// Copyright (c) 2026 SafetyCulture Pty Ltd. All Rights Reserved.

namespace S12.Protobuf.Validation;

/// <summary>Describes the first field of a message that failed validation.</summary>
public sealed class ValidationError
{
    /// <summary>Text of the innermost failure, which names the field it came from.</summary>
    private readonly string _failure;

    /// <summary>Dotted chain of the message fields walked to reach it, empty at the innermost.</summary>
    private readonly string _path;

    private ValidationError(string field, string failure, string path)
    {
        Field = field;
        _failure = failure;
        _path = path;
    }

    /// <summary>Dot-separated path to the field, from the root message being validated.</summary>
    public string Field { get; }

    /// <summary>Human-readable statement of the rule that was not met.</summary>
    public string Message => _path.Length == 0 ? _failure : $"invalid field {_path}: {_failure}";

    /// <summary>An error for a field that did not meet <paramref name="description"/>.</summary>
    public static ValidationError Create(string field, string description) =>
        new(field, $"{field}: {description}", "");

    /// <summary>An error for a field that carries no value.</summary>
    public static ValidationError Required(string field) =>
        new(field, $"field {field} is required", "");

    /// <summary>Returns the error re-rooted under <paramref name="parentField"/>.</summary>
    public ValidationError Nest(string parentField) =>
        new($"{parentField}.{Field}", _failure,
            _path.Length == 0 ? parentField : $"{parentField}.{_path}");

    public override string ToString() => Message;
}
