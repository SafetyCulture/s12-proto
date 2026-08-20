// Copyright (c) 2026 SafetyCulture Pty Ltd. All Rights Reserved.

namespace S12.Protobuf.Validation;

/// <summary>Describes the first field of a message that failed validation.</summary>
public sealed class ValidationError
{
    private ValidationError(string field, string description)
    {
        Field = field;
        Description = description;
    }

    /// <summary>Dot-separated path to the field, from the root message being validated.</summary>
    public string Field { get; }

    /// <summary>Human-readable statement of the rule that was not met.</summary>
    public string Description { get; }

    public static ValidationError Create(string field, string description) => new(field, description);

    /// <summary>Returns the error re-rooted under <paramref name="parentField"/>.</summary>
    public ValidationError Nest(string parentField) => new($"{parentField}.{Field}", Description);

    public override string ToString() => $"invalid field {Field}: {Description}";
}
